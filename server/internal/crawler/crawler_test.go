package crawler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/krzysu/website-analyzer/internal/models"
)

func TestCrawl_BasicExtraction(t *testing.T) {
	// Create a mock HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		// Simulate HTML5 doctype
		_, err := w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>Test Page</title></head>
<body>
<h1>Heading 1</h1>
<h2>Heading 2</h2>
<a href="/internal">Internal Link</a>
<a href="http://external.com">External Link</a>
<form><input type="password" name="password"></form>
</body>
</html>`))
		assert.NoError(t, err)
	}))
	defer ts.Close()

	result := &models.CrawlResult{
		URL:       ts.URL,
		Status:    "queued",
		Headings:  make(map[string]int),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := Crawl(result)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, "Test Page", result.PageTitle)
	assert.Equal(t, "HTML5", result.HTMLVersion)
	assert.Equal(t, 1, result.Headings["h1"])
	assert.Equal(t, 1, result.Headings["h2"])
	assert.Equal(t, 1, result.InternalLinksCount)
	assert.Equal(t, 1, result.ExternalLinksCount)
	assert.True(t, result.HasLoginForm)
	assert.Equal(t, "completed", result.Status)
}

func TestCrawl_InaccessibleLinks(t *testing.T) {
	// Mock server for the main page
	mainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, err := w.Write([]byte(`<!DOCTYPE html>
<html>
<body>
<a href="/broken-404">Broken Link 404</a>
<a href="/broken-500">Broken Link 500</a>
<a href="/ok">OK Link</a>
</body>
</html>`))
		assert.NoError(t, err)
	}))
	defer mainServer.Close()

	// Mock server for the links
	linkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/broken-404" {
			w.WriteHeader(http.StatusNotFound)
		} else if r.URL.Path == "/broken-500" {
			w.WriteHeader(http.StatusInternalServerError)
		} else if r.URL.Path == "/ok" {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer linkServer.Close()

	// Replace the links in the main page HTML to point to the mock link server
	mainPageHTML := `<!DOCTYPE html>
<html>
<body>
<a href="` + linkServer.URL + `/broken-404">Broken Link 404</a>
<a href="` + linkServer.URL + `/broken-500">Broken Link 500</a>
<a href="` + linkServer.URL + `/ok">OK Link</a>
</body>
</html>`

	// Create a new main server with the updated HTML
	mainServerWithMockLinks := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, err := w.Write([]byte(mainPageHTML))
		assert.NoError(t, err)
	}))
	defer mainServerWithMockLinks.Close()

	result := &models.CrawlResult{
		URL:       mainServerWithMockLinks.URL,
		Status:    "queued",
		Headings:  make(map[string]int),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := Crawl(result)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, 2, result.InaccessibleLinksCount)
	assert.Len(t, result.BrokenLinks, 2)

	// Check if the broken links are correctly reported
	found404 := false
	found500 := false
	for _, bl := range result.BrokenLinks {
		statusCode, ok := bl["status_code"].(int)
		if !ok {
			// If it's not an int, try float64 (common with JSON unmarshalling)
			statusCodeFloat, okFloat := bl["status_code"].(float64)
			if okFloat {
				statusCode = int(statusCodeFloat)
			}
		}

		if statusCode == 404 {
			found404 = true
		} else if statusCode == 500 {
			found500 = true
		}
	}
	assert.True(t, found404)
	assert.True(t, found500)
}

func TestCrawl_ErrorHandling(t *testing.T) {
	// Test with an invalid URL
	result := &models.CrawlResult{
		URL:       "invalid-url",
		Status:    "queued",
		Headings:  make(map[string]int),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := Crawl(result)
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "error", result.Status)
	assert.Contains(t, result.ErrorMessage, "unsupported protocol scheme")

	// Test with a server that returns an error (e.g., connection refused)
	// This is harder to mock directly with httptest.NewServer as it implies
	// the server itself is not reachable. For a real-world scenario, you'd
	// need to simulate network issues or use a custom http.Client with a Transport
	// that can inject errors.

	// For now, we can simulate a non-existent domain
	result = &models.CrawlResult{
		URL:       "http://nonexistent-domain-12345.com",
		Status:    "queued",
		Headings:  make(map[string]int),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = Crawl(result)
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "error", result.Status)
	assert.Contains(t, result.ErrorMessage, "no such host")
}
