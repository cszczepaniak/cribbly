.PHONY: generate
generate: tailwind generate-templ

.PHONY: tailwind
tailwind:
	npx @tailwindcss/cli -i internal/ui/components/css/input.css -o public/output.css

.PHONY: generate-templ
generate-templ:
	go tool templ generate

.PHONY: format
format: format-go format-templ

.PHONY: format-go
format-go:
	./scripts/format.sh

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: format-templ
format-templ:
	go tool templ fmt .

.PHONY: dev
dev:
	./scripts/run_local.sh
