# Windows Development Guide for go-strutil

This guide provides Windows-specific instructions for building and testing go-strutil.

## Prerequisites

1. **Go Installation**: Install Go from https://golang.org/dl/
2. **Git**: Install Git from https://git-scm.com/download/win
3. **PowerShell**: Use PowerShell or Command Prompt

## Quick Start (Windows)

### Using Batch Script (Recommended)
```batch
# Build everything
build_windows.bat all

# Run tests only
build_windows.bat test

# Run benchmarks
build_windows.bat bench

# Format code
build_windows.bat fmt
```

### Using Windows Makefile
```powershell
# If you have make installed (e.g., through Chocolatey)
make -f Makefile.windows build
make -f Makefile.windows test
make -f Makefile.windows bench
```

### Using Go Commands Directly
```powershell
# Download dependencies
go mod download

# Format code
go fmt ./...

# Run tests
go test -v ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# Build
go build -v ./...
```

## Go Version Support on Windows

### Fully Supported Versions
- Go 1.11+ (recommended for Windows)
- Go 1.18+ (with generics support)
- Go 1.20+ (latest features)

### Limited Support
- Go 1.6-1.10 (legacy support, may have CI limitations on Windows)

## Windows-Specific Features

### Build Tags Detection
The build system automatically detects your Go version and applies appropriate build tags:

- **Go 1.18+**: Uses `generics` tag for enhanced performance
- **Go 1.11-1.17**: Standard build without generics
- **Go 1.6-1.10**: Legacy build mode

### Path Handling
All file paths in tests have been updated to be cross-platform compatible:
- Uses relative paths instead of Unix-style absolute paths
- Handles Windows path separators correctly

## Build Artifacts

After running builds, you'll find:
- `build/`: Compiled binaries
- `reports/benchmarks/`: Benchmark results (`bench-windows.txt`)
- `reports/coverage/`: Coverage reports (HTML format)

## Common Issues and Solutions

### Issue: "go: command not found"
**Solution**: Ensure Go is installed and added to your PATH environment variable.

### Issue: Git commands fail
**Solution**: Install Git and ensure it's in your PATH, or use GitHub Desktop.

### Issue: PowerShell execution policy
**Solution**: Run in PowerShell as Administrator:
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Issue: Module download fails
**Solution**: Configure Go proxy:
```powershell
go env -w GOPROXY=https://proxy.golang.org,direct
go env -w GOSUMDB=sum.golang.org
```

## IDE Support

### VS Code
1. Install the Go extension
2. Open the project folder
3. The extension will auto-detect the module

### GoLand
1. Open the project folder
2. GoLand will automatically detect Go modules
3. Run/Debug configurations work out of the box

### Notepad++ / Other Editors
Use the command line tools provided in this guide.

## CI/CD on Windows

The GitHub Actions CI includes Windows builds:
- Tests run on `windows-latest`
- Supports Go 1.11+ (older versions excluded due to GitHub Actions limitations)
- Automatic cross-platform testing

## Performance on Windows

Benchmark results are typically within 5-10% of Linux performance. Windows-specific optimizations include:
- Efficient string processing
- Memory allocation optimization
- Cross-platform file handling

## Troubleshooting

### Enable Verbose Logging
```powershell
go test -v -x ./...
go build -v -x ./...
```

### Check Environment
```powershell
go version
go env
```

### Clean Build Cache
```powershell
go clean -cache
go clean -modcache
go clean -testcache
```

## Contributing on Windows

1. Fork the repository
2. Create a feature branch
3. Use `build_windows.bat all` to test your changes
4. Commit and push your changes
5. Create a pull request

The CI system will automatically test your changes on Windows, Linux, and macOS.