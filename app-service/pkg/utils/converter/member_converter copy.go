package converter

import (
	"errors"

	"app-service/internal/dto"
	"app-service/internal/model"
	"github.com/jinzhu/copier"
)

func MemberModelToDTO(dataModel *model.Member) (*dto.Member, error) {
	dataDTO := &dto.Member{}
	err := copier.Copy(&dataDTO, &dataModel)
	if err != nil {
		return nil, errors.New("error converting Member model to dto")
	}
	return dataDTO, nil
}

func MemberDTOToModel(dataDTO *dto.Member) (*model.Member, error) {
	dataModel := &model.Member{}
	err := copier.CopyWithOption(&dataModel, &dataDTO, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errors.New("error converting member dto to model")
	}
	return dataModel, nil
}
