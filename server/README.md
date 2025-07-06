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

Set the following environment variables for database connection. You can create a `.env` file in the project root, and the application will automatically load these variables.

- `DB_USER`: Your MySQL username.
- `DB_PASSWORD`: Your MySQL password.
- `DB_HOST`: Your MySQL host (e.g., `127.0.0.1` or `localhost`).
- `DB_PORT`: Your MySQL port (e.g., `3306`).
- `DB_NAME`: The name of your database (e.g., `crawler_db`).
- `PORT`: The port the application will run on (e.g., `8080`).
- `API_KEY`: A secret key required for authenticating API requests. Generate a strong, random key.

### 3. Running the Application

#### 1. Using Docker Compose

For easier setup with MySQL, a `docker-compose.yml` file is provided in the project root.

1.  **Run Docker Compose (Production Mode):**

    ```bash
    docker-compose up --build backend
    ```

    This will build the backend image, start the MySQL container, and then start the backend application.

2.  **Run Docker Compose (Development Mode with Live-Reloading):**

    For development, you can use the `backend-dev` service which includes live-reloading. Any changes to the Go source code in the `server` directory will automatically trigger a recompile and restart of the backend.

    ```bash
    docker-compose up --build backend-dev
    ```

### 4. API Endpoints

The backend exposes the following RESTful API endpoints:

- **`POST /urls`**

  - **Description:** Adds a new URL to the queue for analysis.
  - **Request Body:** `{"url": "http://example.com"}`
  - **Example:** `curl -X POST -H "Content-Type: application/json" -d '{"url": "http://example.com"}' http://localhost:8080/urls`

- **`GET /urls`**

  - **Description:** Retrieves a paginated, sortable, and filterable list of all analyzed URLs and their crawl results.
  - **Example:** `curl http://localhost:8080/urls`

- **`GET /urls/:id`**

  - **Description:** Retrieves detailed information for a single crawl result by its ID.
  - **Example:** `curl http://localhost:8080/urls/123`

- **`DELETE /urls`**

  - **Description:** Deletes multiple crawl results.
  - **Request Body:** `{"ids": [1, 2, 3]}`
  - **Example:** `curl -X DELETE -H "Content-Type: application/json" -d '{"ids": [1, 2]}' http://localhost:8080/urls`

- **`POST /urls/rerun`**
  - **Description:** Re-runs analysis on multiple URLs by their IDs.
  - **Request Body:** `{"ids": [1, 2, 3]}`
  - **Example:** `curl -X POST -H "Content-Type: application/json" -d '{"ids": [1, 2]}' http://localhost:8080/urls/rerun`

### 5. Testing

To run the tests for the backend, navigate to the `server` directory and execute:

```bash
go test ./...
```

### 6. Linting

This project uses `golangci-lint` for linting. To run the linter, ensure you have it installed and then run:

```bash
golangci-lint run
```
