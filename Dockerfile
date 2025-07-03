# Stage 1: Build stage
FROM golang:1.24.4-alpine AS builder


# Set the working directory
WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Optional: Set build args
ARG VERSION
ARG COMMIT
ARG DATE


# Build the Go application
# Get commit hash if VERSION is not passed
#RUN apk update && apk upgrade && apk add bash && apk add git && chmod +x /app/build.sh && /app/build.sh
RUN apk update && apk upgrade && apk add git && \
    VERSION=$(git describe --tags $(git rev-list --tags --max-count=1 2>/dev/null) 2>/dev/null || git rev-parse --short HEAD) && \
    COMMIT=$(git rev-parse --short HEAD) && \
    DATE=$(date +%Y-%m-%d) && \
    CGO_ENABLED=0 GOOS=linux go build -o ./savebot -ldflags "-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE" cmd/savebot/main.go

# Stage 2: Final stage
FROM alpine:edge

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/savebot .

# Set the timezone and install CA certificates
RUN apk --no-cache add ca-certificates tzdata

# Set the entrypoint command
ENTRYPOINT ["/app/savebot","-c","/app/config.yaml"]