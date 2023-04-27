package model

type Episode struct {
	Id       uint `gorm:"type:int;primary_key"`
	SeasonId uint
	Name     string `gorm:"type:varchar(255)"`
	Number   uint   `gorm:"type:int"`
	VideoUrl string `gorm:"type:varchar(255)"`
}
