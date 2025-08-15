# --- Build Stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the binary for the indexing-job
RUN CGO_ENABLED=0 GOOS=linux go build -o /indexing-job ./cmd/indexing-job

# --- Final Stage ---
FROM alpine:latest

COPY --from=builder /indexing-job /indexing-job

ENTRYPOINT ["/indexing-job"]
