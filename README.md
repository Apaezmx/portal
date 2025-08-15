# Portal: An AI-Powered Search Engine

Portal is a next-generation search engine built on a novel architecture of interconnected AI "experts." The goal is to provide users with not only broad search results but also the ability to dive deep into specific topics and websites with conversational AI agents.

## Vision

The internet is a vast repository of knowledge, but traditional search engines often only scratch the surface. Portal's vision is to create a more intelligent way to explore the web. We do this by creating a specialized AI expert for every webpage. These experts can be queried directly, providing in-depth, context-aware answers that are impossible to achieve with a standard search query.

By organizing these leaf experts into a hierarchy of "middlemen" experts, Portal can also answer broad questions by intelligently aggregating and synthesizing information from multiple sources. This creates a powerful, scalable system for knowledge discovery.

## Core Concepts

*   **AI Experts:** The fundamental building block of Portal. Each expert is an AI agent with deep knowledge of a specific webpage or a broader topic.
*   **Leaf Experts:** An expert for a single webpage. For small pages, this is a Gemini model with the page's content as context. For larger pages, it's a RAG-powered agent.
*   **Middlemen Experts:** Thematic experts that know about a group of leaf experts or other middlemen experts. They delegate and synthesize information to answer broader queries.
*   **Hierarchical Structure:** A tree of experts with a single root, allowing for efficient query routing and aggregation.
*   **Asynchronous Indexing:** A robust system of background jobs for crawling the web, creating experts, and keeping them up-to-date.

## Documentation

This repository contains the initial design and architecture documentation for the Portal project. The documents herein describe the proposed services, data models, APIs, and potential challenges.

*   `architecture.md`: The high-level system architecture.
*   `services.md`: Detailed descriptions of each microservice.
*   `jobs.md`: Information on the asynchronous background jobs.
*   `interfaces.md`: API contracts and data formats.
*   `data_models.md`: The proposed database schema.
*   `pros_and_cons.md`: An analysis of the pros and cons of this approach.
*   `risks_and_tradeoffs.md`: A discussion of technical risks and tradeoffs.
