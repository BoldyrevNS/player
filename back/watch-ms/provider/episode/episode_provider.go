package episode

import (
	"gorm.io/gorm"
	"watch-ms/model"
)

type Provider interface {
	FindAll() ([]model.Episode, error)
	FindBySeasonId(seasonId uint) ([]model.Episode, error)
	FindById(episodeId uint) (model.Episode, error)
	Create(data model.Episode) error
	Delete(episodeId uint) error
}

type providerImpl struct {
	db *gorm.DB
}

func NewEpisodeProvider(db *gorm.DB) Provider {
	return &providerImpl{
		db: db,
	}
}

func (p *providerImpl) FindAll() ([]model.Episode, error) {
	var episodes []model.Episode
	res := p.db.Find(&episodes)
	return episodes, res.Error
}

func (p *providerImpl) FindById(episodeId uint) (model.Episode, error) {
	var episode model.Episode
	res := p.db.Find(&episode, episodeId)
	return episode, res.Error
}

func (p *providerImpl) FindBySeasonId(seasonId uint) ([]model.Episode, error) {
	var episodes []model.Episode
	res := p.db.Where("season_id = ?", seasonId).Find(&episodes)
	return episodes, res.Error
}

func (p *providerImpl) Create(episode model.Episode) error {
	res := p.db.Create(&episode)
	return res.Error
}

func (p *providerImpl) Delete(episodeId uint) error {
	var episode model.Episode
	res := p.db.Where("id = ?", episodeId).Delete(&episode)
	return res.Error
}
