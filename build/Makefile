VERSION				:= $(shell git describe --tags --always --dirty="-dev")
DATE				:= $(shell date -u '+%Y-%m-%d-%H%M UTC')
VERSION_FLAGS		:= -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'

CMD_GO				:= $(shell which go)
CMD_GOLINT			:= $(shell which golint)
CMD_GOIMPORTS		:= $(shell which goimports)
CMD_GOFMT			:= $(shell which gofmt)
CMD_ECHO			:= $(shell which echo)
CMD_RM				:= $(shell which rm)
CMD_MKDIR			:= $(shell which mkdir)

GO_VERSION			:= $(shell go version | awk '{print $$3}' | sed 's/go//')
GO_VERSION_NUM		:= $(shell echo $(GO_VERSION) | sed 's/\.//g' | cut -c1-3)

PATH_RACE_REPORT	:= "/tmp/golang-race.report"

## Setup Environment
setup::
	@$(CMD_ECHO) -e "\033[1;40;32mSetup Build Environment.\033[01;m\x1b[0m"
	@$(CMD_GO) mod download
	@$(CMD_GO) mod verify
	@$(CMD_GO) mod tidy
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run Linter
lint:: setup
	@$(CMD_ECHO) -e "\033[1;40;32mRun Linter.\033[01;m\x1b[0m"
	@$(CMD_GO) vet ./...
	@$(CMD_GOFMT) -l . | grep -v '^$$' || echo "All files formatted"
	@if [ -n "$(CMD_GOLINT)" ]; then $(CMD_GOLINT) ./...; fi
	@if [ -n "$(CMD_GOIMPORTS)" ]; then $(CMD_GOIMPORTS) -w .; fi
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run Go Test
test:: 
	@$(CMD_ECHO) -e "\033[1;40;32mRun Go Test.\033[01;m\x1b[0m"
	@$(CMD_GO) test -v ./...
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run Go Test with Race Detection
test-race:: 
	@$(CMD_ECHO) -e "\033[1;40;32mRun Go Test with Race Detection.\033[01;m\x1b[0m"
	@$(CMD_ECHO) -e "\033[1;40;36mRace detection report: $(PATH_RACE_REPORT)\033[01;m\x1b[0m"
	@GORACE="log_path=$(PATH_RACE_REPORT)" $(CMD_GO) test -race -v ./...
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run Go Test with Coverage
test-coverage::
	@$(CMD_ECHO) -e "\033[1;40;32mRun Go Test with Coverage.\033[01;m\x1b[0m"
	@$(CMD_GO) test -coverprofile=coverage.out -covermode=atomic ./...
	@$(CMD_GO) tool cover -html=coverage.out -o coverage.html
	@$(CMD_GO) tool cover -func=coverage.out
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run Benchmarks
bench::
	@$(CMD_ECHO) -e "\033[1;40;32mRun Benchmarks.\033[01;m\x1b[0m"
	@$(CMD_GO) test -bench=. -benchmem -run=^$$ ./...
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Build
build:: setup
	@$(CMD_ECHO) -e "\033[1;40;32mBuild.\033[01;m\x1b[0m"
	@$(CMD_GO) build $(VERSION_FLAGS) ./...
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Clean-up
clean::
	@$(CMD_ECHO) -e "\033[1;40;32mClean-up.\033[01;m\x1b[0m"
	@$(CMD_RM) -rf *.coverprofile *.out *.html *.test *.prof
	@$(CMD_GO) clean -cache -testcache
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Show current version
version::
	@$(CMD_ECHO) "Version: $(VERSION)"
	@$(CMD_ECHO) "Go Version: $(GO_VERSION)"
	@$(CMD_ECHO) "Build Date: $(DATE)"

## Release workflow
release:: clean setup lint test-race test-coverage
	@$(CMD_ECHO) -e "\033[1;40;32mRelease Checks Complete.\033[01;m\x1b[0m"
	@if [ -z "$(TAG)" ]; then \
		$(CMD_ECHO) "Usage: make release TAG=v1.2.3"; \
		exit 1; \
	fi
	@if [ -n "$$(git status --porcelain)" ]; then \
		$(CMD_ECHO) "Working directory is dirty. Please commit changes first."; \
		exit 1; \
	fi
	@git tag -a $(TAG) -m "Release $(TAG)"
	@$(CMD_ECHO) -e "\033[1;40;36mTagged: $(TAG)\033[01;m\x1b[0m"
	@$(CMD_ECHO) "Push with: git push origin $(TAG)"

## Cross-platform testing (requires Docker)
test-cross-platform::
	@$(CMD_ECHO) -e "\033[1;40;32mRun Cross-Platform Tests.\033[01;m\x1b[0m"
	@if [ -f scripts/test-all-platforms.sh ]; then \
		bash scripts/test-all-platforms.sh; \
	else \
		$(CMD_ECHO) -e "\033[1;40;31mCross-platform test script not found.\033[01;m\x1b[0m"; \
	fi
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Test on specific OS (Podman-based)
test-ubuntu:: 
	@$(CMD_ECHO) -e "\033[1;40;32mTest on Ubuntu.\033[01;m\x1b[0m"
	@podman run --rm -v $(PWD):/app -w /app golang:latest bash -c "go mod download && go test -v ./..."

test-alpine::
	@$(CMD_ECHO) -e "\033[1;40;32mTest on Alpine Linux.\033[01;m\x1b[0m"
	@podman run --rm -v $(PWD):/app -w /app golang:alpine sh -c "apk add --no-cache git gcc musl-dev && go mod download && go test -v ./..."

test-centos::
	@$(CMD_ECHO) -e "\033[1;40;32mTest on CentOS.\033[01;m\x1b[0m"
	@podman run --rm -v $(PWD):/app -w /app quay.io/centos/centos:stream9 bash -c "dnf install -y golang git && go mod download && go test -v ./..."

## Multi-architecture testing
test-multiarch::
	@$(CMD_ECHO) -e "\033[1;40;32mRun Multi-Architecture Tests.\033[01;m\x1b[0m"
	@$(CMD_ECHO) "Testing AMD64..."
	@GOARCH=amd64 $(CMD_GO) test -v ./...
	@$(CMD_ECHO) "Testing 386..."
	@GOARCH=386 $(CMD_GO) test -v ./... || $(CMD_ECHO) "386 tests failed (may not be supported)"
	@if [ "$$(uname -m)" = "arm64" ] || [ "$$(uname -m)" = "aarch64" ]; then \
		$(CMD_ECHO) "Testing ARM64..."; \
		GOARCH=arm64 $(CMD_GO) test -v ./...; \
	fi
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Test with different Go versions (requires Podman)
test-go-versions::
	@$(CMD_ECHO) -e "\033[1;40;32mTest with Different Go Versions.\033[01;m\x1b[0m"
	@for version in 1.18 1.19 1.20 1.21; do \
		$(CMD_ECHO) "Testing with Go $$version..."; \
		podman run --rm -v $(PWD):/app -w /app golang:$$version bash -c "go mod download && go test -v ./..." || \
		$(CMD_ECHO) "Go $$version tests failed"; \
	done
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Performance benchmarks across platforms
bench-cross-platform::
	@$(CMD_ECHO) -e "\033[1;40;32mRun Cross-Platform Benchmarks.\033[01;m\x1b[0m"
	@$(CMD_MKDIR) -p bench-results
	@$(CMD_ECHO) "Local benchmarks..."
	@$(CMD_GO) test -bench=. -benchmem -run=^$$ ./... > bench-results/local.txt
	@$(CMD_ECHO) "Ubuntu benchmarks..."
	@podman run --rm -v $(PWD):/app -w /app golang:latest bash -c "go mod download && go test -bench=. -benchmem -run=^$$ ./..." > bench-results/ubuntu.txt || true
	@$(CMD_ECHO) "Alpine benchmarks..."
	@podman run --rm -v $(PWD):/app -w /app golang:alpine sh -c "apk add --no-cache git gcc musl-dev && go mod download && go test -bench=. -benchmem -run=^$$ ./..." > bench-results/alpine.txt || true
	@$(CMD_ECHO) -e "\033[1;40;36mBenchmark results saved to bench-results/\033[01;m\x1b[0m"

## Full CI pipeline
ci:: clean setup lint test-race test-coverage bench
	@$(CMD_ECHO) -e "\033[1;40;32mCI Pipeline Complete.\033[01;m\x1b[0m"

## Extended CI with cross-platform testing
ci-extended:: clean setup lint test-race test-coverage bench test-multiarch
	@$(CMD_ECHO) -e "\033[1;40;32mExtended CI Pipeline Complete.\033[01;m\x1b[0m"

## Show Help
help::
	@$(CMD_ECHO) "Available targets:"
	@$(CMD_ECHO) "  setup        - Setup build environment"
	@$(CMD_ECHO) "  lint         - Run linter"
	@$(CMD_ECHO) "  test         - Run tests"
	@$(CMD_ECHO) "  test-race    - Run tests with race detection"
	@$(CMD_ECHO) "  test-coverage - Run tests with coverage"
	@$(CMD_ECHO) "  bench        - Run benchmarks"
	@$(CMD_ECHO) "  build        - Build the project"
	@$(CMD_ECHO) "  clean        - Clean build artifacts"
	@$(CMD_ECHO) "  version      - Show version information"
	@$(CMD_ECHO) "  release TAG=vX.Y.Z - Create a release tag"
	@$(CMD_ECHO) "  ci           - Run full CI pipeline"
	@$(CMD_ECHO) "  ci-extended  - Run extended CI with cross-platform tests"
	@$(CMD_ECHO) ""
	@$(CMD_ECHO) "Cross-platform testing:"
	@$(CMD_ECHO) "  test-cross-platform - Run all cross-platform tests"
	@$(CMD_ECHO) "  test-ubuntu     - Test on Ubuntu (Docker)"
	@$(CMD_ECHO) "  test-alpine     - Test on Alpine Linux (Docker)"  
	@$(CMD_ECHO) "  test-centos     - Test on CentOS (Docker)"
	@$(CMD_ECHO) "  test-multiarch  - Test multiple architectures"
	@$(CMD_ECHO) "  test-go-versions - Test with different Go versions"
	@$(CMD_ECHO) "  bench-cross-platform - Run benchmarks across platforms"
	@$(CMD_ECHO) ""
	@$(CMD_ECHO) "  help         - Show this help"

.PHONY: setup lint test test-race test-coverage bench build clean version release ci ci-extended help
.PHONY: test-cross-platform test-ubuntu test-alpine test-centos test-multiarch test-go-versions bench-cross-platform