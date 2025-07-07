# Website Analyzer Client

This is the client-side application for the Website Analyzer, built with React and TypeScript. It provides a user interface to submit URLs for analysis, view the crawling results, and inspect detailed reports for each analyzed website.

## Features

*   **URL Submission:** Easily submit single URLs or multiple URLs in bulk for analysis.
*   **Crawl Results Overview:** View a list of all submitted URLs along with their high-level crawl status.
*   **Detailed Analysis:** Dive deep into individual website reports, including:
    *   HTML version detected.
    *   Presence of login forms.
    *   Counts of different heading tags (H1, H2, etc.).
    *   Distribution of internal vs. external links.
    *   Identification and status codes of broken links.
*   **Responsive Design:** Built with Tailwind CSS and Radix UI for a modern and responsive user experience.

## Technologies Used

*   **React 19:** A JavaScript library for building user interfaces.
*   **TypeScript:** A typed superset of JavaScript that compiles to plain JavaScript.
*   **Vite:** A fast build tool that provides an extremely fast development experience.
*   **Tailwind CSS:** A utility-first CSS framework for rapidly building custom designs.
*   **Radix UI:** A collection of unstyled, accessible UI components.
*   **React Router DOM:** For declarative routing in React applications.
*   **TanStack Query (React Query):** For data fetching, caching, and synchronization.
*   **Chart.js & React-Chartjs-2:** For rendering data visualizations, specifically pie charts for link distribution.
*   **Vitest:** A fast unit test framework powered by Vite.
*   **Biome:** A tool for formatting and linting code.

## Getting Started

### Prerequisites

*   Node.js (LTS version recommended)
*   npm (comes with Node.js)

### Installation

1.  Navigate to the `client` directory:
    ```bash
    cd client
    ```
2.  Install the dependencies:
    ```bash
    npm install
    ```

### Running the Development Server

To start the development server:

```bash
npm run dev
```

This will typically start the application on `http://localhost:5173` (or another available port).

### Building for Production

To build the application for production:

```bash
npm run build
```

This command compiles the application into the `dist` directory, ready for deployment.

## Linting and Type Checking

This project uses Biome for linting and formatting, and TypeScript for type checking.

*   **Run Linter and Type Checker:**
    ```bash
    npm run lint
    ```
*   **Run Linter and Fix Issues:**
    ```bash
    npm run lint:fix
    ```
*   **Run Type Checker only:**
    ```bash
    npm run typecheck
    ```

## Testing

Tests are written using Vitest.

*   **Run Tests:**
    ```bash
    npm run test
    ```
*   **Run Tests in Watch Mode:**
    ```bash
    npm run test:watch
    ```