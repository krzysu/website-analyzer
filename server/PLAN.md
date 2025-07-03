# Web Crawler Backend Development Plan

This document outlines the plan for developing the Go-based backend for the web crawler application.

## 1. Core Technology Stack

*   **Language:** Go (Golang)
*   **Web Framework:** Gin
    *   *Reasoning:* Gin is a high-performance, minimalist framework with a robust middleware ecosystem, making it ideal for building efficient APIs.
*   **Database:** MySQL
    *   *Driver:* `go-sql-driver/mysql`
*   **HTML Parsing:** `golang.org/x/net/html`
    *   *Reasoning:* A standard, powerful, and robust library for parsing HTML documents.
*   **Concurrency:** Go's built-in goroutines and channels for background crawling tasks.

## 2. Database Schema

We will use a single table to store the URL crawl data.

**Table: `crawl_results`**

| Column Name                 | Data Type      | Description                                                                                             |
| --------------------------- | -------------- | ------------------------------------------------------------------------------------------------------- |
| `id`                        | `VARCHAR(36)`  | **Primary Key.** UUID for the record.                                                                   |
| `url`                       | `TEXT`         | The target URL submitted for crawling.                                                                  |
| `status`                    | `VARCHAR(20)`  | The current status of the crawl (`queued`, `running`, `completed`, `error`).                            |
| `page_title`                | `VARCHAR(255)` | The extracted `<title>` of the page.                                                                      |
| `html_version`              | `VARCHAR(50)`  | The HTML version (e.g., "HTML 5", "XHTML 1.0").                                                         |
| `headings_json`             | `JSON`         | A JSON object storing the count of each heading level (e.g., `{"h1": 1, "h2": 4}`).                     |
| `internal_links_count`      | `INT`          | The total number of internal links found.                                                               |
| `external_links_count`      | `INT`          | The total number of external links found.                                                               |
| `inaccessible_links_count`  | `INT`          | The number of links that returned a 4xx or 5xx status code.                                             |
| `broken_links_json`         | `JSON`         | A JSON array of objects for inaccessible links (e.g., `[{"url": "...", "status_code": 404}]`).          |
| `has_login_form`            | `BOOLEAN`      | `true` if a form with a password input field is detected.                                               |
| `error_message`             | `TEXT`         | Stores any error message if the crawl fails.                                                            |
| `created_at`                | `TIMESTAMP`    | Timestamp when the URL was added.                                                                       |
| `updated_at`                | `TIMESTAMP`    | Timestamp when the record was last updated.                                                             |

## 3. Project Structure

```
/web-crawler-backend
├── cmd/
│   └── server/
│       └── main.go         // Entry point, Gin server setup
├── internal/
│   ├── api/                // Gin handlers, routes, and middleware
│   │   ├── handlers.go
│   │   └── routes.go
│   ├── crawler/            // Core crawling logic
│   │   └── crawler.go
│   ├── database/           // Database connection and queries
│   │   └── mysql.go
│   └── models/             // Go structs for our data
│       └── crawl.go
├── pkg/
│   └── utils/              // Shared utility functions
├── go.mod
├── go.sum
└── Dockerfile              // For containerization
```

## 4. Background Processing

Crawling is a long-running task and must not block API requests.
- A simple in-memory job queue will be implemented using a buffered channel.
- A pool of worker goroutines will listen on this channel.
- When a new URL is submitted, it's added to the job queue. A worker picks it up, performs the crawl, and updates the database with the results.

## 5. API Endpoints

The API will be designed to support all specified frontend features.

| Method | Path                       | Description                                                                                             |
| ------ | -------------------------- | ------------------------------------------------------------------------------------------------------- |
| `POST` | `/urls`                    | Adds a new URL to the processing queue. Body: `{"url": "..."}`.                                         |
| `GET`  | `/urls`                    | **(For Dashboard)** Retrieves a paginated, sortable, and filterable list of all crawl results.          |
| `GET`  | `/urls/{id}`               | **(For Details View)** Retrieves the full details for a single URL, including the broken links list.    |
| `DELETE`| `/urls`                   | **(For Bulk Actions)** Deletes multiple URLs. Body: `{"ids": ["id1", "id2"]}`.                          |
| `POST` | `/urls/rerun`              | **(For Bulk Actions)** Re-runs analysis on multiple URLs. Body: `{"ids": ["id1", "id2"]}`.               |
| `GET`  | `/urls/status`             | **(For Real-time Progress)** A potential endpoint for the frontend to poll for status updates on multiple URLs. |

*Note: True real-time updates would require WebSockets, which can be added as a future enhancement. Polling the main `/api/v1/urls` endpoint will provide sufficient status updates for the initial version.*

## 6. Implementation Steps

1.  **Setup Project:**
    *   Initialize Go module: `go mod init github.com/your-user/web-crawler`.
    *   Create the project directory structure as outlined above.

2.  **Database Layer (`internal/database`):**
    *   Implement functions to connect to MySQL.
    *   Write CRUD (Create, Read, Update, Delete) functions for the `crawl_results` table.
    *   The `List` function must support pagination, sorting, and filtering query parameters from the API.

3.  **Model (`internal/models`):**
    *   Define the `CrawlResult` struct that maps directly to the `crawl_results` database table.

4.  **Crawler Logic (`internal/crawler`):**
    *   Create a `Crawl` function that accepts a URL string.
    *   It will perform the following steps:
        a. Fetch the URL content using `net/http`.
        b. Parse the HTML using `golang.org/x/net/html`.
        c. Traverse the parsed document tree to extract:
            - Doctype (for HTML version).
            - Page title.
            - All heading tags (`<h1>`, `<h2>`, etc.).
            - All anchor tags (`<a>`).
        d. For each link, classify it as internal or external.
        e. Concurrently check each link's status code to identify inaccessible links.
        f. Check for the presence of `<form>` tags containing `<input type="password">`.
        g. Aggregate all collected data into a `CrawlResult` struct.
    *   This component will be called by the background workers.

5.  **API Layer (`internal/api`):**
    *   Set up a Gin router in `routes.go`.
    *   Define all endpoints from the table above.
    *   Implement handler functions in `handlers.go` for each route. These handlers will validate input and interact with the database layer and the job queue.

6.  **Main Application (`cmd/server/main.go`):**
    *   Initialize the database connection.
    *   Initialize and start the background worker pool and job queue.
    *   Set up and run the Gin HTTP server.

7.  **Containerization:**
    *   Write a `Dockerfile` to build a container image for the application, enabling easy deployment.
