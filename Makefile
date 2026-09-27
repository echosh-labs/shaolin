# Shaolin Academy - Build & Development Automation

.PHONY: all dev test build seed clean help

help:
	@echo "Shaolin Academy Platform - Available Commands:"
	@echo "  make dev          Start both Go API server (8081) and Next.js frontend (3000)"
	@echo "  make test         Execute full backend tests, frontend typecheck, and e2e API suite"
	@echo "  make build        Compile Go binary and Next.js production bundle"
	@echo "  make seed         Re-migrate and re-seed clean SQLite database (shaolin.db)"
	@echo "  make clean        Remove build artifacts and temporary databases"

dev:
	@bash ./scripts/dev.sh

test:
	@bash ./scripts/test.sh

build:
	@bash ./scripts/build.sh

seed:
	@bash ./scripts/seed.sh

clean:
	@rm -rf frontend/.next frontend/out backend/apiServer backend/shaolin.db
	@echo "✅ Cleaned build artifacts and temporary database."
