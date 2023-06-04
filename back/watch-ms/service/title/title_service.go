package title

import (
	"watch-ms/DTO"
	"watch-ms/model"
	"watch-ms/provider/title"
)

type Service interface {
	CreateNewTitle(data DTO.CreateTitleDTO) error
	GetTitlesExcludeWatched(userId uint) ([]DTO.TitleDTO, error)
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

func (s *serviceImpl) GetTitlesExcludeWatched(userId uint) ([]DTO.TitleDTO, error) {
	var titles []DTO.TitleDTO
	titlesRaw, err := s.titleProvider.GetTitlesExcludeWatched(userId)
	if err != nil {
		return nil, err
	}
	for _, t := range titlesRaw {
		titles = append(titles, DTO.TitleDTO{
			Id:           t.Id,
			Name:         t.Name,
			ThumbnailUrl: t.ThumbnailUrl,
			CategoryName: t.CategoryName,
			CategoryId:   t.CategoryId,
		})
	}
	return titles, nil
}
