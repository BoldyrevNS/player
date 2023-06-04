package title

import (
	"gorm.io/gorm"
	"watch-ms/model"
)

type titleWithCategory struct {
	Id           uint
	Name         string
	ThumbnailUrl string
	CategoryName string
	CategoryId   uint
}

type Provider interface {
	Create(title model.Title) error
	FindAll() ([]model.Title, error)
	Delete(titleId uint) error
	GetTitlesExcludeWatched(userId uint) ([]titleWithCategory, error)
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

func (p *providerImpl) GetTitlesExcludeWatched(userId uint) ([]titleWithCategory, error) {
	var titles []titleWithCategory
	query := p.db.
		Table("titles").
		Joins("INNER JOIN categories ON categories.id = titles.category_id").
		Joins("INNER JOIN seasons ON titles.id = seasons.title_id").
		Joins("INNER JOIN episodes ON seasons.id = episodes.season_id").
		Joins("LEFT JOIN watches ON episodes.id = watches.episode_id AND watches.user_id = ?", userId).
		Select("DISTINCT titles.id AS id, titles.name AS name, titles.thumbnail_url AS thumbnail_url, categories.name AS category_name, categories.id AS category_id").
		Where("watches.id IS NULL")
	if err := query.Find(&titles).Error; err != nil {
		return nil, err
	}
	return titles, nil
}
