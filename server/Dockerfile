# Use the official Go image as a base image
FROM golang:1.24.4-alpine AS base

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod .
COPY go.sum .

# Download Go modules
RUN go mod download

# Production stage
FROM base AS production

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o /website-analyzer-server ./cmd/server

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["/website-analyzer-server"]

# Development stage
FROM base AS development

# Install air for live reloading
RUN go install github.com/air-verse/air@latest

# Copy the rest of the application source code for development
COPY . .

# Expose the port the app runs on
EXPOSE 8080

# Command to run air for live reloading
CMD ["air"]
