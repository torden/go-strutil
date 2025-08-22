//go:build !go1.18
// +build !go1.18

package strutils_test

import (
	"testing"

	strutils "github.com/torden/go-strutil"
)

func Test_VersionCompatibility_Legacy(t *testing.T) {
	t.Parallel()

	t.Run("Go Version Detection", func(t *testing.T) {
		version, err := strutils.GetGoVersion()
		if err != nil {
			t.Fatalf("Failed to get Go version: %v", err)
		}

		if version.Major != 1 {
			t.Errorf("Expected Go major version 1, got %d", version.Major)
		}

		if version.Minor < 6 {
			t.Errorf("Expected Go minor version >= 6, got %d", version.Minor)
		}

		t.Logf("Detected Go version: %s", version.Original)
	})

	t.Run("Feature Detection Legacy", func(t *testing.T) {
		version, err := strutils.GetGoVersion()
		if err != nil {
			t.Fatalf("Failed to get Go version: %v", err)
		}

		// For Go < 1.18, check that we're using legacy features
		if version.Minor >= 18 {
			t.Skip("Skipping legacy test on Go 1.18+")
		}

		t.Logf("Running legacy feature tests for Go %s", version.Original)
	})

	t.Run("Basic StringProc Legacy", func(t *testing.T) {
		strproc := strutils.NewStringProc()

		// Test basic number formatting (interface-based)
		result, err := strproc.NumberFmt("12345")
		if err != nil {
			t.Errorf("NumberFmt failed: %v", err)
		}
		if result != "12,345" {
			t.Errorf("NumberFmt result mismatch: got %q, expected %q",
				result, "12,345")
		}

		// Test padding
		padded, err := strproc.PaddingBoth("test", " ", 10)
		if err != nil {
			t.Errorf("PaddingBoth failed: %v", err)
		}
		if len(padded) != 10 {
			t.Errorf("PaddingBoth length mismatch: got %d, expected 10", len(padded))
		}

		t.Logf("Legacy tests completed successfully")
	})
}