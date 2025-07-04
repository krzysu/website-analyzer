# Use the official Go image as a base image
FROM golang:1.24.4-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod .
COPY go.sum .

# Download Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o /website-analyzer-server ./cmd/server

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["/website-analyzer-server"]
