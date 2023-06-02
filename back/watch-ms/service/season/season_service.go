package season

import (
	"watch-ms/DTO"
	"watch-ms/model"
	"watch-ms/provider/season"
)

type Service interface {
	CreateSeason(data DTO.CreateSeasonDTO) error
}

type serviceImpl struct {
	seasonProvider season.Provider
}

func NewSeasonService(seasonProvider season.Provider) Service {
	return &serviceImpl{
		seasonProvider: seasonProvider,
	}
}

func (s *serviceImpl) CreateSeason(data DTO.CreateSeasonDTO) error {
	err := s.seasonProvider.Create(model.Season{
		TitleId: data.TitleId,
		Number:  data.Number,
	})
	if err != nil {
		return err
	}
	return nil
}
