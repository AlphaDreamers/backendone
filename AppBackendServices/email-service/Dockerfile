# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy the Go workspace and module files
COPY go.mod go.sum ./

# Ensure that Go modules are downloaded
RUN go mod tidy

# Copy the entire project
COPY . .
COPY .env /app/.env

# Vendor the dependencies
RUN go mod vendor

# Build the binary
RUN go build -mod=vendor -o /bin/email-service

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/.env .
COPY --from=builder /bin/email-service .

# Command to run the service
CMD ["./email-service"]
