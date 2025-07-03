package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/krzysu/web-crawler/internal/api"
	"github.com/krzysu/web-crawler/internal/database"
	"github.com/krzysu/web-crawler/internal/worker"
)

func main() {
	// Initialize the database connection
	if err := database.InitDB("user:password@tcp(127.0.0.1:3306)/crawler_db"); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize and start the worker dispatcher
	dispatcher := worker.NewDispatcher(5) // 5 workers
	dispatcher.Run()

	// Set up the Gin router
	router := gin.Default()
	api.SetupRoutes(router)

	// Start the server
	log.Println("Server starting on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
