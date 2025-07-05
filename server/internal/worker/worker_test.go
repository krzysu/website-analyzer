package worker

import (
	"sync"
	"testing"
	"time"

	"github.com/krzysu/website-analyzer/internal/database"
	"github.com/krzysu/website-analyzer/internal/models"
	"github.com/krzysu/website-analyzer/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorker_RecrawlResetsCountableFields(t *testing.T) {
	// 1. Create a mock HTTP server
	ts := testutils.NewComplexWebsite() // Use the more complex fixture here
	defer ts.Close()

	db, err := database.NewDBForTest()
	require.NoError(t, err)
	defer db.Close()

	// 2. Create a CrawlResult with initial values
	initialResult := &models.CrawlResult{
		URL:                  ts.URL,
		Status:               "completed",
		InternalLinksCount:   10, // Initial value
		ExternalLinksCount:   20, // Initial value
		InaccessibleLinksCount: 5,  // Initial value
		Headings:             map[string]int{"h1": 2},
	}
	err = db.CreateCrawlResult(initialResult)
	require.NoError(t, err)

	var wg sync.WaitGroup
	dispatcher := NewDispatcher(1, db, &wg)
	dispatcher.Run()

	// 4. Enqueue a re-crawl job
	job := Job{ID: initialResult.ID, URL: initialResult.URL}
	dispatcher.JobQueue <- job

	// 5. Wait for the worker to finish by polling the database
	var finalResult *models.CrawlResult
	require.Eventually(t, func() bool {
		var err error
		finalResult, err = db.GetCrawlResult(initialResult.ID)
		if err != nil {
			return false
		}
		return finalResult.Status == "completed"
	}, 5*time.Second, 50*time.Millisecond, "worker did not complete in time")

	// 7. Assert that the countable fields have been overwritten
	assert.Equal(t, "completed", finalResult.Status)
	assert.Equal(t, 3, finalResult.InternalLinksCount, "InternalLinksCount should be overwritten, not accumulated")
	assert.Equal(t, 0, finalResult.ExternalLinksCount, "ExternalLinksCount should be overwritten, not accumulated")
	assert.Equal(t, 2, finalResult.InaccessibleLinksCount, "InaccessibleLinksCount should be reset")
	assert.Equal(t, 1, finalResult.Headings["h1"], "Headings should be overwritten, not accumulated")
	assert.Equal(t, 1, finalResult.Headings["h2"], "Headings should be overwritten, not accumulated")
}

func TestWorker(t *testing.T) {
	// Use the simple website fixture
	ts := testutils.NewSimpleWebsite()
	defer ts.Close()

	db, err := database.NewDBForTest()
	require.NoError(t, err)
	defer db.Close()

	var wg sync.WaitGroup
	dispatcher := NewDispatcher(1, db, &wg)
	dispatcher.Run()

	// Test creating a new crawl result
	job := Job{URL: ts.URL}
	dispatcher.JobQueue <- job

	// Test re-crawling an existing result
	result := &models.CrawlResult{URL: ts.URL + "/recrawl"} // Use a different URL for the second job
	err = db.CreateCrawlResult(result)
	require.NoError(t, err)
	job = Job{ID: result.ID, URL: result.URL}
	dispatcher.JobQueue <- job

	// Wait for both jobs to complete by polling the database
	require.Eventually(t, func() bool {
		results, err := db.GetCrawlResults(2, 0, "", "")
		if err != nil {
			return false
		}
		if len(results) != 2 {
			return false
		}
		// Check that both have a final status. The actual status doesn't matter as much
		// as the fact that they are no longer "running" or "queued".
		return (results[0].Status == "completed" || results[0].Status == "error") &&
			(results[1].Status == "completed" || results[1].Status == "error")
	}, 5*time.Second, 50*time.Millisecond)

	// Verify the results in the database
	results, err := db.GetCrawlResults(2, 0, "", "")
	require.NoError(t, err)
	assert.Len(t, results, 2)
}
