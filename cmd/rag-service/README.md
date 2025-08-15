# RAG Service

The RAG Service is a gRPC service responsible for handling all operations related to Retrieval-Augmented Generation.

## Responsibilities

-   Receives content from the Indexing Job.
-   Performs text chunking and vectorization (currently mocked).
-   Stores text chunks and their vector embeddings in the PostgreSQL database.
-   Provides a gRPC endpoint for other services to retrieve relevant context chunks for a given query.

## Running the Service

To run the service locally (assuming the PostgreSQL database is running via Docker Compose):

```sh
go run ./cmd/rag-service -grpc-port=50051 -db-conn="host=localhost user=postgres password=postgres dbname=portal sslmode=disable"
```

Note: The `-db-conn` flag uses `localhost` because the service is running on the host machine, not within the Docker network. When run inside Docker Compose, it uses the default value which points to the `postgres` container.

## Building the Service

To build the binary:

```sh
go build -o rag-service ./cmd/rag-service
```
