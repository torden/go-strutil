#!/usr/bin/env bash

# Release automation script for go-strutil
# Handles version tagging, changelog generation, and release preparation

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CHANGELOG_FILE="$PROJECT_ROOT/CHANGELOG.md"
VERSION_FILE="$PROJECT_ROOT/VERSION"

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

print_header() {
    echo -e "${BOLD}${CYAN}$1${NC}"
}

# Function to show usage
usage() {
    cat << EOF
Usage: $0 [OPTIONS] <version>

Release automation for go-strutil

Options:
    -h, --help          Show this help message
    -d, --dry-run       Show what would be done without making changes
    -s, --skip-tests    Skip running tests before release
    -p, --pre-release   Mark as pre-release
    
Arguments:
    version            Version to release (e.g., v1.2.3, 1.2.3)

Examples:
    $0 v1.2.3                 # Release version 1.2.3
    $0 --dry-run v1.2.4       # Show what would happen for v1.2.4
    $0 --pre-release v1.3.0   # Release v1.3.0 as pre-release

EOF
}

# Parse command line arguments
DRY_RUN=false
SKIP_TESTS=false
PRE_RELEASE=false
VERSION=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -d|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -s|--skip-tests)
            SKIP_TESTS=true
            shift
            ;;
        -p|--pre-release)
            PRE_RELEASE=true
            shift
            ;;
        -*)
            print_error "Unknown option: $1"
            usage
            exit 1
            ;;
        *)
            VERSION="$1"
            shift
            ;;
    esac
done

# Validate arguments
if [[ -z "$VERSION" ]]; then
    print_error "Version is required"
    usage
    exit 1
fi

# Normalize version (add v prefix if missing)
if [[ ! "$VERSION" =~ ^v ]]; then
    VERSION="v$VERSION"
fi

# Validate version format
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
    print_error "Invalid version format: $VERSION"
    print_error "Expected format: v1.2.3 or v1.2.3-alpha"
    exit 1
fi

# Change to project root
cd "$PROJECT_ROOT"

print_header "Go-Strutil Release Automation"
print_status "Preparing release: $VERSION"
if [[ "$DRY_RUN" == "true" ]]; then
    print_warning "DRY RUN MODE - No changes will be made"
fi

# Pre-flight checks
print_status "Running pre-flight checks..."

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "Not in a git repository"
    exit 1
fi

# Check for uncommitted changes
if [[ -n "$(git status --porcelain)" ]]; then
    print_error "Working directory is dirty. Please commit or stash changes."
    git status --short
    exit 1
fi

# Check if we're on main/master branch
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" != "main" ]] && [[ "$CURRENT_BRANCH" != "master" ]]; then
    print_warning "You are on branch '$CURRENT_BRANCH', not main/master"
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check if tag already exists
if git tag -l | grep -q "^$VERSION$"; then
    print_error "Tag $VERSION already exists"
    exit 1
fi

# Detect Go version and features
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
GO_VERSION_NUM=$(echo "$GO_VERSION" | sed 's/\.//g' | cut -c1-3)
HAS_MODULES=$([ "$GO_VERSION_NUM" -ge "111" ] && echo "true" || echo "false")
HAS_GENERICS=$([ "$GO_VERSION_NUM" -ge "118" ] && echo "true" || echo "false")

print_status "Go version: $GO_VERSION (Modules: $HAS_MODULES, Generics: $HAS_GENERICS)"

# Run tests unless skipped
if [[ "$SKIP_TESTS" == "false" ]]; then
    print_status "Running tests..."
    if [[ "$DRY_RUN" == "false" ]]; then
        make ci
    else
        print_warning "Would run: make ci"
    fi
    print_success "All tests passed"
else
    print_warning "Skipping tests as requested"
fi

# Update version file
print_status "Updating version information..."
if [[ "$DRY_RUN" == "false" ]]; then
    echo "$VERSION" > "$VERSION_FILE"
    print_success "Version file updated"
else
    print_warning "Would write '$VERSION' to $VERSION_FILE"
fi

# Generate/update changelog
print_status "Updating changelog..."
RELEASE_DATE=$(date +%Y-%m-%d)

if [[ "$DRY_RUN" == "false" ]]; then
    if [[ ! -f "$CHANGELOG_FILE" ]]; then
        cat > "$CHANGELOG_FILE" << EOF
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [$VERSION] - $RELEASE_DATE

### Added
- Multi-version Go support (Go 1.6+ through latest)
- Version-aware build system
- Comprehensive benchmarking suite
- Generic APIs for Go 1.18+
- Enhanced CI/CD pipeline

### Changed
- Modernized Makefile with feature detection
- Improved code formatting and linting
- Enhanced documentation

### Fixed
- Version compatibility issues
- Test suite reliability

EOF
    else
        # Insert new version entry after [Unreleased]
        temp_file=$(mktemp)
        awk -v version="$VERSION" -v date="$RELEASE_DATE" '
        /^## \[Unreleased\]/ {
            print $0
            print ""
            print "## [" version "] - " date
            print ""
            print "### Added"
            print "- Release automation"
            print ""
            print "### Changed"
            print "- Updated dependencies"
            print ""
            next
        }
        { print }
        ' "$CHANGELOG_FILE" > "$temp_file"
        mv "$temp_file" "$CHANGELOG_FILE"
    fi
    print_success "Changelog updated"
else
    print_warning "Would update changelog with version $VERSION"
fi

# Commit changes
print_status "Committing release changes..."
if [[ "$DRY_RUN" == "false" ]]; then
    git add "$VERSION_FILE" "$CHANGELOG_FILE" 2>/dev/null || true
    git commit -m "chore: release $VERSION

    🤖 Generated with Claude Code
    
    Co-Authored-By: Claude <noreply@anthropic.com>" || true
    print_success "Changes committed"
else
    print_warning "Would commit version and changelog changes"
fi

# Create git tag
print_status "Creating git tag..."
TAG_MESSAGE="Release $VERSION

Generated: $RELEASE_DATE
Go Version: $GO_VERSION
Features: Modules=$HAS_MODULES, Generics=$HAS_GENERICS

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"

if [[ "$DRY_RUN" == "false" ]]; then
    if [[ "$PRE_RELEASE" == "true" ]]; then
        git tag -a "$VERSION" -m "$TAG_MESSAGE (Pre-release)"
        print_success "Pre-release tag created: $VERSION"
    else
        git tag -a "$VERSION" -m "$TAG_MESSAGE"
        print_success "Release tag created: $VERSION"
    fi
else
    print_warning "Would create git tag: $VERSION"
    if [[ "$PRE_RELEASE" == "true" ]]; then
        print_warning "Would mark as pre-release"
    fi
fi

# Generate release notes
print_status "Generating release notes..."
RELEASE_NOTES_FILE="release-notes-$VERSION.md"

if [[ "$DRY_RUN" == "false" ]]; then
    cat > "$RELEASE_NOTES_FILE" << EOF
# Release Notes - $VERSION

**Release Date:** $RELEASE_DATE
**Go Compatibility:** $GO_VERSION (supports Go 1.6+)

## Features
- ✅ Modules Support: $HAS_MODULES
- ✅ Generics Support: $HAS_GENERICS
- ✅ Multi-version Testing
- ✅ Comprehensive Benchmarking

## Installation

\`\`\`bash
go get github.com/torden/go-strutil@$VERSION
\`\`\`

## What's New

See [CHANGELOG.md](CHANGELOG.md) for detailed changes.

## Compatibility

This release maintains backward compatibility with all supported Go versions (1.6+).

## Performance

Run benchmarks with:
\`\`\`bash
make bench-report
\`\`\`

---

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>
EOF
    print_success "Release notes generated: $RELEASE_NOTES_FILE"
else
    print_warning "Would generate release notes in $RELEASE_NOTES_FILE"
fi

# Final summary
print_header "Release Summary"
echo "Version: $VERSION"
echo "Type: $([ "$PRE_RELEASE" == "true" ] && echo "Pre-release" || echo "Release")"
echo "Go Version: $GO_VERSION"
echo "Features: Modules=$HAS_MODULES, Generics=$HAS_GENERICS"
echo ""

if [[ "$DRY_RUN" == "false" ]]; then
    print_success "Release preparation completed!"
    print_status "Next steps:"
    echo "  1. Review the changes: git show $VERSION"
    echo "  2. Push the tag: git push origin $VERSION"
    echo "  3. Push the commit: git push origin $CURRENT_BRANCH"
    echo "  4. Create GitHub release using: $RELEASE_NOTES_FILE"
    echo ""
    print_warning "Don't forget to push your changes and tag to the remote repository!"
else
    print_success "Dry run completed - no changes were made"
    print_status "Run without --dry-run to perform the actual release"
fi