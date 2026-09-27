#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "🔨 [1/2] Building Go API Binary..."
cd "$DIR/backend"
/usr/local/go/bin/go build -o apiServer ./cmd/api/main.go
echo "✅ Go API binary compiled: backend/apiServer"

echo "⚡ [2/2] Building Next.js Production Bundle..."
cd "$DIR/frontend"
/usr/bin/npm run build
echo "✅ Next.js production bundle compiled successfully."

echo "🎉 Build finished successfully!"
