# Website Analyzer

Website Analyzer is a full-stack application designed to crawl and analyze websites, providing detailed reports on various aspects such as HTML version, link distribution, heading structure, and broken links.

## Features

- **URL Submission:** Submit single or bulk URLs for comprehensive analysis.
- **Detailed Crawl Reports:** Get insights into:
  - HTML version and page title.
  - Presence of login forms.
  - Counts of different heading tags (H1, H2, etc.).
  - Distribution of internal and external links.
  - Identification of broken links with their status codes.
- **User-Friendly Interface:** A responsive web interface for easy interaction and visualization of results.

## Architecture

The project consists of two main components:

1.  **Client (Frontend):** A React application that provides the user interface for submitting URLs and viewing analysis reports.
2.  **Server (Backend):** A Go application that handles URL crawling, data processing, and serves the API for the client.
3.  **Database:** A MySQL database for storing crawl results.

## Technologies Used

### Client

- **React 19**
- **TypeScript**
- **Vite**
- **Tailwind CSS**
- **Radix UI**
- **TanStack Query**
- **Chart.js**

### Server

- **Go**
- **Gin Web Framework**
- **MySQL**
- **Docker**

## Authentication

API requests to the backend are authenticated using an API Key. The server expects an `X-API-Key` header with a valid API key. This key is configured via the `API_KEY` environment variable.

## Further Improvements / Missing Features

The following features or areas could be further improved or are currently missing:

- **URL Management:** Implementation for starting/stopping processing on selected URLs. <- server is so fast that I don't see a need for this
- **Real-Time Progress:** Displaying real-time crawl status updates (queued → running → done / error) for URLs in the dashboard. <- polling every 5sec takes care of status updates
- **Dashboard Features:** Full implementation of sorting and column filters.
- **Automated Front-end Tests:** Expanding test coverage for happy-path scenarios in the client application.

## Getting Started

To get the Website Analyzer up and running on your local machine, follow these steps:

### Prerequisites

- [Docker](https://www.docker.com/get-started) (includes Docker Compose)

### Setup and Run

1.  Clone the repository:

    ```bash
    git clone https://github.com/krzysu/website-analyzer.git
    cd website-analyzer
    ```

2.  Start the Docker containers:

    - **Development Mode (with live reloading for server):**

      ```bash
      docker-compose up --build backend-dev
      ```

      This will start the database and the server in development mode with live reloading.

    - **Production Mode (for client and server):**

      ```bash
      docker-compose up --build backend
      ```

      This will start the database and the server without live reloading.

3.  Access the application:
    - First, navigate to the `client` directory:
      ```bash
      cd client
      ```
    - Then, install the client dependencies:
      ```bash
      npm install
      ```
    - Finally, start the client development server:
      ```bash
      npm run dev
      ```
    - The client application will then be available at `http://localhost:5173` (or another available port).
    - The server API will be available at `http://localhost:8080`.

## More Information

For more detailed information on the client and server applications, including development, testing, and linting instructions, please refer to their respective README files:

- [`client/README.md`](./client/README.md)
- [`server/README.md`](./server/README.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
