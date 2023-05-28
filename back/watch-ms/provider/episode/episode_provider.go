package episode

import (
	"gorm.io/gorm"
	"watch-ms/model"
)

type Provider interface {
	GetAll() ([]model.Episode, error)
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

func (p *providerImpl) GetAll() ([]model.Episode, error) {
	var episodes []model.Episode
	res := p.db.Find(&episodes)
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
