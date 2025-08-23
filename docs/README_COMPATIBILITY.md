# Go Version Compatibility Guide

## Overview

The `go-strutil` library is designed to work across a wide range of Go versions, from legacy versions (Go 1.6+) to the latest releases (Go 1.24+). This document explains the compatibility strategy and how to use version-specific features.

## Supported Go Versions

### Minimum Requirements
- **Go 1.6**: Minimum supported version
- **Go 1.7**: Context package support
- **Go 1.11**: Go modules support
- **Go 1.18**: Generics support (recommended for new projects)

### Version Support Matrix

| Go Version | Support Level | Features Available |
|------------|---------------|-------------------|
| 1.6 - 1.10 | Legacy Support | Core functionality, manual dependency management |
| 1.11 - 1.17 | Modern Support | Full functionality with Go modules |
| 1.18+ | Enhanced Support | All features + generics + performance optimizations |

## Feature Availability by Version

### Core Features (All Versions)
- String processing functions (`AddSlashes`, `StripSlashes`, `Nl2Br`, etc.)
- String validation functions (`IsValidEmail`, `IsValidDomain`, etc.)
- Number formatting (`NumberFmt`)
- String utilities (`ReverseStr`, `Padding`, etc.)

### Go 1.7+ Features
- Context support for cancellation and timeouts
- Enhanced error handling

### Go 1.11+ Features
- Go modules support
- Improved dependency management
- Module-aware builds

### Go 1.18+ Features
- Generic type-safe APIs
- Enhanced performance with type constraints
- Modern error handling with error wrapping

## Usage Examples

### Basic Usage (All Versions)
```go
package main

import (
    "fmt"
    "github.com/torden/go-strutil"
)

func main() {
    strproc := strutils.NewStringProc()

    // Basic string processing - works in all versions
    result := strproc.AddSlashes("test\\string")
    fmt.Println(result) // Output: test\\\\string

    // Number formatting - works in all versions
    formatted, err := strproc.NumberFmt(12345)
    if err != nil {
        panic(err)
    }
    fmt.Println(formatted) // Output: 12,345
}
```

### Generic Usage (Go 1.18+)
```go
package main

import (
    "fmt"
    "github.com/torden/go-strutil"
)

func main() {
    strproc := strutils.NewStringProc()

    // Type-safe generic number formatting
    result, err := strproc.NumberFmtGeneric(12345)
    if err != nil {
        panic(err)
    }
    fmt.Println(result) // Output: 12,345

    // Generic slice processing
    processor := strutils.NewSliceProcessor[string]()
    unique := processor.UniqueSlice([]string{"a", "b", "a", "c"})
    fmt.Println(unique) // Output: [a b c]
}
```

### Context Usage (Go 1.7+)
```go
package main

import (
    "fmt"
    "github.com/torden/go-strutil"
)

func main() {
    // Context support (available in Go 1.7+)
    ctx := strutils.Background()

    // Use context for cancellation in long-running operations
    // (Implementation would depend on specific functions supporting context)
}
```

## Build Tags and Conditional Compilation

The library uses Go build tags to provide version-specific implementations:

### Available Build Tags
- `legacy`: For Go versions < 1.7
- `generics`: For Go versions >= 1.18

### Building for Specific Versions

```bash
# Build for legacy Go versions (1.6-1.10)
go build -tags "legacy" ./...

# Build for modern Go versions without generics (1.11-1.17)
go build ./...

# Build for Go versions with generics (1.18+)
go build -tags "generics" ./...
```

## Testing Across Versions

### Running Tests
```bash
# Test with current Go version
go test ./...

# Test with race detection (Go 1.1+)
go test -race ./...

# Run benchmarks
go test -bench=. -benchmem ./...
```

### CI/CD Testing
The project includes a comprehensive CI matrix that tests across multiple Go versions:

```yaml
# See .github/workflows/ci-matrix.yml for full configuration
strategy:
  matrix:
    go-version: ['1.6', '1.11', '1.18', '1.21', '1.24']
    os: [ubuntu-latest, windows-latest, macos-latest]
```

## Migration Guide

### From Go 1.6-1.10 to 1.11+
1. Enable Go modules: `go mod init`
2. Update imports if needed
3. Remove vendor directories
4. Use `go mod tidy` to manage dependencies

### From Go 1.11-1.17 to 1.18+
1. Update `go.mod` to require Go 1.18+
2. Start using generic APIs for better type safety:
   ```go
   // Before (interface-based)
   result, err := strproc.NumberFmt(value)

   // After (generic)
   result, err := strproc.NumberFmtGeneric(value)
   ```
3. Use generic slice processors for better performance

## Performance Considerations

### Version-Specific Optimizations
- **Go 1.18+**: Generic implementations provide better performance and type safety
- **Go 1.13+**: Enhanced error handling with error wrapping
- **Go 1.10+**: Improved string building performance

### Benchmarking
Run version-specific benchmarks to compare performance:

```bash
# Compare generic vs interface implementations (Go 1.18+)
go test -bench=BenchmarkVersionSpecific ./...
```

## Troubleshooting

### Common Issues

#### Module Resolution (Go 1.11+)
```bash
# Clear module cache if needed
go clean -modcache
go mod download
```

#### Build Tags Not Working
```bash
# Ensure proper spacing in build tags
// +build go1.18

# Use correct tag syntax for newer Go versions
//go:build go1.18
```

#### Legacy Version Support
For Go versions < 1.11, ensure dependencies are in your `GOPATH`:
```bash
go get github.com/torden/go-strutil
```

## Contributing

When contributing to the library:

1. Test across multiple Go versions using the CI matrix
2. Use appropriate build tags for version-specific code
3. Maintain backward compatibility
4. Document any new version requirements

## Future Compatibility

The library will continue to support:
- At least 3 major Go versions back from the latest
- All features available in the minimum supported version
- Gradual adoption of new Go features with fallbacks

For the latest compatibility information, see the [CI test results](https://github.com/torden/go-strutil/actions).