#!make

# Include the envfile that contains all the metadata about the app
include build/noteapp/envfile
export $(shell sed 's/=.*//' build/noteapp/envfile)

APP_NAME:=noteapp

define LOCAL_SERVER_HELP_INFO
# Use to run noteapp server in local machine.
#
# Example:
# 	make local-server
endef
local-server: build-noteapp
	./bin/${APP_NAME}.linux

define BUILD_HELP_INFO
# Build run the build step process of Noteapp Docker image.
#
# Example:
# 	make build
endef
build: clean docker-build-noteapp clean

define BUILD_NOTEAPP_HELP_INFO
# Use to build the executable file of noteapp.
# The executable will store in backend/bin/
#
# Example:
# 	make build-noteapp
endef
.PHONY: build-noteapp
build-noteapp: # Use to build the executable file of the noteapp. The executable will store in ./bin/ directory
ifdef NOTEAPP_VERSION
	@echo "🛠 Building Noteapp version: ${NOTEAPP_VERSION}"
endif
ifeq ($(wildcard ./bin/.*),)
	@echo " ---> Creating bin directory"
	@mkdir ./bin
endif
	@echo " ---> Building Noteapp"
	@CGO_ENABLED=0 go build \
		-a \
		-tags netgo \
		-ldflags '-w -extldflags "-static"' \
		-ldflags '-X "main.Version=${NOTEAPP_VERSION}"' \
		-o ./bin/noteapp.linux \
		./cmd/noteapp_server/main.go

define DOCKER_BUILD_NOTEAPP_HELP_INFO
# Use to build the Docker image of noteapp.
# This will tag the image with latest an its version.
#
#	Example:
# 	make docker-build-noteapp
endef
.PHONY: docker-build-noteapp
docker-build-noteapp:
	@echo "🛠 Building Noteapp Docker Image"
	docker build -t ${APP_NAME} -f ./build/noteapp/docker/Dockerfile .
	docker tag ${APP_NAME} jayvib/${APP_NAME}:latest
	docker tag ${APP_NAME} jayvib/${APP_NAME}:${NOTEAPP_VERSION}

define FMT_HELP_INFO
# Use to format the Go source code.
#
# Example:
# 	make fmt
endef
.PHONY: fmt
fmt:
	@echo "🧹 Formatting source code"
	@go fmt ./...

define UNIT_TESTS_HELP_INFO
# Use to run unit testing in noteapp source code.
#
# Example:
# 	make unit-test
endef
.PHONY: unit-test
unit-test: lint
	@echo "🏃 Running unit tests"
	@go test -short ./... -tags=unit-tests -race | grep -v '^?'

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
ifdef GOPROXY
	@echo "👉 Go proxy setting found: GOPROXY=${GOPROXY}"
endif
	@echo "📥 Downloading Dependencies"
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

define INSTALL_NOTEAPP_CLI_HELP_INFO
# Use to install the noteapp cli tool.
#
# Example:
# 	make install-noteapp-cli
endef
install-noteapp-cli:
	@echo " 👉 Installing noteapp CLI ⚙"
	@go install ./cmd/noteapp_cli/

define CLEAN_HELP_INFO
# Use to clean the image layers after building the Docker image
#
# Example:
# 	make clean
endef
clean:
	@echo "🧹️ Cleaning up resources"
	@docker images -q --filter "dangling=true" | xargs docker rmi || true