FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o app ./cmd/main.go

EXPOSE 8006
CMD ["./app"]