# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**CodeMate** is a personal AI assistant for code documentation that allows querying codebases in natural language using local LLM models. Built as an HTTP gateway service using Go and Gin framework, integrating Qdrant (vector database) and Ollama (local LLM server) for RAG (Retrieval-Augmented Generation) pipeline.

### Key Features

- Codebase indexing (Go, Python, JS, etc.)
- Code search via vector database (RAG)
- Natural language Q&A about code using local LLM
- Documentation and explanation generation
- Code improvement suggestions

## Development Commands

### Building and Running

- `make run` - Run the application (starts server on localhost:8080)
- `make build` - Build binary to `bin/gateway`
- `make deps` - Install/update dependencies
- `make clean` - Remove build artifacts and coverage files

### Testing and Quality

- `make test` - Run all tests with race detection and coverage
- `make lint` - Run golangci-lint
- `make fmt` - Format code with `go fmt`

### Docker Infrastructure in future

- `make docker-up` - Start all Docker services (PostgreSQL, Qdrant, Ollama)
- `make docker-down` - Stop all Docker services
- `make docker-status` - Check service status
- `make docker-logs` - View service logs
- `make pull-model` - Download LLM model (llama3.2:3b)
- `make pull-embed-model` - Download embeddings model (nomic-embed-text)
- `./scripts/docker.sh debug` - Start services with Adminer web UI (port 8081)

## Architecture

### Project Structure

```text
cmd/gateway/          # Application entry point (main.go)
internal/             # Private application code
├── config/          # Configuration loading (Viper)
├── logger/          # Logging setup (zap)
└── gateway/         # HTTP layer
    ├── handler/     # Request handlers
    ├── middleware/  # HTTP middleware
    ├── router/      # Route setup
    └── dto/         # Data transfer objects
pkg/                 # Public reusable libraries
├── ollama/          # Ollama API client (generate, embed, ping)
└── qdrant/          # Qdrant client (upsert, search, delete, count)
configs/             # Configuration files (config.yaml)
deployments/docker/  # Docker Compose configuration
scripts/             # Automation scripts (docker.sh)
```

### Application Layers

**Entry Point** ([cmd/gateway/main.go](cmd/gateway/main.go)):

- Loads configuration from `configs/config.yaml`
- Initializes structured logger (zap)
- Sets up router with handlers and middleware
- Runs HTTP server with graceful shutdown (10s timeout)
- Handles SIGINT/SIGTERM signals

**Configuration** ([internal/config/config.go](internal/config/config.go)):

- Uses Viper for config management
- Loads from `configs/config.yaml` or environment variables
- Environment variables use `CODEMATE_` prefix (e.g., `CODEMATE_SERVER_PORT`)
- Sections: `server`, `logger`, `qdrant` (future), `ollama` (future)
- All settings have sensible defaults

**Router** ([internal/gateway/router/router.go](internal/gateway/router/router.go)):

- Gin framework in debug/release mode
- Middleware chain order is critical: Recovery → Logger → CORS
- Routes:
  - Health checks: `GET /health`, `GET /ready`
  - API v1: `/api/v1/query` (POST)
  - Planned: `/api/v1/index`, `/api/v1/status/:id`

**Handlers** ([internal/gateway/handler/](internal/gateway/handler/)):

- Constructor pattern: `NewXxxHandler(logger)` returns handler struct
- All handlers receive `*zap.Logger` for structured logging
- Use DTOs from `internal/gateway/dto` for request/response
- Gin binding validates required fields automatically
- Health handler checks service status
- Query handler: currently returns placeholder response, will implement RAG pipeline

**Middleware** ([internal/gateway/middleware/](internal/gateway/middleware/)):

- **Recovery**: Catches panics, logs them, returns 500 (must be first)
- **Logger**: Logs all requests with method, path, status, latency, IP, user agent
- **CORS**: Enables cross-origin requests

**Logger** ([internal/logger/logger.go](internal/logger/logger.go)):

- Built on uber-go/zap for structured, high-performance logging
- Levels: debug, info, warn, error
- Encodings: json (production), console (development)
- Outputs to stdout by default, configurable to files

**External Services** ([pkg/](pkg/)):

- **Ollama Client** ([pkg/ollama/client.go](pkg/ollama/client.go)): Generate text, create embeddings, health checks
- **Qdrant Client** ([pkg/qdrant/client.go](pkg/qdrant/client.go)): Vector CRUD operations, cosine similarity search
- All clients have unit tests and error handling

### Current Architecture (MVP)

```text
┌─────────────┐
│  Go Gateway │  ← HTTP API (Gin)
│   (main)    │
└──┬─────┬────┘
   │     │
   │     └──────────┐
   │                │
┌──▼──────┐   ┌────▼────┐
│ Ollama  │   │ Qdrant  │
│ (LLM)   │   │(Vectors)│
└─────────┘   └─────────┘
```

### Planned Microservices Architecture

```text
┌────────────────┐
│   OpenWebUI    │  ← Web interface
└───────┬────────┘
        │ HTTP
┌───────▼─────────┐
│  API Gateway    │  ← Routing, auth
└─┬──────┬───┬────┘
  │gRPC  │   │gRPC
┌─▼──┐ ┌─▼─┐ ┌▼───┐
│Idx │ │Qry│ │Chat│  ← Microservices
└─┬──┘ └─┬─┘ └┬───┘
  │   ┌──▼───▼──┐
  │   │ Qdrant  │
  │   └─────────┘
  └──►┌──────────┐
      │PostgreSQL│  ← Metadata
      └──────────┘
```

### Key Patterns

1. **Dependency Injection**: Logger is injected into all handlers and middleware
2. **Constructor Functions**: All components use `NewXxx()` constructors
3. **Graceful Shutdown**: 10-second timeout for in-flight requests
4. **Middleware Ordering**: Recovery must be first to catch panics in other middleware
5. **RESTful Versioning**: API routes under `/api/v1`
6. **DTO Pattern**: Separate request/response types in `dto` package
7. **Error Wrapping**: Use `fmt.Errorf("context: %w", err)` for error chains
8. **Context Propagation**: Pass `context.Context` through service layers for cancellation

### Configuration

Default values are set in [internal/config/config.go:69](internal/config/config.go#L69):

- Server: localhost:8080, debug mode
- Logger: info level, json encoding, stdout
- Qdrant: <http://localhost:6333>, collection "codemate_code"
- Ollama: <http://localhost:11434>, model llama3.2:3b

Override via environment variables:

```bash
CODEMATE_SERVER_PORT=9090
CODEMATE_LOGGER_LEVEL=debug
CODEMATE_QDRANT_URL=http://qdrant:6333
CODEMATE_OLLAMA_MODEL=llama3.2:latest
```

### Service Ports

- **8080** - Go API Gateway
- **5432** - PostgreSQL (metadata storage)
- **6333** - Qdrant HTTP API
- **6334** - Qdrant gRPC API
- **11434** - Ollama API
- **8081** - Adminer (debug mode only)

## Technology Stack

### Backend

- **Go 1.25+** with Gin (HTTP framework)
- **Viper** for configuration (YAML + env vars)
- **Zap** for structured logging
- **GORM** (planned) for ORM
- **golang-migrate** (planned) for database migrations
- **Wire** (planned) for dependency injection

### AI/ML

- **Ollama** - Local LLM server (llama3.2:3b model, 3B parameters)
- **nomic-embed-text** - Embeddings model
- **Qdrant** - Vector database for code embeddings
- **Cosine similarity** for vector search

### Infrastructure

- **Docker + Docker Compose** for service orchestration
- **PostgreSQL 16** for metadata (connected but not yet used in code)
- **Adminer** for database management (debug mode)

## Implementation Status

### Completed (Weeks 1-2)

- Project structure following Go best practices
- Configuration system (Viper with YAML and env vars)
- Structured logging (Zap with JSON/console formats)
- HTTP API (Gin with middleware: Recovery, Logger, CORS)
- Health check endpoints (`/health`, `/ready`)
- Graceful shutdown with signal handling

### In Progress (Week 3)

- Docker Compose infrastructure (PostgreSQL, Qdrant, Ollama)
- Ollama client ([pkg/ollama/client.go](pkg/ollama/client.go)) with Generate, Embed, Ping
- Qdrant client ([pkg/qdrant/client.go](pkg/qdrant/client.go)) with Upsert, Search, Delete, Count
- Unit tests for both clients
- Docker management script ([scripts/docker.sh](scripts/docker.sh))
- Docker infrastructure setup and testing

### Then

- Code parser (Go AST - Abstract Syntax Tree)
- Indexer service for code chunking

### Planned

- **Week 4**: RAG Pipeline (Retriever, Prompt engineering, Query handler integration)
- **Weeks 5-7**: Microservices split (Gateway, Indexer, Query, Chat services via gRPC)
- **Weeks 8-9**: OpenWebUI integration, n8n automation
- **Weeks 10-12**: Production features (Keycloak auth, Swagger docs, monitoring)

### RAG Pipeline Implementation

Query handler at [internal/gateway/handler/query.go:54](internal/gateway/handler/query.go#L54) currently returns placeholder. Future implementation:

1. **Retrieve**: User question → Ollama embedding → Qdrant search → Top-K code chunks
2. **Augment**: Format retrieved code + question into prompt
3. **Generate**: Send prompt to Ollama → Stream/return answer with sources

## API Testing

### Health Checks

```bash
# Liveness probe
curl http://localhost:8080/health
# Response: {"status":"ok","version":"1.0.0"}

# Readiness probe
curl http://localhost:8080/ready
# Response: {"status":"ok","version":"1.0.0"}
```

### Query Endpoint (Placeholder)

```bash
curl -X POST http://localhost:8080/api/v1/query \
  -H "Content-Type: application/json" \
  -d '{"question": "How does authentication work?"}'

# Current response (placeholder):
# {
#   "answer": "THIS IS A PLACEHOLDER RESPONSE, RAG PIPELINE WILL BE IMPLEMENTED LATER",
#   "sources": [{"file_path":"example/main.go","name":"main","type":"function","score":0.95}]
# }
```

### Direct Service Testing

```bash
# Test Ollama
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2:3b",
  "prompt": "Say hello",
  "stream": false
}'

# Test Qdrant (dashboard)
# Open http://localhost:6333/dashboard in browser
```

## Common Issues

- **Lint cache conflicts**: Makefile removes `golangci-lint.lock` before linting to avoid lock issues
- **Port conflicts**: Change `CODEMATE_SERVER_PORT` if 8080 is in use
- **Cache directories**: Uses user's home directory for Go caches (`~/.cache/go-build`, `~/go/pkg/mod`)
- **GPU Memory**: GTX 1050 Ti (4GB VRAM) limits model size - 3B parameter models work, 7B+ will use CPU
- **Model downloads**: First `make pull-model` downloads ~2GB, takes time depending on connection
