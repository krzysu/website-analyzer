package crawler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/krzysu/web-crawler/internal/models"
	"golang.org/x/net/html"
)

// Crawl performs the crawling of a single URL.
func Crawl(targetURL string) (*models.CrawlResult, error) {
	// Create a new CrawlResult
	result := &models.CrawlResult{
		ID:        uuid.New().String(),
		URL:       targetURL,
		Status:    "running",
		Headings:  make(map[string]int),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Fetch the URL
	resp, err := http.Get(targetURL)
	if err != nil {
		result.Status = "error"
		result.ErrorMessage = err.Error()
		return result, err
	}
	defer resp.Body.Close()

	// Get the HTML version
	result.HTMLVersion = getHTMLVersion(resp)

	// Parse the HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		result.Status = "error"
		result.ErrorMessage = err.Error()
		return result, err
	}

	// Extract information from the parsed HTML
	links := extractInfo(doc, result)

	// Check the status of the links concurrently
	checkLinks(links, result)

	// Set the status to completed
	result.Status = "completed"
	result.UpdatedAt = time.Now()

	return result, nil
}

// extractInfo traverses the HTML document and extracts the required information.
func extractInfo(n *html.Node, result *models.CrawlResult) []string {
	var links []string
	if n.Type == html.ElementNode {
		switch n.Data {
		case "title":
			if n.FirstChild != nil {
				result.PageTitle = n.FirstChild.Data
			}
		case "h1", "h2", "h3", "h4", "h5", "h6":
			result.Headings[n.Data]++
		case "a":
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}
					baseURL, _ := url.Parse(result.URL)
					resolvedLink := baseURL.ResolveReference(link)
					links = append(links, resolvedLink.String())

					if link.Host == "" || link.Host == baseURL.Host {
						result.InternalLinksCount++
					} else {
						result.ExternalLinksCount++
					}
				}
			}
		case "form":
			checkForLoginForm(n, result)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, extractInfo(c, result)...)
	}
	return links
}

// checkForLoginForm checks if a form contains a password input field.
func checkForLoginForm(n *html.Node, result *models.CrawlResult) {
	if n.Type == html.ElementNode && n.Data == "input" {
		for _, attr := range n.Attr {
			if attr.Key == "type" && attr.Val == "password" {
				result.HasLoginForm = true
				return
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		checkForLoginForm(c, result)
	}
}

// getHTMLVersion tries to determine the HTML version from the doctype.
func getHTMLVersion(resp *http.Response) string {
	// This is a simplified approach. A more robust solution would involve
	// more complex parsing of the doctype.
	// For now, we will just check for HTML5 doctype.
	// A proper implementation would require a more sophisticated check.
	return "HTML5"
}

// checkLinks checks the status of a list of links concurrently.
func checkLinks(links []string, result *models.CrawlResult) {
	var wg sync.WaitGroup
	brokenLinksChan := make(chan map[string]any, len(links))

	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			resp, err := http.Head(link)
			if err != nil {
				// Handle error, maybe add to a separate list of failed requests
				return
			}
			if resp.StatusCode >= 400 {
				brokenLinksChan <- map[string]any{"url": link, "status_code": resp.StatusCode}
			}
		}(link)
	}

	wg.Wait()
	close(brokenLinksChan)

	for brokenLink := range brokenLinksChan {
		result.BrokenLinks = append(result.BrokenLinks, brokenLink)
	}
	result.InaccessibleLinksCount = len(result.BrokenLinks)
}