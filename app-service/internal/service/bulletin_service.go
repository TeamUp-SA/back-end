package service

import (
	"mime/multipart"

	"app-service/internal/dto"
	"app-service/internal/model"
	"app-service/internal/repository"
	"app-service/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IBulletinService interface {
	GetBulletins() ([]dto.Bulletin, error)
	GetBulletinByID(bulletinID primitive.ObjectID) (*dto.Bulletin, error)
	GetBulletinsByAuthorID(authorID primitive.ObjectID) ([]dto.Bulletin, error)
	GetBulletinsByGroupID(groupID primitive.ObjectID) ([]dto.Bulletin, error)
	CreateBulletin(bulletin *model.Bulletin, imageFile multipart.File, fileHeader *multipart.FileHeader) (*dto.Bulletin, error)
	UpdateBulletin(bulletinID primitive.ObjectID, updatedBulletin *dto.BulletinUpdateRequest, imageFile multipart.File, fileHeader *multipart.FileHeader) (*dto.Bulletin, error)
	DeleteBulletin(bulletinID primitive.ObjectID) error
}

type BulletinService struct {
	bulletinRepository repository.IBulletinRepository
}

func NewBulletinService(r repository.IBulletinRepository) IBulletinService {
	return BulletinService{
		bulletinRepository: r,
	}
}

func (s BulletinService) GetBulletins() ([]dto.Bulletin, error) {
	bulletins, err := s.bulletinRepository.GetBulletins()
	if err != nil {
		return nil, err
	}
	return bulletins, nil
}

func (s BulletinService) GetBulletinByID(bulletinID primitive.ObjectID) (*dto.Bulletin, error) {
	bulletinDTO, err := s.bulletinRepository.GetBulletinByID(bulletinID)
	if err != nil {
		return nil, err
	}
	return bulletinDTO, nil
}

func (s BulletinService) GetBulletinsByAuthorID(authorID primitive.ObjectID) ([]dto.Bulletin, error) {
	bulletins, err := s.bulletinRepository.GetBulletinsByAuthorID(authorID)
	if err != nil {
		return nil, err
	}
	return bulletins, nil
}

func (s BulletinService) GetBulletinsByGroupID(groupID primitive.ObjectID) ([]dto.Bulletin, error) {
	bulletins, err := s.bulletinRepository.GetBulletinsByGroupID(groupID)
	if err != nil {
		return nil, err
	}
	return bulletins, nil
}

func (s BulletinService) CreateBulletin(bulletin *model.Bulletin, imageFile multipart.File, fileHeader *multipart.FileHeader) (*dto.Bulletin, error) {
	if imageFile != nil && fileHeader != nil {
		s3URL, err := utils.UploadFileToS3(imageFile, fileHeader)
		if err != nil {
			return nil, err
		}
		bulletin.Image = s3URL
	}

	newBulletin, err := s.bulletinRepository.CreateBulletin(bulletin)
	if err != nil {
		return nil, err
	}

	return newBulletin, nil
}

func (s BulletinService) UpdateBulletin(bulletinID primitive.ObjectID, updatedBulletin *dto.BulletinUpdateRequest, imageFile multipart.File, fileHeader *multipart.FileHeader) (*dto.Bulletin, error) {
	if imageFile != nil && fileHeader != nil {
		s3URL, err := utils.UploadFileToS3(imageFile, fileHeader)
		if err != nil {
			return nil, err
		}
		updatedBulletin.Image = &s3URL
	}

	updatedDTO, err := s.bulletinRepository.UpdateBulletin(bulletinID, updatedBulletin)
	if err != nil {
		return nil, err
	}

	return updatedDTO, nil
}

func (s BulletinService) DeleteBulletin(bulletinID primitive.ObjectID) error {

	err := s.bulletinRepository.DeleteBulletin(bulletinID)
	if err != nil {
		return err
	}

	return nil
}
