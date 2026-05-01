.PHONY: test lint build run-example

# Run unit tests across all sub-modules
test:
	@echo "==> Running standard tests..."
	@go test -v -cover ./...

# Vet and fmt
lint:
	@echo "==> Vetting code..."
	@go vet ./...
	@echo "==> Formatting code..."
	@go fmt ./...

# Run the example
run-example:
	@echo "==> Running example implementation..."
	@cd examples/gin && go run main.go

# Mod tidy
tidy:
	@echo "==> Tidying dependencies..."
	@go mod tidy
	@cd examples/gin && go mod tidy
