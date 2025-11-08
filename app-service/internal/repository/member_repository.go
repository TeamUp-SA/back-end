package repository

import (
	"context"
	"time"

	"app-service/internal/dto"
	"app-service/internal/model"
	"app-service/pkg/utils/converter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IMemberRepository interface {
	GetMember() ([]dto.Member, error)
	GetMemberByID(memberID primitive.ObjectID) (*dto.Member, error)
	GetMembersByIDs(memberIDs []primitive.ObjectID) ([]dto.Member, error)
	CreateMemberData(member *model.Member) (*dto.Member, error)
	GetMemberByUsername(username string) (*model.Member, error)
	UpdateMemberData(memberID primitive.ObjectID, updatedMember *model.Member) (*dto.Member, error)
}

type MemberRepository struct {
	memberCollection *mongo.Collection
}

func NewMemberRepository(db *mongo.Database, collectionName string) IMemberRepository {
	return &MemberRepository{
		memberCollection: db.Collection(collectionName),
	}
}

func (r *MemberRepository) GetMemberByID(memberID primitive.ObjectID) (*dto.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var member *model.Member

	err := r.memberCollection.FindOne(ctx, bson.M{"_id": memberID}).Decode(&member)
	if err != nil {
		return nil, err
	}
	return converter.MemberModelToDTO(member)
}

func (r *MemberRepository) GetMember() ([]dto.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var memberList []dto.Member

	dataList, err := r.memberCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer dataList.Close(ctx)

	for dataList.Next(ctx) {
		var memberModel *model.Member
		if err = dataList.Decode(&memberModel); err != nil {
			return nil, err
		}
		memberDTO, memberErr := converter.MemberModelToDTO(memberModel)
		if memberErr != nil {
			return nil, err
		}
		memberList = append(memberList, *memberDTO)
	}

	return memberList, nil
}

func (r *MemberRepository) GetMembersByIDs(memberIDs []primitive.ObjectID) ([]dto.Member, error) {
	if len(memberIDs) == 0 {
		return []dto.Member{}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": memberIDs}}
	cursor, err := r.memberCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	members := make([]dto.Member, 0, len(memberIDs))
	for cursor.Next(ctx) {
		var memberModel model.Member
		if err := cursor.Decode(&memberModel); err != nil {
			return nil, err
		}

		memberDTO, convertErr := converter.MemberModelToDTO(&memberModel)
		if convertErr != nil {
			return nil, convertErr
		}

		members = append(members, *memberDTO)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *MemberRepository) GetMemberByUsername(username string) (*model.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var member *model.Member

	err := r.memberCollection.FindOne(ctx, bson.M{"username": username}).Decode(&member)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func (r *MemberRepository) CreateMemberData(member *model.Member) (*dto.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	member.MemberID = primitive.NewObjectID()
	// member.Cart = []model.OrderProduct{}
	result, err := r.memberCollection.InsertOne(ctx, member)
	if err != nil {
		return nil, err
	}
	var newMember *model.Member
	err = r.memberCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&newMember)

	if err != nil {
		return nil, err
	}

	return converter.MemberModelToDTO(newMember)
}

func (r *MemberRepository) UpdateMemberData(memberID primitive.ObjectID, updatedMember *model.Member) (*dto.Member, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	data, err := bson.Marshal(updatedMember)
	if err != nil {
		return nil, err
	}
	var update bson.M
	err = bson.Unmarshal(data, &update)
	if err != nil {
		return nil, err
	}
	for key, value := range update {
		if value == "" || value == nil || key == "_id" {
			delete(update, key)
		}
	}

	filter := bson.M{"_id": memberID}
	_, err = r.memberCollection.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}

	var newUpdatedMember *model.Member
	err = r.memberCollection.FindOne(ctx, filter).Decode(&newUpdatedMember)

	if err != nil {
		return nil, err
	}

	return converter.MemberModelToDTO(newUpdatedMember)
}
