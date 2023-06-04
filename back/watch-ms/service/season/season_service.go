package season

import (
	"watch-ms/DTO"
	"watch-ms/model"
	"watch-ms/provider/season"
)

type Service interface {
	CreateSeason(data DTO.CreateSeasonDTO) error
	GetAllTitleSeasons(data DTO.GetTitleSeasonsDTO) ([]DTO.TitleSeasonDTO, error)
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

func (s *serviceImpl) GetAllTitleSeasons(data DTO.GetTitleSeasonsDTO) ([]DTO.TitleSeasonDTO, error) {
	var titleSeasons []DTO.TitleSeasonDTO
	seasons, err := s.seasonProvider.FindTitleSeasons(data.TitleId)
	if err != nil {
		return nil, err
	}
	for _, s := range seasons {
		titleSeasons = append(titleSeasons, DTO.TitleSeasonDTO{
			SeasonId: s.Id,
			Number:   s.Number,
		})
	}
	return titleSeasons, err
}
