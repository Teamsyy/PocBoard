FROM golang:1.22-alpine AS development

# Install dependencies for health checks and development
RUN apk add --no-cache wget git

WORKDIR /app

# Install air for hot reloading (compatible with Go 1.22)
RUN go install github.com/air-verse/air@v1.52.3

# Copy go mod files first for better Docker layer caching
COPY go.mod ./
COPY go.sum* ./

# Download dependencies
RUN go mod download
RUN go mod tidy

# Copy source code
COPY . .

# Create uploads directory with proper permissions
RUN mkdir -p uploads/boards && chmod -R 755 uploads

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=10s --timeout=5s --start-period=20s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:8080/health || exit 1

# For development, use air for hot reloading if available, otherwise go run
CMD ["sh", "-c", "if command -v air >/dev/null 2>&1; then air; else go run .; fi"]

# Production stage
FROM golang:1.22-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod ./
COPY go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final production stage
FROM alpine:latest AS production

# Install ca-certificates for HTTPS and wget for health checks
RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/main .

# Create uploads directory
RUN mkdir -p uploads/boards && chmod -R 755 uploads

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# Change ownership of the app directory
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=10s --timeout=5s --start-period=20s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:8080/health || exit 1

# Run the application
CMD ["./main"]