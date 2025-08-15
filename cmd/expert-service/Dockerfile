# --- Build Stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary for the expert-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /expert-service ./cmd/expert-service

# --- Final Stage ---
FROM alpine:latest

COPY --from=builder /expert-service /expert-service

EXPOSE 50052

ENTRYPOINT ["/expert-service"]
