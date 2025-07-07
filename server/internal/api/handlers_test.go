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

	"github.com/krzysu/website-analyzer/internal/database"
	"github.com/krzysu/website-analyzer/internal/models"
	"github.com/krzysu/website-analyzer/internal/worker"
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
	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/urls", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
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
	var req *http.Request
	req, err = http.NewRequest("POST", "/urls", bytes.NewBuffer([]byte(`{"url":}`)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetURLs_Success(t *testing.T) {
	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	// Create some test data
	err = db.CreateCrawlResult(&models.CrawlResult{URL: "http://example.com/1"})
	assert.NoError(t, err)
	err = db.CreateCrawlResult(&models.CrawlResult{URL: "http://example.com/2"})
	assert.NoError(t, err)

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	w := httptest.NewRecorder()
	var req *http.Request
	req, err = http.NewRequest("GET", "/urls?limit=2&offset=0", nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Results []models.CrawlResult `json:"results"`
		Total   int64                `json:"total"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Results, 2)
	assert.Equal(t, int64(2), response.Total)
}

func TestGetURL_Success(t *testing.T) {
	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	// Create a test entry
	result := &models.CrawlResult{URL: "http://example.com/1"}
	err = db.CreateCrawlResult(result)
	assert.NoError(t, err)

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	w := httptest.NewRecorder()
	var req *http.Request
	req, err = http.NewRequest("GET", "/urls/"+strconv.FormatUint(uint64(result.ID), 10), nil)
	assert.NoError(t, err)

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
	err = db.CreateCrawlResult(result)
	assert.NoError(t, err)

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	body := map[string][]uint{"ids": {result.ID}}
	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	var req *http.Request
	req, err = http.NewRequest("DELETE", "/urls", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
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
	err = db.CreateCrawlResult(result)
	assert.NoError(t, err)

	router := setupRouter()
	jobQueue := make(chan worker.Job, 1)
	SetupRoutes(router, db, jobQueue)

	body := map[string][]uint{"ids": {result.ID}}
	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	var req *http.Request
	req, err = http.NewRequest("POST", "/urls/rerun", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	assert.NoError(t, err)
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
	var req *http.Request
	req, err = http.NewRequest("GET", "/urls/999", nil)
	assert.NoError(t, err)

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
	var req *http.Request
	req, err = http.NewRequest("DELETE", "/urls", bytes.NewBuffer([]byte(`{"ids":`)))
	assert.NoError(t, err)
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
	var req *http.Request
	req, err = http.NewRequest("POST", "/urls/rerun", bytes.NewBuffer([]byte(`{"ids":`)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
