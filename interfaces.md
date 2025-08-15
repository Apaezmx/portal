# API and Data Interfaces

This document defines the API contracts and data models for communication between the services in the Portal ecosystem. While internal communication will be implemented using gRPC for performance, the interfaces are described here in a REST-like fashion for clarity.

---

## 1. Public-Facing API (via API Gateway)

These are the endpoints exposed to the public internet.

### `POST /search`
*   **Description:** The main endpoint for submitting a search query to the Portal engine.
*   **Request Body:** `SearchQuery` object.
*   **Response Body:** `SearchResponse` object.

### `POST /e/{url}`
*   **Description:** Allows a user to interact directly with a specific Leaf Expert for a given URL. The URL should be Base64 encoded to be URL-safe.
*   **Request Body:** `ExpertQuery` object.
*   **Response Body:** `ExpertResponse` object.

---

## 2. Internal Service APIs (gRPC)

These are the interfaces used for communication between microservices.

### Expert Service

#### `rpc QueryExpert(ExpertQueryRequest) returns (ExpertResponse)`
*   **Equivalent to:** `POST /internal/experts/{url}/query`
*   **Description:** Called by the **Query Orchestrator** to get a response from a specific Leaf Expert.
*   **Request Body:** `ExpertQuery` object.
*   **Response Body:** `ExpertResponse` object.

#### `rpc CreateOrUpdateExpert(ExpertCreationRequest) returns (Empty)`
*   **Equivalent to:** `POST /internal/experts`
*   **Description:** Called by the **Indexing Job** to create a new expert or update an existing one.
*   **Request Body:** `ExpertCreationRequest` object.
*   **Response Body:** An acknowledgment of success or failure.

### RAG Service

#### `rpc IndexContent(IndexRequest) returns (Empty)`
*   **Equivalent to:** `POST /internal/rag/index`
*   **Description:** Called by the **Indexing Job** to start the RAG pipeline for a large piece of content.
*   **Request Body:** `IndexRequest` object.
*   **Response Body:** An acknowledgment that the indexing process has begun.

#### `rpc RetrieveContext(RetrievalRequest) returns (RetrievalResponse)`
*   **Equivalent to:** `POST /internal/rag/retrieve`
*   **Description:** Called by the **Expert Service** to get relevant text chunks for a query from a RAG-based expert.
*   **Request Body:** `RetrievalRequest` object.
*   **Response Body:** `RetrievalResponse` object.

---

## 3. Asynchronous Communication (Message Queue)

### `crawled-content` Topic
*   **Description:** A message queue topic where the **Crawler/Discovery Service** publishes content for the **Indexing Job** to consume.
*   **Message Body:** `CrawledContentMessage` object.

---

## 4. Core Data Objects

Below are the primary data structures used in the API calls.

### `SearchQuery`
```json
{
  "query": "What are the best practices for writing microservices in Golang?",
  "user_id": "user-12345"
}
```

### `SearchResponse`
```json
{
  "summary": "Golang is well-suited for microservices due to its performance, concurrency model, and strong standard library. Best practices include single responsibility, decentralized data management, and using gRPC for communication.",
  "sources": [
    {
      "url": "https://awesome-go.com/",
      "title": "Awesome Go",
      "snippet": "A curated list of awesome Go frameworks, libraries and software."
    },
    {
      "url": "https://go.dev/blog/using-go-modules",
      "title": "Using Go Modules",
      "snippet": "This post is an introduction to the basics of using Go modules."
    }
  ]
}
```

### `ExpertQuery`
```json
{
  "query": "How do I use the `pgvector` extension with this library?",
  "conversation_history": [
    { "role": "user", "content": "What is this page about?" },
    { "role": "assistant", "content": "This page is the official documentation for the GORM database library for Golang." }
  ]
}
```

### `ExpertResponse`
```json
{
  "answer": "To use `pgvector` with GORM, you'll need to define a custom data type for the `vector` type and then use it in your model struct. You can then use raw SQL queries with the `<=>` operator for similarity search.",
  "source_url": "https://gorm.io/docs/generic_interface.html"
}
```

### `CrawledContentMessage`
```json
{
  "url": "https://example.com/some-article",
  "content": "This is the full, extracted text content of the article...",
  "crawled_at": "2024-10-26T10:00:00Z"
}
```

### `ExpertCreationRequest`
```json
{
  "url": "https://example.com/some-article",
  "content": "This is the full, extracted text content of the article...",
  "expert_type": "simple" // or "rag"
}
```

### `IndexRequest`
```json
{
  "url": "https://example.com/large-article",
  "content": "This is the very large body of text to be indexed..."
}
```

### `RetrievalRequest`
```json
{
  "url": "https://example.com/large-article",
  "query": "What does the author say about performance?"
}
```

### `RetrievalResponse`
```json
{
  "context_chunks": [
    "chunk 1 text...",
    "chunk 5 text...",
    "chunk 12 text..."
  ]
}
```
