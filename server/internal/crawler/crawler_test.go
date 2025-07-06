package crawler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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
	// Test with an invalid URL (unsupported protocol scheme)
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

	// Test with a mock server that immediately closes the connection
	// This simulates a connection refused error.
	closeConnServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Close the underlying TCP connection immediately
		hijacker, ok := w.(http.Hijacker)
		assert.True(t, ok)
		conn, _, err := hijacker.Hijack()
		assert.NoError(t, err)
		conn.Close()
	}))
	defer closeConnServer.Close()

	result = &models.CrawlResult{
		URL:       closeConnServer.URL,
		Status:    "queued",
		Headings:  make(map[string]int),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = Crawl(result)
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "error", result.Status)
	assert.True(t, strings.Contains(result.ErrorMessage, "connection reset by peer") || strings.Contains(result.ErrorMessage, "EOF") || strings.Contains(result.ErrorMessage, "connection refused"))

	// Test with a mock server that returns a non-200 status code for HEAD request
	// This simulates an inaccessible link during checkLinks
	brokenLinkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer brokenLinkServer.Close()

	mainServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<!DOCTYPE html>
<html>
<body>
<a href="%s">Broken Link</a>
</body>
</html>`, brokenLinkServer.URL)
	}))
	defer mainServer.Close()

	result = &models.CrawlResult{
		URL:       mainServer.URL,
		Status:    "queued",
		Headings:  make(map[string]int),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = Crawl(result)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "completed", result.Status)
	assert.Equal(t, 1, result.InaccessibleLinksCount)
}

func TestGetHTMLVersion(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "HTML5 Doctype",
			html:     "<!DOCTYPE html>\n<html><body></body></html>",
			expected: "HTML5",
		},
		{
			name:     "HTML5 Doctype with leading spaces",
			html:     "  <!DOCTYPE html>\n<html><body></body></html>",
			expected: "HTML5",
		},
		{
			name:     "HTML 4.01 Strict Doctype",
			html:     "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\">\n<html><body></body></html>",
			expected: "HTML 4.01 Strict",
		},
		{
			name:     "XHTML 1.0 Strict Doctype",
			html:     "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\">\n<html><body></body></html>",
			expected: "XHTML 1.0 Strict",
		},
		{
			name:     "No Doctype",
			html:     "<html><body></body></html>",
			expected: "Unknown",
		},
		{
			name:     "Empty HTML",
			html:     "",
			expected: "Unknown",
		},
		{
			name:     "Doctype on second line",
			html:     "\n<!DOCTYPE html>\n<html><body></body></html>",
			expected: "HTML5",
		},
		{
			name:     "HTML 2.0 Doctype",
			html:     "<!DOCTYPE html PUBLIC \"-//IETF//DTD HTML 2.0//EN\">\n<html><body></body></html>",
			expected: "HTML 2.0",
		},
		{
			name:     "HTML 3.2 Doctype",
			html:     "<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 3.2 Final//EN\">\n<html><body></body></html>",
			expected: "HTML 3.2",
		},
		{
			name:     "XHTML 1.1 Doctype",
			html:     "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">\n<html><body></body></html>",
			expected: "XHTML 1.1",
		},
		{
			name:     "HTML5 with profile",
			html:     "<!DOCTYPE html profile=\"http://www.w3.org/2000/svg\">\n<html><body></body></html>",
			expected: "HTML5 with profile",
		},
		{
			name:     "HTML5 Doctype and html tag on same line",
			html:     "<!DOCTYPE html><html lang=\"en\"><head><meta charSet=\"utf-8\"/>",
			expected: "HTML5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy http.Response
			version := getHTMLVersion([]byte(tt.html))
			assert.Equal(t, tt.expected, version)
		})
	}
}
