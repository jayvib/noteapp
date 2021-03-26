
APP_NAME:=noteapp

.PHONY: build-noteapp
build-noteapp: # Use to build the executable file of the noteapp. The executable will store in ./bin/ directory
ifeq ($(wildcard ./bin/.*),)
	@echo " ---> Creating bin directory"
	@mkdir ./bin
endif
	@echo " ---> Building Noteapp"
	@CGO_ENABLED=0 go build \
		-a \
		-tags netgo \
		-ldflags '-w -extldflags "-static"' \
		-o ./bin/noteapp.linux \
		./cmd/noteapp/main.go

.PHONY: docker-build-noteapp
docker-build-noteapp:
	docker build -t ${APP_NAME} ./build/noteapp

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: unit-tests
unit-tests:
	@echo "Running unit tests"
	go test -short ./... -tags=unit-tests -race | grep -v '^?'

.PHONY: lint
lint: lint-check-deps
	@echo "[golangci-lint] linting sources"
	@golangci-lint run \
		-E misspell \
		-E golint \
		-E gofmt \
		-E unconvert \
		--exclude-use-default=false \
		./...

.PHONY: lint-check-deps
lint-check-deps:
	@if [ -z `which golangci-lint` ]; then \
		echo "[go get] installing golangci-lint"; \
  fi