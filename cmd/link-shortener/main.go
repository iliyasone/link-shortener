package main

import (
	"flag"
	"log"

	"link-shortener/internal/handlers"
	"link-shortener/internal/storage"

	"github.com/gin-gonic/gin"

	_ "link-shortener/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Link Shortener API
// @version 0.1.0
// @description A simple URL shortening service in Go.
// @host localhost:8080
// @BasePath /
func main() {
	storageType := flag.String("storage", "ram", "Storage backend to use: 'ram' or 'postgres'")
	flag.Parse()

	var store storage.Storage
	if *storageType == "postgres" {
		var err error
		store, err = storage.NewPostgresStorage()
		if err != nil {
			log.Fatalf("Failed to initialize postgres storage: %v", err)
		}
	} else if *storageType == "ram" {
		store = storage.NewRAMStorage()
	} else {
		log.Fatal("Wrong storage type")
	}

	r := gin.Default()
	r.POST("/shorten", handlers.PostShortenURL(store))
	r.GET("/:shortURL", handlers.GetOriginalURL(store))

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server is running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
