# AI Agents CLI Docker Image

# Build stage
FROM golang:1.25.2-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o ai-agents-cli \
    -ldflags="-s -w -X main.Version=${VERSION:-dev} -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.GitCommit=${GIT_COMMIT:-unknown}" \
    .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/ai-agents-cli /app/ai-agents-cli

# Copy example config
COPY --from=builder /build/env.example /app/env.example

# Create non-root user
RUN addgroup -g 1000 cli && \
    adduser -D -u 1000 -G cli cli && \
    chown -R cli:cli /app

USER cli

# Expose CLI
ENTRYPOINT ["/app/ai-agents-cli"]

