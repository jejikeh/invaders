MAIN_PACKAGE_PATH := ./
BINARY_NAME := invaders

.PHONY: build
build:
	go build -o $=${MAIN_PACKAGE_PATH}bin/${BINARY_NAME} -ldflags "-X main.MajorVersion=0 -X main.MinorVersion=`date -u +%Y%m%d.%H%M%S`" ${MAIN_PACKAGE_PATH}

.PHONY: run
run:
	go run ${MAIN_PACKAGE_PATH}