#!/usr/bin/env bash

set -euo pipefail

cleanup() {
	if [[ -n "${TAILWIND_PID:-}" ]] && kill -0 "$TAILWIND_PID" 2>/dev/null; then
		kill "$TAILWIND_PID" || true
	fi
}

trap cleanup EXIT INT TERM

echo "Starting Tailwind CSS watcher..."
npx @tailwindcss/cli -i internal/ui/components/css/input.css -o public/output.css --watch &
TAILWIND_PID=$!

echo "Starting Go app (with templ hot reload)..."
# Generate templ files, and on change, recompile + run the Go binary. Using --proxy sets up a
# websocket to notify the browser to reload without needing a manual refresh.
go tool templ generate \
	--watch \
	--cmd "go run main.go" \
	--proxy "http://localhost:8080"
