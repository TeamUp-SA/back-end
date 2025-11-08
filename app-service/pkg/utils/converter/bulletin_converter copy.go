package converter

import (
	"errors"

	"app-service/internal/dto"
	"app-service/internal/model"

	"github.com/jinzhu/copier"
)

func BulletinModelToDTO(dataModel *model.Bulletin) (*dto.Bulletin, error) {
	dataDTO := &dto.Bulletin{}
	err := copier.Copy(&dataDTO, &dataModel)
	if err != nil {
		return nil, errors.New("error converting Bulletin model to dto")
	}
	return dataDTO, nil
}

func BulletinDTOToModel(dataDTO *dto.Bulletin) (*model.Bulletin, error) {
	dataModel := &model.Bulletin{}
	err := copier.CopyWithOption(&dataModel, &dataDTO, copier.Option{DeepCopy: true})
	if err != nil {
		return nil, errors.New("error converting bulletin dto to model")
	}
	return dataModel, nil
}
