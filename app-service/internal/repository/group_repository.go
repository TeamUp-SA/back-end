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

type IGroupRepository interface {
	GetGroupByID(groupID primitive.ObjectID) (*dto.Group, error)
	GetGroupsByOwnerID(ownerID primitive.ObjectID) ([]dto.Group, error)
	GetGroups() ([]dto.Group, error)
	CreateGroup(group *model.Group) (*dto.Group, error)
	UpdateGroup(groupID primitive.ObjectID, updatedGroup *model.Group) (*dto.Group, error)
	DeleteGroup(groupID primitive.ObjectID) error
}

type GroupRepository struct {
	groupCollection *mongo.Collection
}

func NewGroupRepository(db *mongo.Database, collectionName string) IGroupRepository {
	return &GroupRepository{
		groupCollection: db.Collection(collectionName),
	}
}

func (r *GroupRepository) GetGroupByID(groupID primitive.ObjectID) (*dto.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var group *model.Group

	err := r.groupCollection.FindOne(ctx, bson.M{"_id": groupID}).Decode(&group)
	if err != nil {
		return nil, err
	}
	return converter.GroupModelToDTO(group)
}

func (r GroupRepository) GetGroupsByOwnerID(ownerID primitive.ObjectID) ([]dto.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var groupList []dto.Group

	dataList, err := r.groupCollection.Find(ctx, bson.M{"owner_id": ownerID})
	if err != nil {
		return nil, err
	}
	defer dataList.Close(ctx)
	for dataList.Next(ctx) {
		var groupModel *model.Group
		if err = dataList.Decode(&groupModel); err != nil {
			return nil, err
		}
		groupDTO, groupErr := converter.GroupModelToDTO(groupModel)
		if groupErr != nil {
			return nil, err
		}
		groupList = append(groupList, *groupDTO)
	}

	return groupList, nil
}

func (r *GroupRepository) GetGroups() ([]dto.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var groupList []dto.Group

	dataList, err := r.groupCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer dataList.Close(ctx)

	for dataList.Next(ctx) {
		var groupModel *model.Group
		if err = dataList.Decode(&groupModel); err != nil {
			return nil, err
		}
		groupDTO, groupErr := converter.GroupModelToDTO(groupModel)
		if groupErr != nil {
			return nil, err
		}
		groupList = append(groupList, *groupDTO)
	}

	return groupList, nil
}

func (r *GroupRepository) CreateGroup(group *model.Group) (*dto.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	group.GroupID = primitive.NewObjectID()
	result, err := r.groupCollection.InsertOne(ctx, group)
	if err != nil {
		return nil, err
	}
	var newGroup *model.Group
	err = r.groupCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&newGroup)

	if err != nil {
		return nil, err
	}

	return converter.GroupModelToDTO(newGroup)
}

func (r *GroupRepository) UpdateGroup(groupID primitive.ObjectID, updatedGroup *model.Group) (*dto.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	updateData, err := bson.Marshal(updatedGroup)
	if err != nil {
		return nil, err
	}

	var update bson.M
	if err := bson.Unmarshal(updateData, &update); err != nil {
		return nil, err
	}

	update = bson.M{"$set": update}

	filter := bson.M{"_id": groupID}
	_, err = r.groupCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	var newUpdatedGroup *model.Group
	err = r.groupCollection.FindOne(ctx, filter).Decode(&newUpdatedGroup)

	if err != nil {
		return nil, err
	}

	return converter.GroupModelToDTO(newUpdatedGroup)
}

func (r *GroupRepository) DeleteGroup(groupID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := r.groupCollection.DeleteOne(ctx, bson.M{"_id": groupID})
	return err
}
