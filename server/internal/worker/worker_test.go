
package worker

import (
	"sync"
	"testing"

	"github.com/krzysu/web-crawler/internal/database"
	"github.com/krzysu/web-crawler/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	db, err := database.NewDBForTest()
	assert.NoError(t, err)
	defer db.Close()

	var wg sync.WaitGroup
	dispatcher := NewDispatcher(1, db, &wg)
	dispatcher.Run()

	// Test creating a new crawl result
	job := Job{URL: "http://example.com"}
	wg.Add(1)
	dispatcher.JobQueue <- job

	// Test re-crawling an existing result
	result := &models.CrawlResult{URL: "http://example.com/recrawl"}
	err = db.CreateCrawlResult(result)
	assert.NoError(t, err)
	job = Job{ID: result.ID, URL: result.URL}
	wg.Add(1)
	dispatcher.JobQueue <- job

	wg.Wait()

	// Verify the results in the database
	results, err := db.GetCrawlResults(2, 0, "", "")
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "completed", results[0].Status)
	assert.Equal(t, "completed", results[1].Status)
}
