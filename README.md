# Website Analyzer

This project consists of a client-side application and a Go backend for analyzing websites.

## Running the Backend

### Production Mode

To run the backend in production mode (without live-reloading):

```bash
docker-compose up --build backend
```

### Development Mode (with Live-Reloading)

For development, you can use the `backend-dev` service which includes live-reloading. Any changes to the Go source code in the `server` directory will automatically trigger a recompile and restart of the backend.

```bash
docker-compose up --build backend-dev
```

## Database

The project uses a MySQL database. The `docker-compose.yml` sets up a `db` service for this purpose.

## Client

(Add client-side instructions here)
