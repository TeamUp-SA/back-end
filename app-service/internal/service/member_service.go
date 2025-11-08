package service

import (
	"mime/multipart"

	"app-service/internal/dto"
	"app-service/internal/model"
	"app-service/internal/repository"
	"app-service/internal/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMemberService interface {
	GetMemberByID(memberID primitive.ObjectID) (*dto.Member, error)
	GetMember() ([]dto.Member, error)
	CreateMemberData(member *model.Member) (*dto.Member, error)
	UpdateMemberData(memberID primitive.ObjectID, updatedMember *model.Member, imageFile multipart.File, fileHeader *multipart.FileHeader) (*dto.Member, error)
}

type MemberService struct {
	memberRepository repository.IMemberRepository
}

func NewMemberService(r repository.IMemberRepository) IMemberService {
	return MemberService{
		memberRepository: r,
	}
}

func (s MemberService) CreateMemberData(member *model.Member) (*dto.Member, error) {
	newMember, err := s.memberRepository.CreateMemberData(member)

	if err != nil {
		return nil, err
	}

	return newMember, nil
}

func (s MemberService) GetMemberByID(memberID primitive.ObjectID) (*dto.Member, error) {
	memberDTO, err := s.memberRepository.GetMemberByID(memberID)
	if err != nil {
		return nil, err
	}
	return memberDTO, nil
}

func (s MemberService) GetMember() ([]dto.Member, error) {
	members, err := s.memberRepository.GetMember()
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (s MemberService) UpdateMemberData(memberID primitive.ObjectID, updatedMember *model.Member, imageFile multipart.File, fileHeader *multipart.FileHeader) (*dto.Member, error) {
	if imageFile != nil && fileHeader != nil {
		s3URL, err := utils.UploadFileToS3(imageFile, fileHeader)
		if err != nil {
			return nil, err
		}
		updatedMember.ProfileImage = s3URL
	}

	updatedMemberDTO, err := s.memberRepository.UpdateMemberData(memberID, updatedMember)
	if err != nil {
		return nil, err
	}

	return updatedMemberDTO, nil
}
