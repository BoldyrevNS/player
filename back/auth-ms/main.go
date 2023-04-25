package main

import (
	"auth-ms/db"
	_ "auth-ms/docs"
	"auth-ms/model"
	"auth-ms/router"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func migrate() error {
	err := db.Instance.AutoMigrate(&model.User{})
	return err
}

// @title		Player auth-service
// @version		1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @host 		localhost:8080
// @BasePath	/api/v1/
func main() {
	err := godotenv.Load(".env")
	db.NewDatabase().Init()
	err = migrate()
	if err != nil {
		log.Fatalf("Migration error")
	}
	routes := router.NewRouter()

	server := &http.Server{
		Addr:    ":8080",
		Handler: routes,
	}

	err = server.ListenAndServe()

	if err != nil {
		log.Fatalf("Serve error: %v", err)
	}
}
