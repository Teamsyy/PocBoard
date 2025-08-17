FROM golang:1.22-alpine AS development

WORKDIR /app

# Install air for hot reloading (compatible with Go 1.22)
RUN go install github.com/air-verse/air@v1.52.3

# Copy go mod files first
COPY go.mod ./
COPY go.sum* ./

# Download dependencies
RUN go mod download
RUN go mod tidy

# Copy source code
COPY . .

# Create uploads directory
RUN mkdir -p uploads

# Expose port
EXPOSE 8080

# For development, we'll use a simple approach - run go directly
CMD ["go", "run", "."]