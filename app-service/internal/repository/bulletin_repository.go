package repository

import (
	"context"
	"fmt"
	"time"

	"app-service/internal/dto"
	"app-service/internal/model"
	"app-service/pkg/utils/converter"

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
	UpdateBulletin(bulletinID primitive.ObjectID, updatedBulletin *dto.BulletinUpdateRequest) (*dto.Bulletin, error)
	DeleteBulletin(bulletinID primitive.ObjectID) error
	DeleteBulletinsByGroupID(groupID primitive.ObjectID) (int64, error)
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

func (r BulletinRepository) UpdateBulletin(bulletinID primitive.ObjectID, req *dto.BulletinUpdateRequest) (*dto.Bulletin, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	updateFields := bson.M{}
	if req.Title != nil {
		updateFields["title"] = *req.Title
	}
	if req.Description != nil {
		updateFields["description"] = *req.Description
	}
	if req.GroupID != nil {
		updateFields["group_id"] = *req.GroupID
	}
	if req.Date != nil {
		updateFields["date"] = *req.Date
	}
	if req.Image != nil {
		updateFields["image"] = *req.Image
	}
	if req.Tags != nil {
		updateFields["tags"] = *req.Tags
	}

	if len(updateFields) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	update := bson.M{"$set": updateFields}
	filter := bson.M{"_id": bulletinID}

	_, err := r.bulletinCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	var updated model.Bulletin
	if err := r.bulletinCollection.FindOne(ctx, filter).Decode(&updated); err != nil {
		return nil, err
	}

	return converter.BulletinModelToDTO(&updated)
}

func (r BulletinRepository) DeleteBulletin(bulletinID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := r.bulletinCollection.DeleteOne(ctx, bson.M{"_id": bulletinID})
	return err
}

func (r BulletinRepository) DeleteBulletinsByGroupID(groupID primitive.ObjectID) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := r.bulletinCollection.DeleteMany(ctx, bson.M{"group_id": groupID})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
