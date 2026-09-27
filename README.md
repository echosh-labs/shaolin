# Shaolin Platform — Martial Arts Dojo Management & Portal System

Full-stack martial arts school management system, student grading portal, curriculum scheduler, and community hub for Shaolin Academy.

---

## Architecture Overview

- **Backend (`/backend`)**: High-performance Go (1.22+) backend service using Go standard library multiplexer with modern wildcard path parameters, SQLite relational persistence, JWT authentication with bcrypt hashing, and in-memory telemetry analytics. Runs on **Port 8081**.
- **Frontend (`/frontend`)**: Modern Next.js (App Router), TypeScript, and TailwindCSS responsive user portal for student bookings, training curriculum, disciple certification, and administrative controls. Runs on **Port 3000**.
- **Scraper / Ingestion (`/scraper`)**: Automated scripts and data extraction utilities for curriculum imports and roster syncing.
- **Automation (`/scripts`, `Makefile`)**: Unified lifecycle management for development, builds, database migrations, and testing.

---

## Quickstart

### Prerequisites
- Go 1.22+
- Node.js 20+ (npm / pnpm)

### Development
Start both Go backend (Port 8081) and Next.js frontend (Port 3000):
```bash
make dev
```

### Build Production Binaries
Compile Go binary and Next.js bundle:
```bash
make build
```

### Run Test Suite
Run backend unit tests, frontend typechecking, and integration tests:
```bash
make test
```

### Seed Database
Migrate and re-seed clean SQLite database (`shaolin.db`):
```bash
make seed
```
