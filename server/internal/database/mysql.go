package database

import (
	"fmt"
	"os"

	"github.com/krzysu/web-crawler/internal/models"
	gorm_mysql "gorm.io/driver/mysql"
	gorm_sqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the database connection pool.
type DB struct {
	db *gorm.DB
}

// NewDB creates a new DB instance with GORM.
func NewDB() (*DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	gormDB, err := gorm.Open(gorm_mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate will create or update the table based on the model.
	err = gormDB.AutoMigrate(&models.CrawlResult{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate database: %w", err)
	}

	return &DB{db: gormDB}, nil
}

// NewDBForTest creates a new DB instance for testing with SQLite.
func NewDBForTest() (*DB, error) {
	gormDB, err := gorm.Open(gorm_sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate will create or update the table based on the model.
	err = gormDB.AutoMigrate(&models.CrawlResult{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate database: %w", err)
	}

	return &DB{db: gormDB}, nil
}

// Close closes the database connection (not typically needed for GORM, but good practice).
func (d *DB) Close() error {
	// GORM manages its own connection pool, so explicit close might not be necessary
	// depending on the driver. For SQLite in-memory, it's effectively closed when process exits.
	// For MySQL, it might close the underlying sql.DB connection.
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// CreateCrawlResult inserts a new CrawlResult into the database.
func (d *DB) CreateCrawlResult(result *models.CrawlResult) error {
	return d.db.Create(result).Error
}

// GetCrawlResult retrieves a CrawlResult from the database by ID.
func (d *DB) GetCrawlResult(id uint) (*models.CrawlResult, error) {
	result := &models.CrawlResult{}
	err := d.db.First(result, "id = ?", id).Error
	return result, err
}

// UpdateCrawlResult updates an existing CrawlResult in the database.
func (d *DB) UpdateCrawlResult(result *models.CrawlResult) error {
	return d.db.Save(result).Error
}

// DeleteCrawlResult deletes a CrawlResult from the database.
func (d *DB) DeleteCrawlResult(id uint) error {
	return d.db.Delete(&models.CrawlResult{}, "id = ?", id).Error
}

// DeleteCrawlResults deletes multiple CrawlResults from the database.
func (d *DB) DeleteCrawlResults(ids []uint) error {
	return d.db.Delete(&models.CrawlResult{}, "id IN ?", ids).Error
}

// GetCrawlResults retrieves a paginated list of CrawlResults.
func (d *DB) GetCrawlResults(limit, offset int, sortBy, filterBy string) ([]*models.CrawlResult, error) {
	var results []*models.CrawlResult
	query := d.db.Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(sortBy)
	}

	if filterBy != "" {
		query = query.Where("url LIKE ?", "%"+filterBy+"%")
	}

	err := query.Find(&results).Error
	return results, err
}