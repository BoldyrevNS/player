package watch

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"watch-ms/model"
	"watch-ms/provider/episode"
)

type withVideoUrl struct {
	WatchId  uint
	VideoUrl string
}

type Provider interface {
	Create(data model.Watch) (*model.Watch, error)
	CreateAndReturnJoinVideoUrl(data model.Watch) (withVideoUrl, error)
	Delete(watchId uint) error
	GetAllUserWatch(userId uint) ([]model.Watch, error)
	UpdateThumbnailUrl(watchId uint, newThumbnailUrl string) error
}

type providerImpl struct {
	db              *gorm.DB
	episodeProvider episode.Provider
}

func NewWatchProvider(db *gorm.DB, episodeProvider episode.Provider) Provider {
	return &providerImpl{
		db:              db,
		episodeProvider: episodeProvider,
	}
}

func (p *providerImpl) Create(data model.Watch) (*model.Watch, error) {
	watch := model.Watch{
		UserId:       data.UserId,
		EpisodeId:    data.EpisodeId,
		ThumbnailUrl: data.ThumbnailUrl,
	}
	res := p.db.Clauses(clause.Returning{}).Create(&watch)
	return &watch, res.Error
}

func (p *providerImpl) CreateAndReturnJoinVideoUrl(requestData model.Watch) (withVideoUrl, error) {
	var responseData withVideoUrl
	err := p.db.Transaction(func(tx *gorm.DB) error {
		createData, err := p.Create(requestData)
		if err != nil {
			return err
		}
		episodeData, err := p.episodeProvider.FindById(requestData.EpisodeId)
		if err != nil {
			return err
		}
		responseData = withVideoUrl{
			WatchId:  createData.Id,
			VideoUrl: episodeData.VideoUrl,
		}
		return nil
	})
	if err != nil {
		return withVideoUrl{}, err
	}
	return responseData, err
}

func (p *providerImpl) Delete(watchId uint) error {
	var watch model.Watch
	res := p.db.Where("id = ?", watchId).Delete(&watch)
	return res.Error
}

func (p *providerImpl) GetAllUserWatch(userId uint) ([]model.Watch, error) {
	var watch []model.Watch
	res := p.db.Where("user_id = ?", userId).Find(&watch)
	return watch, res.Error
}

func (p *providerImpl) UpdateThumbnailUrl(watchId uint, newThumbnailUrl string) error {
	res := p.db.Where("id = ?", watchId).Updates(&model.Watch{
		ThumbnailUrl: newThumbnailUrl,
	})
	return res.Error
}
