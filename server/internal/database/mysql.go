package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/krzysu/web-crawler/internal/models"
)

var db *sql.DB

// InitDB initializes the database connection.
func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}

// Close closes the database connection.
func Close() {
	db.Close()
}

// CreateCrawlResult inserts a new CrawlResult into the database.
func CreateCrawlResult(result *models.CrawlResult) error {
	query := `INSERT INTO crawl_results (id, url, status, page_title, html_version, headings_json, internal_links_count, external_links_count, inaccessible_links_count, broken_links_json, has_login_form, error_message, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := db.Exec(query, result.ID, result.URL, result.Status, result.PageTitle, result.HTMLVersion, result.Headings, result.InternalLinksCount, result.ExternalLinksCount, result.InaccessibleLinksCount, result.BrokenLinks, result.HasLoginForm, result.ErrorMessage, result.CreatedAt, result.UpdatedAt)
	return err
}

// GetCrawlResult retrieves a CrawlResult from the database by ID.
func GetCrawlResult(id string) (*models.CrawlResult, error) {
	result := &models.CrawlResult{}
	query := `SELECT id, url, status, page_title, html_version, headings_json, internal_links_count, external_links_count, inaccessible_links_count, broken_links_json, has_login_form, error_message, created_at, updated_at FROM crawl_results WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&result.ID, &result.URL, &result.Status, &result.PageTitle, &result.HTMLVersion, &result.Headings, &result.InternalLinksCount, &result.ExternalLinksCount, &result.InaccessibleLinksCount, &result.BrokenLinks, &result.HasLoginForm, &result.ErrorMessage, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateCrawlResult updates an existing CrawlResult in the database.
func UpdateCrawlResult(result *models.CrawlResult) error {
	query := `UPDATE crawl_results SET status = ?, page_title = ?, html_version = ?, headings_json = ?, internal_links_count = ?, external_links_count = ?, inaccessible_links_count = ?, broken_links_json = ?, has_login_form = ?, error_message = ?, updated_at = ? WHERE id = ?`
	_, err := db.Exec(query, result.Status, result.PageTitle, result.HTMLVersion, result.Headings, result.InternalLinksCount, result.ExternalLinksCount, result.InaccessibleLinksCount, result.BrokenLinks, result.HasLoginForm, result.ErrorMessage, time.Now(), result.ID)
	return err
}

// DeleteCrawlResult deletes a CrawlResult from the database.
func DeleteCrawlResult(id string) error {
	query := `DELETE FROM crawl_results WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// DeleteCrawlResults deletes multiple CrawlResults from the database.
func DeleteCrawlResults(ids []string) error {
	query := fmt.Sprintf("DELETE FROM crawl_results WHERE id IN (?%s)", strings.Repeat(",?", len(ids)-1))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	_, err := db.Exec(query, args...)
	return err
}

// GetCrawlResults retrieves a paginated list of CrawlResults.
func GetCrawlResults(limit, offset int, sortBy, filterBy string) ([]*models.CrawlResult, error) {
	query := fmt.Sprintf("SELECT id, url, status, page_title, html_version, headings_json, internal_links_count, external_links_count, inaccessible_links_count, broken_links_json, has_login_form, error_message, created_at, updated_at FROM crawl_results WHERE url LIKE ? ORDER BY %s LIMIT ? OFFSET ?", sortBy)
	rows, err := db.Query(query, "%"+filterBy+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []*models.CrawlResult{}
	for rows.Next() {
		result := &models.CrawlResult{}
		err := rows.Scan(&result.ID, &result.URL, &result.Status, &result.PageTitle, &result.HTMLVersion, &result.Headings, &result.InternalLinksCount, &result.ExternalLinksCount, &result.InaccessibleLinksCount, &result.BrokenLinks, &result.HasLoginForm, &result.ErrorMessage, &result.CreatedAt, &result.UpdatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}