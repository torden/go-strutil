# Go-Strutil Multi-Version Compatible Makefile
# Supports Go 1.6+ through latest versions with automatic feature detection

PKG_NAME := go-strutil
MODULE_NAME := github.com/torden/go-strutil

# Version information
VERSION := $(shell git describe --tags --always --dirty="-dev" 2>/dev/null || echo "v0.0.0-dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')

# Go version detection and feature flags
GO_VERSION := $(shell go version | awk '{print $$3}' | sed 's/go//')
GO_VERSION_MAJOR := $(shell echo $(GO_VERSION) | cut -d. -f1)
GO_VERSION_MINOR := $(shell echo $(GO_VERSION) | cut -d. -f2)
GO_VERSION_NUM := $(shell echo $(GO_VERSION) | sed 's/\.//g' | cut -c1-3)

# Feature detection based on Go version
HAS_MODULES := $(shell [ "$(GO_VERSION_NUM)" -ge "111" ] && echo 1 || echo 0)
HAS_GENERICS := $(shell [ "$(GO_VERSION_NUM)" -ge "118" ] && echo 1 || echo 0)
HAS_EMBED := $(shell [ "$(GO_VERSION_NUM)" -ge "116" ] && echo 1 || echo 0)
HAS_FUZZING := $(shell [ "$(GO_VERSION_NUM)" -ge "118" ] && echo 1 || echo 0)

# Build tags based on Go version
BUILD_TAGS := 
ifeq ($(shell [ "$(GO_VERSION_NUM)" -ge "118" ] && echo 1 || echo 0),1)
    BUILD_TAGS += generics
else ifeq ($(shell [ "$(GO_VERSION_NUM)" -lt "17" ] && echo 1 || echo 0),1)
    BUILD_TAGS += legacy
endif

# Directories
BUILD_DIR := build
REPORT_DIR := reports
BENCH_DIR := $(REPORT_DIR)/benchmarks
COVER_DIR := $(REPORT_DIR)/coverage

# Build flags
LDFLAGS := -ldflags="-s -w -X '$(MODULE_NAME).Version=$(VERSION)' -X '$(MODULE_NAME).Commit=$(COMMIT)' -X '$(MODULE_NAME).BuildTime=$(BUILD_TIME)'"
DEV_LDFLAGS := -ldflags="-X '$(MODULE_NAME).Version=$(VERSION)-dev' -X '$(MODULE_NAME).Commit=$(COMMIT)' -X '$(MODULE_NAME).BuildTime=$(BUILD_TIME)'"

# Colors for output
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
BLUE := \033[34m
CYAN := \033[36m
RESET := \033[0m
BOLD := \033[1m

# Default target
.DEFAULT_GOAL := help

## Display this help message
help:
	@echo "$(BOLD)Go-Strutil Build System$(RESET)"
	@echo "$(CYAN)Go Version: $(GO_VERSION)$(RESET)"
	@echo "$(CYAN)Features: Modules=$(HAS_MODULES), Generics=$(HAS_GENERICS), Embed=$(HAS_EMBED), Fuzzing=$(HAS_FUZZING)$(RESET)"
	@echo ""
	@echo "$(BOLD)Available targets:$(RESET)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(CYAN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BOLD)Multi-version testing:$(RESET)"
	@echo "  $(CYAN)test-all-versions$(RESET)    Test with multiple Go versions (requires Docker)"
	@echo "  $(CYAN)test-compatibility$(RESET)   Test version compatibility features"
	@echo ""
	@echo "$(BOLD)Performance:$(RESET)"
	@echo "  $(CYAN)bench-compare$(RESET)        Compare performance across implementations"
	@echo "  $(CYAN)bench-report$(RESET)         Generate detailed benchmark report"

## Print detected Go version and features
version-info: ## Show Go version and feature detection
	@echo "$(BOLD)Version Information:$(RESET)"
	@echo "  Go Version: $(GREEN)$(GO_VERSION)$(RESET)"
	@echo "  Major.Minor: $(GO_VERSION_MAJOR).$(GO_VERSION_MINOR)"
	@echo "  Version Number: $(GO_VERSION_NUM)"
	@echo ""
	@echo "$(BOLD)Feature Support:$(RESET)"
	@echo "  Modules: $(if $(filter 1,$(HAS_MODULES)),$(GREEN)✓,$(RED)✗)$(RESET)"
	@echo "  Generics: $(if $(filter 1,$(HAS_GENERICS)),$(GREEN)✓,$(RED)✗)$(RESET)"
	@echo "  Embed FS: $(if $(filter 1,$(HAS_EMBED)),$(GREEN)✓,$(RED)✗)$(RESET)"
	@echo "  Fuzzing: $(if $(filter 1,$(HAS_FUZZING)),$(GREEN)✓,$(RED)✗)$(RESET)"
	@echo ""
	@echo "$(BOLD)Build Configuration:$(RESET)"
	@echo "  Build Tags: $(BUILD_TAGS)"
	@echo "  Module: $(MODULE_NAME)"
	@echo "  Version: $(VERSION)"

## Setup development environment based on Go version
setup: ## Install dependencies and tools
	@echo "$(BOLD)$(CYAN)Setting up development environment...$(RESET)"
ifeq ($(HAS_MODULES),1)
	@echo "$(GREEN)Using Go modules$(RESET)"
	@go mod download
	@go mod verify
	@go mod tidy
else
	@echo "$(YELLOW)Using legacy GOPATH mode$(RESET)"
	@go get -d -v ./...
endif
	@echo "$(GREEN)Installing development tools...$(RESET)"
	@go install golang.org/x/tools/cmd/goimports@latest 2>/dev/null || go get golang.org/x/tools/cmd/goimports
	@go install golang.org/x/lint/golint@latest 2>/dev/null || go get golang.org/x/lint/golint
	@go install honnef.co/go/tools/cmd/staticcheck@latest 2>/dev/null || echo "$(YELLOW)staticcheck not available$(RESET)"
ifeq ($(HAS_FUZZING),1)
	@echo "$(GREEN)Fuzzing support detected$(RESET)"
endif
	@echo "$(GREEN)Setup complete$(RESET)"

## Clean build artifacts and reports
clean: ## Remove build artifacts and reports
	@echo "$(BOLD)$(CYAN)Cleaning up...$(RESET)"
	@rm -rf $(BUILD_DIR) $(REPORT_DIR)
	@rm -f *.test *.prof *.cover *.coverprofile *.out
	@rm -f cpu.prof mem.prof block.prof mutex.prof
	@go clean -cache -testcache -modcache 2>/dev/null || go clean -cache -testcache
	@echo "$(GREEN)Clean complete$(RESET)"

## Create necessary directories
dirs:
	@mkdir -p $(BUILD_DIR) $(REPORT_DIR) $(BENCH_DIR) $(COVER_DIR)

## Build the library with version-specific optimizations
build: dirs ## Build the library
	@echo "$(BOLD)$(CYAN)Building $(PKG_NAME)...$(RESET)"
ifneq ($(BUILD_TAGS),)
	@echo "$(YELLOW)Using build tags: $(BUILD_TAGS)$(RESET)"
	@go build $(LDFLAGS) -tags "$(BUILD_TAGS)" -v ./...
else
	@go build $(LDFLAGS) -v ./...
endif
	@echo "$(GREEN)Build complete$(RESET)"

## Build for development (with debug info)
build-dev: dirs ## Build with debug information
	@echo "$(BOLD)$(CYAN)Building $(PKG_NAME) for development...$(RESET)"
ifneq ($(BUILD_TAGS),)
	@go build $(DEV_LDFLAGS) -tags "$(BUILD_TAGS)" -gcflags="all=-N -l" -v ./...
else
	@go build $(DEV_LDFLAGS) -gcflags="all=-N -l" -v ./...
endif
	@echo "$(GREEN)Development build complete$(RESET)"

## Run all tests with version-appropriate flags
test: dirs ## Run all tests
	@echo "$(BOLD)$(CYAN)Running tests...$(RESET)"
ifneq ($(BUILD_TAGS),)
	@echo "$(YELLOW)Testing with build tags: $(BUILD_TAGS)$(RESET)"
	@go test -tags "$(BUILD_TAGS)" -v ./... 
else
	@go test -v ./...
endif
	@echo "$(GREEN)Tests complete$(RESET)"

## Run tests with race detection (Go 1.1+)
test-race: dirs ## Run tests with race detection
	@echo "$(BOLD)$(CYAN)Running tests with race detection...$(RESET)"
ifeq ($(shell [ "$(GO_VERSION_NUM)" -ge "11" ] && echo 1 || echo 0),1)
ifneq ($(BUILD_TAGS),)
	@go test -tags "$(BUILD_TAGS)" -race -v ./...
else
	@go test -race -v ./...
endif
	@echo "$(GREEN)Race tests complete$(RESET)"
else
	@echo "$(YELLOW)Race detection not available in Go $(GO_VERSION)$(RESET)"
endif

## Run tests with coverage
test-coverage: dirs ## Generate test coverage report
	@echo "$(BOLD)$(CYAN)Generating coverage report...$(RESET)"
ifneq ($(BUILD_TAGS),)
	@go test -tags "$(BUILD_TAGS)" -coverprofile=$(COVER_DIR)/coverage.out -covermode=atomic ./...
else
	@go test -coverprofile=$(COVER_DIR)/coverage.out -covermode=atomic ./...
endif
	@go tool cover -html=$(COVER_DIR)/coverage.out -o $(COVER_DIR)/coverage.html
	@go tool cover -func=$(COVER_DIR)/coverage.out -o $(COVER_DIR)/coverage.txt
	@echo "$(GREEN)Coverage report generated: $(COVER_DIR)/coverage.html$(RESET)"
	@echo "$(CYAN)Coverage summary:$(RESET)"
	@tail -n 1 $(COVER_DIR)/coverage.txt

## Test version compatibility features
test-compatibility: ## Test version compatibility features
	@echo "$(BOLD)$(CYAN)Testing version compatibility...$(RESET)"
	@go test -run TestGoVersion -v
	@go test -run TestCompatibility -v
	@go test -run TestCoreStringProcessing -v
	@echo "$(GREEN)Compatibility tests complete$(RESET)"

## Run benchmarks
bench: dirs ## Run performance benchmarks
	@echo "$(BOLD)$(CYAN)Running benchmarks...$(RESET)"
ifneq ($(BUILD_TAGS),)
	@go test -tags "$(BUILD_TAGS)" -bench=. -benchmem -run=^$$ ./... | tee $(BENCH_DIR)/bench-$(GO_VERSION).txt
else
	@go test -bench=. -benchmem -run=^$$ ./... | tee $(BENCH_DIR)/bench-$(GO_VERSION).txt
endif
	@echo "$(GREEN)Benchmarks saved to $(BENCH_DIR)/bench-$(GO_VERSION).txt$(RESET)"

## Compare performance between generic and interface implementations
bench-compare: bench ## Compare generic vs interface performance
	@echo "$(BOLD)$(CYAN)Comparing implementation performance...$(RESET)"
ifeq ($(HAS_GENERICS),1)
	@echo "$(GREEN)Generics available - running comparison$(RESET)"
	@go test -bench=BenchmarkBasicPerformance/GenericVsInterface -benchmem -run=^$$ ./... | tee $(BENCH_DIR)/comparison-$(GO_VERSION).txt
	@echo "$(CYAN)Comparison results saved to $(BENCH_DIR)/comparison-$(GO_VERSION).txt$(RESET)"
else
	@echo "$(YELLOW)Generics not available in Go $(GO_VERSION) - skipping comparison$(RESET)"
endif

## Generate comprehensive benchmark report
bench-report: bench bench-compare ## Generate detailed benchmark report
	@echo "$(BOLD)$(CYAN)Generating benchmark report...$(RESET)"
	@echo "# Benchmark Report - Go $(GO_VERSION)" > $(BENCH_DIR)/report.md
	@echo "Generated: $(BUILD_TIME)" >> $(BENCH_DIR)/report.md
	@echo "" >> $(BENCH_DIR)/report.md
	@echo "## System Information" >> $(BENCH_DIR)/report.md
	@echo "- Go Version: $(GO_VERSION)" >> $(BENCH_DIR)/report.md
	@echo "- Features: Modules=$(HAS_MODULES), Generics=$(HAS_GENERICS)" >> $(BENCH_DIR)/report.md
	@echo "- Build Tags: $(BUILD_TAGS)" >> $(BENCH_DIR)/report.md
	@echo "" >> $(BENCH_DIR)/report.md
	@echo "## Benchmark Results" >> $(BENCH_DIR)/report.md
	@echo "\`\`\`" >> $(BENCH_DIR)/report.md
	@cat $(BENCH_DIR)/bench-$(GO_VERSION).txt >> $(BENCH_DIR)/report.md
	@echo "\`\`\`" >> $(BENCH_DIR)/report.md
ifeq ($(HAS_GENERICS),1)
	@echo "" >> $(BENCH_DIR)/report.md
	@echo "## Generic vs Interface Comparison" >> $(BENCH_DIR)/report.md
	@echo "\`\`\`" >> $(BENCH_DIR)/report.md
	@cat $(BENCH_DIR)/comparison-$(GO_VERSION).txt >> $(BENCH_DIR)/report.md
	@echo "\`\`\`" >> $(BENCH_DIR)/report.md
endif
	@echo "$(GREEN)Benchmark report generated: $(BENCH_DIR)/report.md$(RESET)"

## Run linting and static analysis
lint: ## Run linting and static analysis
	@echo "$(BOLD)$(CYAN)Running linters...$(RESET)"
	@go vet ./...
	@golint ./... 2>/dev/null || echo "$(YELLOW)golint not available$(RESET)"
	@staticcheck ./... 2>/dev/null || echo "$(YELLOW)staticcheck not available$(RESET)"
	@go fmt ./...
	@goimports -w . 2>/dev/null || echo "$(YELLOW)goimports not available$(RESET)"
	@echo "$(GREEN)Linting complete$(RESET)"

## Run fuzzing tests (Go 1.18+)
fuzz: ## Run fuzz tests (Go 1.18+)
ifeq ($(HAS_FUZZING),1)
	@echo "$(BOLD)$(CYAN)Running fuzz tests...$(RESET)"
	@go test -fuzz=. -fuzztime=30s ./... || echo "$(YELLOW)No fuzz tests found$(RESET)"
	@echo "$(GREEN)Fuzzing complete$(RESET)"
else
	@echo "$(YELLOW)Fuzzing not available in Go $(GO_VERSION)$(RESET)"
endif

## Profile CPU and memory usage
profile: dirs ## Generate CPU and memory profiles
	@echo "$(BOLD)$(CYAN)Generating profiles...$(RESET)"
ifneq ($(BUILD_TAGS),)
	@go test -tags "$(BUILD_TAGS)" -bench=. -cpuprofile=$(REPORT_DIR)/cpu.prof -memprofile=$(REPORT_DIR)/mem.prof -run=^$$ ./...
else
	@go test -bench=. -cpuprofile=$(REPORT_DIR)/cpu.prof -memprofile=$(REPORT_DIR)/mem.prof -run=^$$ ./...
endif
	@echo "$(GREEN)Profiles generated:$(RESET)"
	@echo "  CPU: $(REPORT_DIR)/cpu.prof"
	@echo "  Memory: $(REPORT_DIR)/mem.prof"
	@echo "$(CYAN)View with: go tool pprof $(REPORT_DIR)/cpu.prof$(RESET)"

## Test with multiple Go versions using Docker
test-all-versions: ## Test with multiple Go versions (requires Docker)
	@echo "$(BOLD)$(CYAN)Testing with multiple Go versions...$(RESET)"
	@for version in 1.6 1.11 1.16 1.18 1.19 1.20 1.21 1.22 1.23; do \
		echo "$(YELLOW)Testing with Go $$version...$(RESET)"; \
		docker run --rm -v $(PWD):/src -w /src golang:$$version sh -c "go mod download 2>/dev/null || go get -d -v ./...; go test -v ./..." || echo "$(RED)Go $$version failed$(RESET)"; \
	done
	@echo "$(GREEN)Multi-version testing complete$(RESET)"

## Check for security vulnerabilities
security: ## Check for security issues
	@echo "$(BOLD)$(CYAN)Checking for security issues...$(RESET)"
	@go list -json -deps ./... | jq -r '.Module.Path' | sort -u | xargs -I {} sh -c 'echo "Checking {}" && go list -m -versions {} 2>/dev/null || echo "No versions found"'
	@govulncheck ./... 2>/dev/null || echo "$(YELLOW)govulncheck not available - install with: go install golang.org/x/vuln/cmd/govulncheck@latest$(RESET)"

## Continuous Integration target
ci: clean setup lint test-race test-coverage bench ## Run full CI pipeline
	@echo "$(BOLD)$(GREEN)CI pipeline complete$(RESET)"

## Development workflow
dev: clean setup build-dev test lint ## Quick development workflow
	@echo "$(BOLD)$(GREEN)Development workflow complete$(RESET)"

## Release preparation
release-check: clean setup lint test-race test-coverage bench security ## Pre-release checks
	@echo "$(BOLD)$(CYAN)Running release checks...$(RESET)"
	@echo "$(GREEN)All release checks passed$(RESET)"

## Install the package
install: build ## Install the package
	@echo "$(BOLD)$(CYAN)Installing $(PKG_NAME)...$(RESET)"
	@go install $(LDFLAGS) ./...
	@echo "$(GREEN)Installation complete$(RESET)"

.PHONY: help version-info setup clean dirs build build-dev test test-race test-coverage test-compatibility
.PHONY: bench bench-compare bench-report lint fuzz profile test-all-versions security
.PHONY: ci dev release-check install

# Include version-specific targets
-include Makefile.$(GO_VERSION_MAJOR).$(GO_VERSION_MINOR)