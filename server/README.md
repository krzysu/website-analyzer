# Web Crawler Backend

This is the backend service for a web crawler application, built with Go and Gin, using MySQL for data storage.

## Features

- Accepts a website URL for crawling.
- Extracts key information from crawled pages:
  - HTML version
  - Page title
  - Count of heading tags (H1, H2, etc.)
  - Number of internal vs. external links
  - Number of inaccessible links (4xx or 5xx status codes)
  - Presence of a login form
- Provides RESTful API endpoints for:
  - Adding new URLs for analysis.
  - Retrieving paginated, sortable, and filterable crawl results.
  - Retrieving detailed information for a single crawl result.
  - Deleting multiple crawl results.
  - Re-running analysis on multiple URLs.
- Background processing of crawl jobs using a worker pool.

## Technologies Used

- **Language:** Go (Golang)
- **Web Framework:** Gin
- **Database:** MySQL (with GORM)
- **HTML Parsing:** `golang.org/x/net/html`
- **Concurrency:** Go Goroutines and Channels

## Setup and Running

### 1. Prerequisites

- Go (1.16 or higher)
- MySQL Server
- Docker (optional, for containerized deployment)

### 2. Database Setup

First, you need to set up your MySQL database. Create a database named `crawler_db` (or whatever you configure in `main.go`). GORM will handle the table creation and migration automatically based on the `models.CrawlResult` struct.

Set the following environment variables for database connection. You can create a `.env` file in the project root and use a tool like `godotenv` to load these variables.

- `DB_USER`: Your MySQL username.
- `DB_PASSWORD`: Your MySQL password.
- `DB_HOST`: Your MySQL host (e.g., `127.0.0.1` or `localhost`).
- `DB_PORT`: Your MySQL port (e.g., `3306`).
- `DB_NAME`: The name of your database (e.g., `crawler_db`).

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

    _Note: Ensure your MySQL server is accessible from within the Docker container if it's running on a different host or network. You might need to link containers or use the host network mode._
