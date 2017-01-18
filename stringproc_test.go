// Package strutils made by torden <https://github.com/torden/go-strutil>
// license that can be found in the LICENSE file.
package strutils

import (
	"strings"
	"testing"

	"github.com/dustin/go-humanize"
)

func TestAddSlashes(t *testing.T) {

	strproc := NewStringProc()
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

	strproc := NewStringProc()
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

	strproc := NewStringProc()
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

	strproc := NewStringProc()
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
	strproc := NewStringProc()

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

		retval := strproc.WordWrapSimple(v.str, v.wd, v.breakstr)
		if v.okstr != retval {
			t.Errorf("Original Value : %v\n", v.str)
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v.okstr)
		}
	}
}

func TestWordWrapAround(t *testing.T) {
	strproc := NewStringProc()

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

		retval := strproc.WordWrapAround(v.str, v.wd, v.breakstr)
		if v.okstr != retval {
			t.Errorf("Original Value : %v\n", v.str)
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v.okstr)
		}
	}
}

func TestNumbertFmt(t *testing.T) {

	strproc := NewStringProc()
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
		retval, err := strproc.NumberFmt(k)
		if v != retval {
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		}
		if err != nil {
			t.Errorf("Return Error : %v", err)
		}
	}
}

func BenchmarkTestNumbertFmt(b *testing.B) {

	strproc := NewStringProc()
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

	strproc := NewStringProc()
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

	strproc := NewStringProc()
	dataset := make(map[int]paddingTestVal)

	dataset[0] = paddingTestVal{"Life isn't always what one like.", "*", padBoth, 38, "***Life isn't always what one like.***"}
	dataset[1] = paddingTestVal{"Life isn't always what one like.", "*", padLeft, 38, "******Life isn't always what one like."}
	dataset[2] = paddingTestVal{"Life isn't always what one like.", "*", padRight, 38, "Life isn't always what one like.******"}
	dataset[3] = paddingTestVal{"Life isn't always what one like.", "*-=", padBoth, 37, "*-Life isn't always what one like.*-="}
	dataset[4] = paddingTestVal{"Life isn't always what one like.", "*-=", padLeft, 37, "*-=*-Life isn't always what one like."}
	dataset[5] = paddingTestVal{"Life isn't always what one like.", "*-=", padRight, 37, "Life isn't always what one like.*-=*-"}

	dataset[6] = paddingTestVal{"가나다라마바사아자차카타파하", "*", padBoth, 48, "***가나다라마바사아자차카타파하***"}
	dataset[7] = paddingTestVal{"가나다라마바사아자차카타파하", "*", padLeft, 48, "******가나다라마바사아자차카타파하"}
	dataset[8] = paddingTestVal{"가나다라마바사아자차카타파하", "*", padRight, 48, "가나다라마바사아자차카타파하******"}
	dataset[9] = paddingTestVal{"가나다라마바사아자차카타파하", "*-=", padBoth, 47, "*-가나다라마바사아자차카타파하*-="}
	dataset[10] = paddingTestVal{"가나다라마바사아자차카타파하", "*-=", padLeft, 47, "*-=*-가나다라마바사아자차카타파하"}
	dataset[11] = paddingTestVal{"가나다라마바사아자차카타파하", "*-=", padRight, 47, "가나다라마바사아자차카타파하*-=*-"}

	for _, v := range dataset {

		retval := strproc.padding(v.str, v.fill, v.m, v.mx)
		if v.okstr != retval {
			t.Errorf("Original Value : %v\n", v.str)
			t.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v.okstr)
		}
	}
}

func TestUppercaseFirstWords(t *testing.T) {

	strproc := NewStringProc()
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

	strproc := NewStringProc()
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

	strproc := NewStringProc()
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
