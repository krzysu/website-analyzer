package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krzysu/web-crawler/internal/database"
	"github.com/krzysu/web-crawler/internal/worker"
)

func AddURL(c *gin.Context) {
	var json struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Submit job to the worker queue
	worker.JobQueue <- worker.Job{URL: json.URL}

	c.JSON(http.StatusOK, gin.H{"message": "URL submitted for crawling"})
}

func GetURLs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sortBy", "created_at")
	filterBy := c.DefaultQuery("filterBy", "")

	results, err := database.GetCrawlResults(limit, offset, sortBy, filterBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

func GetURL(c *gin.Context) {
	id := c.Param("id")
	result, err := database.GetCrawlResult(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Result not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func DeleteURLs(c *gin.Context) {
	var json struct {
		IDs []string `json:"ids"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DeleteCrawlResults(json.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "URLs deleted successfully"})
}

func RerunURLs(c *gin.Context) {
	var json struct {
		IDs []string `json:"ids"`
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
