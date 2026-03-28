# Build stage
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# Disable CGO
ENV CGO_ENABLED=0

# Install build dependencies
RUN apk add --no-cache git
ENV GOPRIVATE=github.com/aceextension/*

# 1. Mirror Host Structure for Workspace
# Copy go.work and module definitions first for caching
COPY sass-platform-go-api/go.work sass-platform-go-api/go.work.sum ./sass-platform-go-api/
COPY sass-platform-go-api/common/go.mod sass-platform-go-api/common/go.sum ./sass-platform-go-api/common/
COPY sass-platform-go-api/core/go.mod sass-platform-go-api/core/go.sum ./sass-platform-go-api/core/
COPY sass-platform-go-api/ace_packages/identity/go.mod sass-platform-go-api/ace_packages/identity/go.sum ./sass-platform-go-api/ace_packages/identity/
COPY sass-platform-go-api/ace_packages/notification/go.mod sass-platform-go-api/ace_packages/notification/go.sum ./sass-platform-go-api/ace_packages/notification/
COPY sass-platform-go-api/ace_packages/subscription/go.mod sass-platform-go-api/ace_packages/subscription/go.sum ./sass-platform-go-api/ace_packages/subscription/
COPY sass-platform-go-api/api/go.mod sass-platform-go-api/api/go.sum ./sass-platform-go-api/api/


# Sync workspace (this uses the mirrored structure)
WORKDIR /app/sass-platform-go-api
RUN go work sync

# 2. Copy Full Source
WORKDIR /app
COPY sass-platform-go-api ./sass-platform-go-api

# 3. Build the application
WORKDIR /app/sass-platform-go-api
RUN go build -buildvcs=false -o /app/api-bin ./api/cmd/api/main.go

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
FROM golang:1.25.6-alpine AS development

# Mirror Host Structure for Development
WORKDIR /app
COPY sass-platform-go-api ./sass-platform-go-api

WORKDIR /app/sass-platform-go-api

# Install build dependencies
RUN apk add --no-cache git

# Install Air for hot-reloading and Swag for docs
RUN go install github.com/air-verse/air@v1.61.7
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Expose port
EXPOSE 3001

# Command to run the application with Air
CMD ["air", "-c", "air.toml"]
