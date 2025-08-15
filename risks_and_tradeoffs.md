# Risks and Technical Tradeoffs

This document discusses the key technical risks and the important architectural and design tradeoffs that the Portal project will face during its implementation.

---

## 1. Risk: Query Latency

High query latency is the most significant threat to user experience. The multi-layered, sequential nature of the query process is inherently slow.

*   **Mitigation Strategies:**
    *   **Aggressive Caching:** Caching will be critical at multiple levels.
        *   **Leaf Expert Responses:** Cache responses from Leaf Experts for common questions.
        *   **Synthesized Middleman Responses:** Cache the aggregated responses from Middlemen Experts.
        *   **RAG Context:** Cache the retrieved context from the vector database for frequent query patterns.
    *   **Extreme Parallelization:** The fan-out from Middlemen to Leaf Experts must be highly parallel. Golang's concurrency features are well-suited for this.
    *   **Optimized Communication:** Use gRPC for all internal service-to-service communication to minimize network overhead.
    *   **Model Selection:** For routing and simple synthesis tasks, consider using smaller, faster, and cheaper LLMs. Reserve the most powerful models for the final user-facing response generation.

## 2. Tradeoff: Cost vs. Quality

The extensive use of LLMs and embedding models will be expensive. Managing this cost requires careful tradeoffs.

*   **Decision: LLM Model Selection**
    *   **Tradeoff:** Using state-of-the-art models (like Gemini Advanced) for all tasks will produce the highest quality results but will be prohibitively expensive. Using smaller models (like Gemini) or even fine-tuned open-source models for specific tasks (like query routing) can dramatically reduce cost, but may lower accuracy.
    *   **Recommendation:** Implement a flexible system that allows different models to be used for different tasks. Start with cheaper models and strategically upgrade them where quality is most critical.

*   **Decision: RAG Content Threshold**
    *   **Tradeoff:** At what point does a page become "large" enough to warrant the cost and complexity of a RAG expert? Setting the threshold too low increases indexing costs (embedding and storage). Setting it too high risks feeding too much context into the LLM for "simple" experts, which is also expensive and can reduce quality.
    *   **Recommendation:** This threshold should be a configurable value that can be tuned based on cost-benefit analysis after launch.

## 3. Tradeoff: Vector Database: `pgvector` vs. Dedicated Solution

The choice of where to store vector embeddings is a key architectural decision.

*   **Option A: PostgreSQL with `pgvector`**
    *   **Pros:** Simplifies the tech stack (one less system to manage). Reduces operational overhead. Perfect for getting started and for moderately-sized datasets.
    *   **Cons:** May not offer the same performance or advanced features (like filtering and hybrid search) as dedicated solutions at extreme scale (billions of vectors).

*   **Option B: Dedicated Vector Database (e.g., Pinecone, Weaviate, Milvus)**
    *   **Pros:** Built from the ground up for vector search. Offers superior performance, scalability, and features for very large-scale applications.
    *   **Cons:** Adds another piece of infrastructure to the stack, increasing operational complexity and cost.

*   **Recommendation:** **Start with `pgvector`**. The `RAG Service` should be designed with a clear interface that abstracts the underlying vector store. This allows the team to begin with a simpler architecture and migrate to a dedicated solution in the future if performance benchmarks indicate it's necessary, without changing the other services.

## 4. Risk: Data Quality and LLM Hallucination

The system's credibility depends on the accuracy of its responses. There are two primary risks: low-quality source content and LLM "hallucinations."

*   **Mitigation Strategies:**
    *   **Source Vetting:** In the long run, the `Crawler/Discovery Service` may need to incorporate a "source quality" score to prioritize indexing authoritative domains.
    *   **Strict Prompt Engineering:** The prompts sent to Leaf Experts must be heavily constrained, instructing the model to *only* use the provided context (the webpage content) to answer the question.
    *   **Synthesis Guardrails:** The synthesis prompts for Middlemen Experts must be carefully crafted to prioritize faithful aggregation over creative "reasoning" that might introduce inaccuracies.
    *   **Emphasize Source Linking:** The UI must always make it easy for users to see the source URLs for the information in the response. This is the ultimate guardrail, as it allows users to verify the information themselves.

## 5. Tradeoff: Manual vs. Automated Hierarchy Management

The expert hierarchy is the backbone of the query routing system. How it's built and maintained is a critical decision.

*   **Option A: Manual/Curated Hierarchy**
    *   **Pros:** High accuracy and quality. The hierarchy will be logical and well-structured.
    *   **Cons:** Does not scale. Becomes a bottleneck as the number of indexed sites grows.

*   **Option B: Automated ML-based Hierarchy**
    *   **Pros:** The only viable long-term solution for scalability. Can adapt to the changing web automatically.
    *   **Cons:** A very complex ML project in its own right. Requires a separate pipeline for training and running classification models to place new Leaf Experts into the tree. Will make mistakes that could lead to poor query routing.

*   **Recommendation:** A hybrid approach. **Start with a manual, curated hierarchy** for a set of core, high-quality topics. Use this initial dataset to **bootstrap the development of an automated system**. The system could start as a "recommendation engine" for human curators before being allowed to operate fully autonomously.
