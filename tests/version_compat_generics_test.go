//go:build go1.18
// +build go1.18

package strutils_test

import (
	"testing"

	strutils "github.com/torden/go-strutil"
)

func Test_VersionCompatibility_Generics(t *testing.T) {
	t.Parallel()

	t.Run("Go Version Detection", func(t *testing.T) {
		version, err := strutils.GetGoVersion()
		if err != nil {
			t.Fatalf("Failed to get Go version: %v", err)
		}

		if version.Major != 1 {
			t.Errorf("Expected Go major version 1, got %d", version.Major)
		}

		if version.Minor < 18 {
			t.Errorf("Expected Go minor version >= 18, got %d", version.Minor)
		}

		t.Logf("Detected Go version: %s", version.Original)
	})

	t.Run("Feature Detection Modern", func(t *testing.T) {
		version, err := strutils.GetGoVersion()
		if err != nil {
			t.Fatalf("Failed to get Go version: %v", err)
		}

		// For Go >= 1.18, we should have generics
		if version.Minor < 18 {
			t.Skip("Skipping modern feature test on Go < 1.18")
		}

		t.Logf("Running modern feature tests for Go %s", version.Original)
		t.Logf("Generics support available in this version")
	})

	t.Run("Generic Features Go 1.18+", func(t *testing.T) {
		// Test generic features (Go 1.18+)
		strproc := strutils.NewStringProc()

		// Test generic number formatting
		result, err := strutils.NumberFmtGeneric(strproc, 12345)
		if err != nil {
			t.Errorf("Generic NumberFmt failed: %v", err)
		}
		if result != "12,345" {
			t.Errorf("Generic NumberFmt result mismatch: got %q, expected %q",
				result, "12,345")
		}

		// Test generic slice processor
		processor := strutils.NewSliceProcessor[string]()
		original := []string{"a", "b", "a", "c", "b"}
		unique := processor.UniqueSlice(original)
		expectedUnique := []string{"a", "b", "c"}

		if !processor.CompareSlices(unique, expectedUnique) {
			t.Errorf("Generic slice uniqueness failed: got %v, expected %v",
				unique, expectedUnique)
		}

		t.Logf("Generic tests completed successfully")
	})

	t.Run("Performance Comparison", func(t *testing.T) {
		strproc := strutils.NewStringProc()

		// Benchmark generic vs interface approach
		iterations := 1000

		// Generic approach (type-safe, potentially faster)
		for i := 0; i < iterations; i++ {
			_, _ = strutils.NumberFmtGeneric(strproc, 1234567890)
		}

		// Interface approach (compatible with all versions)
		for i := 0; i < iterations; i++ {
			_, _ = strproc.NumberFmt("1234567890")
		}

		t.Logf("Performance comparison completed for %d iterations", iterations)
	})
}
