
build:
	@echo Installing build tooling...
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.0

	@echo Linting go code...
	@golangci-lint run ./...
	
	@echo Running tests...
	@go test -coverpkg ./... -race ./... -coverprofile coverage.out
	
	# this removes the cmd package from coverage analysis
	@cp coverage.out coverage.out.tmp && \
		grep -v 'cmd/go-coverage-analysis' coverage.out.tmp > coverage.out && \
		rm coverage.out.tmp
	
	@go run cmd/go-coverage-analysis/*.go run --standard 90
.PHONY: build
