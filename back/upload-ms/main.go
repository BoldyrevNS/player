package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"upload-ms/router"
)

func main() {
	_ = godotenv.Load(".env")

	routes := router.NewRouter()
	server := &http.Server{
		Addr:    ":8089",
		Handler: routes,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("serve error: %v", err)
	}
}
