# Indexing Job

The Indexing Job is a worker service that listens for messages from the Crawler Service and processes them for indexing.

## Responsibilities

-   Connects to the NATS message queue and subscribes to the `crawled-content` subject.
-   Receives messages containing crawled webpage content.
-   Determines whether the content is suitable for a "simple" expert or a "RAG" expert based on its length.
-   Calls the `CreateOrUpdateExpert` RPC on the Expert Service to trigger the creation or update of the corresponding AI expert.

## Running the Service

To run the service locally:

```sh
go run ./cmd/indexing-job -nats-url="nats://localhost:4222" -expert-svc-addr="localhost:50052"
```

When run inside Docker Compose, it uses the default values which point to the `nats` and `expert-service` containers.

## Building the Service

To build the binary:

```sh
go build -o indexing-job ./cmd/indexing-job
```
