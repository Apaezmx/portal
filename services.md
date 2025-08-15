# Portal Microservices

This document provides a detailed description of each microservice within the Portal architecture. All services will be developed in Golang to ensure high performance and concurrency.

---

## 1. API Gateway

*   **Purpose:** To provide a single, unified entry point for all client requests. It simplifies the client-side by abstracting the internal microservice architecture.

*   **Key Features:**
    *   **Request Routing:** Routes incoming HTTP requests to the appropriate backend service (e.g., `/search` to Query Orchestrator, `/e/{url}` to Expert Service).
    *   **Authentication & Authorization:** Validates user credentials or API keys.
    *   **Rate Limiting:** Protects the system from abuse and ensures fair usage.
    *   **SSL Termination:** Handles HTTPS and decrypts traffic before forwarding it to internal services.
    *   **Request/Response Logging:** Centralized logging for all incoming and outgoing traffic.

*   **Interactions:**
    *   Receives requests from the public internet (users).
    *   Forwards requests to the **Query Orchestrator** and **Expert Service**.

*   **Technology Stack:**
    *   A dedicated API Gateway solution like Kong, Traefik, or a custom-built Go service using a library like `chi` or `gorilla/mux` with appropriate middleware.

---

## 2. Query Orchestrator

*   **Purpose:** To act as the brain of the search process, embodying the Root and Middlemen Experts. It interprets user queries and orchestrates the process of gathering and synthesizing information.

*   **Key Features:**
    *   **Query Analysis:** Uses an LLM to understand the user's query and identify the relevant topics.
    *   **Query Routing:** Determines which Middlemen or Leaf Experts to query based on the analysis.
    *   **Parallel Fan-out:** Sends multiple concurrent requests to downstream experts to minimize latency.
    *   **Response Synthesis:** Receives responses from multiple experts and uses an LLM to aggregate and synthesize them into a single, coherent answer.
    *   **State Management:** Manages the hierarchy of experts (which middleman knows about which topics/experts).

*   **Interactions:**
    *   Receives requests from the **API Gateway**.
    *   Sends requests to other instances of itself (for Middleman-to-Middleman communication) and to the **Expert Service**.

*   **Technology Stack:**
    *   **Golang:** For the core application logic.
    *   **gRPC:** For efficient, low-latency communication with the Expert Service.
    *   **Gemini API:** For query analysis and response synthesis.

---

## 3. Expert Service

*   **Purpose:** To manage the lifecycle and querying of all Leaf Experts. It abstracts the complexity of whether an expert is simple or RAG-based.

*   **Key Features:**
    *   **Expert Retrieval:** Fetches expert metadata from the database based on a URL.
    *   **Query Execution:**
        *   For **simple experts**, it formats a prompt with the full page content and the user's query.
        *   For **RAG experts**, it first queries the **RAG Service** to get relevant context before formatting the prompt.
    *   **LLM Interaction:** Calls the Gemini API with the prepared prompt to get an answer.
    *   **Expert Creation/Update:** Exposes an internal API that is called by the **Indexing Job**. This endpoint takes page content and creates the necessary expert metadata and, if needed, triggers the RAG pipeline.

*   **Interactions:**
    *   Receives requests from the **Query Orchestrator** and **API Gateway** (for direct expert access).
    *   Communicates with the **RAG Service** to get context for RAG-based experts.
    *   Interacts with the **Postgres Database** to fetch expert metadata.
    *   Receives instructions from the **Indexing Job** to create/update experts.

*   **Technology Stack:**
    *   **Golang**
    *   **gRPC** for its internal API.
    *   **Gemini API**

---

## 4. RAG Service

*   **Purpose:** To handle all the operations required for Retrieval-Augmented Generation. This service is a specialized backend for the **Expert Service**.

*   **Key Features:**
    *   **Text Processing:** Chunks large documents into smaller, manageable pieces.
    *   **Vectorization:** Converts text chunks into vector embeddings using a text embedding model.
    *   **Vector Storage:** Stores the embeddings in a Postgres database with the `pgvector` extension.
    *   **Similarity Search:** Given a query vector, it performs a similarity search on the stored vectors to find the most relevant text chunks.

*   **Interactions:**
    *   Receives requests from the **Expert Service** to retrieve relevant chunks for a given query.
    *   Receives instructions from the **Indexing Job** to process and store content for a new or updated RAG expert.
    *   Heavily interacts with the **Postgres Database**.

*   **Technology Stack:**
    *   **Golang**
    *   **gRPC** for its internal API.
    *   **Text Embedding Model API** (e.g., Google's `text-embedding-004`).
    *   **PostgreSQL** with `pgvector` extension.

---

## 5. Crawler/Discovery Service

*   **Purpose:** To discover and fetch content from the web to be indexed by Portal.

*   **Key Features:**
    *   **URL Frontier Management:** Manages a queue of URLs to be crawled.
    *   **Webpage Fetching:** Downloads the HTML content of webpages.
    *   **Content Extraction:** Parses the HTML to extract the main text content, stripping out boilerplate like ads, nav bars, and footers.
    *   **Link Discovery:** Discovers new URLs from the links on crawled pages and adds them to the frontier.
    *   **Politeness Policy:** Respects `robots.txt` and implements rate limiting to avoid overwhelming servers.

*   **Interactions:**
    *   Publishes crawled content (URL and extracted text) to a **Message Queue** (e.g., RabbitMQ, Kafka).
    *   Does not receive direct requests from other services; it runs continuously in the background.

*   **Technology Stack:**
    *   **Golang** (with libraries like `colly` for crawling).
    *   **Message Queue** (RabbitMQ, Kafka, or NATS).
