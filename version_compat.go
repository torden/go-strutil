//go:build go1.6
// +build go1.6

package strutils

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// GoVersionInfo holds Go version information
type GoVersionInfo struct {
	Major    int
	Minor    int
	Patch    int
	Original string
}

// GetGoVersion returns the current Go version information
func GetGoVersion() (*GoVersionInfo, error) {
	version := runtime.Version()

	// Remove "go" prefix if present
	if strings.HasPrefix(version, "go") {
		version = version[2:]
	}

	// Split version string
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid Go version format: %s", runtime.Version())
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %v", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %v", err)
	}

	var patch int
	if len(parts) >= 3 {
		// Handle beta/rc versions like "1.18beta1" or "1.18rc1"
		patchStr := parts[2]
		// Extract numeric part only
		var numStr strings.Builder
		for _, r := range patchStr {
			if r >= '0' && r <= '9' {
				numStr.WriteRune(r)
			} else {
				break
			}
		}
		if numStr.Len() > 0 {
			patch, _ = strconv.Atoi(numStr.String())
		}
	}

	return &GoVersionInfo{
		Major:    major,
		Minor:    minor,
		Patch:    patch,
		Original: runtime.Version(),
	}, nil
}

// IsVersionAtLeast checks if the current Go version is at least the specified version
func (v *GoVersionInfo) IsVersionAtLeast(major, minor int) bool {
	if v.Major > major {
		return true
	}
	if v.Major == major && v.Minor >= minor {
		return true
	}
	return false
}

// SupportsContext returns true if the Go version supports context package
func (v *GoVersionInfo) SupportsContext() bool {
	return v.IsVersionAtLeast(1, 7)
}

// SupportsModules returns true if the Go version supports Go modules
func (v *GoVersionInfo) SupportsModules() bool {
	return v.IsVersionAtLeast(1, 11)
}

// SupportsGenerics returns true if the Go version supports generics
func (v *GoVersionInfo) SupportsGenerics() bool {
	return v.IsVersionAtLeast(1, 18)
}

// SupportsEmbedFS returns true if the Go version supports embed.FS
func (v *GoVersionInfo) SupportsEmbedFS() bool {
	return v.IsVersionAtLeast(1, 16)
}

// String returns a string representation of the version
func (v *GoVersionInfo) String() string {
	return fmt.Sprintf("Go %d.%d.%d (%s)", v.Major, v.Minor, v.Patch, v.Original)
}
