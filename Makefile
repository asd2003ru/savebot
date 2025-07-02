

# Define variables
BIN_PATH := build/savebot

GO_BUILD_FLAGS := -v
CGO_ENABLED=0

# Default target
.PHONY: all
all: build

# Build target
.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o $(BIN_PATH) ./cmd/savebot


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