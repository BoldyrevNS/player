package model

type Episode struct {
	Id            uint `gorm:"type:int;primary_key"`
	SeasonId      uint
	EpisodeName   string `gorm:"type:varchar(255)"`
	EpisodeNumber uint   `gorm:"type:int"`
	VideoUrl      string `gorm:"type:varchar(255)"`
	ThumbnailUrl  string `gorm:"type:varchar(255)"`
}
