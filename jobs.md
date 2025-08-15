# Portal Asynchronous Jobs

This document describes the background jobs that run asynchronously to support the Portal ecosystem. These jobs are essential for content ingestion, expert creation, and data maintenance. They are designed to be resilient and scalable.

---

## 1. Indexing Job

*   **Purpose:** To process newly crawled web content and create or update the corresponding Leaf Experts. This job is the bridge between the **Crawler/Discovery Service** and the **Expert Service**.

*   **Trigger:** This job is triggered by messages appearing in a specific topic on the **Message Queue**. The **Crawler/Discovery Service** places a message on this queue every time it successfully fetches and extracts content from a webpage.

*   **Workflow:**
    1.  **Consume Message:** The job worker pulls a message from the queue. The message contains the URL and the extracted text content of a webpage.
    2.  **Determine Expert Type:** The job analyzes the size of the content.
        *   If the content is small (below a certain threshold), it's marked for a "simple expert."
        *   If the content is large, it's marked for a "RAG expert."
    3.  **Create/Update Simple Expert:** If the expert is simple, the job makes an API call to an internal endpoint on the **Expert Service**. This call includes the URL and the full content. The Expert Service then saves this information to the database.
    4.  **Trigger RAG Pipeline:** If the expert is to be RAG-based, the job makes an API call to the **RAG Service**. This call includes the URL and the content. The RAG Service then initiates the process of chunking, vectorizing, and storing the content in the vector database.
    5.  **Update Expert Metadata:** Once the RAG pipeline is complete (which might be an asynchronous process itself), the **Expert Service** is notified to update the expert's metadata in the database, marking it as an active RAG expert.
    6.  **Acknowledge Message:** After successfully processing the content and creating/updating the expert, the job acknowledges the message on the queue to remove it.

*   **Interactions:**
    *   **Message Queue:** Consumes messages produced by the **Crawler/Discovery Service**.
    *   **Expert Service:** Calls its internal API to create/update expert metadata and store content for simple experts.
    *   **RAG Service:** Calls its internal API to initiate the RAG indexing pipeline for large content.

---

## 2. RAG Refresh Job

*   **Purpose:** To periodically update the content of existing RAG-based Leaf Experts to prevent the information from becoming stale. Data freshness is critical for a search engine.

*   **Trigger:** This job runs on a schedule (e.g., once a day, once a week). It can be managed by a cron scheduler like Kubernetes CronJob or a similar system.

*   **Workflow:**
    1.  **Identify Stale Experts:** The job queries the **Postgres Database** to find RAG experts that haven't been updated in a certain amount of time (e.g., `last_updated_at < now() - 7 days`).
    2.  **Request Re-crawl:** For each stale expert identified, the job sends a message to the **Crawler/Discovery Service** (potentially via a dedicated message queue topic) requesting a high-priority re-crawl of the expert's URL.
    3.  **Crawler Fetches New Content:** The crawler picks up the request, re-crawls the page, and publishes the new content to the indexing queue, just like it would for a newly discovered page.
    4.  **Indexing Job Updates Expert:** The **Indexing Job** consumes the message with the updated content and follows its standard workflow. This will transparently update the RAG expert with the new information.

*   **Interactions:**
    *   **Postgres Database:** Queries the database to find stale experts.
    *   **Crawler/Discovery Service:** Sends requests (likely via a message queue) to re-crawl specific URLs.

This two-job system decouples the process of finding content from the process of indexing it, allowing each part to be scaled and managed independently.
