# =============================================================================
# Default target
# =============================================================================
all: build


# =============================================================================
# Build
# =============================================================================
.PHONY: build
build: format lint mod-update


# =============================================================================
# Format
# =============================================================================
.PHONY: format
format:
	@echo "Formatting ..."
	@gofmt -s -l -w .


# =============================================================================
# Lint
# =============================================================================
.PHONY: lint-install
lint-install:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -d -b $(GOPATH)/bin v1.18.0

.PHONY: lint
lint:
	@echo "Linting ..."
	@golangci-lint run


# =============================================================================
# Modules
# =============================================================================
.PHONY: mod-update
mod-update:
	@echo "Updating modules ..."
	@go get -v all
	@go mod tidy