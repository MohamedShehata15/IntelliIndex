FROM golang:1.24

WORKDIR /app

# Install development tools
RUN go install github.com/cosmtrek/air@v1.49.0 && \
    go install github.com/go-delve/delve/cmd/dlv@latest

# Install system dependencies including PostgreSQL client
RUN apt-get update && apt-get install -y --no-install-recommends \
    postgresql-client \
    curl \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# Set environment variables
ENV CGO_ENABLED=1 \
    GO111MODULE=on

# Create .air.toml config in /tmp directory instead of /app
RUN mkdir -p /tmp/air && \
    echo '[build]\n  cmd = "go build -o /app/tmp/main ./cmd/api"\n  bin = "/app/tmp/main"\n  delay = 1000\n  exclude_dir = ["assets", "tmp", "vendor"]\n  include_ext = ["go", "tpl", "tmpl", "html"]\n  exclude_regex = ["_test\\.go"]\n[misc]\n  clean_on_exit = true' > /tmp/air/.air.example.toml

# Make sure tmp directory exists
RUN mkdir -p /app/tmp

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# HEALTHCHECK for container health monitoring
HEALTHCHECK --interval=30s --timeout=3s CMD curl -f http://localhost:8080/health || exit 1

# air will watch for file changes and automatically rebuild
CMD ["sh", "-c", "if [ ! -f .air.toml ]; then cp /tmp/air/.air.example.toml .air.toml; fi && air -c .air.toml"]
