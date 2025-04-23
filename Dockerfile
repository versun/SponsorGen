# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Ensure config.yaml exists (for debugging)
RUN ls -la && echo "Checking for config.yaml:" && ls -la config.yaml

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o sponsorgen .

# Final stage - using smaller alpine image
FROM alpine:3.19

# Install necessary packages (ca-certificates for HTTPS connections)
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user to run the application
RUN adduser -D -h /app appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/sponsorgen .

# Copy configuration and asset files
COPY --from=builder /app/config.yaml ./config.yaml
COPY --from=builder /app/assets ./assets

# Create necessary directories with proper permissions
RUN mkdir -p ./output ./cache && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose the default port (match with docker-compose and workflow)
EXPOSE 5000

# Add health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:5000/sponsors.json || exit 1

# Command to run
ENTRYPOINT ["./sponsorgen"]
# Default arguments, can be overridden at runtime
CMD ["-config", "./config.yaml", "-port", "5000"]