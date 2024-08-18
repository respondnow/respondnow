# Makefile for building respond server
# Reference Guide - https://www.gnu.org/software/make/manual/make.html

DOCKER_REGISTRY ?= docker.io
DOCKER_REPO ?= respondnow
DOCKER_IMAGE ?= respond
DOCKER_TAG ?= ci
IS_DOCKER_INSTALLED = $(shell which docker >> /dev/null 2>&1; echo $$?)
HOME = $(shell echo $$HOME)

.PHONY: help
help:
	@echo ""
	@echo "Usage:-"
	@echo "\tmake deps                        -- sets up dependencies for server image build"
	@echo "\tmake gofmt-check                 -- run golang formatting checks for server"
	@echo "\tmake build-server                -- builds server binary"
	@echo "\tmake run-dev-server              -- run server locally for dev purposes"
	@echo "\tmake gen-server-api-doc          -- generate server's swagger API doc"
	@echo "\tmake build-server-image          -- builds server's amd64 image"
	@echo "\tmake push-server-image           -- pushes the server's amd64 image"
	@echo ""

.PHONY: all
all: deps build-server test

.PHONY: deps
deps: _build_check_docker go-dep-update

_build_check_docker:
	@if [ $(IS_DOCKER_INSTALLED) -eq 1 ]; \
		then echo "" \
		&& echo "ERROR:\tdocker is not installed. Please install it before build." \
		&& echo "" \
		&& exit 1; \
		fi;

.PHONY: test
test:
	@echo "------------------"
	@echo "--> Run Go Test"
	@echo "------------------"
	@go test ./... -coverprofile=coverage.txt -covermode=atomic -v -count=1


.PHONY: go-dep-update
go-dep-update:
	go mod tidy
	go mod download


.PHONY: gen-server-api-doc
gen-server-api-doc:
	swag init -g cmd/main.go


.PHONY: build-server
build-server: go-dep-update gen-server-api-doc
	rm -f bin/server
	CGO_ENABLED=0 go build -o bin/server cmd/main.go

.PHONY: run-dev-server
run-dev-server: go-dep-update gen-server-api-doc
	MONGO_SERVER_URL=mongodb://127.0.0.1:27018,127.0.0.1:27019,127.0.0.1:27020 \
	MONGO_DB_NAME=respondnow \
	MONGO_DB_USERNAME='' \
	MONGO_DB_PASSWORD='' \
	ZOOM_LINK='http://my-example-zoom-link' \
	go run ./cmd/ \
		--debug=true \
		--config-file=./cmd/config.yaml

.PHONY: build-server-image
build-server-image:
	@echo "--------------------------------------"
	@echo "--> Build RespondNow server image for amd64" 
	@echo "--------------------------------------"
	@docker build -f package/Dockerfile  --no-cache -t $(DOCKER_REGISTRY)/$(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) .  --build-arg TARGETPLATFORM="linux/amd64"

.PHONY: push-server-image
push-server-image:
	@echo "--------------------------------------"
	@echo "--> Push RespondNow server image for amd64" 
	@echo "--------------------------------------"
	@docker push $(DOCKER_REGISTRY)/$(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG)


.PHONY: gofmt-check
gofmt-check:
	@echo "------------------"
	@echo "--> Run gofmt for RespondNow server"
	@echo "------------------"
	@gfmt=$$(gofmt -s -l . | wc -l); \
	if [ "$${gfmt}" -ne 0 ]; then \
		echo "The following files were found to be not go formatted:"; \
   		gofmt -s -l .; \
   		exit 1; \
  	fi
