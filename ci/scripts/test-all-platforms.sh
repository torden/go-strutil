#!/bin/bash
# 모든 플랫폼에서 테스트 실행 스크립트

set -e

# 색상 정의
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 로그 함수
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Podman이 설치되어 있는지 확인
check_podman() {
    if ! command -v podman &> /dev/null; then
        log_error "Podman is not installed. Please install Podman first."
        exit 1
    fi
    
    if ! podman info &> /dev/null; then
        log_error "Podman is not running properly."
        exit 1
    fi
    
    log_info "Podman is available"
}

# 테스트 결과 저장 디렉터리 생성
setup_results_dir() {
    RESULTS_DIR="test-results/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$RESULTS_DIR"
    log_info "Test results will be saved to: $RESULTS_DIR"
}

# Ubuntu 테스트
test_ubuntu() {
    local version=$1
    log_info "Testing on Ubuntu $version..."
    
    cat << EOF > Dockerfile.ubuntu-$version
FROM ubuntu:$version

RUN apt-get update && apt-get install -y \\
    golang-go \\
    git \\
    wget \\
    curl \\
    build-essential

# Go 최신 버전 설치 (시스템 Go가 오래된 경우)
RUN if [ "$version" = "18.04" ] || [ "$version" = "20.04" ]; then \\
        wget -O - https://dl.google.com/go/go1.21.0.linux-amd64.tar.gz | tar -C /usr/local -xz && \\
        echo 'export PATH=/usr/local/go/bin:\$PATH' >> /etc/profile; \\
    fi

WORKDIR /app
COPY . .

CMD ["bash", "-c", "source /etc/profile && go version && go mod download && go test -v -race ./... && go test -bench=. -benchmem ./..."]
EOF

    if podman build -f Dockerfile.ubuntu-$version -t strutil-test-ubuntu-$version .; then
        if podman run --rm strutil-test-ubuntu-$version > "$RESULTS_DIR/ubuntu-$version.log" 2>&1; then
            log_success "Ubuntu $version tests passed"
        else
            log_error "Ubuntu $version tests failed"
            cat "$RESULTS_DIR/ubuntu-$version.log"
        fi
    else
        log_error "Failed to build Ubuntu $version image"
    fi
    
    rm -f Dockerfile.ubuntu-$version
}

# Rocky Linux 테스트
test_rocky() {
    local version=$1
    log_info "Testing on Rocky Linux $version..."
    
    cat << EOF > Dockerfile.rocky-$version
FROM rockylinux:$version

RUN dnf update -y && dnf install -y \\
    golang \\
    git \\
    wget \\
    gcc \\
    gcc-c++ \\
    make

WORKDIR /app
COPY . .

CMD ["bash", "-c", "go version && go mod download && go test -v -race ./... && go test -bench=. -benchmem ./..."]
EOF

    if podman build -f Dockerfile.rocky-$version -t strutil-test-rocky-$version .; then
        if podman run --rm strutil-test-rocky-$version > "$RESULTS_DIR/rocky-$version.log" 2>&1; then
            log_success "Rocky Linux $version tests passed"
        else
            log_error "Rocky Linux $version tests failed"
            cat "$RESULTS_DIR/rocky-$version.log"
        fi
    else
        log_error "Failed to build Rocky Linux $version image"
    fi
    
    rm -f Dockerfile.rocky-$version
}

# CentOS Stream 테스트
test_centos_stream() {
    log_info "Testing on CentOS Stream..."
    
    cat << EOF > Dockerfile.centos-stream
FROM quay.io/centos/centos:stream9

RUN dnf update -y && dnf install -y \\
    golang \\
    git \\
    wget \\
    gcc \\
    gcc-c++ \\
    make

WORKDIR /app
COPY . .

CMD ["bash", "-c", "go version && go mod download && go test -v -race ./... && go test -bench=. -benchmem ./..."]
EOF

    if podman build -f Dockerfile.centos-stream -t strutil-test-centos-stream .; then
        if podman run --rm strutil-test-centos-stream > "$RESULTS_DIR/centos-stream.log" 2>&1; then
            log_success "CentOS Stream tests passed"
        else
            log_error "CentOS Stream tests failed"
            cat "$RESULTS_DIR/centos-stream.log"
        fi
    else
        log_error "Failed to build CentOS Stream image"
    fi
    
    rm -f Dockerfile.centos-stream
}

# Alpine Linux 테스트 (경량화)
test_alpine() {
    log_info "Testing on Alpine Linux..."
    
    cat << EOF > Dockerfile.alpine
FROM golang:alpine

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app
COPY . .

CMD ["sh", "-c", "go version && go mod download && go test -v ./... && go test -bench=. -benchmem ./..."]
EOF

    if podman build -f Dockerfile.alpine -t strutil-test-alpine .; then
        if podman run --rm strutil-test-alpine > "$RESULTS_DIR/alpine.log" 2>&1; then
            log_success "Alpine Linux tests passed"
        else
            log_error "Alpine Linux tests failed"
            cat "$RESULTS_DIR/alpine.log"
        fi
    else
        log_error "Failed to build Alpine Linux image"
    fi
    
    rm -f Dockerfile.alpine
}

# macOS 테스트 (로컬에서만 실행 가능)
test_macos() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        log_info "Testing on macOS (local)..."
        
        if go test -v -race ./... > "$RESULTS_DIR/macos.log" 2>&1; then
            log_success "macOS tests passed"
        else
            log_error "macOS tests failed"
            cat "$RESULTS_DIR/macos.log"
        fi
        
        # macOS 벤치마크
        go test -bench=. -benchmem ./... > "$RESULTS_DIR/macos-benchmark.log" 2>&1
    else
        log_warning "macOS tests skipped (not running on macOS)"
    fi
}

# Windows 테스트 (WSL에서 실행 가능)
test_windows_wsl() {
    if command -v cmd.exe &> /dev/null; then
        log_info "Testing on Windows (via WSL)..."
        
        # Windows PowerShell을 통해 테스트 실행
        cmd.exe /c "go test -v ./..." > "$RESULTS_DIR/windows-wsl.log" 2>&1
        
        if [ $? -eq 0 ]; then
            log_success "Windows (WSL) tests passed"
        else
            log_error "Windows (WSL) tests failed"
            cat "$RESULTS_DIR/windows-wsl.log"
        fi
    else
        log_warning "Windows tests skipped (not running in WSL)"
    fi
}

# 멀티 아키텍처 테스트
test_multi_arch() {
    log_info "Testing multiple architectures..."
    
    # AMD64
    log_info "Testing AMD64..."
    GOARCH=amd64 go test -v ./... > "$RESULTS_DIR/amd64.log" 2>&1
    
    # ARM64 (Apple Silicon이나 ARM 서버에서)
    if [[ $(uname -m) == "arm64" ]] || [[ $(uname -m) == "aarch64" ]]; then
        log_info "Testing ARM64..."
        GOARCH=arm64 go test -v ./... > "$RESULTS_DIR/arm64.log" 2>&1
    fi
    
    # 32bit (호환성 테스트)
    log_info "Testing 386..."
    GOARCH=386 go test -v ./... > "$RESULTS_DIR/386.log" 2>&1 || log_warning "386 architecture test failed (may not be supported)"
}

# 성능 비교 리포트 생성
generate_performance_report() {
    log_info "Generating performance comparison report..."
    
    cat << EOF > "$RESULTS_DIR/performance-summary.md"
# Performance Test Summary

Generated on: $(date)

## Test Results by Platform

EOF

    for logfile in "$RESULTS_DIR"/*.log; do
        if [[ -f "$logfile" ]]; then
            platform=$(basename "$logfile" .log)
            echo "### $platform" >> "$RESULTS_DIR/performance-summary.md"
            echo "\`\`\`" >> "$RESULTS_DIR/performance-summary.md"
            
            # 벤치마크 결과 추출
            grep -E "Benchmark|PASS|FAIL" "$logfile" | head -20 >> "$RESULTS_DIR/performance-summary.md" || echo "No benchmark data found" >> "$RESULTS_DIR/performance-summary.md"
            
            echo "\`\`\`" >> "$RESULTS_DIR/performance-summary.md"
            echo "" >> "$RESULTS_DIR/performance-summary.md"
        fi
    done
    
    log_success "Performance report saved to: $RESULTS_DIR/performance-summary.md"
}

# 메인 실행 함수
main() {
    log_info "Starting cross-platform tests..."
    
    check_podman
    setup_results_dir
    
    # 각 플랫폼별 테스트 실행
    test_ubuntu "20.04"
    test_ubuntu "22.04"
    test_rocky "8"
    test_rocky "9"
    test_centos_stream
    test_alpine
    test_macos
    test_windows_wsl
    test_multi_arch
    
    # 리포트 생성
    generate_performance_report
    
    log_success "All cross-platform tests completed!"
    log_info "Results available in: $RESULTS_DIR"
}

# 스크립트 실행
main "$@"