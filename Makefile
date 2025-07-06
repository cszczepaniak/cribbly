.PHONY: tailwind
tailwind:
	npx @tailwindcss/cli -i internal/ui/components/css/input.css -o public/output.css
	npx prettier public/output.css -w

.PHONY: tailwind-watch
tailwind-watch:
	npx @tailwindcss/cli -i internal/ui/components/css/input.css -o public/output.css --watch

.PHONY: generate-templ
generate-templ:
	go tool templ generate
	@make format-go

.PHONY: format
format: format-go format-templ

.PHONY: format-go
format-go:
	go tool goimports -local github.com/cszczepaniak/cribbly -w .

.PHONY: format-templ
format-templ:
	go tool templ fmt .
