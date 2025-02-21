# Language Learning Backend

This is a Go-based backend server for a language learning application that serves as:
- Vocabulary inventory
- Learning Record Store (LRS)
- Unified launch pad for learning apps

## Prerequisites

- Go 1.21 or later
- SQLite3
- Mage (Go task runner)

## Installation

1. Install Go:
```bash
sudo apt update
sudo apt install golang-go
```

2. Install SQLite3:
```bash
sudo apt install sqlite3
```

3. Install project dependencies:
```bash
go mod download
```

4. Install Mage:
```bash
go install github.com/magefile/mage@latest
```

## Project Structure

```
backend_go/
├── main.go              # Main application entry point
├── go.mod              # Go module file
├── magefile.go         # Mage task definitions
├── internal/           # Internal packages
│   ├── database/       # Database related code
│   ├── db/            # Database migrations and seeds
│   ├── handlers/      # HTTP handlers
│   ├── models/        # Data models
│   ├── repositories/  # Data access layer
│   └── services/      # Business logic
└── pkg/               # Public packages
    ├── ginutils/      # Gin framework utilities
    └── sqliteutils/   # SQLite utilities
```

## Development

To run the server:
```bash
mage run
```

To build the application:
```bash
mage build
```

To run tests:
```bash
mage test
```

## API Endpoints

- GET /api/dashboard/last_study_session
- GET /api/dashboard/study_progress
- GET /api/dashboard/quick_stats
- GET /api/study_activities/:id
- GET /api/study_activities/:id/study_sessions
- POST /api/study_activities
- GET /api/words

For detailed API documentation, refer to the Backend-Technical-Specs.md file.
