MAIN_PACKAGE_PATH := ./
BINARY_NAME := invaders
MAJOR_VERSION := 0
MINOR_VERSION := $(shell date -u +%Y%m%d.%H%M%S)

.PHONY: build
build:
	go build -o $=${MAIN_PACKAGE_PATH}bin/${BINARY_NAME}.${MAJOR_VERSION}.${MINOR_VERSION} -ldflags "-X main.MajorVersion=${MAJOR_VERSION} -X main.MinorVersion=${MINOR_VERSION}" ${MAIN_PACKAGE_PATH}

.PHONY: run
run:
	go run ${MAIN_PACKAGE_PATH}

.DEFAULT_GOAL := all
BUILD_DIR := bin

TITLE := wooff

fmt:
	@gofmt -w .

build: fmt
	@go build -o $(BUILD_DIR)/$(TITLE) .

test: fmt
	@go test -v ./...

all: test build

clean:
	rm -rf $(BUILD_DIR)

%/.:
	mkdir -p $@
