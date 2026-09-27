#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$DIR/backend"

echo "🥋 Initializing & Seeding Shaolin Academy SQLite Database..."
rm -f shaolin.db
/usr/local/go/bin/go run ./cmd/seed/main.go

echo "✅ Database initialized with pristine seed records."
