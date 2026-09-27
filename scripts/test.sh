#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "🧪 [1/3] Running Go Unit & Repository Test Suites..."
cd "$DIR/backend"
/usr/local/go/bin/go test ./...

echo "🔍 [2/3] Checking Frontend TypeScript & Production Build..."
cd "$DIR/frontend"
/usr/bin/npm run build

echo "⚡ [3/3] Running Full-Stack E2E API Integration Suite..."
# Check if API server is already listening on 8081
if curl -s http://localhost:8081/api/event-types > /dev/null 2>&1; then
    echo "API server already active on http://localhost:8081. Running verification suite..."
    cd "$DIR"
    python3 scripts/verify_api.py
else
    echo "Spawning test API server on http://localhost:8081..."
    cd "$DIR/backend"
    /usr/local/go/bin/go run ./cmd/api/main.go > /tmp/shaolin_test_api.log 2>&1 &
    SERVER_PID=$!
    trap "kill $SERVER_PID 2>/dev/null || true" EXIT
    sleep 2
    cd "$DIR"
    python3 scripts/verify_api.py
fi

echo "🎉 All test suites passed with 100% success!"
