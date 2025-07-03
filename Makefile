VERSION=$(shell git describe --tags 2>/dev/null || git rev-parse HEAD)
COMMIT=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date -u +"%Y-%m-%d %H:%M:%S")

.PHONY: all
all: build

# Build target

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -o ./build/savebot  -ldflags "-X 'main.version=${VERSION}' -X 'main.commit=${COMMIT}' -X 'main.date=${BUILD_DATE}'" cmd/savebot/main.go


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