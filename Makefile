.PHONY: help setup build docker-build docker-run

REGISTRY := docker.io/pellepedro
BINARY   := halo
VERSION  := $(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")
BUILD    := $(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
LDFLAGS  := -ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"

CURRENT_DIR = $(shell pwd)

help:	## - Show help message
	@printf "\033[32m\xE2\x9c\x93 usage: make [target]\n\n\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


build:	## - Building Application Binary Locally
	@printf "\033[32m\xE2\x9c\x93 Building Application Binary Locally \033[0m"
	@go build -o $(BINARY) $(LDFLAGS) cmd/app/main.go

generate:	## - Generate Protocol Buffer Stubs
	@printf "\033[32m\xE2\x9c\x93 Generate Protocol Buffer Stubs \033[0m"
	@docker build --target=artifact --output type=local,dest=${CURRENT_DIR}/protos/ -f Dockerfile.protoc  .

docker-build:	## - Building Application Container
	@printf "\033[32m\xE2\x9c\x93 Building Application Container ${REGISTRY}/${BINARY} \033[0m"
	@DOCKER_BUILDKIT=1 docker build \
		-t ${REGISTRY}/${BINARY}:latest \
		-t ${REGISTRY}/${BINARY}:${VERSION} \
		--build-arg binary=$(BINARY) --build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile .

docker-run:	## - Running Application Container
	@printf "\033[32m\xE2\x9c\x93 Running Application Container ${REGISTRY}/${BINARY} \033[0m"
	@docker run ${REGISTRY}/${BINARY}

run-grpc-server:	## - Running GRPC Server
	@printf "\033[32m\xE2\x9c\x93 Running GRPC server \033[0m\n"
	@go run cmd/grpc/server/*.go

run-grpc-client:	## - Running GRPC Client
	@printf "\033[32m\xE2\x9c\x93 Running GRPC client \033[0m\n"
	@go run cmd/grpc/client/main.go