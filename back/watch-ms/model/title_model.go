package model

type Title struct {
	Id           uint `gorm:"type:int;primary_key"`
	CategoryId   uint
	Name         string   `gorm:"type:varchar(255);unique"`
	ThumbnailUrl string   `gorm:"type:varchar(255)"`
	Category     Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
