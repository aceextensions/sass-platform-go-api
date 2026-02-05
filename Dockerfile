# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy the workspace and modules
COPY go.work go.work.sum ./
COPY common/go.mod common/go.sum ./common/
COPY core/go.mod core/go.sum ./core/
COPY identity/go.mod identity/go.sum ./identity/
COPY api/go.mod api/go.sum ./api/

# Download dependencies
RUN go work sync

# Copy source code
COPY . .

# Build the application
RUN go build -o /app/api-bin ./api/cmd/api/main.go

# Production stage
FROM alpine:latest AS production

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/api-bin /app/api

# Set execution permissions
RUN chmod +x /app/api

# Expose port
EXPOSE 3001

# Command to run the application
CMD ["/app/api"]

# Development stage
FROM golang:1.24-alpine AS development

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Install Air for hot-reloading and Swag for docs
RUN go install github.com/air-verse/air@v1.61.7
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Expose port
EXPOSE 3001

# Command to run the application with Air
CMD ["air", "-c", "air.toml"]
