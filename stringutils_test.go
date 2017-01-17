package strutils

import (
	"strings"
	"testing"

	"github.com/dustin/go-humanize"
)

func TestAddSlashes(t *testing.T) {

	strutil := NewStringUtils()
	dataset := map[string]string{
		`대한민국만세`:     `대한민국만세`,
		`대한\민국만세`:    `대한\\민국만세`,
		`대한\\민국만세`:   `대한\\민국만세`,
		"abcdefgz":   "abcdefgz",
		`a\bcdefgz`:  `a\\bcdefgz`,
		`a\\bcdefgz`: `a\\bcdefgz`,
	}

	var retval string
	for k, v := range dataset {
		retval = strutil.AddSlashes(k)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
	}
}

func TestStripSlashes(t *testing.T) {

	strutil := NewStringUtils()
	dataset := map[string]string{
		`대한민국만세`:       `대한민국만세`,
		`대한\\민국만세`:     `대한\민국만세`,
		`대한\\\\민국만세`:   `대한\\민국만세`,
		"abcdefgz":     "abcdefgz",
		`a\\bcdefgz`:   `a\bcdefgz`,
		`a\\\\bcdefgz`: `a\\bcdefgz`,
	}

	var retval string
	for k, v := range dataset {
		retval = strutil.StripSlashes(k)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
	}
}

func TestNl2Br(t *testing.T) {

	strutil := NewStringUtils()
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

	var retval string
	for k, v := range dataset {
		retval = strutil.Nl2Br(k)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
	}
}

func BenchmarkNl2Br(b *testing.B) {

	strutil := NewStringUtils()
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

	for i := 0; i < b.N; i++ {
		var retval string
		for k, v := range dataset {
			retval = strutil.Nl2Br(k)
			if v != retval {
				b.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
			}
		}
	}
}

func BenchmarkNl2BrUseStringReplace(b *testing.B) {

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

func TestNumbertFmt(t *testing.T) {

	strutil := NewStringUtils()
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

	for k, v := range dataset {
		retval, err := strutil.NumberFmt(k)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
		if err != nil {
			t.Errorf("Return Error : %v", err)
		}
	}
}

func BenchmarkTestNumbertFmt(b *testing.B) {

	strutil := NewStringUtils()
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

	for i := 0; i < b.N; i++ {
		for k, v := range dataset {
			retval, err := strutil.NumberFmt(k)
			if v != retval {
				b.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
			}
			if err != nil {
				b.Errorf("Return Error : %v", err)
			}
		}
	}

}

//BenchmarkTestNumbertFmtInt64-8                	 2000000	       712 ns/op
//BenchmarkTestNumbertFmtInt64UseHumanUnits-8   	 2000000	       761 ns/op
func BenchmarkTestNumbertFmtInt64(b *testing.B) {

	strutil := NewStringUtils()
	dataset := map[interface{}]string{
		123456789101112: "123,456,789,101,112",
	}

	for i := 0; i < b.N; i++ {
		for k, v := range dataset {
			retval, err := strutil.NumberFmt(k)
			if v != retval {
				b.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
			}
			if err != nil {
				b.Errorf("Return Error : %v", err)
			}
		}
	}

}

func BenchmarkTestNumbertFmtInt64UseHumanUnits(b *testing.B) {

	dataset := map[int64]string{
		123456789101112: "123,456,789,101,112",
	}

	for i := 0; i < b.N; i++ {
		for k, v := range dataset {
			retval := humanize.Comma(k)
			if v != retval {
				b.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
			}
		}
	}
}

type paddingTestVal struct {
	str   string
	fill  string
	m     int
	mx    int
	okstr string
}

func TestPadding(t *testing.T) {

	strutil := NewStringUtils()
	dataset := make(map[int]paddingTestVal)

	dataset[0] = paddingTestVal{"Life isn't always what one like.", "*", PAD_BOTH, 38, "***Life isn't always what one like.***"}
	dataset[1] = paddingTestVal{"Life isn't always what one like.", "*", PAD_LEFT, 38, "******Life isn't always what one like."}
	dataset[2] = paddingTestVal{"Life isn't always what one like.", "*", PAD_RIGHT, 38, "Life isn't always what one like.******"}
	dataset[3] = paddingTestVal{"Life isn't always what one like.", "*-=", PAD_BOTH, 37, "*-Life isn't always what one like.*-="}
	dataset[4] = paddingTestVal{"Life isn't always what one like.", "*-=", PAD_LEFT, 37, "*-=*-Life isn't always what one like."}
	dataset[5] = paddingTestVal{"Life isn't always what one like.", "*-=", PAD_RIGHT, 37, "Life isn't always what one like.*-=*-"}

	dataset[6] = paddingTestVal{"가나다라마바사아자차카타파하", "*", PAD_BOTH, 48, "***가나다라마바사아자차카타파하***"}
	dataset[7] = paddingTestVal{"가나다라마바사아자차카타파하", "*", PAD_LEFT, 48, "******가나다라마바사아자차카타파하"}
	dataset[8] = paddingTestVal{"가나다라마바사아자차카타파하", "*", PAD_RIGHT, 48, "가나다라마바사아자차카타파하******"}
	dataset[9] = paddingTestVal{"가나다라마바사아자차카타파하", "*-=", PAD_BOTH, 47, "*-가나다라마바사아자차카타파하*-="}
	dataset[10] = paddingTestVal{"가나다라마바사아자차카타파하", "*-=", PAD_LEFT, 47, "*-=*-가나다라마바사아자차카타파하"}
	dataset[11] = paddingTestVal{"가나다라마바사아자차카타파하", "*-=", PAD_RIGHT, 47, "가나다라마바사아자차카타파하*-=*-"}

	for _, v := range dataset {

		retval := strutil.padding(v.str, v.fill, v.m, v.mx)
		if v.okstr != retval {
			t.Errorf("INPUT TEXT : %v\n", v.str)
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v.okstr)
		}
	}
}
