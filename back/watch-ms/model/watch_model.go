package model

type Watch struct {
	Id           uint   `gorm:"type:int;primary_key"`
	UserId       uint   `gorm:"type:int"`
	EpisodeId    uint   `gorm:"type:int"`
	ThumbnailUrl string `gorm:"type:varchar(255)"`
	Episode      Episode
}
