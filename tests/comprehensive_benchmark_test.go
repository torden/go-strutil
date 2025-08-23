package strutils_test

import (
	"math"
	"strings"
	"testing"

	strutils "github.com/torden/go-strutil"
)

// StringProc 함수들의 종합 벤치마크
func BenchmarkStringProc_AddSlashes(b *testing.B) {
	strproc := strutils.NewStringProc()
	testCases := map[string]string{
		"Small":  strings.Repeat("abc'def", 10),
		"Medium": strings.Repeat("abc'def\"ghi\\", 100),
		"Large":  strings.Repeat("test'data\"with\\slashes", 1000),
	}

	for name, testData := range testCases {
		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(testData)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := strproc.AddSlashes(testData)
				_ = result // 컴파일러 최적화 방지
			}
		})
	}
}

func BenchmarkStringProc_StripSlashes(b *testing.B) {
	strproc := strutils.NewStringProc()
	testCases := map[string]string{
		"Small":  strings.Repeat("abc\\'def", 10),
		"Medium": strings.Repeat("abc\\'def\\\"ghi\\\\", 100),
		"Large":  strings.Repeat("test\\'data\\\"with\\\\slashes", 1000),
	}

	for name, testData := range testCases {
		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(testData)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := strproc.StripSlashes(testData)
				_ = result
			}
		})
	}
}

func BenchmarkStringProc_UpperCaseFirstWords(b *testing.B) {
	strproc := strutils.NewStringProc()
	testCases := map[string]string{
		"Small":  "hello world test case",
		"Medium": strings.Repeat("hello world test case ", 25),
		"Large":  strings.Repeat("hello world test case with many words ", 100),
	}

	for name, testData := range testCases {
		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(testData)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := strproc.UpperCaseFirstWords(testData)
				_ = result
			}
		})
	}
}

func BenchmarkStringProc_LowerCaseFirstWords(b *testing.B) {
	strproc := strutils.NewStringProc()
	testCases := map[string]string{
		"Small":  "HELLO WORLD TEST CASE",
		"Medium": strings.Repeat("HELLO WORLD TEST CASE ", 25),
		"Large":  strings.Repeat("HELLO WORLD TEST CASE WITH MANY WORDS ", 100),
	}

	for name, testData := range testCases {
		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(testData)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := strproc.LowerCaseFirstWords(testData)
				_ = result
			}
		})
	}
}

func BenchmarkStringProc_ReverseStr(b *testing.B) {
	strproc := strutils.NewStringProc()
	testCases := map[string]string{
		"ASCII_Small":     "abcdefghijklmnopqrstuvwxyz",
		"ASCII_Medium":    strings.Repeat("abcdefghijklmnopqrstuvwxyz", 10),
		"ASCII_Large":     strings.Repeat("abcdefghijklmnopqrstuvwxyz", 100),
		"Unicode_Small":   "가나다라마바사아자차카타파하",
		"Unicode_Medium":  strings.Repeat("가나다라마바사아자차카타파하", 10),
		"Unicode_Large":   strings.Repeat("가나다라마바사아자차카타파하", 100),
		"Mixed_Small":     "Hello 안녕 World 世界",
		"Mixed_Medium":    strings.Repeat("Hello 안녕 World 世界 ", 25),
	}

	for name, testData := range testCases {
		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(testData)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := strproc.ReverseStr(testData)
				_ = result
			}
		})
	}
}

func BenchmarkStringProc_ReverseNormalStr(b *testing.B) {
	strproc := strutils.NewStringProc()
	testCases := map[string]string{
		"Small":  "abcdefghijklmnopqrstuvwxyz",
		"Medium": strings.Repeat("abcdefghijklmnopqrstuvwxyz", 10),
		"Large":  strings.Repeat("abcdefghijklmnopqrstuvwxyz", 100),
	}

	for name, testData := range testCases {
		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(testData)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := strproc.ReverseNormalStr(testData)
				_ = result
			}
		})
	}
}

func BenchmarkStringProc_ReverseUnicode(b *testing.B) {
	strproc := strutils.NewStringProc()
	testCases := map[string]string{
		"Korean":   strings.Repeat("가나다라마바사", 50),
		"Japanese": strings.Repeat("あいうえおかきくけこ", 50),
		"Chinese":  strings.Repeat("天地玄黃宇宙洪荒", 50),
		"Mixed":    strings.Repeat("Hello안녕世界", 50),
	}

	for name, testData := range testCases {
		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(testData)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := strproc.ReverseUnicode(testData)
				_ = result
			}
		})
	}
}

func BenchmarkStringProc_NumberFmt(b *testing.B) {
	strproc := strutils.NewStringProc()
	
	testCases := []struct {
		name  string
		value interface{}
	}{
		{"Int_Small", 12345},
		{"Int_Large", 1234567890123},
		{"Float_Small", 123.45},
		{"Float_Large", 1234567890.123456},
		{"Float_Scientific", 1.23456789e+10},
		{"String_Int", "1234567890"},
		{"String_Float", "123456.789"},
		{"String_Scientific", "1.23e+10"},
		{"Negative_Int", -1234567890},
		{"Negative_Float", -123456.789},
	}

	for _, testCase := range testCases {
		b.Run(testCase.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result, err := strproc.NumberFmt(testCase.value)
				if err != nil {
					b.Fatalf("NumberFmt error: %v", err)
				}
				_ = result
			}
		})
	}
}

func BenchmarkStringProc_Padding(b *testing.B) {
	strproc := strutils.NewStringProc()
	testData := "test string for padding"
	
	paddingTests := []struct {
		name   string
		method func(string, string, int) string
	}{
		{"PaddingLeft", strproc.PaddingLeft},
		{"PaddingRight", strproc.PaddingRight},
		{"PaddingBoth", strproc.PaddingBoth},
	}

	for _, test := range paddingTests {
		b.Run(test.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := test.method(testData, "*", 50)
				_ = result
			}
		})
	}
}

func BenchmarkStringProc_WordWrap(b *testing.B) {
	strproc := strutils.NewStringProc()
	testData := "The quick brown fox jumps over the lazy dog. This is a test sentence for word wrapping functionality."
	
	b.Run("WordWrapSimple", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := strproc.WordWrapSimple(testData, 10, "\n")
			if err != nil {
				b.Fatalf("WordWrapSimple error: %v", err)
			}
			_ = result
		}
	})

	b.Run("WordWrapAround", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result, err := strproc.WordWrapAround(testData, 10, "\n")
			if err != nil {
				b.Fatalf("WordWrapAround error: %v", err)
			}
			_ = result
		}
	})
}

func BenchmarkStringProc_MD5Hash(b *testing.B) {
	strproc := strutils.NewStringProc()
	testCases := map[string]string{
		"Small":  "hello world",
		"Medium": strings.Repeat("test data for hashing ", 10),
		"Large":  strings.Repeat("large test data for MD5 hashing performance ", 100),
	}

	for name, testData := range testCases {
		b.Run(name, func(b *testing.B) {
			b.SetBytes(int64(len(testData)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result, err := strproc.MD5Hash(testData)
				if err != nil {
					b.Fatalf("MD5Hash error: %v", err)
				}
				_ = result
			}
		})
	}
}

func BenchmarkStringProc_HumanByteSize(b *testing.B) {
	strproc := strutils.NewStringProc()
	
	testCases := []struct {
		name  string
		bytes interface{}
	}{
		{"Bytes", 512},
		{"Kilobytes", 2048},
		{"Megabytes", 5242880},
		{"Gigabytes", 3221225472},
		{"Large", int64(math.MaxInt32)},
	}

	for _, testCase := range testCases {
		b.Run(testCase.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result, err := strproc.HumanByteSize(testCase.bytes, 2, strutils.CamelCaseLong)
				if err != nil {
					b.Fatalf("HumanByteSize error: %v", err)
				}
				_ = result
			}
		})
	}
}

// StringValidator 함수들의 종합 벤치마크
func BenchmarkStringValidator_IsValidEmail(b *testing.B) {
	strvalidator := strutils.NewStringValidator()
	
	testCases := map[string]string{
		"Simple":      "test@example.com",
		"Complex":     "user.name+tag@example-domain.co.uk",
		"Long":        "very.long.email.address.with.many.parts@very-long-domain-name.example.org",
		"Invalid_Simple": "@invalid.com",
		"Invalid_Complex": "user@@domain..com",
	}

	for name, email := range testCases {
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := strvalidator.IsValidEmail(email)
				_ = result
			}
		})
	}
}

func BenchmarkStringValidator_IsValidDomain(b *testing.B) {
	strvalidator := strutils.NewStringValidator()
	
	testCases := map[string]string{
		"Simple":         "example.com",
		"Subdomain":      "subdomain.example.com",
		"Long":           "very-long-domain-name-with-many-hyphens.example.org",
		"International":  "xn--example-domain.com",
		"Invalid":        "invalid..domain.com",
	}

	for name, domain := range testCases {
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := strvalidator.IsValidDomain(domain)
				_ = result
			}
		})
	}
}

func BenchmarkStringValidator_IsValidURL(b *testing.B) {
	strvalidator := strutils.NewStringValidator()
	
	testCases := map[string]string{
		"HTTP_Simple":  "http://example.com",
		"HTTPS_Path":   "https://example.com/path/to/page",
		"Complex":      "https://user:pass@subdomain.example.com:8080/path?query=value&other=test#fragment",
		"Long":         "https://very-long-domain.example.org/very/long/path/with/many/segments?param1=value1&param2=value2&param3=value3#section",
		"Invalid":      "not-a-url",
	}

	for name, url := range testCases {
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := strvalidator.IsValidURL(url)
				_ = result
			}
		})
	}
}

func BenchmarkStringValidator_IsValidIPAddr(b *testing.B) {
	strvalidator := strutils.NewStringValidator()
	
	testCases := map[string]struct {
		ip     string
		ipType int
	}{
		"IPv4_Simple":       {"192.168.1.1", strutils.IPv4},
		"IPv4_Loopback":     {"127.0.0.1", strutils.IPv4},
		"IPv4_Broadcast":    {"255.255.255.255", strutils.IPv4},
		"IPv6_Simple":       {"2001:db8::1", strutils.IPv6},
		"IPv6_Loopback":     {"::1", strutils.IPv6},
		"IPv6_Full":         {"2001:0db8:85a3:0000:0000:8a2e:0370:7334", strutils.IPv6},
		"IPv4_CIDR":         {"192.168.1.0/24", strutils.IPv4CIDR},
		"IPv6_CIDR":         {"2001:db8::/32", strutils.IPv6CIDR},
		"IPv4MappedIPv6":    {"2001:470:1f09:495::3:217.126.185.21", strutils.IPv4MappedIPv6},
		"Invalid_IPv4":      {"256.256.256.256", strutils.IPv4},
		"Invalid_IPv6":      {"2001:db8::1::2", strutils.IPv6},
	}

	for name, testCase := range testCases {
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result, err := strvalidator.IsValidIPAddr(testCase.ip, testCase.ipType)
				if err != nil {
					// 에러가 있어도 벤치마크는 계속
				}
				_ = result
			}
		})
	}
}

func BenchmarkStringValidator_IsValidMACAddr(b *testing.B) {
	strvalidator := strutils.NewStringValidator()
	
	testCases := map[string]string{
		"Colon":     "00:11:22:33:44:55",
		"Dash":      "00-11-22-33-44-55",
		"Uppercase": "AA:BB:CC:DD:EE:FF",
		"Mixed":     "aA:bB:cC:dD:eE:fF",
		"Invalid":   "invalid-mac-address",
	}

	for name, mac := range testCases {
		b.Run(name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := strvalidator.IsValidMACAddr(mac)
				_ = result
			}
		})
	}
}

func BenchmarkStringValidator_IsPureText(b *testing.B) {
	strvalidator := strutils.NewStringValidator()
	
	testCases := map[string]string{
		"Plain":           "This is just plain text",
		"WithNumbers":     "Text with numbers 123 and symbols !@#",
		"LongText":        strings.Repeat("This is a long text for testing performance. ", 20),
		"HTML":            "<script>alert('test')</script>",
		"Mixed":           "Normal text with <b>some</b> HTML tags",
		"Unicode":         "Unicode text: 안녕하세요 こんにちは 你好",
		"ControlChars":    "Text with control\x00chars\x01",
	}

	b.Run("PureTextNormal", func(b *testing.B) {
		for name, text := range testCases {
			b.Run(name, func(b *testing.B) {
				b.SetBytes(int64(len(text)))
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					result, err := strvalidator.IsPureTextNormal(text)
					if err != nil {
						// 에러가 있어도 벤치마크는 계속
					}
					_ = result
				}
			})
		}
	})

	b.Run("PureTextStrict", func(b *testing.B) {
		for name, text := range testCases {
			b.Run(name, func(b *testing.B) {
				b.SetBytes(int64(len(text)))
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					result, err := strvalidator.IsPureTextStrict(text)
					if err != nil {
						// 에러가 있어도 벤치마크는 계속
					}
					_ = result
				}
			})
		}
	})
}

func BenchmarkStringValidator_FilePath(b *testing.B) {
	strvalidator := strutils.NewStringValidator()
	
	testCases := map[string]string{
		"Simple":         "file.txt",
		"WithPath":       "/path/to/file.txt",
		"Windows":        "C:\\Path\\To\\File.txt",
		"Relative":       "./relative/path/file.txt",
		"Long":           strings.Repeat("dir/", 10) + "very_long_filename.extension",
		"Special":        "file-name_with.special.chars.txt",
		"Invalid":        "invalid<>chars.txt",
	}

	b.Run("IsValidFilePath", func(b *testing.B) {
		for name, path := range testCases {
			b.Run(name, func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					result := strvalidator.IsValidFilePath(path)
					_ = result
				}
			})
		}
	})

	b.Run("IsValidFilePathWithRelativePath", func(b *testing.B) {
		for name, path := range testCases {
			b.Run(name, func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					result := strvalidator.IsValidFilePathWithRelativePath(path)
					_ = result
				}
			})
		}
	})
}

// 복합 작업 벤치마크
func BenchmarkStringProc_CompositeOperations(b *testing.B) {
	strproc := strutils.NewStringProc()
	testData := "The quick brown fox jumps over the lazy dog"

	operations := []struct {
		name string
		fn   func(string) string
	}{
		{"AddSlashes_UpperCase", func(s string) string {
			return strproc.UpperCaseFirstWords(strproc.AddSlashes(s))
		}},
		{"Reverse_PadBoth", func(s string) string {
			return strproc.PaddingBoth(strproc.ReverseStr(s), "*", 60)
		}},
		{"Chain_Multiple", func(s string) string {
			s = strproc.AddSlashes(s)
			s = strproc.UpperCaseFirstWords(s)
			s = strproc.PaddingLeft(s, " ", len(s)+10)
			return s
		}},
	}

	for _, op := range operations {
		b.Run(op.name, func(b *testing.B) {
			b.SetBytes(int64(len(testData)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result := op.fn(testData)
				_ = result
			}
		})
	}
}

func BenchmarkStringValidator_CompositeValidations(b *testing.B) {
	strvalidator := strutils.NewStringValidator()
	
	b.Run("MultipleValidations", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = strvalidator.IsValidEmail("test@example.com")
			_ = strvalidator.IsValidDomain("example.com")
			_ = strvalidator.IsValidURL("https://example.com")
			ipResult, _ := strvalidator.IsValidIPAddr("192.168.1.1", strutils.IPv4)
			_ = ipResult
		}
	})
}

// 메모리 할당 벤치마크
func BenchmarkStringProc_MemoryAllocation(b *testing.B) {
	strproc := strutils.NewStringProc()
	testData := strings.Repeat("test data ", 1000) // 큰 문자열

	b.Run("AddSlashes_LargeString", func(b *testing.B) {
		b.SetBytes(int64(len(testData)))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := strproc.AddSlashes(testData)
			_ = result
		}
	})

	b.Run("ReverseStr_LargeString", func(b *testing.B) {
		b.SetBytes(int64(len(testData)))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := strproc.ReverseStr(testData)
			_ = result
		}
	})

	b.Run("UpperCaseFirstWords_LargeString", func(b *testing.B) {
		b.SetBytes(int64(len(testData)))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := strproc.UpperCaseFirstWords(testData)
			_ = result
		}
	})
}

// 병렬 처리 벤치마크
func BenchmarkStringProc_Parallel(b *testing.B) {
	strproc := strutils.NewStringProc()
	testData := "test data for parallel processing"

	b.Run("AddSlashes_Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				result := strproc.AddSlashes(testData)
				_ = result
			}
		})
	})

	b.Run("UpperCaseFirstWords_Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				result := strproc.UpperCaseFirstWords(testData)
				_ = result
			}
		})
	})

	b.Run("NumberFmt_Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				result, _ := strproc.NumberFmt(123456)
				_ = result
			}
		})
	})
}

func BenchmarkStringValidator_Parallel(b *testing.B) {
	strvalidator := strutils.NewStringValidator()

	b.Run("IsValidEmail_Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				result := strvalidator.IsValidEmail("test@example.com")
				_ = result
			}
		})
	})

	b.Run("IsValidURL_Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				result := strvalidator.IsValidURL("https://example.com")
				_ = result
			}
		})
	})

	b.Run("IsValidIPAddr_Parallel", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				result, _ := strvalidator.IsValidIPAddr("192.168.1.1", strutils.IPv4)
				_ = result
			}
		})
	})
}

// 특수 케이스 벤치마크
func BenchmarkStringProc_EdgeCases(b *testing.B) {
	strproc := strutils.NewStringProc()

	b.Run("EmptyString", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := strproc.AddSlashes("")
			_ = result
		}
	})

	b.Run("SingleChar", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := strproc.UpperCaseFirstWords("a")
			_ = result
		}
	})

	b.Run("VeryLongString", func(b *testing.B) {
		longString := strings.Repeat("abcdefghijklmnopqrstuvwxyz", 1000)
		b.SetBytes(int64(len(longString)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := strproc.AddSlashes(longString)
			_ = result
		}
	})

	b.Run("UnicodeString", func(b *testing.B) {
		unicodeString := strings.Repeat("안녕하세요世界こんにちは", 100)
		b.SetBytes(int64(len(unicodeString)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := strproc.ReverseStr(unicodeString)
			_ = result
		}
	})
}

// 함수 호출 오버헤드 벤치마크
func BenchmarkFunction_CallOverhead(b *testing.B) {
	strproc := strutils.NewStringProc()
	strvalidator := strutils.NewStringValidator()

	b.Run("StringProc_Creation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			proc := strutils.NewStringProc()
			_ = proc
		}
	})

	b.Run("StringValidator_Creation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			validator := strutils.NewStringValidator()
			_ = validator
		}
	})

	b.Run("Method_Call_Overhead", func(b *testing.B) {
		testData := "test"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// 매우 간단한 작업으로 순수 메소드 호출 오버헤드 측정
			result := strproc.AddSlashes(testData)
			_ = result
		}
	})

	b.Run("Validation_Call_Overhead", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := strvalidator.IsValidEmail("test@example.com")
			_ = result
		}
	})
}