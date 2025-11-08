package converter

import (
	"errors"

	"app-service/internal/dto"
	"app-service/internal/model"
	"github.com/jinzhu/copier"
)

func EducationModelToDTO(dataModel *model.Education) (*dto.Education, error) {
	dataDTO := &dto.Education{}
	err := copier.Copy(&dataDTO, &dataModel)
	if err != nil {
		return nil, errors.New("error converting Education model to dto")
	}
	return dataDTO, nil
}

func EducationDTOToModel(dataDTO *dto.Education) (*model.Education, error) {
	dataModel := &model.Education{}
	err := copier.CopyWithOption(&dataModel, &dataDTO, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errors.New("error converting education dto to model")
	}
	return dataModel, nil
}
