package strutils_test

import (
	"fmt"
	"strings"
	"testing"
)

func Test_AddSlashes_StripSlashes_Roundtrip(t *testing.T) {
	t.Parallel()

	// Test cases that should roundtrip perfectly
	testCases := []string{
		"simple string",
		"string with\ttabs",
		"string with\nnewlines",
		"string with\rcarriage returns",
		"string with\bbackspaces",
		"string with\fform feeds",
		"string with\vvertical tabs",
		"string with\abell",
		"string with\x00null",
		`string with\backslashes`,
		`string with\\double backslashes`,
		`string with"double quotes`,
		`string with'single quotes`,
		"complex\tstring\nwith\rmultiple\bescape\fcharacters\vand\\backslashes",
		"korea\ttest\nstring",
		`C:\Windows\System32\cmd.exe`,
		`SELECT * FROM table WHERE name = "John's \"special\" item"`,
		"",                     // empty string
		"no escape characters", // no escapes needed

		// All ASCII control characters
		"test\x01\x02\x03data",
		"file\x04\x05\x06sep",
		"string\x0e\x0f\x10more",
		"data\x11\x12\x13\x14end",
		"control\x15\x16\x17\x18chars",
		"final\x19\x1a\x1b\x1c\x1d\x1e\x1ftest",
		"delete\x7fchar",

		// Mixed content
		"normal text with \x01 control and \t tab",
		"path/to/file\x1fseparator\x7fdelete",
		"unicode\x00null\x1fcontrol",
	}

	for i, original := range testCases {
		t.Run(fmt.Sprintf("Case_%d", i), func(t *testing.T) {
			// Add slashes
			escaped := strproc.AddSlashes(original)

			// Strip slashes should return to original
			restored := strproc.StripSlashes(escaped)

			if restored != original {
				t.Errorf("Roundtrip failed for case %d:\nOriginal:  %q\nEscaped:   %q\nRestored:  %q",
					i, original, escaped, restored)
			}

			// Log the transformation for verification
			t.Logf("Case %d: %q -> %q -> %q", i, original, escaped, restored)
		})
	}
}

func Test_AddSlashes_EdgeCases(t *testing.T) {
	t.Parallel()

	edgeCases := map[string]string{
		// Single characters - named escapes
		"\t": "\\t",
		"\n": "\\n",
		"\r": "\\r",
		"\b": "\\b",
		"\f": "\\f",
		"\v": "\\v",
		"\a": "\\a",
		"\\": "\\\\",
		"\"": "\\\"",
		"'":  "\\'",

		// Null character
		"\x00": "\\0",

		// Multiple same characters
		"\t\t":   "\\t\\t",
		"\n\n\n": "\\n\\n\\n",
		"\\\\\\": "\\\\\\\\\\\\",

		// Mixed escape sequences
		"\t\n\r": "\\t\\n\\r",
		"\\t":    "\\\\t",     // literal \t
		"\\\\n":  "\\\\\\\\n", // literal \\n

		// Control characters (hex escapes)
		"\x01": "\\x01", // SOH
		"\x02": "\\x02", // STX
		"\x03": "\\x03", // ETX
		"\x04": "\\x04", // EOT
		"\x05": "\\x05", // ENQ
		"\x06": "\\x06", // ACK
		"\x0e": "\\x0e", // SO
		"\x0f": "\\x0f", // SI
		"\x10": "\\x10", // DLE
		"\x11": "\\x11", // DC1
		"\x12": "\\x12", // DC2
		"\x13": "\\x13", // DC3
		"\x14": "\\x14", // DC4
		"\x15": "\\x15", // NAK
		"\x16": "\\x16", // SYN
		"\x17": "\\x17", // ETB
		"\x18": "\\x18", // CAN
		"\x19": "\\x19", // EM
		"\x1a": "\\x1a", // SUB
		"\x1b": "\\x1b", // ESC
		"\x1c": "\\x1c", // FS
		"\x1d": "\\x1d", // GS
		"\x1e": "\\x1e", // RS
		"\x1f": "\\x1f", // US
		"\x7f": "\\x7f", // DEL

		// Unicode with escapes
		"korea\ttest":   "korea\\ttest",
		"unicode\ntest": "unicode\\ntest",

		// Mixed control and normal characters
		"test\x01data":    "test\\x01data",
		"file\x1f\x7fsep": "file\\x1f\\x7fsep",
	}

	for input, expected := range edgeCases {
		result := strproc.AddSlashes(input)
		if result != expected {
			t.Errorf("AddSlashes failed for input %q:\nExpected: %q\nActual:   %q",
				input, expected, result)
		}
	}
}

func Test_StripSlashes_EdgeCases(t *testing.T) {
	t.Parallel()

	edgeCases := map[string]string{
		// Single escape sequences - named
		"\\t":  "\t",
		"\\n":  "\n",
		"\\r":  "\r",
		"\\b":  "\b",
		"\\f":  "\f",
		"\\v":  "\v",
		"\\a":  "\a",
		"\\\\": "\\",
		"\\\"": "\"",
		"\\'":  "'",
		"\\0":  "\x00",

		// Multiple escape sequences
		"\\t\\n\\r": "\t\n\r",
		"\\\\\\\\":  "\\\\",

		// Hexadecimal escape sequences
		"\\x01": "\x01",
		"\\x02": "\x02",
		"\\x0a": "\n", // \x0a is newline
		"\\x0d": "\r", // \x0d is carriage return
		"\\x1f": "\x1f",
		"\\x7f": "\x7f",
		"\\xff": "\xff", // max byte value
		"\\x00": "\x00", // null byte

		// Mixed hex digits case
		"\\x0A": "\n",    // uppercase hex
		"\\xaB": "\xab",  // mixed case
		"\\XFF": "\\XFF", // invalid \X (should keep as literal)

		// Invalid/incomplete hex escapes (should keep backslash)
		"\\x":   "\\x",   // no hex digits
		"\\x1":  "\\x1",  // only one hex digit
		"\\xg1": "\\xg1", // invalid hex digit
		"\\x1g": "\\x1g", // invalid hex digit

		// Invalid/unrecognized escapes (should keep backslash)
		"\\z":   "\\z",
		"\\123": "\\123",
		"\\u":   "\\u", // not supported unicode escape

		// Trailing backslash
		"text\\": "text\\",

		// Empty and edge cases
		"":           "",
		"\\":         "\\", // single trailing backslash
		"no_escapes": "no_escapes",

		// Complex mixed cases
		"test\\x01\\tdata\\x7f": "test\x01\tdata\x7f",
		"path\\\\file\\x1fname": "path\\file\x1fname",
	}

	for input, expected := range edgeCases {
		result := strproc.StripSlashes(input)
		if result != expected {
			t.Errorf("StripSlashes failed for input %q:\nExpected: %q\nActual:   %q",
				input, expected, result)
		}
	}
}

func Benchmark_AddSlashes_Performance(b *testing.B) {

	testStrings := []string{
		"simple string without escapes",
		"string\twith\nmultiple\rescapes\band\\backslashes",
		`complex "quoted" string with 'mixed' quotes and \backslashes`,
		"very long string " + strings.Repeat("with\tescapes\nand\rspecial\bcharacters\\", 100),
	}

	b.ResetTimer()

	for _, str := range testStrings {
		b.Run(fmt.Sprintf("len_%d", len(str)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = strproc.AddSlashes(str)
			}
		})
	}
}

func Benchmark_StripSlashes_Performance(b *testing.B) {

	// Pre-escaped strings for benchmarking
	testStrings := []string{
		"simple string without escapes",
		"string\\twith\\nmultiple\\rescapes\\band\\\\backslashes",
		"complex \\\"quoted\\\" string with \\'mixed\\' quotes and \\\\backslashes",
		"very long string " + strings.Repeat("with\\tescapes\\nand\\rspecial\\bcharacters\\\\", 100),
	}

	b.ResetTimer()

	for _, str := range testStrings {
		b.Run(fmt.Sprintf("len_%d", len(str)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = strproc.StripSlashes(str)
			}
		})
	}
}
