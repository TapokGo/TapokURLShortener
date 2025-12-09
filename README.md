# TapokURLShortener

A simple URL shortener service written in Go.  
Personal learning project — not intended for production use.

## Features

- Shorten long URLs via HTTP API
- Redirect from short code to original URL
- Data persistence with SQLite
- Structured logging to file
- Configurable via YAML

## Quick Start

### Prerequisites

- Go 1.22+

### Build & Run

```bash
make run
```
The server starts on http://localhost:8082.

## Usage
### Shorten a URL:

```bash
curl -X POST http://localhost:8082/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com"}'
```

Response:

```json
{"short_url":"http://localhost:8082/abc123"}
```

Visiting http://localhost:8082/abc123 redirects to the original URL.

## Configuration
The app loads configuration from:

- dev.yaml — local development config (git-ignored)
- .env.example.yaml — example production config

## Project Structure
```bash
internal/
├── config/      # Configuration loading
├── handler/     # HTTP handlers (v1 API)
├── logger/      # Structured logging
├── repo/        # SQLite storage
└── service/     # Business logic
```

## Makefile Commands
- make build — build binary
- make run — build and run
- make test — run tests
- make lint — run golangci-lint
- make help - show all commands
