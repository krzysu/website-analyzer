# Web Crawler Backend

This is the backend service for a web crawler application, built with Go and Gin, using MySQL for data storage.

## Features

*   Accepts a website URL for crawling.
*   Extracts key information from crawled pages:
    *   HTML version
    *   Page title
    *   Count of heading tags (H1, H2, etc.)
    *   Number of internal vs. external links
    *   Number of inaccessible links (4xx or 5xx status codes)
    *   Presence of a login form
*   Provides RESTful API endpoints for:
    *   Adding new URLs for analysis.
    *   Retrieving paginated, sortable, and filterable crawl results.
    *   Retrieving detailed information for a single crawl result.
    *   Deleting multiple crawl results.
    *   Re-running analysis on multiple URLs.
*   Background processing of crawl jobs using a worker pool.

## Technologies Used

*   **Language:** Go (Golang)
*   **Web Framework:** Gin
*   **Database:** MySQL
*   **HTML Parsing:** `golang.org/x/net/html`
*   **Concurrency:** Go Goroutines and Channels

## Setup and Running

### 1. Prerequisites

*   Go (1.16 or higher)
*   MySQL Server
*   Docker (optional, for containerized deployment)

### 2. Database Setup

First, you need to set up your MySQL database. Create a database named `crawler_db` (or whatever you configure in `main.go`).

```sql
CREATE DATABASE IF NOT EXISTS crawler_db;
USE crawler_db;

CREATE TABLE IF NOT EXISTS crawl_results (
    id VARCHAR(36) PRIMARY KEY,
    url TEXT NOT NULL,
    status VARCHAR(20) NOT NULL,
    page_title VARCHAR(255),
    html_version VARCHAR(50),
    headings_json JSON,
    internal_links_count INT,
    external_links_count INT,
    inaccessible_links_count INT,
    broken_links_json JSON,
    has_login_form BOOLEAN,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

Update the database connection string in `cmd/server/main.go`:

```go
// Initialize the database connection
if err := database.InitDB("user:password@tcp(127.0.0.1:3306)/crawler_db"); err != nil {
    log.Fatalf("Failed to connect to database: %v", err)
}
```

Replace `user:password@tcp(127.0.0.1:3306)/crawler_db` with your MySQL credentials and host.

### 3. Running the Application

#### A. Locally (without Docker)

1.  **Navigate to the project directory:**
    ```bash
    cd /Users/krzysu/work/web-crawler/server
    ```

2.  **Download Go modules:**
    ```bash
    go mod tidy
    ```

3.  **Run the application:**
    ```bash
    go run cmd/server/main.go
    ```

    The server will start on `http://localhost:8080`.

#### B. Using Docker

1.  **Build the Docker image:**
    ```bash
    docker build -t web-crawler-backend .
    ```

2.  **Run the Docker container:**
    ```bash
    docker run -p 8080:8080 --name web-crawler-app web-crawler-backend
    ```

    *Note: Ensure your MySQL server is accessible from within the Docker container if it's running on a different host or network. You might need to link containers or use the host network mode.*