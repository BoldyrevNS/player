package model

type WatchThumbnail struct {
	WatchId    uint   `gorm:"type:int;unique"`
	BucketName string `gorm:"type:varchar(255);"`
	Filename   string `gorm:"type:varchar(255);"`
	VideoUrl   string `gorm:"type:varchar(255);"`
}
