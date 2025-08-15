# Query Orchestrator Service

The Query Orchestrator is the central "brain" of the Portal search engine. It is a gRPC service that receives search queries and coordinates with the network of experts to generate an answer.

## Responsibilities

-   Exposes a gRPC `Search` endpoint.
-   Receives search queries from the API Gateway.
-   Interprets the query and determines which experts to consult (currently a simplified, hardcoded logic).
-   Calls the `QueryExpert` RPC on the Expert Service for each required expert.
-   Synthesizes the answers from the experts into a final response (currently a simplified, passthrough logic).

## Running the Service

To run the service locally:

```sh
go run ./cmd/query-orchestrator -grpc-port=50053 -expert-svc-addr="localhost:50052"
```

When run inside Docker Compose, it uses the default values which point to the `expert-service` container.

## Building the Service

To build the binary:

```sh
go build -o query-orchestrator ./cmd/query-orchestrator
```
