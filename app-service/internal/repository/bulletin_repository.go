package repository

import (
	"context"
	"time"

	"github.com/Ntchah/TeamUp-application-service/internal/dto"
	"github.com/Ntchah/TeamUp-application-service/internal/model"
	"github.com/Ntchah/TeamUp-application-service/pkg/utils/converter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IBulletinRepository interface {
	GetBulletins() ([]dto.Bulletin, error)
	GetBulletinByID(bulletinID primitive.ObjectID) (*dto.Bulletin, error)
	GetBulletinsByAuthorID(authorID primitive.ObjectID) ([]dto.Bulletin, error)
	GetBulletinsByGroupID(groupID primitive.ObjectID) ([]dto.Bulletin, error)
	CreateBulletin(bulletin *model.Bulletin) (*dto.Bulletin, error)
	UpdateBulletin(bulletinID primitive.ObjectID, updatedBulletin *model.Bulletin) (*dto.Bulletin, error)
	DeleteBulletin(bulletinID primitive.ObjectID) error
}

type BulletinRepository struct {
	bulletinCollection *mongo.Collection
}

func NewBulletinRepository(db *mongo.Database, bulletincollectionName string) IBulletinRepository {
	return BulletinRepository{
		bulletinCollection: db.Collection(bulletincollectionName),
	}
}

func (r BulletinRepository) GetBulletins() ([]dto.Bulletin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var bulletinList []dto.Bulletin

	dataList, err := r.bulletinCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer dataList.Close(ctx)
	for dataList.Next(ctx) {
		var bulletinModel *model.Bulletin
		if err = dataList.Decode(&bulletinModel); err != nil {
			return nil, err
		}
		bulletinDTO, bulletinErr := converter.BulletinModelToDTO(bulletinModel)
		if bulletinErr != nil {
			return nil, err
		}
		bulletinList = append(bulletinList, *bulletinDTO)
	}

	return bulletinList, nil
}

func (r BulletinRepository) GetBulletinByID(bulletinID primitive.ObjectID) (*dto.Bulletin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var bulletin *model.Bulletin

	err := r.bulletinCollection.FindOne(ctx, bson.M{"_id": bulletinID}).Decode(&bulletin)
	if err != nil {
		return nil, err
	}
	return converter.BulletinModelToDTO(bulletin)
}

func (r BulletinRepository) GetBulletinsByAuthorID(authorID primitive.ObjectID) ([]dto.Bulletin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var bulletinList []dto.Bulletin

	dataList, err := r.bulletinCollection.Find(ctx, bson.M{"author_id": authorID})
	if err != nil {
		return nil, err
	}
	defer dataList.Close(ctx)
	for dataList.Next(ctx) {
		var bulletinModel *model.Bulletin
		if err = dataList.Decode(&bulletinModel); err != nil {
			return nil, err
		}
		bulletinDTO, bulletinErr := converter.BulletinModelToDTO(bulletinModel)
		if bulletinErr != nil {
			return nil, err
		}
		bulletinList = append(bulletinList, *bulletinDTO)
	}

	return bulletinList, nil
}

func (r BulletinRepository) GetBulletinsByGroupID(groupID primitive.ObjectID) ([]dto.Bulletin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var bulletinList []dto.Bulletin

	dataList, err := r.bulletinCollection.Find(ctx, bson.M{"group_id": groupID})
	if err != nil {
		return nil, err
	}
	defer dataList.Close(ctx)
	for dataList.Next(ctx) {
		var bulletinModel *model.Bulletin
		if err = dataList.Decode(&bulletinModel); err != nil {
			return nil, err
		}
		bulletinDTO, bulletinErr := converter.BulletinModelToDTO(bulletinModel)
		if bulletinErr != nil {
			return nil, err
		}
		bulletinList = append(bulletinList, *bulletinDTO)
	}

	return bulletinList, nil
}

func (r BulletinRepository) CreateBulletin(bulletin *model.Bulletin) (*dto.Bulletin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	bulletin.BulletinID = primitive.NewObjectID()
	bulletin.CreatedAt = time.Now()
	result, err := r.bulletinCollection.InsertOne(ctx, bulletin)
	if err != nil {
		return nil, err
	}
	var newBulletin *model.Bulletin
	err = r.bulletinCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&newBulletin)
	if err != nil {
		return nil, err
	}
	return converter.BulletinModelToDTO(newBulletin)
}

func (r BulletinRepository) UpdateBulletin(bulletinID primitive.ObjectID, updatedBulletin *model.Bulletin) (*dto.Bulletin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title": updatedBulletin.Title,
			"description":  updatedBulletin.Description,
			"group_id":  updatedBulletin.GroupID,
			"date": updatedBulletin.Date,
			"image": updatedBulletin.Image,
			"tags": updatedBulletin.Tags,
			"createdAt": time.Now(),
		},
	}

	filter := bson.M{"_id": bulletinID}
	_, err := r.bulletinCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	var newUpdatedBulletin *model.Bulletin
	err = r.bulletinCollection.FindOne(ctx, filter).Decode(&newUpdatedBulletin)
	if err != nil {
		return nil, err
	}

	return converter.BulletinModelToDTO(newUpdatedBulletin)
}

func (r BulletinRepository) DeleteBulletin(bulletinID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := r.bulletinCollection.DeleteOne(ctx, bson.M{"_id": bulletinID})
	return err
}
