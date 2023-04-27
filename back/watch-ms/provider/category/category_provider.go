package category

import (
	"gorm.io/gorm"
	"watch-ms/model"
)

type Provider interface {
	GetAll() ([]model.Category, error)
	Create(category model.Category) error
	Delete(categoryId uint) error
}

type providerImpl struct {
	db *gorm.DB
}

func NewCategoryProvider(db *gorm.DB) Provider {
	return &providerImpl{db: db}
}

func (p *providerImpl) GetAll() ([]model.Category, error) {
	var categories []model.Category
	res := p.db.Find(&categories)
	return categories, res.Error
}

func (p *providerImpl) Create(category model.Category) error {
	res := p.db.Create(&category)
	return res.Error
}

func (p *providerImpl) Delete(categoryId uint) error {
	var category model.Category
	res := p.db.Where("id = ?", categoryId).Delete(&category)
	return res.Error
}
