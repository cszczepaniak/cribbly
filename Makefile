.PHONY: tailwind
tailwind:
	npx @tailwindcss/cli -i internal/ui/components/css/input.css -o public/output.css
	@make format-prettier

.PHONY: tailwind-watch
tailwind-watch:
	npx @tailwindcss/cli -i internal/ui/components/css/input.css -o public/output.css --watch

.PHONY: generate-templ
generate-templ:
	go tool templ generate
	@make format-go

.PHONY: format
format: format-prettier format-go format-templ

.PHONY: format-prettier
format-prettier:
	npx prettier . --write --log-level=warn

.PHONY: format-go
format-go:
	go tool goimports -local github.com/cszczepaniak/cribbly -w .

.PHONY: format-templ
format-templ:
	go tool templ fmt .
