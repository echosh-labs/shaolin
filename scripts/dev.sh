#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "🥋 Starting Shaolin Academy Full-Stack Development Environment..."

# Trap SIGINT and SIGTERM to kill child processes
cleanup() {
    echo -e "\n🛑 Shutting down development servers..."
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null || true
    exit 0
}
trap cleanup SIGINT SIGTERM EXIT

# 1. Start Go API Backend (Port 8081)
echo "🚀 Starting Go API Server on http://localhost:8081..."
cd "$DIR/backend"
/usr/local/go/bin/go run ./cmd/api/main.go &
BACKEND_PID=$!

# 2. Start Next.js Frontend (Port 3000)
echo "⚡ Starting Next.js Dev Server on http://localhost:3000..."
cd "$DIR/frontend"
/usr/bin/npm run dev &
FRONTEND_PID=$!

echo "🌟 System is online!"
echo "   - Frontend UI:  http://localhost:3000"
echo "   - Go API:       http://localhost:8081"
echo "   - Press Ctrl+C to terminate all services."

# Wait for background processes
wait $BACKEND_PID $FRONTEND_PID
