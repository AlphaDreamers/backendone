FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy all files (including config.yaml)
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/auth-service ./cmd/main.go

# Final stage
FROM alpine:3.18
WORKDIR /app

# Copy certificates and config
COPY --from=builder /app/certificates /app/certificates
COPY --from=builder /app/config.yaml /app/config.yaml
COPY --from=builder /app/bin/auth-service /app/auth-service

# Install timezone data if needed
RUN apk add --no-cache tzdata

EXPOSE 8004
CMD ["/app/auth-service"]