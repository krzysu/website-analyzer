
package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/krzysu/web-crawler/internal/database"
	"github.com/krzysu/web-crawler/internal/models"
	"github.com/krzysu/web-crawler/internal/worker"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestAddURL_Success(t *testing.T) {
	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	body := map[string]string{"url": "http://example.com"}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "URL submitted for crawling", response["message"])

	// Check if the job was added to the queue
	job := <-jobQueue
	assert.Equal(t, "http://example.com", job.URL)
}

func TestAddURL_InvalidJSON(t *testing.T) {
	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/urls", bytes.NewBuffer([]byte(`{"url":}`)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetURLs_Success(t *testing.T) {
	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	// Create some test data
	db.CreateCrawlResult(&models.CrawlResult{URL: "http://example.com/1"})
	db.CreateCrawlResult(&models.CrawlResult{URL: "http://example.com/2"})

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/urls?limit=2&offset=0", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var results []models.CrawlResult
	err = json.Unmarshal(w.Body.Bytes(), &results)
	assert.NoError(t, err)
	assert.Len(t, results, 2)
}

func TestGetURL_Success(t *testing.T) {
	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	// Create a test entry
	result := &models.CrawlResult{URL: "http://example.com/1"}
	db.CreateCrawlResult(result)

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/urls/"+strconv.FormatUint(uint64(result.ID), 10), nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var returnedResult models.CrawlResult
	err = json.Unmarshal(w.Body.Bytes(), &returnedResult)
	assert.NoError(t, err)
	assert.Equal(t, result.URL, returnedResult.URL)
}

func TestDeleteURLs_Success(t *testing.T) {
	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	// Create a test entry
	result := &models.CrawlResult{URL: "http://example.com/1"}
	db.CreateCrawlResult(result)

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	body := map[string][]uint{"ids": {result.ID}}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/urls", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the entry was deleted
	_, err = db.GetCrawlResult(result.ID)
	assert.Error(t, err)
}

func TestRerunURLs_Success(t *testing.T) {
	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	// Create a test entry
	result := &models.CrawlResult{URL: "http://example.com/1"}
	db.CreateCrawlResult(result)

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	body := map[string][]uint{"ids": {result.ID}}
	jsonBody, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/urls/rerun", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the job was added to the queue
	job := <-jobQueue
	assert.Equal(t, result.ID, job.ID)
}

func TestGetURL_NotFound(t *testing.T) {
	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/urls/999", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteURLs_InvalidJSON(t *testing.T) {
	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/urls", bytes.NewBuffer([]byte(`{"ids":`)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRerunURLs_InvalidJSON(t *testing.T) {
	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/urls/rerun", bytes.NewBuffer([]byte(`{"ids":`)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
