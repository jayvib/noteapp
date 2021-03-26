APP_NAME:=noteapp
VERSION:=0.0.0

define LOCAL_SERVER_HELP_INFO
# Use to run noteapp server in local machine.
#
# Example:
# 	make local-server
endef
local-server: build-noteapp
	./bin/${APP_NAME}.linux

define BUILD_NOTEAPP_HELP_INFO
# Use to build the executable file of noteapp.
# The executable will store in backend/bin/
#
# Example:
# 	make build-noteapp
endef
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

define DOCKER_BUILD_NOTEAPP_HELP_INFO
# Use to build the Docker image of noteapp.
# This will tag the image with latest an its version.
#
#	Example:
# 	make docker-build-noteapp
endef
.PHONY: docker-build-noteapp
docker-build-noteapp:
	docker build -t ${APP_NAME} ./build/noteapp
	docker tag ${APP_NAME} jayvib/${APP_NAME}:latest
	docker tag ${APP_NAME} jayvib/${APP_NAME}:${VERSION}

define FMT_HELP_INFO
# Use to format the Go source code.
#
# Example:
# 	make fmt
endef
.PHONY: fmt
fmt:
	go fmt ./...

define UNIT_TESTS_HELP_INFO
# Use to run unit testing in noteapp source code.
#
# Example:
# 	make unit-test
endef
.PHONY: unit-test
unit-test: lint
	@echo "Running unit tests"
	go test -short ./... -tags=unit-tests -race | grep -v '^?'

define LINT_HELP_INFO
# Use to lint the noteapp source code.
#
# Example:
# 	make lint
endef
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

define LINT_CHECK_DEPS_HELP_INFO
# Use to check the lint executable.
#
# Example:
# 	make lint-check-deps
endef
.PHONY: lint-check-deps
lint-check-deps:
	@if [ -z `which golangci-lint` ]; then \
		echo "[go get] installing golangci-lint"; \
  fi

define MOD_HELP_INFO
## Use to download the dependencies.
#
# Example:
#		make mod
endef
.PHONY: mod
mod:
	@go mod download

define GENERATE_HELP_INFO
# Use to run the go:generate in source code.
#
# Example:
# 	make generate
endef
.PHONY: generate
generate:
	go generate ./...