
build:
	@echo Linting go code...
	@golangci-lint run ./...
	
	@echo Running tests...
	@go test -coverpkg ./... -race ./... -coverprofile coverage.out
	
	@go run cmd/coverage/main.go run --standard 0
.PHONY: build
