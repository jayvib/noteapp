
fmt:
	go fmt ./...

unit-tests: lint
	@echo "Running unit tests"
	go test -short ./... -tags=unit-tests -race | grep -v '^?'

lint: lint-check-deps
	@echo "[golangci-lint] linting sources"
	@golangci-lint run \
		-E misspell \
		-E golint \
		-E gofmt \
		-E unconvert \
		--exclude-use-default=false \
		./...

lint-check-deps:
	@if [ -z `which golangci-lint` ]; then \
		echo "[go get] installing golangci-lint"; \
    GO111MODULE=on go get -u github.com/golangci/golangci-lint/cmd/golangci-lint; \
  fi