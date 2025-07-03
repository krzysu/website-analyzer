package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	
	"github.com/krzysu/web-crawler/internal/database"
	"github.com/krzysu/web-crawler/internal/models"
	"github.com/krzysu/web-crawler/internal/worker"
)

func AddURL(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json struct {
			URL string `json:"url"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create a new CrawlResult and save it with "queued" status
		result := &models.CrawlResult{
			
			URL:       json.URL,
			Status:    "queued",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.CreateCrawlResult(result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Submit job to the worker queue
		worker.JobQueue <- worker.Job{ID: result.ID, URL: result.URL}

		c.JSON(http.StatusOK, gin.H{"message": "URL submitted for crawling", "id": strconv.FormatUint(uint64(result.ID), 10)})
	}
}

func GetURLs(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
		offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
		sortBy := c.DefaultQuery("sortBy", "created_at")
		filterBy := c.DefaultQuery("filterBy", "")

		results, err := db.GetCrawlResults(limit, offset, sortBy, filterBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}

func GetURL(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
			return
		}
		result, err := db.GetCrawlResult(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Result not found"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func DeleteURLs(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json struct {
			IDs []uint `json:"ids"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.DeleteCrawlResults(json.IDs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "URLs deleted successfully"})
	}
}

func RerunURLs(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json struct {
			IDs []uint `json:"ids"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, id := range json.IDs {
			// Submit re-crawl job to the worker queue
			worker.JobQueue <- worker.Job{ID: id}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Re-crawl initiated for selected URLs"})
	}
}
