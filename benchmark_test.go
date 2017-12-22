package strutils_test

import (
	"strings"
	"testing"

	humanize "github.com/dustin/go-humanize"
	"github.com/torden/go-strutil"
)

/*
BenchmarkReverseStrSwap-8                     	 1000000	      1137 ns/op
BenchmarkReverseStrUseReverseLoop-8           	 1000000	      1608 ns/op
*/

func Benchmark_strutils_ReverseStr(b *testing.B) {

	benchMarkStr1 := strings.Repeat("0123456789", 100)

	strproc := strutils.NewStringProc()
	for i := 0; i < b.N; i++ {
		strproc.ReverseStr(benchMarkStr1)
	}
}

func Benchmark_strutils_ReverseNormalStr(b *testing.B) {

	benchMarkStr1 := strings.Repeat("0123456789", 100)

	strproc := strutils.NewStringProc()
	for i := 0; i < b.N; i++ {
		strproc.ReverseNormalStr(benchMarkStr1)
	}
}

func Benchmark_strutils_ReverseReverseUnicode(b *testing.B) {

	benchMarkStr1 := strings.Repeat("0123456789", 100)

	strproc := strutils.NewStringProc()
	for i := 0; i < b.N; i++ {
		strproc.ReverseUnicode(benchMarkStr1)
	}
}

func Benchmark_strutils_ReverseStrSwap(b *testing.B) {

	benchMarkStr1 := strings.Repeat("0123456789", 100)

	strproc := strutils.NewStringProc()
	for i := 0; i < b.N; i++ {
		strproc.ReverseNormalStr(benchMarkStr1)
	}
}

func Benchmark_strutils_Nl2Br(b *testing.B) {

	strproc := strutils.NewStringProc()
	dataset := map[string]string{
		"대한\n민국만세":     "대한<br />민국만세",
		"대한\r\n민국만세":   "대한<br />민국만세",
		"대한민국만세\r\n":   "대한민국만세<br />",
		"대한민국만세\n\r":   "대한민국만세<br />",
		"대한민국만세\n":     "대한민국만세<br />",
		"abcdefgh":     "abcdefgh",
		"abc\ndefgh":   "abc<br />defgh",
		"abcde\r\nfgh": "abcde<br />fgh",
		"abcdefgh\r\n": "abcdefgh<br />",
		"abcdefgh\n\r": "abcdefgh<br />",
	}

	//check : common
	for i := 0; i < b.N; i++ {
		var retval string
		for k, v := range dataset {
			retval = strproc.Nl2Br(k)
			if v != retval {
				b.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
			}
		}
	}
}

func Benchmark_strutils_Nl2BrUseStringReplace(b *testing.B) {

	dataset := map[string]string{
		"대한\n민국만세":     "대한<br />민국만세",
		"대한\r\n민국만세":   "대한<br />민국만세",
		"대한민국만세\r\n":   "대한민국만세<br />",
		"대한민국만세\n\r":   "대한민국만세<br />",
		"대한민국만세\n":     "대한민국만세<br />",
		"abcdefgh":     "abcdefgh",
		"abc\ndefgh":   "abc<br />defgh",
		"abcde\r\nfgh": "abcde<br />fgh",
		"abcdefgh\r\n": "abcdefgh<br />",
		"abcdefgh\n\r": "abcdefgh<br />",
	}

	//check : common
	for i := 0; i < b.N; i++ {
		var retval string
		for k, v := range dataset {
			retval = strings.Replace(k, "\r\n", "<br />", -1)
			retval = strings.Replace(retval, "\n\r", "<br />", -1)
			retval = strings.Replace(retval, "\n", "<br />", -1)
			if v != retval {
				b.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
			}
		}
	}
}
func Benchmark_strutils_TestNumbertFmt(b *testing.B) {

	strproc := strutils.NewStringProc()
	dataset := map[interface{}]string{
		123456789101112: "123,456,789,101,112",
		123456.1234:     "123,456.1234",
		-123456.1234:    "-123,456.1234",
		1.1234561e+06:   "1.1234561e+06",
		1234.1234:       "1,234.1234",
		12345.1234:      "12,345.1234",
		-1.1234561e+06:  "-1.1234561e+06",
		-12345.16:       "-12,345.16",
		12345.16:        "12,345.16",
		1234:            "1,234",
		12.12123098123:  "12.12123098123",
		1.212e+24:       "1.212e+24",
		123456789:       "123,456,789",
	}

	//benchmark : common
	for i := 0; i < b.N; i++ {
		for k, v := range dataset {
			retval, err := strproc.NumberFmt(k)
			if v != retval {
				b.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
			}
			if err != nil {
				b.Errorf("Return Error : %v", err)
			}
		}
	}
}

func Benchmark_strutils_TestNumbertFmtInt64(b *testing.B) {
	//BenchmarkTestNumbertFmtInt64-8                	 2000000	       712 ns/op
	//BenchmarkTestNumbertFmtInt64UseHumanUnits-8   	 2000000	       761 ns/op

	strproc := strutils.NewStringProc()
	dataset := map[interface{}]string{
		123456789101112: "123,456,789,101,112",
	}

	//benchmark : common
	for i := 0; i < b.N; i++ {
		for k, v := range dataset {
			retval, err := strproc.NumberFmt(k)
			if v != retval {
				b.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
			}
			if err != nil {
				b.Errorf("Return Error : %v", err)
			}
		}
	}

}

func Benchmark_strutils_TestNumbertFmtInt64UseHumanUnits(b *testing.B) {

	dataset := map[int64]string{
		123456789101112: "123,456,789,101,112",
	}

	//benchmark : common
	for i := 0; i < b.N; i++ {
		for k, v := range dataset {
			retval := humanize.Comma(k)
			if v != retval {
				b.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
			}
		}
	}
}
