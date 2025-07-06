.PHONY: tailwind
tailwind:
	tailwindcss -i internal/ui/components/css/input.css -o public/output.css

.PHONY: tailwind-watch
tailwind-watch:
	tailwindcss -i internal/ui/components/css/input.css -o public/output.css --watch
