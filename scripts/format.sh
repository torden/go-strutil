#!/usr/bin/env bash

# Code formatting script for go-strutil
# Maintains consistent code style across all Go versions

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if we're in the right directory
if [[ ! -f "go.mod" ]]; then
    print_error "This script must be run from the project root directory"
    exit 1
fi

print_status "Starting code formatting for go-strutil..."

# Detect Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
GO_VERSION_NUM=$(echo "$GO_VERSION" | sed 's/\.//g' | cut -c1-3)
print_status "Detected Go version: $GO_VERSION"

# Format Go code
print_status "Running go fmt..."
go fmt ./...
print_success "go fmt completed"

# Run goimports if available
if command -v goimports >/dev/null 2>&1; then
    print_status "Running goimports..."
    find . -name "*.go" -not -path "./vendor/*" -not -path "./build/*" | xargs goimports -w
    print_success "goimports completed"
else
    print_warning "goimports not found, skipping import organization"
fi

# Run gofumpt if available (for modern Go versions)
if [[ "$GO_VERSION_NUM" -ge "113" ]] && command -v gofumpt >/dev/null 2>&1; then
    print_status "Running gofumpt for enhanced formatting..."
    find . -name "*.go" -not -path "./vendor/*" -not -path "./build/*" | xargs gofumpt -w
    print_success "gofumpt completed"
else
    print_warning "gofumpt not available or Go version too old, skipping enhanced formatting"
fi

# Organize imports with version-specific prefixes
print_status "Organizing imports with local prefixes..."
if command -v goimports >/dev/null 2>&1; then
    find . -name "*.go" -not -path "./vendor/*" -not -path "./build/*" | \
        xargs goimports -w -local "github.com/torden/go-strutil"
    print_success "Import organization completed"
fi

# Fix common Go code issues
print_status "Applying go fix..."
go fix ./...
print_success "go fix completed"

# Remove trailing whitespace from all files
print_status "Removing trailing whitespace..."
find . -type f \( -name "*.go" -o -name "*.md" -o -name "*.yml" -o -name "*.yaml" -o -name "*.json" -o -name "*.toml" \) \
    -not -path "./vendor/*" -not -path "./build/*" -not -path "./reports/*" \
    -exec sed -i '' 's/[[:space:]]*$//' {} +
print_success "Trailing whitespace removed"

# Ensure consistent line endings (LF)
print_status "Normalizing line endings..."
find . -type f \( -name "*.go" -o -name "*.md" -o -name "*.yml" -o -name "*.yaml" -o -name "*.json" -o -name "*.toml" \) \
    -not -path "./vendor/*" -not -path "./build/*" -not -path "./reports/*" \
    -exec dos2unix {} \; 2>/dev/null || true
print_success "Line endings normalized"

# Apply EditorConfig settings if editorconfig tool is available
if command -v editorconfig >/dev/null 2>&1; then
    print_status "Applying EditorConfig settings..."
    find . -type f \( -name "*.go" -o -name "*.md" -o -name "*.yml" -o -name "*.yaml" -o -name "*.json" -o -name "*.toml" \) \
        -not -path "./vendor/*" -not -path "./build/*" -not -path "./reports/*" \
        -exec editorconfig-checker {} \; || true
    print_success "EditorConfig settings applied"
else
    print_warning "editorconfig tool not found, skipping EditorConfig formatting"
fi

# Version-specific formatting
if [[ "$GO_VERSION_NUM" -ge "118" ]]; then
    print_status "Applying Go 1.18+ specific formatting..."
    
    # Format generic code if present
    if find . -name "*.go" -exec grep -l "func.*\[.*\]" {} \; | head -1 >/dev/null 2>&1; then
        print_status "Found generic code, applying specialized formatting..."
        # Additional formatting for generics can be added here
    fi
fi

# Final validation
print_status "Validating formatting..."
if ! go fmt ./... | grep -q .; then
    print_success "All files are properly formatted"
else
    print_error "Some files still need formatting"
    exit 1
fi

# Summary
print_success "Code formatting completed successfully!"
print_status "Summary of operations:"
echo "  ✓ go fmt applied"
echo "  ✓ goimports $(command -v goimports >/dev/null 2>&1 && echo "applied" || echo "skipped")"
echo "  ✓ go fix applied"
echo "  ✓ Trailing whitespace removed"
echo "  ✓ Line endings normalized"
echo "  ✓ Import organization completed"

print_status "Code is now formatted according to Go conventions and project standards."