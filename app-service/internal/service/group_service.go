package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"app-service/internal/dto"
	"app-service/internal/kafka"
	"app-service/internal/model"
	"app-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IGroupService interface {
	GetGroupByID(groupID primitive.ObjectID) (*dto.Group, error)
	GetGroupsByOwnerID(ownerID primitive.ObjectID) ([]dto.Group, error)
	GetGroups() ([]dto.Group, error)
	CreateGroup(group *model.Group) (*dto.Group, error)
	UpdateGroup(groupID primitive.ObjectID, updatedGroup *dto.GroupUpdateRequest) (*dto.Group, error)
	DeleteGroup(groupID primitive.ObjectID, requesterID primitive.ObjectID) error
	NotifyGroupMembers(ctx context.Context, groupID primitive.ObjectID, subject, message string) error
}

type GroupService struct {
	groupRepository    repository.IGroupRepository
	memberRepository   repository.IMemberRepository
	bulletinRepository repository.IBulletinRepository
	producer           kafka.Producer
}

var (
	ErrGroupNotFound         = errors.New("group not found")
	ErrMessageBodyEmpty      = errors.New("message is required")
	ErrProducerNotConfigured = errors.New("notification producer is not configured")
)

const hardcodedNotificationEmail = "dalai2547@gmail.com"

func NewGroupService(groupRepo repository.IGroupRepository, memberRepo repository.IMemberRepository, bulletinRepo repository.IBulletinRepository, producer kafka.Producer) IGroupService {
	return GroupService{
		groupRepository:    groupRepo,
		memberRepository:   memberRepo,
		bulletinRepository: bulletinRepo,
		producer:           producer,
	}
}

func (s GroupService) CreateGroup(group *model.Group) (*dto.Group, error) {
	// You may not need to hash passwords for groups, so you can remove that part.
	newGroup, err := s.groupRepository.CreateGroup(group)
	if err != nil {
		return nil, err
	}
	if newGroup != nil {
		subject := fmt.Sprintf("Group created: %s", strings.TrimSpace(newGroup.Title))
		message := fmt.Sprintf("A new group \"%s\" has been created.\n\nDescription:\n%s", strings.TrimSpace(newGroup.Title), strings.TrimSpace(newGroup.Description))
		s.dispatchGroupNotification(newGroup.GroupID, subject, message)
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
	if updatedGroupDTO != nil {
		subject := fmt.Sprintf("Group updated: %s", strings.TrimSpace(updatedGroupDTO.Title))
		message := fmt.Sprintf("Group \"%s\" has been updated. Visit the group page for the latest details.", strings.TrimSpace(updatedGroupDTO.Title))
		s.dispatchGroupNotification(updatedGroupDTO.GroupID, subject, message)
	}
	return updatedGroupDTO, nil
}

func (s GroupService) DeleteGroup(groupID primitive.ObjectID, requesterID primitive.ObjectID) error {
	_, err := s.groupRepository.GetGroupByID(groupID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrGroupNotFound
		}
		return err
	}

	if err := s.groupRepository.DeleteGroup(groupID); err != nil {
		return err
	}

	if s.bulletinRepository != nil {
		if _, bulletinErr := s.bulletinRepository.DeleteBulletinsByGroupID(groupID); bulletinErr != nil {
			return bulletinErr
		}
	}

	return nil
}

func (s GroupService) NotifyGroupMembers(ctx context.Context, groupID primitive.ObjectID, subject, message string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if s.producer == nil {
		return ErrProducerNotConfigured
	}

	message = strings.TrimSpace(message)
	if message == "" {
		return ErrMessageBodyEmpty
	}

	group, err := s.groupRepository.GetGroupByID(groupID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrGroupNotFound
		}
		return fmt.Errorf("fetch group: %w", err)
	}

	defaultSubject := strings.TrimSpace(subject)
	if defaultSubject == "" {
		defaultSubject = fmt.Sprintf("Update from %s", group.Title)
	}
	emailBody := fmt.Sprintf("%s\n\n— Sent from group: %s", message, group.Title)

	notification := kafka.NotificationMessage{
		Type:    kafka.NotificationTypeEmail,
		To:      hardcodedNotificationEmail,
		Subject: defaultSubject,
		Message: emailBody,
	}

	if err := s.producer.Publish(ctx, []kafka.NotificationMessage{notification}); err != nil {
		return fmt.Errorf("publish notifications: %w", err)
	}

	return nil
}

func (s GroupService) dispatchGroupNotification(groupID primitive.ObjectID, subject, message string) {
	if s.producer == nil || groupID == primitive.NilObjectID {
		return
	}

	subject = strings.TrimSpace(subject)
	message = strings.TrimSpace(message)
	if message == "" {
		return
	}

	go func() {
		err := s.NotifyGroupMembers(context.Background(), groupID, subject, message)
		if err != nil {
			if errors.Is(err, ErrMessageBodyEmpty) {
				log.Printf("group notification skipped for %s: %v", groupID.Hex(), err)
				return
			}
			log.Printf("failed to notify group %s: %v", groupID.Hex(), err)
		}
	}()
}
