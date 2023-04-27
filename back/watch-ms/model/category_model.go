package model

type Category struct {
	Id   uint   `gorm:"type:int;primary_key"`
	Name string `gorm:"type:varchar(255);unique"`
}
