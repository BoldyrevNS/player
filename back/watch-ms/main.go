package main

import (
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"shared/db"
	"watch-ms/DI"
	"watch-ms/docs"
	"watch-ms/model"
	"watch-ms/router"
	"watch-ms/routines"
)

func migrate(dbInstance *gorm.DB) error {
	err := dbInstance.AutoMigrate(&model.Category{})
	err = dbInstance.AutoMigrate(&model.Episode{})
	err = dbInstance.AutoMigrate(&model.Watch{})
	err = dbInstance.AutoMigrate(&model.Title{})
	err = dbInstance.AutoMigrate(&model.Season{})
	return err
}

// @title		Player auth-service
// @version		1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @BasePath	/api/v1/
func main() {
	err := godotenv.Load(".env")
	docs.SwaggerInfo.Host = os.Getenv("DISPLAY_HOST")
	dbInstance := db.NewDatabase(db.DatabaseConnect{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}).Init()
	err = migrate(dbInstance)
	if err != nil {
		log.Fatalf("migration error")
	}

	container := DI.NewContainer(dbInstance)

	routines.Init(container.Services)

	routes := router.NewRouter(container.Controllers)

	server := &http.Server{
		Addr:    ":8080",
		Handler: routes,
	}

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
