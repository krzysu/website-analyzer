package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// CrawlResult represents the data collected from a crawled URL.

type CrawlResult struct {
	ID                     uint `gorm:"primarykey"`
	CreatedAt              time.Time `gorm:"autoCreateTime"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime"`
	URL                    string `gorm:"type:text"`
	Status                 string `gorm:"type:varchar(20)"`
	PageTitle              string `gorm:"type:varchar(255)"`
	HTMLVersion            string `gorm:"type:varchar(50)"`
	Headings               JSONMap `gorm:"type:json"`
	InternalLinksCount     int
	ExternalLinksCount     int
	InaccessibleLinksCount int
	BrokenLinks            JSONArray `gorm:"type:json"`
	HasLoginForm           bool
	ErrorMessage           string `gorm:"type:text"`
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
