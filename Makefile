.PHONY: tailwind
tailwind:
	npx @tailwindcss/cli -i internal/ui/components/css/input.css -o public/output.css
	npx prettier public/output.css -w

.PHONY: tailwind-watch
tailwind-watch:
	npx @tailwindcss/cli -i internal/ui/components/css/input.css -o public/output.css --watch

.PHONY: generate
generate: tailwind generate-templ

.PHONY: generate-templ
generate-templ:
	go tool templ generate

.PHONY: format
format: format-go format-templ

.PHONY: format-go
format-go:
	./scripts/format.sh

.PHONY: format-templ
format-templ:
	go tool templ fmt .
