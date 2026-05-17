BINARY_NAME := awsmgr
BIN_DIR := bin

.PHONY: all
all: fmt vet test build

.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	go build -trimpath -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME) .

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
