package watch_thumbnail

import (
	"gorm.io/gorm"
	"upload-ms/model"
)

type UpdateWatch struct {
	BucketName string
	Filename   string
}

type Provider interface {
	Create(createData model.WatchThumbnail) error
	FindOne(watchId uint) (model.WatchThumbnail, error)
	UpdateOne(watchId uint, updateData UpdateWatch) error
}

type providerImpl struct {
	db *gorm.DB
}

func NewWatchThumbnailProvider(db *gorm.DB) Provider {
	return &providerImpl{
		db: db,
	}
}

func (p *providerImpl) Create(createData model.WatchThumbnail) error {
	res := p.db.Create(createData)
	return res.Error
}

func (p *providerImpl) FindOne(watchId uint) (model.WatchThumbnail, error) {
	var lastUpload model.WatchThumbnail
	res := p.db.Where("watch_id = ?", watchId).First(&lastUpload)
	return lastUpload, res.Error
}

func (p *providerImpl) UpdateOne(watchId uint, updateData UpdateWatch) error {
	res := *p.db.Where("watch_id = ?", watchId).Updates(&model.WatchThumbnail{
		BucketName: updateData.BucketName,
		Filename:   updateData.Filename,
	})
	return res.Error
}
