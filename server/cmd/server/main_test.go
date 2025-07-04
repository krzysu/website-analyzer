
package main

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/krzysu/web-crawler/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestServerStartup(t *testing.T) {
	// Set a test port
	os.Setenv("PORT", "8081")

	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	router := setupServer(db)
	go router.Run(":" + os.Getenv("PORT"))
	time.Sleep(1 * time.Second)

	// Check if the server is running
	_, err = http.Get("http://localhost:" + os.Getenv("PORT") + "/urls")
	assert.NoError(t, err)
}
