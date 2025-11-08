package service

import (
	"app-service/internal/dto"
	"app-service/internal/model"
	"app-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IGroupService interface {
	GetGroupByID(groupID primitive.ObjectID) (*dto.Group, error)
	GetGroupsByOwnerID(ownerID primitive.ObjectID) ([]dto.Group, error)
	GetGroups() ([]dto.Group, error)
	CreateGroup(group *model.Group) (*dto.Group, error)
	UpdateGroup(groupID primitive.ObjectID, updatedGroup *dto.GroupUpdateRequest) (*dto.Group, error)
	DeleteGroup(groupID primitive.ObjectID) error
}

type GroupService struct {
	groupRepository repository.IGroupRepository
}

func NewGroupService(r repository.IGroupRepository) IGroupService {
	return GroupService{
		groupRepository: r,
	}
}

func (s GroupService) CreateGroup(group *model.Group) (*dto.Group, error) {
	// You may not need to hash passwords for groups, so you can remove that part.
	newGroup, err := s.groupRepository.CreateGroup(group)
	if err != nil {
		return nil, err
	}
	return newGroup, nil
}

func (s GroupService) GetGroupByID(groupID primitive.ObjectID) (*dto.Group, error) {
	groupDTO, err := s.groupRepository.GetGroupByID(groupID)
	if err != nil {
		return nil, err
	}
	return groupDTO, nil
}

func (s GroupService) GetGroupsByOwnerID(ownerID primitive.ObjectID) ([]dto.Group, error) {
	groups, err := s.groupRepository.GetGroupsByOwnerID(ownerID)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (s GroupService) GetGroups() ([]dto.Group, error) {
	groups, err := s.groupRepository.GetGroups()
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (s GroupService) UpdateGroup(groupID primitive.ObjectID, updatedGroup *dto.GroupUpdateRequest) (*dto.Group, error) {
	updatedGroupDTO, err := s.groupRepository.UpdateGroup(groupID, updatedGroup)
	if err != nil {
		return nil, err
	}
	return updatedGroupDTO, nil
}

func (s GroupService) DeleteGroup(groupID primitive.ObjectID) error {

	err := s.groupRepository.DeleteGroup(groupID)
	if err != nil {
		return err
	}

	return nil
}
