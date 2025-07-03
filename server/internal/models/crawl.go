package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// CrawlResult represents the data collected from a crawled URL.

type CrawlResult struct {
	ID                     string    `json:"id"`
	URL                    string    `json:"url"`
	Status                 string    `json:"status"`
	PageTitle              string    `json:"page_title"`
	HTMLVersion            string    `json:"html_version"`
	Headings               JSONMap   `json:"headings"`
	InternalLinksCount     int       `json:"internal_links_count"`
	ExternalLinksCount     int       `json:"external_links_count"`
	InaccessibleLinksCount int       `json:"inaccessible_links_count"`
	BrokenLinks            JSONArray `json:"broken_links"`
	HasLoginForm           bool      `json:"has_login_form"`
	ErrorMessage           string    `json:"error_message"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// JSONMap is a custom type for handling JSON map[string]int in MySQL.
type JSONMap map[string]int

// Value implements the driver.Valuer interface for JSONMap.
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONMap.
func (j *JSONMap) Scan(src interface{}) error {
	if src == nil {
		*j = make(map[string]int)
		return nil
	}
	s, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	return json.Unmarshal(s, j)
}

// JSONArray is a custom type for handling JSON array of objects in MySQL.
type JSONArray []map[string]any

// Value implements the driver.Valuer interface for JSONArray.
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONArray.
func (j *JSONArray) Scan(src interface{}) error {
	if src == nil {
		*j = make([]map[string]any, 0)
		return nil
	}
	s, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	return json.Unmarshal(s, j)
}
