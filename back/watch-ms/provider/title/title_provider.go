package title

import (
	"gorm.io/gorm"
	"watch-ms/model"
)

type Provider interface {
	Create(title model.Title) error
	FindAll() ([]model.Title, error)
	Delete(titleId uint) error
}

type providerImpl struct {
	db *gorm.DB
}

func NewTitleProvider(db *gorm.DB) Provider {
	return &providerImpl{
		db: db,
	}
}

func (p *providerImpl) Create(title model.Title) error {
	res := p.db.Create(&title)
	return res.Error
}

func (p *providerImpl) FindAll() ([]model.Title, error) {
	var titles []model.Title
	res := p.db.Find(&titles)
	return titles, res.Error
}

func (p *providerImpl) Delete(titleId uint) error {
	var title model.Title
	res := p.db.Where("id = ?", titleId).Delete(&title)
	return res.Error
}
