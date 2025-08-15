# --- Build Stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary for the api-gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway ./cmd/api-gateway

# --- Final Stage ---
FROM alpine:latest

COPY --from=builder /api-gateway /api-gateway

# Expose the HTTP port
EXPOSE 8080

ENTRYPOINT ["/api-gateway"]
