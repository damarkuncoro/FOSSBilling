# Multi-stage Dockerfile for FOSSBilling Go Backend
# --------------------------------------------------------
# Build Stage
# --------------------------------------------------------
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates build-base

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build optimized binaries
# -ldflags="-w -s" reduces binary size by removing debug info
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/worker ./cmd/worker && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/cli ./cmd/cli

# --------------------------------------------------------
# Runtime Stage
# --------------------------------------------------------
FROM alpine:3.21

# Security: Add a non-root user
RUN addgroup -S fossbilling && adduser -S fossbilling -G fossbilling

# Install essential runtime tools
RUN apk --no-cache add ca-certificates tzdata postgresql-client

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/bin/worker /app/worker
COPY --from=builder /app/bin/cli /app/cli

# Copy essential assets
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/docs /app/docs
COPY --from=builder /app/extensions /app/extensions

# Set permissions
RUN chown -R fossbilling:fossbilling /app
USER fossbilling

# Default port
EXPOSE 8080

# The worker can be started by overriding the CMD in docker-compose
CMD ["/app/api"]
