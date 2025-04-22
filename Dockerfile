# Build stage
FROM golang:1.19-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sponsorgen .

# Final stage
FROM alpine:latest

# Install necessary packages
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/sponsorgen .

# Copy configuration and asset files
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/assets ./assets

# Create necessary directories
RUN mkdir -p ./output ./cache

# Expose the default port
EXPOSE 8000

# Command to run
ENTRYPOINT ["./sponsorgen"]
# Default arguments, can be overridden at runtime
CMD ["-config", "config.yaml", "-port", "8000"]