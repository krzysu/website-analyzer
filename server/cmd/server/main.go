package main

import (
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/krzysu/web-crawler/internal/api"
	"github.com/krzysu/web-crawler/internal/database"
	"github.com/krzysu/web-crawler/internal/worker"
)

func main() {
	// Initialize the database connection
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	var wg sync.WaitGroup // Create a WaitGroup for the application

	// Initialize and start the worker dispatcher
	dispatcher := worker.NewDispatcher(5, db, &wg) // Pass db and wg to dispatcher
	dispatcher.Run()

	// Set up the Gin router
	router := gin.Default()
	api.SetupRoutes(router, db) // Pass db to API setup

	// Start the server
	log.Println("Server starting on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}