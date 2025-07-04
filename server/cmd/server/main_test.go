package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/krzysu/website-analyzer/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestServerStartup(t *testing.T) {
	// Set a test port
	os.Setenv("PORT", "8081")

	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	router := setupServer(db)
	go func() {
		err = router.Run(":" + os.Getenv("PORT"))
		assert.NoError(t, err)
	}()
	// Wait for the server to start
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get("http://localhost:" + os.Getenv("PORT") + "/urls")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Check if the server is running
	_, err = http.Get("http://localhost:" + os.Getenv("PORT") + "/urls")
	assert.NoError(t, err)
}
