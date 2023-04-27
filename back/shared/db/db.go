package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type DatabaseConnect struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type Database interface {
	Init() *gorm.DB
}

type DatabaseImpl struct {
	instance *gorm.DB
}

func NewDatabase(connect DatabaseConnect) Database {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s", connect.Host, connect.User, connect.Name, connect.Password, connect.Port)))
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	return &DatabaseImpl{instance: db}
}

func (d *DatabaseImpl) Init() *gorm.DB {
	return d.instance
}
