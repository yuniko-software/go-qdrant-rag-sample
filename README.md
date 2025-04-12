# go-qdrant-rag-sample

Semantic Search and Retrieval-Augmented Generation (RAG) in Go using Qdrant and OpenAI

This repository contains a Go application that demonstrates semantic product search and Retrieval-Augmented Generation using OpenAI's GPT and embedding models, with Qdrant as the vector database. The solution ingests structured data from a CSV file, generates embeddings, stores them in Qdrant, and provides REST API endpoints for semantic search and LLM-based answers.

This application simulates a product search system for an online catalog using OpenAI's embedding and LLM models. It supports vector-based semantic search, RAG-style question answering, and automatic ingestion of CSV-based product data into Qdrant.

## Technologies Used

- Go 1.24.0 (go version go1.24.0 windows/amd64)
- Qdrant (vector database)
- OpenAI API (text-embedding-ada-002 and gpt-4o-2024-08-06)
- Fiber (HTTP web framework)
- Resty (HTTP client for Go)

## Environment Setup

Create an environment file at `env/.env` with the following content:

OPENAI_API_KEY=`your-api-key-here` 

QDRANT_HOST=`http://localhost:6333`


This file is used to configure the application for accessing OpenAI and Qdrant.

## Running the Application

### Step 1: Start Qdrant

Start only the Qdrant container using the following command:

docker-compose up -d qdrant


### Step 2: Run the Go Application

From the project root, run:

go run ./cmd/api


This command initializes the Qdrant collection if it does not exist, checks whether the collection is empty, ingests the dataset from `data/products.csv` if necessary, and starts the REST API server on `http://localhost:8080`.

## API Endpoints


Set breakpoints and launch the debugger using Run > Start Debugging.

## API Endpoints

| Method | Endpoint             | Description                              |
|--------|----------------------|------------------------------------------|
| GET    | /search?q=...        | Semantic product search                  |
| GET    | /rag?q=...&top=3     | RAG response with OpenAI and Qdrant      |


## Ingestion Behavior

- Products are sourced from `data/products.csv`.
- At startup, the application checks if the collection in Qdrant is empty.
- If empty, products are embedded and inserted.
- If not, ingestion is skipped to avoid duplicates.











