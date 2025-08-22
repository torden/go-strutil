package strutils_test

import (
	"strings"
	"testing"

	strutils "github.com/torden/go-strutil"
)

func TestGoVersionDetection(t *testing.T) {
	version, err := strutils.GetGoVersion()
	if err != nil {
		t.Fatalf("Failed to get Go version: %v", err)
	}

	if version.Major == 0 {
		t.Error("Major version should not be zero")
	}

	// Test that we can parse the current runtime version
	if !strings.HasPrefix(version.Original, "go") {
		t.Errorf("Original version should start with 'go', got: %s", version.Original)
	}

	t.Logf("Detected Go version: %s", version.String())
}

func TestVersionComparison(t *testing.T) {
	version, err := strutils.GetGoVersion()
	if err != nil {
		t.Fatalf("Failed to get Go version: %v", err)
	}

	// Test basic version comparison logic
	if !version.IsVersionAtLeast(1, 0) {
		t.Error("Current Go version should be at least 1.0")
	}

	// Test context support detection
	hasContext := version.SupportsContext()
	expectedContext := version.IsVersionAtLeast(1, 7)
	if hasContext != expectedContext {
		t.Errorf("Context support mismatch: got %v, expected %v", hasContext, expectedContext)
	}

	// Test modules support detection
	hasModules := version.SupportsModules()
	expectedModules := version.IsVersionAtLeast(1, 11)
	if hasModules != expectedModules {
		t.Errorf("Modules support mismatch: got %v, expected %v", hasModules, expectedModules)
	}

	// Test generics support detection
	hasGenerics := version.SupportsGenerics()
	expectedGenerics := version.IsVersionAtLeast(1, 18)
	if hasGenerics != expectedGenerics {
		t.Errorf("Generics support mismatch: got %v, expected %v", hasGenerics, expectedGenerics)
	}

	t.Logf("Feature support - Context: %v, Modules: %v, Generics: %v",
		hasContext, hasModules, hasGenerics)
}

func TestBasicCompatibilityFeatures(t *testing.T) {
	// Test background context creation
	ctx := strutils.Background()
	if ctx == nil {
		t.Error("Background context should not be nil")
	}
}

func TestCoreStringProcessing(t *testing.T) {
	version, err := strutils.GetGoVersion()
	if err != nil {
		t.Fatalf("Failed to get Go version: %v", err)
	}

	t.Run("BasicFeatures", func(t *testing.T) {
		// These should work in all supported versions (1.6+)
		strproc := strutils.NewStringProc()

		result := strproc.AddSlashes(`test\string`)
		if !strings.Contains(result, `\\`) {
			t.Error("AddSlashes should work in all Go versions")
		}

		result = strproc.ReverseStr("hello")
		if result != "olleh" {
			t.Errorf("ReverseStr failed: got %q, expected %q", result, "olleh")
		}

		// Test number formatting
		formatted, err := strproc.NumberFmt(12345)
		if err != nil {
			t.Errorf("NumberFmt failed: %v", err)
		}
		if formatted != "12,345" {
			t.Errorf("NumberFmt result mismatch: got %q, expected %q",
				formatted, "12,345")
		}
	})

	if version.SupportsGenerics() {
		t.Run("GenericsFeatures", func(t *testing.T) {
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
				t.Errorf("Generic unique slice failed: got %v, expected %v",
					unique, expectedUnique)
			}
		})
	}

	t.Logf("Testing completed for Go %s", version.String())
}

// BenchmarkBasicPerformance benchmarks core functionality
func BenchmarkBasicPerformance(b *testing.B) {
	version, err := strutils.GetGoVersion()
	if err != nil {
		b.Fatalf("Failed to get Go version: %v", err)
	}

	b.Run("StringProcessing", func(b *testing.B) {
		strproc := strutils.NewStringProc()

		b.Run("AddSlashes", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = strproc.AddSlashes("test\\string\\with\\slashes")
			}
		})

		b.Run("ReverseStr", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = strproc.ReverseStr("hello world")
			}
		})

		b.Run("NumberFmt", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = strproc.NumberFmt(1234567890)
			}
		})
	})

	if version.SupportsGenerics() {
		b.Run("GenericVsInterface", func(b *testing.B) {
			strproc := strutils.NewStringProc()

			b.Run("Generic", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = strutils.NumberFmtGeneric(strproc, 1234567890)
				}
			})

			b.Run("Interface", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, _ = strproc.NumberFmt(1234567890)
				}
			})
		})
	}
}
