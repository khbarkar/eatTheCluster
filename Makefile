BINARY_NAME := eatthecluster
BUILD_DIR := bin
GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*')
LDFLAGS := -s -w

.PHONY: all build release run clean

all: build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/eatthecluster

release:
	go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/eatthecluster

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)
