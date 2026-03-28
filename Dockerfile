# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o cinemasys ./cmd/cinemasys
# Final stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/cinemasys .
COPY --from=builder /app/web ./web
COPY --from=builder /app/assets ./assets

# Expose port
EXPOSE 8080

# Run the application
CMD ["./cinemasys"]
