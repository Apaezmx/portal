# Portal System Architecture

This document outlines the high-level architecture of the Portal search engine. The system is designed as a set of distributed microservices that work together to create, manage, and query a hierarchy of AI experts.

## 1. The Hierarchical Expert Model

The core of Portal is its unique hierarchical model of AI experts. This model is designed to provide both broad, high-level answers and deep, specific knowledge. The hierarchy consists of three types of experts:

### 1.1. Root Expert
*   **Singleton:** There is only one Root Expert in the system.
*   **Function:** It serves as the main entry point for all search queries. Its primary responsibility is to understand the user's intent and route the query to the most relevant Middlemen Experts.
*   **Knowledge:** The Root Expert has a high-level understanding of all the topics covered by the Middlemen Experts beneath it.

### 1.2. Middlemen Experts
*   **Thematic Grouping:** Middlemen Experts are organized around specific topics or themes (e.g., "Social Networks," "Cloud Computing," "Open Source Projects").
*   **Function:** A Middleman Expert receives a query from the Root Expert (or another Middleman) and breaks it down. It then routes these sub-queries to the relevant Leaf Experts or other, more specialized Middlemen Experts under its purview. Once it receives answers from the downstream experts, it synthesizes them into a coherent response.
*   **Structure:** This creates a tree-like structure, allowing for multiple layers of Middlemen Experts for complex topics.

### 1.3. Leaf Experts
*   **Source of Truth:** Leaf Experts are the foundation of the hierarchy. Each Leaf Expert is an AI agent focused on a single webpage.
*   **Function:** Its purpose is to answer questions about the content of that specific page. Users can also interact with a Leaf Expert directly (e.g., via `portal.com/e/example.com`).
*   **Two Implementations:**
    1.  **Simple Context:** For small, simple pages, the Leaf Expert is a standard Gemini model call where the entire page's content is passed as context.
    2.  **RAG-based:** For larger, more complex pages, the Leaf Expert is powered by a Retrieval-Augmented Generation (RAG) system. The page's content is chunked, vectorized, and stored in a database (Postgres with `pgvector`). When a query is received, the most relevant chunks are retrieved and passed to the Gemini model as context.

## 2. Query Flow

Here is the step-by-step flow of a typical user query:

1.  **Request Entry:** A user submits a query to `portal.com`. The request hits the **API Gateway**, which authenticates and forwards it to the **Query Orchestrator**.
2.  **Root Expert Analysis:** The Query Orchestrator, acting as the **Root Expert**, analyzes the query to determine its primary intent.
3.  **Delegation to Middlemen:** The Root Expert identifies the relevant **Middlemen Experts** (e.g., for a query like "best social network for developers," it might query the "Social Networks" and "Software Development" Middlemen). It forwards the query to them.
4.  **Fan-out to Leaf Experts:** Each Middleman Expert further analyzes the query and determines which of its **Leaf Experts** are most likely to have the answer. It then sends parallel requests to these Leaf Experts.
5.  **Leaf Expert Response:**
    *   If the Leaf Expert is simple, it injects the page content and the query into a prompt for the LLM.
    *   If it's a RAG expert, it first queries the vector database to find relevant text chunks, then injects those chunks and the query into a prompt.
    *   The LLM generates an answer, which is returned by the Leaf Expert.
6.  **Synthesis and Aggregation:**
    *   The Middlemen Experts collect the responses from the Leaf Experts. They synthesize these individual answers into a more comprehensive response.
    *   The Root Expert collects the responses from the Middlemen Experts and performs a final aggregation to create the final search result.
7.  **Response to User:** The final, aggregated response is sent back to the user through the API Gateway.

## 3. System Components

The architecture is composed of several key microservices, each with a distinct responsibility.

*   **API Gateway:** The single public entry point for all incoming traffic. It handles routing, rate limiting, and authentication.
*   **Query Orchestrator:** A service that embodies the Root and Middlemen Experts. It manages the logic of query decomposition, routing, and response synthesis.
*   **Expert Service:** Manages the Leaf Experts. It exposes an interface to query a specific expert (e.g., `/expert/{url}`). It determines whether to use a simple context or RAG-based approach.
*   **RAG Service:** Provides the backend for RAG-based Leaf Experts. It handles text chunking, vectorization, storage in Postgres, and retrieval.
*   **Crawler/Discovery Service:** A background service that crawls the web to find new pages and content updates.
*   **Indexing Job:** An asynchronous job that takes crawled content and triggers the Expert Service to create or update a Leaf Expert.

This decoupled architecture allows for independent scaling and development of each component. Communication between services will be handled via a combination of synchronous (e.g., gRPC/REST) and asynchronous (e.g., message queue) protocols.
