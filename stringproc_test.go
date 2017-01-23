// testing : string processing
package strutils_test

import (
	"io/ioutil"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/dustin/go-humanize"

	"github.com/torden/go-strutil"
)

func TestAddSlashes(t *testing.T) {

	strproc := strutils.NewStringProc()
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
		retval = strproc.AddSlashes(k)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
	}
}

func TestStripSlashes(t *testing.T) {

	strproc := strutils.NewStringProc()
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
		retval = strproc.StripSlashes(k)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
	}
}

func TestNl2Br(t *testing.T) {

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

	var retval string
	for k, v := range dataset {
		retval = strproc.Nl2Br(k)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
	}
}

func BenchmarkNl2Br(b *testing.B) {

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

type wordwrapTestVal struct {
	str      string
	wd       int
	breakstr string
	okstr    string
}

func TestWordWrapSimple(t *testing.T) {
	strproc := strutils.NewStringProc()

	dataset := make(map[int]wordwrapTestVal)

	dataset[1] = wordwrapTestVal{`The quick brown fox jumped over the lazy dog.`, 3, `*`, `The*quick*brown*fox*jumped*over*the*lazy*dog.`}
	dataset[2] = wordwrapTestVal{`A very long woooooooooooord.`, 3, `*`, `A very*long*woooooooooooord.`}
	dataset[3] = wordwrapTestVal{`A very long woooooooooooooooooord. and something`, 3, `*`, `A very*long*woooooooooooooooooord.*and*something`}
	dataset[4] = wordwrapTestVal{`가 나 다 라 마 바 사 아 자 차 카 타 파 하`, 3, `*`, `가*나*다*라*마*바*사*아*자*차*카*타*파*하`}

	dataset[5] = wordwrapTestVal{`The quick brown fox jumped over the lazy dog.`, 5, `-`, `The quick-brown-fox jumped-over the-lazy dog.`}
	dataset[6] = wordwrapTestVal{`A very long woooooooooooord.`, 5, `-`, `A very-long woooooooooooord.`}
	dataset[7] = wordwrapTestVal{`A very long woooooooooooooooooord. and something`, 5, `-`, `A very-long woooooooooooooooooord.-and something`}
	dataset[8] = wordwrapTestVal{`가 나 다 라 마 바 사 아 자 차 카 타 파 하`, 5, `-`, `가 나-다 라-마 바-사 아-자 차-카 타-파 하`}

	dataset[9] = wordwrapTestVal{`The quick brown fox jumped over the lazy dog.`, 8, `+`, `The quick+brown fox+jumped over+the lazy+dog.`}
	dataset[10] = wordwrapTestVal{`A very long woooooooooooord.`, 8, `+`, `A very long+woooooooooooord.`}
	dataset[11] = wordwrapTestVal{`A very long woooooooooooooooooord. and something`, 8, `+`, `A very long+woooooooooooooooooord.+and something`}
	dataset[12] = wordwrapTestVal{`가 나 다 라 마 바 사 아 자 차 카 타 파 하`, 8, `+`, `가 나 다+라 마 바+사 아 자+차 카 타+파 하`}

	dataset[13] = wordwrapTestVal{`The quick brown fox jumped over the lazy dog.`, 3, `!@#$%`, `The!@#$%quick!@#$%brown!@#$%fox!@#$%jumped!@#$%over!@#$%the!@#$%lazy!@#$%dog.`}
	dataset[14] = wordwrapTestVal{`A very long woooooooooooord.`, 3, `!@#$%`, `A very!@#$%long!@#$%woooooooooooord.`}
	dataset[15] = wordwrapTestVal{`A very long woooooooooooooooooord. and something`, 3, `!@#$%`, `A very!@#$%long!@#$%woooooooooooooooooord.!@#$%and!@#$%something`}
	dataset[16] = wordwrapTestVal{`가 나 다 라 마 바 사 아 자 차 카 타 파 하`, 3, `!@#$%`, `가!@#$%나!@#$%다!@#$%라!@#$%마!@#$%바!@#$%사!@#$%아!@#$%자!@#$%차!@#$%카!@#$%타!@#$%파!@#$%하`}

	dataset[17] = wordwrapTestVal{`The quick brown fox jumped over the lazy dog.`, 5, `*-=*-=`, `The quick*-=*-=brown*-=*-=fox jumped*-=*-=over the*-=*-=lazy dog.`}
	dataset[18] = wordwrapTestVal{`A very long woooooooooooord.`, 5, `*-=*-=`, `A very*-=*-=long woooooooooooord.`}
	dataset[19] = wordwrapTestVal{`A very long woooooooooooooooooord. and something`, 5, `*-=*-=`, `A very*-=*-=long woooooooooooooooooord.*-=*-=and something`}
	dataset[20] = wordwrapTestVal{`가 나 다 라 마 바 사 아 자 차 카 타 파 하`, 5, `*-=*-=`, `가 나*-=*-=다 라*-=*-=마 바*-=*-=사 아*-=*-=자 차*-=*-=카 타*-=*-=파 하`}

	dataset[21] = wordwrapTestVal{`The quick brown fox jumped over the lazy dog.`, 8, `_+_+_+_+_+_+_+_+_+_+_+_+`, `The quick_+_+_+_+_+_+_+_+_+_+_+_+brown fox_+_+_+_+_+_+_+_+_+_+_+_+jumped over_+_+_+_+_+_+_+_+_+_+_+_+the lazy_+_+_+_+_+_+_+_+_+_+_+_+dog.`}
	dataset[22] = wordwrapTestVal{`A very long woooooooooooord.`, 8, `_+_+_+_+_+_+_+_+_+_+_+_+`, `A very long_+_+_+_+_+_+_+_+_+_+_+_+woooooooooooord.`}
	dataset[23] = wordwrapTestVal{`A very long woooooooooooooooooord. and something`, 8, `_+_+_+_+_+_+_+_+_+_+_+_+`, `A very long_+_+_+_+_+_+_+_+_+_+_+_+woooooooooooooooooord._+_+_+_+_+_+_+_+_+_+_+_+and something`}
	dataset[24] = wordwrapTestVal{`가 나 다 라 마 바 사 아 자 차 카 타 파 하`, 8, `_+_+_+_+_+_+_+_+_+_+_+_+`, `가 나 다_+_+_+_+_+_+_+_+_+_+_+_+라 마 바_+_+_+_+_+_+_+_+_+_+_+_+사 아 자_+_+_+_+_+_+_+_+_+_+_+_+차 카 타_+_+_+_+_+_+_+_+_+_+_+_+파 하`}

	for _, v := range dataset {

		retval, _ := strproc.WordWrapSimple(v.str, v.wd, v.breakstr)
		if v.okstr != retval {
			t.Errorf("Original Value : %v\n", v.str)
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v.okstr)
		}
	}

	//check : wd = 0
	_, err := strproc.WordWrapSimple("test", 0, "1234")
	if err == nil {
		t.Errorf("Failure : Couldn't check the `wd at least 1`")
	}
}

func TestWordWrapAround(t *testing.T) {
	strproc := strutils.NewStringProc()

	dataset := make(map[int]wordwrapTestVal)

	dataset[1] = wordwrapTestVal{`The quick brown fox jumped over the lazy dog.`, 3, `*`, `The*quick*brown*fox*jumped*over*the*lazy*dog.`}
	dataset[2] = wordwrapTestVal{`A very long woooooooooooord.`, 3, `*`, `A very*long*woooooooooooord.`}
	dataset[3] = wordwrapTestVal{`A very long woooooooooooooooooord. and something`, 3, `*`, `A very*long*woooooooooooooooooord.*and*something`}
	dataset[4] = wordwrapTestVal{`가 나 다 라 마 바 사 아 자 차 카 타 파 하`, 3, `*`, `가*나*다*라*마*바*사*아*자*차*카*타*파*하`}

	dataset[5] = wordwrapTestVal{`The quick brown fox jumped over the lazy dog.`, 5, `-`, `The quick-brown-fox-jumped-over-the-lazy-dog.`}
	dataset[6] = wordwrapTestVal{`A very long woooooooooooord.`, 5, `-`, `A very-long-woooooooooooord.`}
	dataset[7] = wordwrapTestVal{`A very long woooooooooooooooooord. and something`, 5, `-`, `A very-long-woooooooooooooooooord.-and-something`}
	dataset[8] = wordwrapTestVal{`가 나 다 라 마 바 사 아 자 차 카 타 파 하`, 5, `-`, `가 나-다-라-마 바-사-아-자-차 카-타-파-하`}

	dataset[9] = wordwrapTestVal{`The quick brown fox jumped over the lazy dog.`, 8, `+`, `The quick+brown fox+jumped+over the+lazy+dog.`}
	dataset[10] = wordwrapTestVal{`A very long woooooooooooord.`, 8, `+`, `A very long+woooooooooooord.`}
	dataset[11] = wordwrapTestVal{`A very long woooooooooooooooooord. and something`, 8, `+`, `A very long+woooooooooooooooooord.+and+something`}
	dataset[12] = wordwrapTestVal{`가 나 다 라 마 바 사 아 자 차 카 타 파 하`, 8, `+`, `가 나 다+라 마+바 사+아 자+차 카+타 파+하`}

	dataset[13] = wordwrapTestVal{`The quick brown fox jumped over the lazy dog.`, 3, `!@#$%`, `The!@#$%quick!@#$%brown!@#$%fox!@#$%jumped!@#$%over!@#$%the!@#$%lazy!@#$%dog.`}
	dataset[14] = wordwrapTestVal{`A very long woooooooooooord.`, 3, `!@#$%`, `A very!@#$%long!@#$%woooooooooooord.`}
	dataset[15] = wordwrapTestVal{`A very long woooooooooooooooooord. and something`, 3, `!@#$%`, `A very!@#$%long!@#$%woooooooooooooooooord.!@#$%and!@#$%something`}
	dataset[16] = wordwrapTestVal{`가 나 다 라 마 바 사 아 자 차 카 타 파 하`, 3, `!@#$%`, `가!@#$%나!@#$%다!@#$%라!@#$%마!@#$%바!@#$%사!@#$%아!@#$%자!@#$%차!@#$%카!@#$%타!@#$%파!@#$%하`}

	dataset[17] = wordwrapTestVal{`The quick brown fox jumped over the lazy dog.`, 5, `*-=*-=`, `The quick*-=*-=brown*-=*-=fox*-=*-=jumped*-=*-=over*-=*-=the*-=*-=lazy*-=*-=dog.`}
	dataset[18] = wordwrapTestVal{`A very long woooooooooooord.`, 5, `*-=*-=`, `A very*-=*-=long*-=*-=woooooooooooord.`}
	dataset[19] = wordwrapTestVal{`A very long woooooooooooooooooord. and something`, 5, `*-=*-=`, `A very*-=*-=long*-=*-=woooooooooooooooooord.*-=*-=and*-=*-=something`}
	dataset[20] = wordwrapTestVal{`가 나 다 라 마 바 사 아 자 차 카 타 파 하`, 5, `*-=*-=`, `가 나*-=*-=다*-=*-=라*-=*-=마 바*-=*-=사*-=*-=아*-=*-=자*-=*-=차 카*-=*-=타*-=*-=파*-=*-=하`}

	dataset[21] = wordwrapTestVal{`The quick brown fox jumped over the lazy dog.`, 8, `_+_+_+_+_+_+_+_+_+_+_+_+`, `The quick_+_+_+_+_+_+_+_+_+_+_+_+brown fox_+_+_+_+_+_+_+_+_+_+_+_+jumped_+_+_+_+_+_+_+_+_+_+_+_+over the_+_+_+_+_+_+_+_+_+_+_+_+lazy_+_+_+_+_+_+_+_+_+_+_+_+dog.`}
	dataset[22] = wordwrapTestVal{`A very long woooooooooooord.`, 8, `_+_+_+_+_+_+_+_+_+_+_+_+`, `A very long_+_+_+_+_+_+_+_+_+_+_+_+woooooooooooord.`}
	dataset[23] = wordwrapTestVal{`A very long woooooooooooooooooord. and something`, 8, `_+_+_+_+_+_+_+_+_+_+_+_+`, `A very long_+_+_+_+_+_+_+_+_+_+_+_+woooooooooooooooooord._+_+_+_+_+_+_+_+_+_+_+_+and_+_+_+_+_+_+_+_+_+_+_+_+something`}
	dataset[24] = wordwrapTestVal{`가 나 다 라 마 바 사 아 자 차 카 타 파 하`, 8, `_+_+_+_+_+_+_+_+_+_+_+_+`, `가 나 다_+_+_+_+_+_+_+_+_+_+_+_+라 마_+_+_+_+_+_+_+_+_+_+_+_+바 사_+_+_+_+_+_+_+_+_+_+_+_+아 자_+_+_+_+_+_+_+_+_+_+_+_+차 카_+_+_+_+_+_+_+_+_+_+_+_+타 파_+_+_+_+_+_+_+_+_+_+_+_+하`}

	for _, v := range dataset {

		retval, _ := strproc.WordWrapAround(v.str, v.wd, v.breakstr)
		if v.okstr != retval {
			t.Errorf("Original Value : %v\n", v.str)
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v.okstr)
		}
	}

	var err error

	//check : wd = 0
	_, err = strproc.WordWrapAround("test", 0, "1234")
	if err == nil {
		t.Errorf("Failure : Couldn't check the `wd at least 1`")
	}

	//check : lastspc = 1
	_, err = strproc.WordWrapAround("ttttttt tttttttttt", 2, "1111")
	if err != nil {
		t.Errorf("Failure : Couldn't check the `lastspc = 1`")
	}

}

func TestNumbertFmt(t *testing.T) {

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

		int(math.MaxInt8):  "127",
		uint(math.MaxInt8): "127",

		int8(math.MaxInt8):       "127",
		int16(math.MaxInt16):     "32,767",
		int16(math.MinInt16):     "-32,768",
		int32(math.MaxInt32):     "2,147,483,647",
		int32(math.MinInt32):     "-2,147,483,648",
		int64(math.MaxInt64):     "9,223,372,036,854,775,807",
		int64(math.MinInt64):     "-9,223,372,036,854,775,808",
		uint8(math.MaxUint8):     "255",
		uint16(math.MaxUint16):   "65,535",
		uint32(math.MaxUint32):   "4,294,967,295",
		uint64(math.MaxUint64):   "18,446,744,073,709,551,615",
		float32(math.MaxFloat32): "3.4028235e+38",
		float64(math.MaxFloat64): "1.7976931348623157e+308",
		//BUG(r) :
		//int8(math.MinInt8):       "-128",
		//float32(math.SmallestNonzeroFloat32): "1e-45",
		//float64(math.SmallestNonzeroFloat64): "5e-324",

	}

	for k, v := range dataset {
		retval, err := strproc.NumberFmt(k)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
		if err != nil {
			t.Errorf("Return Error : %v", err)
		}
	}

	var err error

	//check : ParseFloat
	_, err = strproc.NumberFmt("12.11111111111111111111111111111111111111111111111111111111111e12e12e1p029ekj12e")
	if err == nil {
		t.Errorf("Failure : Couldn't check the `Not Support strconv.ParseFloat`")
	}

	//check : not support obj
	_, err = strproc.NumberFmt(complex128(123))
	if err == nil {
		t.Errorf("Failure : Couldn't check the `not support obj`")
	}

	//check : not support numric string
	_, err = strproc.NumberFmt("1234===121212")
	if err == nil {
		t.Errorf("Failure : Couldn't check the `not support obj`")
	}

}

func BenchmarkTestNumbertFmt(b *testing.B) {

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

//BenchmarkTestNumbertFmtInt64-8                	 2000000	       712 ns/op
//BenchmarkTestNumbertFmtInt64UseHumanUnits-8   	 2000000	       761 ns/op
func BenchmarkTestNumbertFmtInt64(b *testing.B) {

	strproc := strutils.NewStringProc()
	dataset := map[interface{}]string{
		123456789101112: "123,456,789,101,112",
	}

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

	strproc := strutils.NewStringProc()
	dataset := make(map[int]paddingTestVal)

	dataset[0] = paddingTestVal{"Life isn't always what one like.", "*", strutils.PadBoth, 38, "***Life isn't always what one like.***"}
	dataset[1] = paddingTestVal{"Life isn't always what one like.", "*", strutils.PadLeft, 38, "******Life isn't always what one like."}
	dataset[2] = paddingTestVal{"Life isn't always what one like.", "*", strutils.PadRight, 38, "Life isn't always what one like.******"}
	dataset[3] = paddingTestVal{"Life isn't always what one like.", "*-=", strutils.PadBoth, 37, "*-Life isn't always what one like.*-="}
	dataset[4] = paddingTestVal{"Life isn't always what one like.", "*-=", strutils.PadLeft, 37, "*-=*-Life isn't always what one like."}
	dataset[5] = paddingTestVal{"Life isn't always what one like.", "*-=", strutils.PadRight, 37, "Life isn't always what one like.*-=*-"}

	dataset[6] = paddingTestVal{"가나다라마바사아자차카타파하", "*", strutils.PadBoth, 48, "***가나다라마바사아자차카타파하***"}
	dataset[7] = paddingTestVal{"가나다라마바사아자차카타파하", "*", strutils.PadLeft, 48, "******가나다라마바사아자차카타파하"}
	dataset[8] = paddingTestVal{"가나다라마바사아자차카타파하", "*", strutils.PadRight, 48, "가나다라마바사아자차카타파하******"}
	dataset[9] = paddingTestVal{"가나다라마바사아자차카타파하", "*-=", strutils.PadBoth, 47, "*-가나다라마바사아자차카타파하*-="}
	dataset[10] = paddingTestVal{"가나다라마바사아자차카타파하", "*-=", strutils.PadLeft, 47, "*-=*-가나다라마바사아자차카타파하"}
	dataset[11] = paddingTestVal{"가나다라마바사아자차카타파하", "*-=", strutils.PadRight, 47, "가나다라마바사아자차카타파하*-=*-"}

	for _, v := range dataset {

		retval := strproc.Padding(v.str, v.fill, v.m, v.mx)
		if v.okstr != retval {
			t.Errorf("Original Value : %v\n", v.str)
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v.okstr)
		}
	}

	//check : mx >= byteStrLen
	testStr := "test"
	retval := strproc.Padding(testStr, "*", strutils.PadBoth, 1)
	if retval != testStr {

		t.Errorf("Failure : Couldn't check the `mx >= byteStrLen`")

	}
}

func TestUppercaseFirstWords(t *testing.T) {

	strproc := strutils.NewStringProc()
	dataset := map[string]string{
		"o say, can you see, by the dawn’s early light,":                    "O Say, Can You See, By The Dawn’s Early Light,",
		"what so proudly we hailed at the twilight’s last gleaming,":        "What So Proudly We Hailed At The Twilight’s Last Gleaming,",
		"whose broad stripes and bright stars, through the perilous fight,": "Whose Broad Stripes And Bright Stars, Through The Perilous Fight,",
		"o’er the ramparts we watched, were so gallantly streaming?":        "O’er The Ramparts We Watched, Were So Gallantly Streaming?",
		"and the rockets’ red glare, the bombs bursting in air,":            "And The Rockets’ Red Glare, The Bombs Bursting In Air,",
		"gave proof through the night that our flag was still there;":       "Gave Proof Through The Night That Our Flag Was Still There;",
		"o say, does that star-spangled banner yet wave":                    "O Say, Does That Star-spangled Banner Yet Wave",
		"o’er the land of the free and the home of the brave?":              "O’er The Land Of The Free And The Home Of The Brave?",
		"가나다 라 마 바사아brownd 가나":                                              "가나다 라 마 바사아brownd 가나",
	}

	for k, v := range dataset {
		retval := strproc.UpperCaseFirstWords(k)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
	}
}

func TestLowercaseFirstWords(t *testing.T) {

	strproc := strutils.NewStringProc()
	dataset := map[string]string{
		"O SAY, CAN YOU SEE, BY THE DAWN’S EARLY LIGHT,":                    "o sAY, cAN yOU sEE, bY tHE dAWN’S eARLY lIGHT,",
		"WHAT SO PROUDLY WE HAILED AT THE TWILIGHT’S LAST GLEAMING,":        "wHAT sO pROUDLY wE hAILED aT tHE tWILIGHT’S lAST gLEAMING,",
		"WHOSE BROAD STRIPES AND BRIGHT STARS, THROUGH THE PERILOUS FIGHT,": "wHOSE bROAD sTRIPES aND bRIGHT sTARS, tHROUGH tHE pERILOUS fIGHT,",
		"O’ER THE RAMPARTS WE WATCHED, WERE SO GALLANTLY STREAMING?":        "o’ER tHE rAMPARTS wE wATCHED, wERE sO gALLANTLY sTREAMING?",
		"AND THE ROCKETS’ RED GLARE, THE BOMBS BURSTING IN AIR,":            "aND tHE rOCKETS’ rED gLARE, tHE bOMBS bURSTING iN aIR,",
		"GAVE PROOF THROUGH THE NIGHT THAT OUR FLAG WAS STILL THERE;":       "gAVE pROOF tHROUGH tHE nIGHT tHAT oUR fLAG wAS sTILL tHERE;",
		"O SAY, DOES THAT STAR-SPANGLED BANNER YET WAVE":                    "o sAY, dOES tHAT sTAR-SPANGLED bANNER yET wAVE",
		"O’ER THE LAND OF THE FREE AND THE HOME OF THE BRAVE?":              "o’ER tHE lAND oF tHE fREE aND tHE hOME oF tHE bRAVE?",
		"가나다 라 마 바사아BROWND 가나":                                              "가나다 라 마 바사아BROWND 가나",
	}

	for k, v := range dataset {
		retval := strproc.LowerCaseFirstWords(k)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
	}
}

func TestSwapCaseFirstWords(t *testing.T) {

	strproc := strutils.NewStringProc()
	dataset := map[string]string{
		"O SAY, CAN YOU SEE, BY THE DAWN’S EARLY LIGHT,":                    "o sAY, cAN yOU sEE, bY tHE dAWN’S eARLY lIGHT,",
		"WHAT SO PROUDLY WE HAILED AT THE TWILIGHT’S LAST GLEAMING,":        "wHAT sO pROUDLY wE hAILED aT tHE tWILIGHT’S lAST gLEAMING,",
		"WHOSE BROAD STRIPES AND BRIGHT STARS, THROUGH THE PERILOUS FIGHT,": "wHOSE bROAD sTRIPES aND bRIGHT sTARS, tHROUGH tHE pERILOUS fIGHT,",
		"O’ER THE RAMPARTS WE WATCHED, WERE SO GALLANTLY STREAMING?":        "o’ER tHE rAMPARTS wE wATCHED, wERE sO gALLANTLY sTREAMING?",
		"AND THE ROCKETS’ RED GLARE, THE BOMBS BURSTING IN AIR,":            "aND tHE rOCKETS’ rED gLARE, tHE bOMBS bURSTING iN aIR,",
		"GAVE PROOF THROUGH THE NIGHT THAT OUR FLAG WAS STILL THERE;":       "gAVE pROOF tHROUGH tHE nIGHT tHAT oUR fLAG wAS sTILL tHERE;",
		"O SAY, DOES THAT STAR-SPANGLED BANNER YET WAVE":                    "o sAY, dOES tHAT sTAR-SPANGLED bANNER yET wAVE",
		"O’ER THE LAND OF THE FREE AND THE HOME OF THE BRAVE?":              "o’ER tHE lAND oF tHE fREE aND tHE hOME oF tHE bRAVE?",
		"o Say, Can You See, By The Dawn’s Early Light,":                    "O say, can you see, by the dawn’s early light,",
		"what So Proudly We Hailed At The Twilight’s Last Gleaming,":        "What so proudly we hailed at the twilight’s last gleaming,",
		"whose Broad Stripes And Bright Stars, Through The Perilous Fight,": "Whose broad stripes and bright stars, through the perilous fight,",
		"o’er The Ramparts We Watched, Were So Gallantly Streaming?":        "O’er the ramparts we watched, were so gallantly streaming?",
		"and The Rockets’ Red Glare, The Bombs Bursting In Air,":            "And the rockets’ red glare, the bombs bursting in air,",
		"gave Proof Through The Night That Our Flag Was Still There;":       "Gave proof through the night that our flag was still there;",
		"o Say, Does That Star-spangled Banner Yet Wave":                    "O say, does that star-spangled banner yet wave",
		"o’er The Land Of The Free And The Home Of The Brave?":              "O’er the land of the free and the home of the brave?",
		"가나다 라 마 바사아brownd 가나":                                              "가나다 라 마 바사아brownd 가나",
		"o sAY, cAN yOU sEE, bY tHE dAWN’S eARLY lIGHT,":                    "O SAY, CAN YOU SEE, BY THE DAWN’S EARLY LIGHT,",
		"wHAT sO pROUDLY wE hAILED aT tHE tWILIGHT’S lAST gLEAMING,":        "WHAT SO PROUDLY WE HAILED AT THE TWILIGHT’S LAST GLEAMING,",
		"wHOSE bROAD sTRIPES aND bRIGHT sTARS, tHROUGH tHE pERILOUS fIGHT,": "WHOSE BROAD STRIPES AND BRIGHT STARS, THROUGH THE PERILOUS FIGHT,",
		"o’ER tHE rAMPARTS wE wATCHED, wERE sO gALLANTLY sTREAMING?":        "O’ER THE RAMPARTS WE WATCHED, WERE SO GALLANTLY STREAMING?",
		"aND tHE rOCKETS’ rED gLARE, tHE bOMBS bURSTING iN aIR,":            "AND THE ROCKETS’ RED GLARE, THE BOMBS BURSTING IN AIR,",
		"gAVE pROOF tHROUGH tHE nIGHT tHAT oUR fLAG wAS sTILL tHERE;":       "GAVE PROOF THROUGH THE NIGHT THAT OUR FLAG WAS STILL THERE;",
		"o sAY, dOES tHAT sTAR-SPANGLED bANNER yET wAVE":                    "O SAY, DOES THAT STAR-SPANGLED BANNER YET WAVE",
		"o’ER tHE lAND oF tHE fREE aND tHE hOME oF tHE bRAVE?":              "O’ER THE LAND OF THE FREE AND THE HOME OF THE BRAVE?",
		"가나다 라 마 바사아BROWND 가나":                                              "가나다 라 마 바사아BROWND 가나",
	}

	for k, v := range dataset {
		retval := strproc.SwapCaseFirstWords(k)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
	}
}

func TestHumanByteSize(t *testing.T) {

	strproc := strutils.NewStringProc()
	dataset := map[interface{}]string{
		1.7976931348623157e+308: "152270531428124968725096603469261934082567927321390584004196605238063615198482718997460353589210907119043200911085747810785909744915680620242659147418948017662928903247753430023357200398869394856103928002466673473125884404826265988290381563441726944871732658253337089007918982991007711232.00Yb",
		1170:         "1.14Kb",
		72125099:     "68.78Mb",
		3276537856:   "3.05Gb",
		27:           "27.00B",
		93735736:     "89.39Mb",
		937592:       "915.62Kb",
		6715287:      "6.40Mb",
		2856906752:   "2.66Gb",
		7040152:      "6.71Mb",
		22016:        "21.50Kb",
		"1170":       "1.14Kb",
		"72125099":   "68.78Mb",
		"3276537856": "3.05Gb",
		"27":         "27.00B",
		"93735736":   "89.39Mb",
		"937592":     "915.62Kb",
		"6715287":    "6.40Mb",
		"2856906752": "2.66Gb",
		"7040152":    "6.71Mb",
		"22016":      "21.50Kb",
		3.40282346638528859811704183484516925440e+38: "288230358971842560.00Yb",
	}

	for k, v := range dataset {
		retval, err := strproc.HumanByteSize(k, 2, strutils.CamelCaseDouble)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
		if err != nil {
			t.Errorf("Error : %v", err)
		}
	}

	var err error

	//check : unit < UpperCaseSingle || unit > CamelCaseLong
	_, err = strproc.HumanByteSize(`1234`, 2, 123)
	if err == nil {
		t.Errorf("Failure : Couldn't check the `retval, err := strproc.HumanByteSize(k, 2, strutils.CamelCaseDouble)`")
	}

	//check : numberToString
	_, err = strproc.HumanByteSize(`abc`, 2, strutils.UpperCaseDouble)
	if err == nil {
		t.Errorf("Failure : Couldn't check the `can't convert number to string`")
	}

	//check : ParseFloat
	_, err = strproc.HumanByteSize(`1234.1234+38`, 2, strutils.UpperCaseDouble)
	if err == nil {
		t.Errorf("Failure : Couldn't check the `strconv.ParseFloat(strNum, 64)`")
	}

}

func TestHumanFileSize(t *testing.T) {

	const tmpFilePath = "./filesizecheck.touch"
	const tmpPath = "./testdir"
	var err error

	//generating a touch file
	tmpdata := []byte("123456789")
	ioutil.WriteFile(tmpFilePath, tmpdata, 0750)

	strproc := strutils.NewStringProc()
	_, err = strproc.HumanFileSize(tmpFilePath, 2, strutils.CamelCaseLong)
	if err != nil {
		t.Errorf("Error : %v", err)
	}

	_, err = strproc.HumanFileSize(tmpFilePath, 2, strutils.CamelCaseDouble)
	if err != nil {
		t.Errorf("Error : %v", err)
	}

	os.Remove(tmpFilePath)

	//check : isDir
	err = os.MkdirAll(tmpPath, 0777)
	if err != nil {
		os.Remove(tmpPath)
		t.Errorf("Failuew : Couldn't Mkdir %q: %s", tmpPath, err)
	}

	_, err = strproc.HumanFileSize(tmpPath, 2, strutils.CamelCaseDouble)
	if err == nil {
		os.Remove(tmpPath)
		t.Errorf("Failure : Couldn't check the `stat.IsDir()`")
	}

	os.Remove(tmpPath)
}

func TestAnyCompare(t *testing.T) {

	var retval bool
	var err error

	strproc := strutils.NewStringProc()

	testInt1 := []int{1, 2, 3}
	testInt2 := []int{1, 2, 3}
	retval, err = strproc.AnyCompare(testInt1, testInt2)
	if retval == false {
		t.Errorf("Could not make an accurate comparison : %v", err)
	}

	testIntFalse1 := []int{1, 2, 3}
	testIntFalse2 := []int{1, 2, 1}
	retval, err = strproc.AnyCompare(testIntFalse1, testIntFalse2)
	if retval == true {
		t.Errorf("Could not make an accurate comparison.")
	}

	testMapStr1 := map[string]string{"a": "va", "vb": "vb"}
	testMapStr2 := map[string]string{"a": "va", "vb": "vb"}
	retval, err = strproc.AnyCompare(testMapStr1, testMapStr2)
	if retval == false {
		t.Errorf("Could not make an accurate comparison : %v", err)
	}

	testMapStrFalse1 := map[string]string{"a": "va", "vb": "vb"}
	testMapStrFalse2 := map[string]string{"a": "va", "v": "vb"}
	retval, err = strproc.AnyCompare(testMapStrFalse1, testMapStrFalse2)
	if retval == true {
		t.Errorf("Could not make an accurate comparison.")
	}

	testMapBool1 := map[string]bool{"a": false, "vb": false}
	testMapBool2 := map[string]bool{"a": false, "vb": true}
	retval, err = strproc.AnyCompare(testMapBool1, testMapBool2)
	if retval == true {
		t.Errorf("Could not make an accurate comparison : %v", err)
	}

	testMultipleDepthMap1 := map[string]map[string]string{
		"H": map[string]string{
			"name":  "Hydrogen",
			"state": "gas",
		},
		"He": map[string]string{
			"name":  "Helium",
			"state": "gas",
		},
		"Li": map[string]string{
			"name":  "Lithium",
			"state": "solid",
		},
		"Be": map[string]string{
			"name":  "Beryllium",
			"state": "solid",
		},
		"B": map[string]string{
			"name":  "Boron",
			"state": "solid",
		},
		"C": map[string]string{
			"name":  "Carbon",
			"state": "solid",
		},
		"N": map[string]string{
			"name":  "Nitrogen",
			"state": "gas",
		},
		"O": map[string]string{
			"name":  "Oxygen",
			"state": "gas",
		},
		"F": map[string]string{
			"name":  "Fluorine",
			"state": "gas",
		},
		"Ne": map[string]string{
			"name":  "Neon",
			"state": "gas",
		},
	}

	testMultipleDepthMap2 := map[string]map[string]string{
		"H": map[string]string{
			"name":  "Hydrogen",
			"state": "gas",
		},
		"He": map[string]string{
			"name":  "Helium",
			"state": "gas",
		},
		"Li": map[string]string{
			"name":  "Lithium",
			"state": "solid",
		},
		"Be": map[string]string{
			"name":  "Beryllium",
			"state": "solid",
		},
		"B": map[string]string{
			"name":  "Boron",
			"state": "solid",
		},
		"C": map[string]string{
			"name":  "Carbon",
			"state": "solid",
		},
		"N": map[string]string{
			"name":  "Nitrogen",
			"state": "gas",
		},
		"O": map[string]string{
			"name":  "Oxygen",
			"state": "gas",
		},
		"F": map[string]string{
			"name":  "Fluorine",
			"state": "gas",
		},
		"Ne": map[string]string{
			"name":  "Neon",
			"state": "gas",
		},
	}

	retval, err = strproc.AnyCompare(testMultipleDepthMap1, testMultipleDepthMap2)
	if retval == false {
		t.Errorf("Could not make an accurate comparison : %v", err)
	}

	testMultipleDepthMapFalse1 := map[string]map[string]string{
		"H": map[string]string{
			"name":  "Hydrogen",
			"state": "gas",
		},
		"He": map[string]string{
			"name":  "Helium",
			"state": "gas",
		},
		"Li": map[string]string{
			"name":  "Lithium",
			"state": "solid",
		},
		"Be": map[string]string{
			"name":  "Beryllium",
			"state": "solid",
		},
		"B": map[string]string{
			"name":  "Boron",
			"state": "solid",
		},
		"C": map[string]string{
			"name":  "Carbon",
			"state": "solid",
		},
		"N": map[string]string{
			"name":  "Nitrogen",
			"state": "gas",
		},
		"O": map[string]string{
			"name":  "Oxygen",
			"state": "gas",
		},
		"F": map[string]string{
			"name":  "Fluorine",
			"state": "gas",
		},
		"Ne": map[string]string{
			"name":  "Neon",
			"state": "gas",
		},
	}

	testMultipleDepthMapFalse2 := map[string]map[string]string{
		"H": map[string]string{
			"name":  "Hydrogen",
			"state": "gas",
		},
		"He": map[string]string{
			"name":  "Helium",
			"state": "gas",
		},
		"Li": map[string]string{
			"name":  "Lithium",
			"state": "solid",
		},
		"Be": map[string]string{
			"name":  "Beryllium",
			"state": "solid",
		},
		"B": map[string]string{
			"name":  "Boron",
			"state": "solid",
		},
		"C": map[string]string{
			"name":  "Carbon",
			"state": "solid",
		},
		"N": map[string]string{
			"name":  "Nitrogen",
			"state": "gas",
		},
		"O": map[string]string{
			"name":  "Oxygen1",
			"state": "gas",
		},
		"F": map[string]string{
			"name":  "Fluorine",
			"state": "gas",
		},
		"Ne": map[string]string{
			"name1": "Neon",
			"state": "gas",
		},
	}

	retval, _ = strproc.AnyCompare(testMultipleDepthMapFalse1, testMultipleDepthMapFalse2)
	if retval == true {
		t.Errorf("Could not make an accurate comparison.")
	}

	testComplexMap1 := map[string]map[string]map[string]int{
		"F": map[string]map[string]int{
			"name": map[string]int{
				"first": 1,
				"last":  2,
			},
		},
		"A": map[string]map[string]int{
			"name": map[string]int{
				"first": 11,
				"last":  21,
			},
		},
	}

	testComplexMap2 := map[string]map[string]map[string]int{
		"F": map[string]map[string]int{
			"name": map[string]int{
				"first": 11,
				"last":  12222,
			},
		},
		"A": map[string]map[string]int{
			"name": map[string]int{
				"first": 11,
				"last":  21,
			},
		},
	}

	retval, _ = strproc.AnyCompare(testComplexMap1, testComplexMap2)
	if retval == true {
		t.Errorf("Could not make an accurate comparison.")
	}

	//check : different type
	testDiffType1 := int(32)
	testDiffType2 := int64(32)
	_, err = strproc.AnyCompare(testDiffType1, testDiffType2)
	if err == nil {
		t.Errorf("Failure : Couldn't check the `Different Type`")
	}

	//check : different len
	testDiffLen1 := []int{1, 2, 3, 4, 5, 6, 7}
	testDiffLen2 := []int{1, 2, 3}
	_, err = strproc.AnyCompare(testDiffLen1, testDiffLen2)
	if err == nil {
		t.Errorf("Failure : Couldn't check the `Different Len`")
	}

	//check : not support compre
	testDiffNotSupport1 := paddingTestVal{}
	testDiffNotSupport2 := paddingTestVal{}
	_, err = strproc.AnyCompare(testDiffNotSupport1, testDiffNotSupport2)
	if err == nil {
		t.Errorf("Failure : Couldn't check the `Not Support Compre`")
		t.Errorf("Error : %v", err)
	}

}
