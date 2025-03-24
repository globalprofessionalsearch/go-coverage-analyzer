
build:
	@echo Installing build tooling...
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.0
	@go install github.com/globalprofessionalsearch/go-coverage-analyzer/cmd/go-coverage-analysis@latest

	@echo Linting go code...
	@golangci-lint run ./...
	
	@echo Running tests...
	@go test -coverpkg ./... -race ./... -coverprofile coverage.out
	
	@go run cmd/coverage/main.go run --standard 0
.PHONY: build
