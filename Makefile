

# Define variables
BIN_PATH := build/savebot
GO_BUILD_FLAGS := -v
CGO_ENABLED=0

# Version
VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1` 2>/dev/null)
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date +%Y-%m-%d)

# Default target
.PHONY: all
all: build

# Build target
.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o $(BIN_PATH) -ldflags "-X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'" ./cmd/savebot


.PHONY: clean
clean:
	rm -rf $(BIN_PATH)
	rm -rf debug/Telegram

.PHONY: image
image:
	docker build -t ghcr.io/asd2003ru/savebot:latest .

.PHONY: push
push:
	docker push ghcr.io/asd2003ru/savebot:latest