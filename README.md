# CodeMate - Personal Code Documentation Assistant

AI-powered assistant for querying and understanding your codebase.

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (для БД и сервисов)

### Installation

1. Clone the repository

```bash
git clone https://github.com/shester1kov/codemate.git
cd codemate
```

2. Install dependencies

```bash
make deps
```

3. Run the application

```bash
make run
```

The server will start on `http://localhost:8080`

### Endpoints

- `GET /health` - Health check
- `GET /ready` - Readiness check
- `POST /api/v1/query` - Query the codebase

### Example Request

```bash
curl -X POST http://localhost:8080/api/v1/query \
  -H "Content-Type: application/json" \
  -d '{"question": "How does the authentication work?"}'
```

## Project Structure

```text
codemate/
├── cmd/gateway/          # Application entry point
├── internal/             # Private application code
│   ├── config/          # Configuration
│   ├── logger/          # Logging setup
│   └── gateway/         # HTTP handlers, middleware, routing
├── configs/             # Config files
└── Makefile            # Build automation
```

## Development

- `make run` - Run the application
- `make build` - Build binary
- `make test` - Run tests
- `make lint` - Run linter
- `make fmt` - Format code

## Configuration

Configuration is managed via `configs/config.yaml` or environment variables.

Environment variables use the `CODEMATE_` prefix:

```bash
CODEMATE_SERVER_PORT=9090
CODEMATE_LOGGER_LEVEL=debug
```

## License

MIT
