package converter

import (
	"errors"

	"github.com/Ntchah/TeamUp-application-service/internal/dto"
	"github.com/Ntchah/TeamUp-application-service/internal/model"
	"github.com/jinzhu/copier"
)

func ExperienceModelToDTO(dataModel *model.Experience) (*dto.Experience, error) {
	dataDTO := &dto.Experience{}
	err := copier.Copy(&dataDTO, &dataModel)
	if err != nil {
		return nil, errors.New("error converting Experience model to dto")
	}
	return dataDTO, nil
}

func ExperienceDTOToModel(dataDTO *dto.Experience) (*model.Experience, error) {
	dataModel := &model.Experience{}
	err := copier.CopyWithOption(&dataModel, &dataDTO, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errors.New("error converting experience dto to model")
	}
	return dataModel, nil
}
