package main

import (
	"auth-ms/docs"
	"auth-ms/router"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"shared/db"
	model "shared/model/user"
)

func migrate(dbInstance *gorm.DB) error {
	err := dbInstance.AutoMigrate(&model.User{})
	return err
}

// @title		Player auth-service
// @version		1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @BasePath	/api/v1/
func main() {
	_ = godotenv.Load(".env")
	docs.SwaggerInfo.Host = os.Getenv("DISPLAY_HOST")
	dbInstance := db.NewDatabase(db.DatabaseConnect{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}).Init()
	err := migrate(dbInstance)
	if err != nil {
		log.Fatalf("migration error")
	}
	routes := router.NewRouter(dbInstance)

	server := &http.Server{
		Addr:    ":8080",
		Handler: routes,
	}

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("Serve error: %v", err)
	}
}
