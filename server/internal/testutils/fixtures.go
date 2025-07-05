package testutils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// NewSimpleWebsite is a test fixture that creates a mock server
// with a basic HTML page.
func NewSimpleWebsite() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `<!DOCTYPE html>
		<html>
		<head><title>Simple Page</title></head>
		<body>
			<h1>Main Heading</h1>
			<a href="/internal-link">Internal</a>
			<a href="http://external.com/page">External</a>
		</body>
		</html>`)
	}))
}

// NewComplexWebsite is a test fixture that creates a mock server with a more
// complex HTML page, including a form and broken links.
func NewComplexWebsite() *httptest.Server {
	mux := http.NewServeMux()

	// Handler for the main page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// This check is important to avoid this handler catching all requests
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		// The external link now points to a path on the same mock server
		fmt.Fprintf(w, `<!DOCTYPE html>
		<html>
		<head><title>Complex Page</title></head>
		<body>
			<h1>First Heading</h1>
			<h2>Second Heading</h2>
			<a href="/internal-ok">Working Internal Link</a>
			<a href="/internal-broken">Broken Internal Link</a>
			<a href="/external-broken">"External" Broken Link</a>
		</body>
		</html>`)
	})

	// Handler for a working internal link
	mux.HandleFunc("/internal-ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Handler for a broken internal link
	mux.HandleFunc("/internal-broken", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// Handler for the "external" link, which is also broken
	mux.HandleFunc("/external-broken", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	return httptest.NewServer(mux)
}
