package main

import (
	"log"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/krzysu/web-crawler/internal/api"
	"github.com/krzysu/web-crawler/internal/database"
	"github.com/krzysu/web-crawler/internal/worker"
)

func setupServer(db *database.DB) *gin.Engine {
	var wg sync.WaitGroup // Create a WaitGroup for the application

	dispatcher := worker.NewDispatcher(5, db, &wg) // Pass db and wg to dispatcher
	dispatcher.Run()

	// Set up the Gin router
	router := gin.Default()
	api.SetupRoutes(router, db, dispatcher.JobQueue) // Pass db and JobQueue to API setup

	return router
}

func main() {
	// Initialize the database connection
	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	router := setupServer(db)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}