# --- Build Stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary for the crawler-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /crawler-service ./cmd/crawler-service

# --- Final Stage ---
FROM alpine:latest

COPY --from=builder /crawler-service /crawler-service

ENTRYPOINT ["/crawler-service"]
