# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum ./

COPY authService/go.mod authService/go.sum ./
RUN go mod tidy

COPY common /app/common
COPY db /app/db
COPY authService /app/authService

COPY ./authService/internal/config/config.development.yaml /app/authService/config.yaml
WORKDIR /app/authService


RUN go work vendor
RUN go build -mod=vendor -o /bin/authService

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/authService/config.yaml .

COPY --from=builder /bin/authService .

CMD ["./authService"]
