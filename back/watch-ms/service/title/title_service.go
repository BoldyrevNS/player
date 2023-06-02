package title

import (
	"watch-ms/DTO"
	"watch-ms/model"
	"watch-ms/provider/title"
)

type Service interface {
	CreateNewTitle(data DTO.CreateTitleDTO) error
}

type serviceImpl struct {
	titleProvider title.Provider
}

func NewTitleService(titleProvider title.Provider) Service {
	return &serviceImpl{
		titleProvider: titleProvider,
	}
}

func (s *serviceImpl) CreateNewTitle(data DTO.CreateTitleDTO) error {
	err := s.titleProvider.Create(model.Title{
		CategoryId: data.CategoryId,
		Name:       data.Name,
	})
	if err != nil {
		return err
	}
	return nil
}
