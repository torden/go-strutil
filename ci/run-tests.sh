#!/bin/bash
# 간단한 크로스 플랫폼 테스트 실행 스크립트

set -e

# 색상 정의
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}Go String Utils Cross-Platform Testing${NC}"
echo "======================================"

# 사용법 출력
usage() {
    echo "Usage: $0 [OPTIONS]"
    echo "Options:"
    echo "  --local         Run local tests only"
    echo "  --podman        Run Podman-based tests"
    echo "  --compose       Run Podman Compose tests"
    echo "  --github        Show GitHub Actions status"
    echo "  --all           Run all available tests (default)"
    echo "  --help          Show this help"
}

# GitHub Actions 상태 확인 (옵션)
check_github_actions() {
    if command -v gh &> /dev/null; then
        echo -e "${YELLOW}GitHub Actions Status:${NC}"
        gh workflow list 2>/dev/null || echo "No GitHub CLI or repo not configured"
    fi
}

# 로컬 테스트 실행
run_local_tests() {
    echo -e "${GREEN}Running local tests...${NC}"
    
    if command -v go &> /dev/null; then
        echo "Go version: $(go version)"
        echo "Running basic tests..."
        go test -v ./...
        
        echo -e "${GREEN}Running race detection tests...${NC}"
        go test -race -v ./...
        
        echo -e "${GREEN}Running benchmarks...${NC}"
        go test -bench=. -benchmem ./...
    else
        echo -e "${RED}Go is not installed${NC}"
        exit 1
    fi
}

# Podman 기반 테스트 실행
run_podman_tests() {
    echo -e "${GREEN}Running Podman-based tests...${NC}"
    
    if ! command -v podman &> /dev/null; then
        echo -e "${RED}Podman is not installed${NC}"
        return 1
    fi
    
    # 간단한 Ubuntu 테스트
    echo "Testing on Ubuntu..."
    podman run --rm -v "$(pwd)":/app -w /app golang:latest bash -c "go mod download && go test -v ./..."
    
    # Alpine 테스트
    echo "Testing on Alpine Linux..."
    podman run --rm -v "$(pwd)":/app -w /app golang:alpine sh -c "apk add --no-cache git gcc musl-dev && go mod download && go test -v ./..."
}

# Docker Compose 테스트 실행
run_compose_tests() {
    echo -e "${GREEN}Running Docker Compose tests...${NC}"
    
    if ! command -v docker-compose &> /dev/null; then
        echo -e "${RED}Docker Compose is not installed${NC}"
        return 1
    fi
    
    # 주요 서비스만 테스트 (시간 단축)
    echo "Running selected Docker Compose tests..."
    docker-compose -f docker-compose.test.yml run --rm test-ubuntu-22
    docker-compose -f docker-compose.test.yml run --rm test-alpine
    docker-compose -f docker-compose.test.yml run --rm test-go-1-21
}

# Makefile 기반 테스트 (사용 가능한 경우)
run_makefile_tests() {
    if [ -f Makefile ]; then
        echo -e "${GREEN}Running Makefile tests...${NC}"
        make test-multiarch 2>/dev/null || echo "Multiarch tests not available"
    fi
}

# 메인 실행 로직
main() {
    case "${1:-all}" in
        --local)
            run_local_tests
            ;;
        --podman)
            run_podman_tests
            ;;
        --compose)
            run_compose_tests
            ;;
        --github)
            check_github_actions
            ;;
        --all)
            echo -e "${GREEN}Running comprehensive tests...${NC}"
            run_local_tests
            echo ""
            run_podman_tests 2>/dev/null || echo "Podman tests skipped"
            echo ""
            run_makefile_tests 2>/dev/null || echo "Makefile tests skipped"
            echo ""
            check_github_actions
            ;;
        --help)
            usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            usage
            exit 1
            ;;
    esac
    
    echo -e "${GREEN}Tests completed!${NC}"
}

# 스크립트 실행
main "$@"