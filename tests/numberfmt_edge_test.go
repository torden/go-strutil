package strutils_test

import (
	"math"
	"testing"
)

func Test_NumberFmt_EdgeCases(t *testing.T) {
	t.Parallel()

	// Test extreme numeric values and edge cases
	testCases := []struct {
		name        string
		input       interface{}
		expected    string
		expectError bool
	}{
		// Integer boundaries
		{"MaxInt64", math.MaxInt64, "9,223,372,036,854,775,807", false},
		{"MinInt64", math.MinInt64, "-9,223,372,036,854,775,808", false},
		
		// Floating point special values  
		{"Positive Infinity", math.Inf(1), "+,Inf", false},
		{"Negative Infinity", math.Inf(-1), " -Inf", false},
		{"NaN", math.NaN(), "NaN", false},
		
		// Zero variations
		{"Zero int64", int64(0), "0", false},
		{"Zero float64", float64(0), "0", false},
		
		// Large floating point numbers
		{"Large scientific", 1e100, "1e+,100", false},
		{"Medium scientific", 1.23456789e50, "1.23456789e+50", false},
		{"Small scientific", 1e-50, "1e,-50", false},
		
		// Complex numbers (should error)
		{"Complex64", complex64(1 + 2i), "", true},
		{"Complex128", complex128(3 + 4i), "", true},
		
		// Invalid string inputs (should error)
		{"Invalid string", "not_a_number", "", true},
		{"Multiple decimals", "12.34.56", "", true},
		{"Empty string", "", "", true},
		{"Mixed alphanumeric", "abc123", "", true},
		{"Number with suffix", "123abc", "", true},
		
		// Valid numeric strings (some are supported)
		{"Leading zeros", "000123456", "000,123,456", false},
		{"Plus sign", "+123456", "+,123,456", false},
		{"Whitespace padded", "  123456  ", "", true},
		{"Negative number", "-987654321", " -987,654,321", false},
		
		// Scientific notation strings (some are supported)
		{"Scientific string 1", "1.23e10", "1.23e10", false},
		{"Scientific string 2", "1.5e-5", "1.5e-5", false},
		{"Scientific string 3", "1e+100", "1e+,100", false},
	}

	// Test hashable types
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := strproc.NumberFmt(testCase.input)
			
			if testCase.expectError {
				assert.AssertNotNil(t, err, "Expected error for input %v", testCase.input)
			} else {
				assert.AssertNil(t, err, "Unexpected error for input %v: %v", testCase.input, err)
				assert.AssertEquals(t, testCase.expected, result, 
					"NumberFmt failed for input %v.\nExpected: %q\nActual: %q", 
					testCase.input, testCase.expected, result)
			}
		})
	}
	
	// Test unhashable types separately
	unhashableTypes := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{"Slice", []int{1, 2, 3}, true},
		{"Map", map[string]int{"a": 1}, true},
		{"Struct", struct{}{}, true},
		{"Channel", make(chan int), true},
		{"Function", func() {}, true},
	}
	
	for _, testCase := range unhashableTypes {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := strproc.NumberFmt(testCase.input)
			if testCase.expectError {
				assert.AssertNotNil(t, err, "Expected error for unhashable type %T", testCase.input)
			}
		})
	}
}

func Test_NumberFmt_LargeNumbers(t *testing.T) {
	t.Parallel()

	// Test very large numbers
	largeNumbers := []struct {
		input    interface{}
		expected string
	}{
		// Test maximum values for different integer types
		{int8(127), "127"},
		{int8(-128), "-128"},
		{uint8(255), "255"},
		{int16(32767), "32,767"},
		{int16(-32768), "-32,768"},
		{uint16(65535), "65,535"},
		{int32(2147483647), "2,147,483,647"},
		{int32(-2147483648), "-2,147,483,648"},
		{uint32(4294967295), "4,294,967,295"},
		
		// Very large floating point
		{float32(3.4028235e+38), "3.4028235e+38"},
		{float64(1.7976931348623157e+308), "1.7976931348623157e+308"},
		
		// Very small floating point  
		{float32(1.175494351e-38), "1.175494351e-38"},
		{float64(2.2250738585072014e-308), "2.2250738585072014e-308"},
	}

	for _, testCase := range largeNumbers {
		result, err := strproc.NumberFmt(testCase.input)
		assert.AssertNil(t, err, "Unexpected error for input %v: %v", testCase.input, err)
		assert.AssertEquals(t, testCase.expected, result,
			"NumberFmt failed for large number %v.\nExpected: %q\nActual: %q",
			testCase.input, testCase.expected, result)
	}
}

func Test_NumberFmt_PrecisionEdgeCases(t *testing.T) {
	t.Parallel()

	// Test floating point precision edge cases
	precisionCases := []struct {
		input    float64
		expected string
	}{
		{1.0000000000000002, "1.0000000000000002"},  // Near machine epsilon
		{0.1 + 0.2, "0.30000000000000004"},          // Classic floating point precision
		{1.0 / 3.0, "0.3333333333333333"},           // Repeating decimal
		{math.Pi, "3.141592653589793"},              // Mathematical constant
		{math.E, "2.718281828459045"},               // Euler's number
		{math.Sqrt(2), "1.4142135623730951"},        // Square root of 2
		{1e-15, "0.000000000000001"},                // Very small but not scientific
		{1e-16, "1e-16"},                            // Should use scientific notation
	}

	for _, testCase := range precisionCases {
		result, err := strproc.NumberFmt(testCase.input)
		assert.AssertNil(t, err, "Unexpected error for input %v: %v", testCase.input, err)
		assert.AssertEquals(t, testCase.expected, result,
			"NumberFmt precision test failed for %v.\nExpected: %q\nActual: %q",
			testCase.input, testCase.expected, result)
	}
}

func Test_NumberFmt_StringEdgeCases(t *testing.T) {
	t.Parallel()

	// Test various string number formats
	stringCases := map[string]struct {
		expected    string
		expectError bool
	}{
		// Whitespace handling
		" 123456 ":     {"123,456", false},
		"\t123456\t":   {"123,456", false},
		"\n123456\n":   {"123,456", false},
		"  +123456  ":  {"123,456", false},
		"  -123456  ":  {"-123,456", false},
		
		// Leading zeros
		"000000123":     {"123", false},
		"000000000":     {"0", false},
		"00123.456":     {"123.456", false},
		
		// Different decimal formats
		"123456.0":      {"123,456", false},
		"123456.00":     {"123,456", false},
		"123456.000000": {"123,456", false},
		"0.123456":      {"0.123456", false},
		".123456":       {"0.123456", false},
		
		// Scientific notation strings
		"1e3":      {"1,000", false},
		"1E3":      {"1,000", false},
		"1.5e3":    {"1,500", false},
		"1.5E+3":   {"1,500", false},
		"1.5e-3":   {"0.0015", false},
		"1.5E-3":   {"0.0015", false},
		
		// Invalid formats
		"12,345":       {"", true},  // Already formatted
		"12.34.56":     {"", true},  // Multiple decimals
		"12e":          {"", true},  // Incomplete scientific
		"e10":          {"", true},  // Missing mantissa
		"12.34e":       {"", true},  // Incomplete scientific
		"12.34ef":      {"", true},  // Invalid scientific
		"abc":          {"", true},  // Not a number
		"12abc":        {"", true},  // Mixed
		"abc12":        {"", true},  // Mixed
		"+-123":        {"", true},  // Double sign
		"":             {"", true},  // Empty
		"   ":          {"", true},  // Only whitespace
		"∞":            {"", true},  // Unicode infinity
		"NaN":          {"", true},  // Text NaN (not Go's math.NaN())
	}

	for input, expected := range stringCases {
		result, err := strproc.NumberFmt(input)
		
		if expected.expectError {
			assert.AssertNotNil(t, err, "Expected error for string input %q", input)
		} else {
			assert.AssertNil(t, err, "Unexpected error for string input %q: %v", input, err)
			assert.AssertEquals(t, expected.expected, result,
				"NumberFmt failed for string %q.\nExpected: %q\nActual: %q",
				input, expected.expected, result)
		}
	}
}

func Test_NumberFmt_UnicodeStrings(t *testing.T) {
	t.Parallel()

	// Test Unicode number strings (should all fail)
	unicodeNumbers := []string{
		"１２３４５６",      // Full-width digits
		"१२३४५६",       // Devanagari digits  
		"۱۲۳۴۵۶",      // Persian digits
		"𝟏𝟐𝟑𝟒𝟓𝟔",      // Mathematical bold digits
		"①②③④⑤⑥",      // Circled numbers
		"½",            // Fraction
		"²³",           // Superscript
		"₁₂₃",          // Subscript
	}

	for _, input := range unicodeNumbers {
		_, err := strproc.NumberFmt(input)
		assert.AssertNotNil(t, err, "Expected error for Unicode number %q", input)
	}
}