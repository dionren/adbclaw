BINARY_NAME := adbclaw
BUILD_DIR := bin
MODULE := github.com/llm-net/adbclaw
GO := go

.PHONY: build test lint clean tidy

build: tidy
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) .

tidy:
	$(GO) mod tidy

test:
	$(GO) test ./...

lint:
	$(GO) vet ./...

clean:
	rm -rf $(BUILD_DIR)
