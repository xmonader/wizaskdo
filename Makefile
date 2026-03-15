.PHONY: build test clean run install

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +%Y-%m-%d)

LDFLAGS = -s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.date=$(DATE)

build:
	go build -ldflags "$(LDFLAGS)" -o wizask .
	go build -ldflags "$(LDFLAGS)" -o wizdo ./cmd/wizdo

test:
	go test -v -race ./pkg/...

clean:
	rm -f wizask wizdo

run:
	go run -ldflags "$(LDFLAGS)" . $(ARGS)

run-do:
	go run -ldflags "$(LDFLAGS)" ./cmd/wizdo $(ARGS)

install:
	go install -ldflags "$(LDFLAGS)" .
	go install -ldflags "$(LDFLAGS)" ./cmd/wizdo

# Development builds (without ldflags)
dev-build:
	go build -o wizask .
	go build -o wizdo ./cmd/wizdo

# Release using goreleaser (requires goreleaser installed)
release:
	goreleaser release --clean

# Snapshot release (for testing)
release-snapshot:
	goreleaser release --snapshot --clean

# Check goreleaser config
release-check:
	goreleaser check
