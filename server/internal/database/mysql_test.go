package database

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/krzysu/website-analyzer/internal/models"
	"github.com/stretchr/testify/assert"
	gorm "gorm.io/gorm"
)

func TestCreateCrawlResult(t *testing.T) {
	dbInstance, err := NewDBForTest()
	assert.NoError(t, err)
	defer dbInstance.Close()

	result := &models.CrawlResult{
		URL:       "http://example.com",
		Status:    "completed",
		PageTitle: "Example Page",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = dbInstance.CreateCrawlResult(result)
	assert.NoError(t, err)

	var retrieved models.CrawlResult
	err = dbInstance.db.First(&retrieved, result.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, result.URL, retrieved.URL)
}

func TestGetCrawlResult(t *testing.T) {
	dbInstance, err := NewDBForTest()
	assert.NoError(t, err)
	defer dbInstance.Close()

	// Create a result first
	result := &models.CrawlResult{
		URL:       "http://example.com/get",
		Status:    "completed",
		PageTitle: "Get Test Page",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = dbInstance.CreateCrawlResult(result)
	assert.NoError(t, err)

	// Test successful retrieval
	retrieved, err := dbInstance.GetCrawlResult(result.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, result.URL, retrieved.URL)

	// Test not found
	_, err = dbInstance.GetCrawlResult(uint(99999))
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestUpdateCrawlResult(t *testing.T) {
	dbInstance, err := NewDBForTest()
	assert.NoError(t, err)
	defer dbInstance.Close()

	result := &models.CrawlResult{
		URL:       "http://example.com/update",
		Status:    "queued",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = dbInstance.CreateCrawlResult(result)
	assert.NoError(t, err)

	result.Status = "completed"
	result.PageTitle = "Updated Page Title"
	err = dbInstance.UpdateCrawlResult(result)
	assert.NoError(t, err)

	var retrieved models.CrawlResult
	err = dbInstance.db.First(&retrieved, result.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "completed", retrieved.Status)
	assert.Equal(t, "Updated Page Title", retrieved.PageTitle)
}

func TestDeleteCrawlResult(t *testing.T) {
	dbInstance, err := NewDBForTest()
	assert.NoError(t, err)
	defer dbInstance.Close()

	result := &models.CrawlResult{
		URL: "http://example.com/delete",
	}
	err = dbInstance.CreateCrawlResult(result)
	assert.NoError(t, err)

	err = dbInstance.DeleteCrawlResult(result.ID)
	assert.NoError(t, err)

	var retrieved models.CrawlResult
	err = dbInstance.db.First(&retrieved, result.ID).Error
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}

func TestDeleteCrawlResults(t *testing.T) {
	dbInstance, err := NewDBForTest()
	assert.NoError(t, err)
	defer dbInstance.Close()

	ids := []uint{}
	for i := 0; i < 3; i++ {
		result := &models.CrawlResult{
			URL: fmt.Sprintf("http://example.com/bulk-delete-%d", i),
		}
		err = dbInstance.CreateCrawlResult(result)
		assert.NoError(t, err)
		ids = append(ids, result.ID)
	}

	err = dbInstance.DeleteCrawlResults(ids)
	assert.NoError(t, err)

	var count int64
	dbInstance.db.Model(&models.CrawlResult{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestGetCrawlResults(t *testing.T) {
	dbInstance, err := NewDBForTest()
	assert.NoError(t, err)
	defer dbInstance.Close()

	// Create some test data
	for i := 0; i < 5; i++ {
		result := &models.CrawlResult{
			URL:       fmt.Sprintf("http://example.com/page%d", i),
			Status:    "completed",
			PageTitle: fmt.Sprintf("Page %d", i),
			CreatedAt: time.Now().Add(time.Duration(i) * time.Hour),
			UpdatedAt: time.Now().Add(time.Duration(i) * time.Hour),
		}
		err = dbInstance.CreateCrawlResult(result)
		assert.NoError(t, err)
	}

	// Test pagination
	results, err := dbInstance.GetCrawlResults(2, 0, "", "")
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Page 0", results[0].PageTitle)

	results, err = dbInstance.GetCrawlResults(2, 2, "", "")
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Page 2", results[0].PageTitle)

	// Test sorting
	results, err = dbInstance.GetCrawlResults(5, 0, "page_title desc", "")
	assert.NoError(t, err)
	assert.Len(t, results, 5)
	assert.Equal(t, "Page 4", results[0].PageTitle)

	// Test filtering
	results, err = dbInstance.GetCrawlResults(5, 0, "", "page1")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "Page 1", results[0].PageTitle)
}
