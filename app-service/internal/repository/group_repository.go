package repository

import (
	"context"
	"fmt"
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
	UpdateGroup(groupID primitive.ObjectID, updatedGroup *dto.GroupUpdateRequest) (*dto.Group, error)
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

func (r GroupRepository) UpdateGroup(groupID primitive.ObjectID, req *dto.GroupUpdateRequest) (*dto.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	updateFields := bson.M{}
	if req.Title != nil {
		updateFields["title"] = *req.Title
	}
	if req.Description != nil {
		updateFields["description"] = *req.Description
	}
	if req.Members != nil {
		updateFields["members"] = *req.Members
	}
	if req.Tags != nil {
		updateFields["tags"] = *req.Tags
	}
	if req.Closed != nil {
		updateFields["closed"] = *req.Closed
	}
	if req.Date != nil {
		updateFields["date"] = *req.Date
	}

	if len(updateFields) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	filter := bson.M{"_id": groupID}
	update := bson.M{"$set": updateFields}

	_, err := r.groupCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Fetch updated document
	var updated model.Group
	if err := r.groupCollection.FindOne(ctx, filter).Decode(&updated); err != nil {
		return nil, err
	}

	return converter.GroupModelToDTO(&updated)
}

func (r *GroupRepository) DeleteGroup(groupID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := r.groupCollection.DeleteOne(ctx, bson.M{"_id": groupID})
	return err
}
