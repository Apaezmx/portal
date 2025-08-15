# Crawler Service

The Crawler Service is a standalone binary that discovers and fetches content from the web.

## Responsibilities

-   Starts crawling from a given seed URL.
-   Follows links to discover new pages, staying within a configurable set of allowed domains.
-   Extracts the text content from each page.
-   Publishes the URL and its content to a NATS message queue (`crawled-content` subject) for the Indexing Job to process.

## Running the Service

To run the service locally:

```sh
go run ./cmd/crawler-service -nats-url="nats://localhost:4222" -allowed-domains="gocolly.dev" -start-url="http://gocolly.dev/"
```

When run inside Docker Compose, it uses the default values which point to the `nats` container.

## Building the Service

To build the binary:

```sh
go build -o crawler-service ./cmd/crawler-service
```
