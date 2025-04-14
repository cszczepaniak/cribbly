#!/usr/bin/env bash

set -e

# Generate templ files, and on change, recompile + run the Go binary. Using --proxy sets up a
# websocket to notify the browser to reload without needing a manual refresh. Neat!
go tool templ generate \
	--watch \
	--cmd "go run main.go" \
	--proxy "http://localhost:8080"
