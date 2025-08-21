# Junk Journal Poc

## Features

- **Token-based Access**: No authentication required - boards are accessed via secret tokens in URLs
- **Visual Canvas Editor**: Drag-and-drop interface for text, images, stickers, and shapes
- **Theme Customization**: Apply different visual themes to your boards
- **File Upload**: Support for image uploads with validation
- **Export Functionality**: Export pages as PNG images
- **Recap Views**: Track journaling activity with date filtering
- **Responsive Design**: Works on desktop and mobile devices

## Tech Stack

### Backend
- **Go 1.22+** with Fiber web framework
- **PostgreSQL 15+** with GORM ORM
- **Docker** for containerization
- **Zap** for structured logging

### Frontend
- **Vue 3** with Composition API and TypeScript
- **Vite** for fast development and builds
- **Tailwind CSS** for styling
- **Fabric.js** for canvas manipulation
- **Pinia** for state management

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Make (optional, for convenience commands)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd junk-journal-board
   ```

2. **Copy environment files**
   ```bash
   cp .env.example .env
   cp backend/.env.example backend/.env
   cp frontend/.env.example frontend/.env
   ```

3. **Start development environment**
   ```bash
   # Using Make (recommended)
   make dev
   
   # Or using Docker Compose directly
   docker compose up --build --watch
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database: localhost:5432

### Available Commands

```bash
# Development
make dev              # Start development environment
make stop             # Stop all services
make clean            # Stop services and remove volumes
make logs             # Show logs from all services

# Database
make db-reset         # Reset database
make db-shell         # Connect to database shell

# Backend
make backend-test     # Run backend tests
make backend-shell    # Open shell in backend container

# Frontend
make frontend-test    # Run frontend tests
make frontend-lint    # Run linter
make frontend-build   # Build for production

# Setup
make setup            # Initial setup and start services
```

## Project Structure

```
├── backend/                 # Go backend application
│   ├── internal/
│   │   ├── config/         # Database configuration
│   │   └── models/         # GORM models
│   ├── main.go             # Application entry point
│   ├── go.mod              # Go dependencies
│   └── Dockerfile          # Backend container
├── frontend/               # Vue 3 frontend application
│   ├── src/
│   │   ├── views/          # Vue components/pages
│   │   ├── router/         # Vue Router configuration
│   │   ├── main.ts         # Application entry point
│   │   └── style.css       # Global styles
│   ├── package.json        # Node dependencies
│   └── Dockerfile          # Frontend container
├── database/
│   └── init.sql            # Database initialization
├── docker-compose.yml      # Development services
├── Makefile               # Development commands
└── README.md              # This file
```

## Environment Variables

### Backend (.env)
- `DB_HOST`: Database host (default: postgres)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: postgres)
- `DB_PASSWORD`: Database password (default: postgres)
- `DB_NAME`: Database name (default: junk_journal)
- `PORT`: Server port (default: 8080)

### Frontend (.env)
- `VITE_API_BASE_URL`: Backend API URL (default: http://localhost:8080/api)

## Development Workflow

1. **Start the development environment**: `make dev`
2. **Make changes** to backend or frontend code
3. **Hot reload** will automatically restart services
4. **View logs** with `make logs` if needed
5. **Run tests** with `make backend-test` or `make frontend-test`

