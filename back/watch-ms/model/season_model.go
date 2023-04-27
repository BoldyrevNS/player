package model

type Season struct {
	Id      uint `gorm:"type:int;primary_key"`
	TitleId uint
	Number  uint  `gorm:"type:int;"`
	Title   Title `gorm:"constraint:OnDelete:CASCADE;"`
}
