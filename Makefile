DOCKER_REGISTRY = index.docker.io
IMAGE_NAME = sesr
IMAGE_VERSION = latest
IMAGE_ORG = flaccid
IMAGE_TAG = $(DOCKER_REGISTRY)/$(IMAGE_ORG)/$(IMAGE_NAME):$(IMAGE_VERSION)

WORKING_DIR := $(shell pwd)

.DEFAULT_GOAL := build

.PHONY: build-static-linux

push:: ## Pushes the docker image to the registry
		@docker push $(IMAGE_TAG)

build:: ## Builds the docker image locally
		@echo building $(IMAGE_TAG)
		@docker build --pull \
		-t $(IMAGE_TAG) $(WORKING_DIR)

run:: ## Runs the docker image locally
		@docker run \
			-it \
			$(DOCKER_REGISTRY)/$(IMAGE_ORG)/$(IMAGE_NAME):$(IMAGE_VERSION)

run-daemon:: ## Runs the daemon locally with go directly
		@go run cmd/sesr/sesr.go --daemon

build-static-linux:: ## Builds a static linux binary
		@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
			go build \
			-o bin/sesr \
			-a -ldflags '-extldflags "-static"' \
				cmd/sesr/sesr.go

docker-build:: ## Builds the docker image locally
		@docker build --pull \
		-t $(IMAGE_TAG) $(WORKING_DIR)

docker-push:: ## Pushes the docker image to the registry
		@docker push $(IMAGE_TAG)

docker-release:: docker-build docker-push ## Builds and pushes the docker image to the registry

# A help target including self-documenting targets (see the awk statement)
define HELP_TEXT
Usage: make [TARGET]... [MAKEVAR1=SOMETHING]...

Available targets:
endef
export HELP_TEXT
help: ## This help target
	@cat .banner
	@echo
	@echo "$$HELP_TEXT"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / \
		{printf "\033[36m%-30s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)
