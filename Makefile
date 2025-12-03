PROJECT_NAME := url-shortener

.PHONY: fmt tidy run build test all
run:
	@echo "Running $(PROJECT_NAME)..."
	@go run ./cmd/url-shortener --config "./local.yaml"

fmt:
	@go fmt ./...

tidy:
	@go mod tidy
