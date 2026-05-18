BINARY_NAME := awsmgr
BIN_DIR := bin
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
IS_CLEAN := $(shell test -z "$$(git status --porcelain)" && echo clean)
LATEST_TAG := $(shell git tag --list 'v[0-9]*.[0-9]*.[0-9]*' --sort=-version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' | head -n 1)
CURRENT_TAG := $(shell git tag --points-at HEAD --list 'v[0-9]*.[0-9]*.[0-9]*' --sort=-version:refname | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' | head -n 1)
BASE_VERSION := $(patsubst v%,%,$(if $(LATEST_TAG),$(LATEST_TAG),v0.0.0))
NEXT_PATCH := $(shell echo "$(BASE_VERSION)" | awk -F. '{printf "%d.%d.%d", $$1, $$2, $$3 + 1}')
NEXT_MINOR := $(shell echo "$(BASE_VERSION)" | awk -F. '{printf "%d.%d.0", $$1, $$2 + 1}')
NEXT_MAJOR := $(shell echo "$(BASE_VERSION)" | awk -F. '{printf "%d.0.0", $$1 + 1}')
AUTO_VERSION := $(if $(and $(CURRENT_TAG),$(IS_CLEAN)),$(patsubst v%,%,$(CURRENT_TAG)),$(NEXT_PATCH)-dev.$(COMMIT))
VERSION ?= $(AUTO_VERSION)
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: all
all: fmt vet test build

.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	go build -trimpath -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME) .

.PHONY: version
version:
	@echo "$(VERSION)"

.PHONY: next-patch
next-patch:
	@echo "$(NEXT_PATCH)"

.PHONY: next-minor
next-minor:
	@echo "$(NEXT_MINOR)"

.PHONY: next-major
next-major:
	@echo "$(NEXT_MAJOR)"

.PHONY: release-snapshot
release-snapshot:
	goreleaser release --snapshot --clean --skip=publish

.PHONY: tag-patch
tag-patch: guard-clean
	git tag -a v$(NEXT_PATCH) -m "Release v$(NEXT_PATCH)"
	@echo "Created tag v$(NEXT_PATCH)"

.PHONY: tag-minor
tag-minor: guard-clean
	git tag -a v$(NEXT_MINOR) -m "Release v$(NEXT_MINOR)"
	@echo "Created tag v$(NEXT_MINOR)"

.PHONY: tag-major
tag-major: guard-clean
	git tag -a v$(NEXT_MAJOR) -m "Release v$(NEXT_MAJOR)"
	@echo "Created tag v$(NEXT_MAJOR)"

.PHONY: guard-clean
guard-clean:
	@test -z "$$(git status --porcelain)" || (echo "worktree has uncommitted changes; commit or stash before tagging" >&2; exit 1)

.PHONY: test
test:
	go test -race -v ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
