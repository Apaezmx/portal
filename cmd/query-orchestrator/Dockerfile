# --- Build Stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary for the query-orchestrator
RUN CGO_ENABLED=0 GOOS=linux go build -o /query-orchestrator ./cmd/query-orchestrator

# --- Final Stage ---
FROM alpine:latest

COPY --from=builder /query-orchestrator /query-orchestrator

EXPOSE 50053

ENTRYPOINT ["/query-orchestrator"]
