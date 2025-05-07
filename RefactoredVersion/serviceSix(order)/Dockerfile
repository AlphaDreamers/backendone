FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o app ./cmd/main.go

EXPOSE 8005
CMD ["./app"]