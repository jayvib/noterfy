#!make

DOCKER_IMAGE_BUILD_FLAG="no-force"

# Include the envfile that contains all the metadata about the app
include build/noterfy/build.env
export $(shell sed 's/=.*//' build/noterfy/build.env)

APP_NAME:=noterfy
GIT_COMMIT:=$(shell git rev-parse --short HEAD)

define DOCKER_DEPLOY_DEV_HELP_INFO
# Use to deploy the development stage services to Docker Swarm
# mode.
#
# Make sure that the Docker in swarm mode already otherwise Docker
# will throw an error.
#
# Example:
#		make docker-deploy-dev
endef
.PHONY: docker-deploy-dev
docker-deploy-dev:
	docker stack deploy -c ./build/noterfy/docker/dev/docker-stack.yml $(APP_NAME)

define LOCAL_SERVER_HELP_INFO
# Use to run noterfy server in local machine.
#
# Example:
# 	make local-server
endef
local-server: build-noterfy
	./bin/${APP_NAME}.linux

define START_DEV_SERVICES_HELP_INFO
# Use to spin the development services
# including the noterfy dependencies.
# by default when the image is not yet exists
# this will build the Docker image of the engine.
# When you want to force build the image despite of
# the image is already existing use:
#      DOCKER_IMAGE_BUILD_FLAG=force
#
# Example:
# 	make start-dev-services
#   make start-dev-services DOCKER_IMAGE_BUILD_FLAG=force
endef
.PHONY: start-dev-services
start-dev-services:
ifeq ($(DOCKER_IMAGE_BUILD_FLAG), force)
	@echo "👉 Forcing the Docker to build the image"
	cd ./build/noterfy/docker/dev && docker-compose build
endif
	@echo "👉 Starting the services"
	@cd ./build/noterfy/docker/dev && \
	 	docker-compose up -d && \
	 	docker-compose ps

define STOP_DEV_SERVICES
# Use to stop the development services
#
# Example:
# 	make stop-dev-services
endef
.PHONY: stop-dev-services
stop-dev-services:
	@echo "👉 Stopping the services"
	@cd ./build/noterfy/docker/dev && \
		docker-compose stop && \
		docker-compose ps


define BUILD_HELP_INFO
# Build run the build step process of Noterfy Docker image.
#
# Example:
# 	make build
endef
build: clean docker-build-noterfy clean

define BUILD_NOTERFY_HELP_INFO
# Use to build the executable file of noterfy.
# The executable will store in bin/
#
# Example:
# 	make build-noterfy
endef
.PHONY: build-noterfy
build-noterfy: # Use to build the executable file of the noterfy. The executable will store in ./bin/ directory
ifdef NOTERFY_VERSION
	@echo "👉 Noterfy Version: ${NOTERFY_VERSION}"
endif
ifeq ($(wildcard ./bin/.*),)
	@echo "📂 Creating bin directory"
	@mkdir ./bin
endif
	@echo "🛠 Building Noterfy"
	@CGO_ENABLED=0 go build \
		-a \
		-tags netgo \
		-ldflags '-extldflags "-static" -w -s -X "main.Version=${NOTERFY_VERSION}" -X "main.BuildCommit=${GIT_COMMIT}"' \
		-o ./bin/noterfy.linux \
		./cmd/noterfy_server/main.go

define DOCKER_BUILD_NOTERFY_HELP_INFO
# Use to build the Docker image of noterfy.
# This will tag the image with latest an its version.
#
#	Example:
# 	make docker-build-noterfy
endef
.PHONY: docker-build-noterfy
docker-build-noterfy:
ifeq ($(DOCKER_BUILDKIT), 1)
	@echo "👉 Docker BuildKit enable"
endif
	@echo "🛠 Building Noterfy Docker Image"
	@docker build \
		-t ${APP_NAME} \
		--build-arg NOTERFY_BUILD_COMMIT=${GIT_COMMIT} \
		--build-arg NOTERFY_VERSION=${NOTERFY_VERSION} \
		-f ./build/noterfy/docker/Dockerfile \
		.
	@docker tag ${APP_NAME} jayvib/${APP_NAME}:latest
	@docker tag ${APP_NAME} jayvib/${APP_NAME}:${NOTERFY_VERSION}

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
# Use to run unit testing in noterfy source code.
#
# Example:
# 	make unit-test
endef
.PHONY: unit-test
unit-test: lint
	@echo "🏃 Running unit tests"
	@go test -short ./... -tags=unit-tests -race | grep -v '^?'

define LINT_HELP_INFO
# Use to lint the noterfy source code.
#
# Example:
# 	make lint
endef
.PHONY: lint
lint: lint-check-deps
	@echo "🔎🔎🔎 Linting sources"
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

define INSTALL_NOTERFY_CLI_HELP_INFO
# Use to install the noterfy cli tool.
#
# Example:
# 	make install-noterfy-cli
endef
install-noterfy-cli:
	@echo " 👉 Installing noterfy CLI ⚙"
	@go install ./cmd/noterfy_cli/

define CLEAN_HELP_INFO
# Use to clean the image layers after building the Docker image
#
# Example:
# 	make clean
endef
clean:
	@echo "🧹️ Cleaning up resources"
	@docker images -q --filter "dangling=true" | xargs docker rmi || true
