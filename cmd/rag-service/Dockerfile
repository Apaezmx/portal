# --- Build Stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /rag-service ./cmd/rag-service

# --- Final Stage ---
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /rag-service /rag-service

# Expose the gRPC port
EXPOSE 50051

# Command to run the service
ENTRYPOINT ["/rag-service"]
