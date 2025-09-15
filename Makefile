.PHONY: dev build test clean watch

dev:
	@echo "🚀 Starting development server..."
	@templ generate
	@go run cmd/server/main.go

build:
	@echo "🔨 Building application..."
	@templ generate
	@go build -o bin/server ./cmd/server

test:
	@echo "🧪 Running tests..."
	@go test ./...

clean:
	@echo "🧹 Cleaning up..."
	@rm -rf bin/
	@find . -name "*_templ.go" -delete

fmt:
	@echo "💅 Formatting code..."
	@go fmt ./...
	@templ fmt .

# Simple watch mode (manual restart)
watch:
	@echo "👀 Watching files (Ctrl+C to restart)..."
	@while true; do \
		echo "⚡ Generating templates and starting server..."; \
		templ generate; \
		go run cmd/server/main.go; \
		echo "💥 Server stopped. Press Enter to restart or Ctrl+C to exit."; \
		read; \
	done