# Expert Service

The Expert Service is a gRPC service that manages the lifecycle and querying of all Leaf Experts. It abstracts the complexity of whether an expert is simple or RAG-based.

## Responsibilities

-   Exposes a gRPC API to create, update, and query experts.
-   Receives instructions from the Indexing Job to create or update experts.
-   Coordinates with the RAG Service to index content for large pages.
-   Responds to queries from the Query Orchestrator by either retrieving simple content from the database or by querying the RAG service for context.

## Running the Service

To run the service locally (assuming the database and RAG service are running):

```sh
go run ./cmd/expert-service -grpc-port=50052 -db-conn="host=localhost user=postgres password=postgres dbname=portal sslmode=disable" -rag-svc-addr="localhost:50051"
```

When run inside Docker Compose, it uses the default values which point to the `postgres` and `rag-service` containers.

## Building the Service

To build the binary:

```sh
go build -o expert-service ./cmd/expert-service
```
