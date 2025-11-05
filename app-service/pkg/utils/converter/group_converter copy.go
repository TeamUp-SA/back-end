package converter

import (
	"errors"

	"github.com/Ntchah/TeamUp-application-service/internal/dto"
	"github.com/Ntchah/TeamUp-application-service/internal/model"
	"github.com/jinzhu/copier"
)

func GroupModelToDTO(dataModel *model.Group) (*dto.Group, error) {
	dataDTO := &dto.Group{}
	err := copier.Copy(&dataDTO, &dataModel)
	if err != nil {
		return nil, errors.New("error converting Group model to dto")
	}
	return dataDTO, nil
}

func GroupDTOToModel(dataDTO *dto.Group) (*model.Group, error) {
	dataModel := &model.Group{}
	err := copier.CopyWithOption(&dataModel, &dataDTO, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errors.New("error converting group dto to model")
	}
	return dataModel, nil
}
