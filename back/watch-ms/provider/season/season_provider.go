package season

import (
	"gorm.io/gorm"
	"watch-ms/model"
)

type Provider interface {
	FindAll() ([]model.Season, error)
	FindTitleSeasons(titleId uint) ([]model.Season, error)
	Create(season model.Season) error
	Delete(seasonId uint) error
}

type providerImpl struct {
	db *gorm.DB
}

func NewSeasonProvider(db *gorm.DB) Provider {
	return &providerImpl{
		db: db,
	}
}

func (p *providerImpl) FindAll() ([]model.Season, error) {
	var seasons []model.Season
	res := p.db.Find(&seasons)
	return seasons, res.Error
}

func (p *providerImpl) FindTitleSeasons(titleId uint) ([]model.Season, error) {
	var seasons []model.Season
	res := p.db.Where("title_id = ?", titleId).Find(&seasons)
	return seasons, res.Error
}

func (p *providerImpl) Create(season model.Season) error {
	res := p.db.Create(&season)
	return res.Error
}

func (p *providerImpl) Delete(seasonId uint) error {
	var season model.Season
	res := p.db.Where("id = ?", seasonId).Delete(&season)
	return res.Error
}
