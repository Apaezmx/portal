# API Gateway

The API Gateway is the single public entry point for all traffic to the Portal system. It is a standard HTTP server that translates public REST-like requests into backend gRPC calls.

## Responsibilities

-   Serves the static frontend application (HTML, CSS, JS).
-   Handles incoming HTTP requests for `/search` and forwards them as gRPC calls to the Query Orchestrator Service.
-   Handles incoming HTTP requests for `/e/{url}` and forwards them as gRPC calls to the Expert Service.

## Running the Service

To run the service locally:

```sh
go run ./cmd/api-gateway -http-port=8080 -expert-svc-addr="localhost:50052" -orch-svc-addr="localhost:50053"
```

When run inside Docker Compose, it uses the default values which point to the `expert-service` and `query-orchestrator` containers.

## Building the Service

To build the binary:

```sh
go build -o api-gateway ./cmd/api-gateway
```
