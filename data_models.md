# Portal Database Schema (PostgreSQL)

This document outlines the proposed database schema for Portal, using PostgreSQL. The schema is designed to store expert metadata, the relationships between experts, and the data required for Retrieval-Augmented Generation (RAG).

We recommend using the [`pgvector`](https://github.com/pgvector/pgvector) extension for PostgreSQL to handle the storage and querying of vector embeddings efficiently.

---

## Table `experts`

This is the central table that stores information about every AI expert in the system, whether it's a Root, Middleman, or Leaf expert.

```sql
CREATE TYPE expert_type AS ENUM ('ROOT', 'MIDDLEMAN', 'LEAF');

CREATE TABLE experts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- The type of expert
    type expert_type NOT NULL,
    -- The name of the expert (e.g., "Social Networks Expert" or the URL for a leaf)
    name TEXT NOT NULL,
    -- For LEAF experts, the URL of the page they are an expert on.
    url TEXT UNIQUE,
    -- For LEAF experts, determines if it uses RAG or simple context.
    is_rag_based BOOLEAN NOT NULL DEFAULT FALSE,
    -- For simple LEAF experts, the full content of the page is stored here.
    raw_content TEXT,
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for quick lookup of leaf experts by URL
CREATE INDEX idx_experts_url ON experts(url);
```

**Notes:**
*   A `url` is only present for `LEAF` experts.
*   For simple (non-RAG) `LEAF` experts, the entire page content is stored in `raw_content`. For RAG experts, this field would be `NULL`.

---

## Table `expert_hierarchy`

This table defines the tree structure of the experts. It uses an adjacency list model to store the parent-child relationships.

```sql
CREATE TABLE expert_hierarchy (
    parent_expert_id UUID NOT NULL REFERENCES experts(id) ON DELETE CASCADE,
    child_expert_id UUID NOT NULL REFERENCES experts(id) ON DELETE CASCADE,
    PRIMARY KEY (parent_expert_id, child_expert_id)
);

-- Index for finding children of a given parent quickly
CREATE INDEX idx_expert_hierarchy_parent ON expert_hierarchy(parent_expert_id);
```

**Example:**
If the "Social Networks" Middleman expert has child Leaf Experts for "twitter.com" and "linkedin.com", this table would have two rows linking the parent's ID to each child's ID.

---

## Table `document_chunks`

This table stores the vectorized data for RAG-based Leaf Experts. Each row represents a small chunk of text from the source webpage.

```sql
-- Make sure the pgvector extension is installed
-- CREATE EXTENSION vector;

CREATE TABLE document_chunks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Foreign key linking this chunk to its corresponding LEAF expert
    expert_id UUID NOT NULL REFERENCES experts(id) ON DELETE CASCADE,
    -- The actual text content of the chunk
    chunk_text TEXT NOT NULL,
    -- The vector embedding of the chunk_text. Assuming an embedding dimension of 768.
    embedding vector(768) NOT NULL,
    -- Optional: a sequential index of the chunk within the document
    chunk_index INTEGER NOT NULL,
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index on the expert_id for quick retrieval of all chunks for a page
CREATE INDEX idx_document_chunks_expert_id ON document_chunks(expert_id);

-- An IVFFlat index for fast approximate nearest neighbor search on the embeddings.
-- This is crucial for RAG performance. The number of lists (100) is a parameter
-- that should be tuned based on the size of the dataset.
CREATE INDEX ON document_chunks USING ivfflat (embedding vector_l2_ops) WITH (lists = 100);
```

**Note on Crawled Content:**
We have made a design decision *not* to have a separate, persistent table for all raw crawled content. Raw content is transiently handled by the `Indexing Job`. It is either stored directly in `experts.raw_content` for simple experts or processed and stored in `document_chunks` for RAG experts. This approach avoids data duplication and significantly reduces storage costs.
