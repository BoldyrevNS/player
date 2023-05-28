package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"upload-ms/docs"
	"upload-ms/router"
)

// @title		Player upload-service
// @version		1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @BasePath	/api/v1/
func main() {
	_ = godotenv.Load(".env")
	docs.SwaggerInfo.Host = os.Getenv("DISPLAY_HOST")

	routes := router.NewRouter()
	server := &http.Server{
		Addr:    ":8080",
		Handler: routes,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
