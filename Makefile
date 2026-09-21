# Configuration variables
APP ?= o7k
TAG ?= 0.1.0
ITERATION ?= 1
TARGET_ENV ?= dev
VERSION ?= $(TAG)-$(TARGET_ENV).$(ITERATION)
#VERSION ?= $(TAG)

# Build metadata
BIN_VERSION ?= $(VERSION)
BIN_GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BIN_BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS = -s -w \
	-X github.com/akyriako/o7k/internal/version.Version=$(BIN_VERSION) \
	-X github.com/akyriako/o7k/internal/version.Commit=$(BIN_GIT_COMMIT) \
	-X github.com/akyriako/o7k/internal/version.BuildDate=$(BIN_BUILD_DATE)

build:
	@echo "Building $(APP) $(BIN_VERSION)..."
	@mkdir -p bin
	CGO_ENABLED=0 go build \
		-trimpath \
		-ldflags "$(LDFLAGS)" \
		-o bin/$(APP) \
		./cmd/main.go

run:
	go run ./cmd/main.go

test:
	go test ./...

version: build
	./bin/$(APP) --version

clean:
	@echo "Cleaning up..."
	rm -rf bin

.PHONY: build run test version clean