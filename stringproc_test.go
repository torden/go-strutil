package strutils_test

import (
	"io/ioutil"
	"math"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	strutils "github.com/torden/go-strutil"
)

func Test_strutils_AddSlashes(t *testing.T) {
	t.Parallel()
	assert.TurnOnUnitTestMode()

	dataset := map[string]string{
		`대한민국만세`:     `대한민국만세`,
		`대한\민국만세`:    `대한\\민국만세`,
		`대한\\민국만세`:   `대한\\민국만세`,
		"abcdefgz":   "abcdefgz",
		`a\bcdefgz`:  `a\\bcdefgz`,
		`a\\bcdefgz`: `a\\bcdefgz`,
	}

	// check : common
	var retval string

	for k, v := range dataset {
		retval = strproc.AddSlashes(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func Test_strutils_StripSlashes(t *testing.T) {
	t.Parallel()

	dataset := map[string]string{
		`대한민국만세`:       `대한민국만세`,
		`대한\\민국만세`:     `대한\민국만세`,
		`대한\\\\민국만세`:   `대한\\민국만세`,
		"abcdefgz":     "abcdefgz",
		`a\\bcdefgz`:   `a\bcdefgz`,
		`a\\\\bcdefgz`: `a\\bcdefgz`,
	}

	// check : common
	var retval string

	for k, v := range dataset {
		retval = strproc.StripSlashes(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func Test_strutils_Nl2Br(t *testing.T) {
	t.Parallel()

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

	// check : common
	var retval string
	for k, v := range dataset {
		retval = strproc.Nl2Br(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func Test_strutils_Br2Nl(t *testing.T) {
	t.Parallel()

	dataset := map[string]string{
		"대한<br />민국만세":   "대한\n민국만세",
		"대한민국만세<br />":   "대한민국만세\n",
		"abc<br />defgh": "abc\ndefgh",
		"abcde<br />fgh": "abcde\nfgh",
		"abcdefgh<br />": "abcdefgh\n",

		"대한<br/>민국만세":   "대한\n민국만세",
		"대한민국만세<br/>":   "대한민국만세\n",
		"abc<br/>defgh": "abc\ndefgh",
		"abcde<br/>fgh": "abcde\nfgh",
		"abcdefgh<br/>": "abcdefgh\n",

		"대한<br>민국만세":   "대한\n민국만세",
		"대한민국만세<br>":   "대한민국만세\n",
		"abc<br>defgh": "abc\ndefgh",
		"abcde<br>fgh": "abcde\nfgh",
		"abcdefgh<br>": "abcdefgh\n",

		"abcdefgh": "abcdefgh",
		"대한민국만세":   "대한민국만세",

		"<a href='http://www.president.go.kr/'>대한민국만세</a><br>":   "<a href='http://www.president.go.kr/'>대한민국만세</a>\n",
		"<a href='http://www.president.go.kr/'>abcde</a><br>fgh": "<a href='http://www.president.go.kr/'>abcde</a>\nfgh",

		"world peace!!<a href='http://www.president.go.kr/'>대한민국만세</a><br>":   "world peace!!<a href='http://www.president.go.kr/'>대한민국만세</a>\n",
		"world peace!!<a href='http://www.president.go.kr/'>abcde</a><br>fgh": "world peace!!<a href='http://www.president.go.kr/'>abcde</a>\nfgh",

		"world peace!!<a href='http://www.president.go.kr/'><br />대한민국만세</a><br>":   "world peace!!<a href='http://www.president.go.kr/'>\n대한민국만세</a>\n",
		"world peace!!<a href='http://www.president.go.kr/'><br />abcde</a><br>fgh": "world peace!!<a href='http://www.president.go.kr/'>\nabcde</a>\nfgh",
	}

	// check : common
	var retval string
	for k, v := range dataset {
		retval = strproc.Br2Nl(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

type wordwrapTestVal struct {
	str      string
	wd       int
	breakstr string
	okstr    string
}

func Test_strutils_WordWrapSimple(t *testing.T) {
	t.Parallel()

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

	// check : common
	for _, v := range dataset {

		retval, _ := strproc.WordWrapSimple(v.str, v.wd, v.breakstr)
		assert.AssertEquals(t, v.okstr, retval, "Original Value : %v\nReturn Value mismatch.\nExpected: %v\nActual: %v", v.str, retval, v)
	}

	// check : wd = 0
	_, err := strproc.WordWrapSimple("test", 0, "1234")
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `wd at least 1`")
}

func Test_strutils_WordWrapAround(t *testing.T) {
	t.Parallel()

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

	// check : common
	for _, v := range dataset {

		retval, _ := strproc.WordWrapAround(v.str, v.wd, v.breakstr)
		assert.AssertEquals(t, v.okstr, retval, "Original Value : %v\nReturn Value mismatch.\nExpected: %v\nActual: %v", v.str, retval, v.okstr)
	}

	var err error

	// check : wd = 0
	_, err = strproc.WordWrapAround("test", 0, "1234")
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `wd at least 1`")

	// check : lastspc = 1
	_, err = strproc.WordWrapAround("ttttttt tttttttttt", 2, "1111")
	assert.AssertNil(t, err, "Failure : Couldn't check the `lastspc = 1`")

	// check : except
	_, err = strproc.WordWrapAround("t t", 1, "*")
	assert.AssertNil(t, err, "Failure : Couldn't check the `specific except`")
}

func Test_strutils_NumbertFmt(t *testing.T) {
	t.Parallel()

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
		// BUG(r) :
		// int8(math.MinInt8):       "-128",
		// float32(math.SmallestNonzeroFloat32): "1e-45",
		// float64(math.SmallestNonzeroFloat64): "5e-324",

	}

	// check : common
	for k, v := range dataset {
		retval, err := strproc.NumberFmt(k)

		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		assert.AssertNil(t, err, "Return Error : %v", err)
	}

	var err error

	// check : ParseFloat
	_, err = strproc.NumberFmt("12.11111111111111111111111111111111111111111111111111111111111e12e12e1p029ekj12e")
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `Not Support strconv.ParseFloat`")

	// check : not support obj
	_, err = strproc.NumberFmt(complex128(123))
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `not support obj`")

	// check : not support obj
	_, err = strproc.NumberFmt(complex64(123))
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `not support obj`")

	// check : not support obj
	_, err = strproc.NumberFmt(true)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `not support obj`")

	// check : not support numric string
	_, err = strproc.NumberFmt("1234===121212")
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `not support obj`")
}

type paddingTestVal struct {
	str   string
	fill  string
	m     int
	mx    int
	okstr string
}

func Test_strutils_Padding(t *testing.T) {
	t.Parallel()

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

	// check : common
	for _, v := range dataset {

		retval := strproc.Padding(v.str, v.fill, v.m, v.mx)
		assert.AssertEquals(t, v.okstr, retval, "Original Value : %v\nReturn Value mismatch.\nExpected: %v\nActual: %v", v.str, retval, v.okstr)
	}

	// check : mx >= byteStrLen
	testStr := "test"
	retval := strproc.Padding(testStr, "*", strutils.PadBoth, 1)
	assert.AssertEquals(t, testStr, retval, "Failure : Couldn't check the `mx >= byteStrLen`")
}

func Test_strutils_UppercaseFirstWords(t *testing.T) {
	t.Parallel()

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

	// check : common
	for k, v := range dataset {
		retval := strproc.UpperCaseFirstWords(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func Test_strutils_LowercaseFirstWords(t *testing.T) {
	t.Parallel()

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

	// check : common
	for k, v := range dataset {
		retval := strproc.LowerCaseFirstWords(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func Test_strutils_SwapCaseFirstWords(t *testing.T) {
	t.Parallel()

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

	// check : common
	for k, v := range dataset {
		retval := strproc.SwapCaseFirstWords(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func Test_strutils_HumanByteSize(t *testing.T) {
	t.Parallel()

	dataset := map[interface{}]string{
		1.7976931348623157e+308: "152270531428124968725096603469261934082567927321390584004196605238063615198482718997460353589210907119043200911085747810785909744915680620242659147418948017662928903247753430023357200398869394856103928002466673473125884404826265988290381563441726944871732658253337089007918982991007711232.00Yb",
		1170:                    "1.14Kb",
		72125099:                "68.78Mb",
		3276537856:              "3.05Gb",
		27:                      "27.00B",
		93735736:                "89.39Mb",
		937592:                  "915.62Kb",
		6715287:                 "6.40Mb",
		2856906752:              "2.66Gb",
		7040152:                 "6.71Mb",
		22016:                   "21.50Kb",
		"1170":                  "1.14Kb",
		"72125099":              "68.78Mb",
		"3276537856":            "3.05Gb",
		"27":                    "27.00B",
		"93735736":              "89.39Mb",
		"937592":                "915.62Kb",
		"6715287":               "6.40Mb",
		"2856906752":            "2.66Gb",
		"7040152":               "6.71Mb",
		"22016":                 "21.50Kb",
		3.40282346638528859811704183484516925440e+38:                                   "288230358971842560.00Yb",
		"12121212121212121212121212121212121212121211212121212211212121212121.1234e+3": "0.00NaN",
	}

	// check : common
	for k, v := range dataset {
		retval, err := strproc.HumanByteSize(k, 2, strutils.CamelCaseDouble)

		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		assert.AssertNil(t, err, "Error : %v", err)
	}

	var err error

	// check : unit < UpperCaseSingle || unit > CamelCaseLong
	_, err = strproc.HumanByteSize(`1234`, 2, 123)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `retval, err := strproc.HumanByteSize(k, 2, strutils.CamelCaseDouble)`")

	// check : numberToString
	_, err = strproc.HumanByteSize(`abc`, 2, strutils.UpperCaseDouble)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `can't convert number to string`")

	// check : ParseFloat
	_, err = strproc.HumanByteSize("100.7976931348623157e+308", 2, strutils.UpperCaseDouble)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `strconv.ParseFloat(strNum, 64)`")

	// check : Complex64
	_, err = strproc.HumanByteSize(complex64(13), 2, strutils.UpperCaseDouble)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `not support obj.(complex128)`")

	// check : Complex128
	_, err = strproc.HumanByteSize(complex128(2+3i), 2, strutils.UpperCaseDouble)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `not support obj.(complex128)`")
}

func Test_strutils_HumanFileSize(t *testing.T) {
	t.Parallel()

	const tmpFilePath = "./filesizecheck.touch"
	const tmpPath = "./testdir"

	var retval string
	var err error

	// generating a touch file
	tmpdata := []byte(strings.Repeat("*", 1024*1024*13))
	err = ioutil.WriteFile(tmpFilePath, tmpdata, 0750)
	assert.AssertNil(t, err, "Error : ", err)

	dataset := map[uint8]string{
		strutils.LowerCaseSingle: "13.00m",
		strutils.LowerCaseDouble: "13.00mb",
		strutils.UpperCaseSingle: "13.00M",
		strutils.UpperCaseDouble: "13.00MB",
		strutils.CamelCaseLong:   "13.00MegaByte",
		strutils.CamelCaseDouble: "13.00Mb",
	}

	// check : common
	for k, v := range dataset {
		retval, err = strproc.HumanFileSize(tmpFilePath, 2, k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		assert.AssertNil(t, err, "Error : %v", err)
	}

	// check : lost file description
	go func() {
		defer func() { recover() }()
		time.Sleep(time.Nanosecond * 10)
		err := os.Remove(tmpFilePath)
		assert.AssertLog(t, err, "lost file : %s but it's OK", tmpFilePath)
		return
	}()
	_, err = strproc.HumanFileSize(tmpFilePath, 2, strutils.CamelCaseDouble)
	assert.AssertLog(t, err, "PASS")

	defer func(t *testing.T) {
		err := os.Remove(tmpFilePath)
		assert.AssertLog(t, err, "lost file : %s but it's OK", tmpFilePath)
		return
	}(t)

	// check : isDir
	err = os.MkdirAll(tmpPath, 0777)
	assert.AssertNil(t, err, "Failure : Couldn't Mkdir %q: %s", tmpPath, err)

	_, err = strproc.HumanFileSize(tmpPath, 2, strutils.CamelCaseDouble)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `stat.IsDir()`")

	// check : os.Open
	_, err = strproc.HumanFileSize("/hello_word_txt", 2, strutils.CamelCaseDouble)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `os.Open()`")

	// check : not support obj.(complex128)
	_, err = strproc.HumanByteSize(complex128(1+3i), 2, strutils.CamelCaseLong)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `Not Support obj.(complex129)`")
}

func Test_strutils_AnyCompare(t *testing.T) {
	t.Parallel()

	var retval bool
	var err error

	testStr1 := "ABCD"
	testStr2 := "ABCD"
	retval, err = strproc.AnyCompare(testStr1, testStr2)
	assert.AssertTrue(t, retval, "Couldn't make an accurate comparison : %v", err)

	testUint1 := uint64(1234567)
	testUint2 := uint64(1234567)
	retval, err = strproc.AnyCompare(testUint1, testUint2)
	assert.AssertTrue(t, retval, "Couldn't make an accurate comparison : %v", err)

	testFloat1 := float64(123.123)
	testFloat2 := float64(123.123)
	retval, err = strproc.AnyCompare(testFloat1, testFloat2)
	assert.AssertTrue(t, retval, "Couldn't make an accurate comparison : %v", err)

	testComplex1 := complex64(123.123)
	testComplex2 := complex64(123.123)
	retval, err = strproc.AnyCompare(testComplex1, testComplex2)
	assert.AssertTrue(t, retval, "Couldn't make an accurate comparison : %v", err)

	testInt1 := []int{1, 2, 3}
	testInt2 := []int{1, 2, 3}
	retval, err = strproc.AnyCompare(testInt1, testInt2)
	assert.AssertTrue(t, retval, "Couldn't make an accurate comparison : %v", err)

	testIntFalse1 := []int{1, 2, 3}
	testIntFalse2 := []int{1, 2, 1}
	retval, err = strproc.AnyCompare(testIntFalse1, testIntFalse2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison : %v", err)

	testMultipleDepthMapDiffType1 := map[string]map[string]string{
		"H": {
			"name":  "Hydrogen",
			"state": "gas",
		},
	}

	testMultipleDepthMapDiffType2 := map[string]map[string]int{
		"H": {
			"name":  1,
			"state": 2,
		},
	}
	retval, err = strproc.AnyCompare(testMultipleDepthMapDiffType1, testMultipleDepthMapDiffType2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison : %v", err)

	testMultipleDepthMapDiffType3 := map[string]map[string]int{
		"H": {
			"name":  1,
			"state": 2,
		},
	}

	testMultipleDepthMapDiffType4 := map[string]map[string]uint{
		"H": {
			"name":  1,
			"state": 2,
		},
	}
	retval, err = strproc.AnyCompare(testMultipleDepthMapDiffType3, testMultipleDepthMapDiffType4)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison : %v", err)

	testMultipleDepthMapDiffType5 := map[string]map[string]uint{
		"H": {
			"name":  1,
			"state": 2,
		},
	}

	testMultipleDepthMapDiffType6 := map[string]map[string]int{
		"H": {
			"name":  1,
			"state": 2,
		},
	}
	retval, err = strproc.AnyCompare(testMultipleDepthMapDiffType5, testMultipleDepthMapDiffType6)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison : %v", err)

	testMultipleDepthMapDiffType7 := map[string]map[string]float64{
		"H": {
			"name":  1,
			"state": 2,
		},
	}

	testMultipleDepthMapDiffType8 := map[string]map[string]int{
		"H": {
			"name":  1,
			"state": 2,
		},
	}
	retval, err = strproc.AnyCompare(testMultipleDepthMapDiffType7, testMultipleDepthMapDiffType8)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison : %v", err)

	testMultipleDepthMapDiffType9 := map[string]map[string]complex64{
		"H": {
			"name":  1,
			"state": 2,
		},
	}

	testMultipleDepthMapDiffType10 := map[string]map[string]int{
		"H": {
			"name":  1,
			"state": 2,
		},
	}
	retval, err = strproc.AnyCompare(testMultipleDepthMapDiffType9, testMultipleDepthMapDiffType10)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison : %v", err)

	testMapStr1 := map[string]string{"a": "va", "vb": "vb"}
	testMapStr2 := map[string]string{"a": "va", "vb": "vb"}
	retval, err = strproc.AnyCompare(testMapStr1, testMapStr2)
	assert.AssertTrue(t, retval, "Couldn't make an accurate comparison : %v", err)

	testMapStrDiff1 := map[string]string{"a": "va", "vb": "v"}
	testMapStrDiff2 := map[string]string{"a": "va", "vb": "vb"}
	retval, err = strproc.AnyCompare(testMapStrDiff1, testMapStrDiff2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison : %v", err)

	testMapStrFalse1 := map[string]string{"a": "va", "vb": "vb"}
	testMapStrFalse2 := map[string]string{"a": "va", "v": "vb"}
	retval, err = strproc.AnyCompare(testMapStrFalse1, testMapStrFalse2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison : %v", err)

	testMapBool1 := map[string]bool{"a": false, "vb": false}
	testMapBool2 := map[string]bool{"a": false, "vb": true}
	retval, err = strproc.AnyCompare(testMapBool1, testMapBool2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison : %v", err)

	testMultipleDepthMap1 := map[string]map[string]string{
		"H": {
			"name":  "Hydrogen",
			"state": "gas",
		},
		"He": {
			"name":  "Helium",
			"state": "gas",
		},
		"Li": {
			"name":  "Lithium",
			"state": "solid",
		},
		"Be": {
			"name":  "Beryllium",
			"state": "solid",
		},
		"B": {
			"name":  "Boron",
			"state": "solid",
		},
		"C": {
			"name":  "Carbon",
			"state": "solid",
		},
		"N": {
			"name":  "Nitrogen",
			"state": "gas",
		},
		"O": {
			"name":  "Oxygen",
			"state": "gas",
		},
		"F": {
			"name":  "Fluorine",
			"state": "gas",
		},
		"Ne": {
			"name":  "Neon",
			"state": "gas",
		},
	}

	testMultipleDepthMap2 := map[string]map[string]string{
		"H": {
			"name":  "Hydrogen",
			"state": "gas",
		},
		"He": {
			"name":  "Helium",
			"state": "gas",
		},
		"Li": {
			"name":  "Lithium",
			"state": "solid",
		},
		"Be": {
			"name":  "Beryllium",
			"state": "solid",
		},
		"B": {
			"name":  "Boron",
			"state": "solid",
		},
		"C": {
			"name":  "Carbon",
			"state": "solid",
		},
		"N": {
			"name":  "Nitrogen",
			"state": "gas",
		},
		"O": {
			"name":  "Oxygen",
			"state": "gas",
		},
		"F": {
			"name":  "Fluorine",
			"state": "gas",
		},
		"Ne": {
			"name":  "Neon",
			"state": "gas",
		},
	}

	retval, err = strproc.AnyCompare(testMultipleDepthMap1, testMultipleDepthMap2)
	assert.AssertTrue(t, retval, "Couldn't make an accurate comparison : %v", err)

	testMultipleDepthMapFalse1 := map[string]map[string]string{
		"H": {
			"name":  "Hydrogen",
			"state": "gas",
		},
		"He": {
			"name":  "Helium",
			"state": "gas",
		},
		"Li": {
			"name":  "Lithium",
			"state": "solid",
		},
		"Be": {
			"name":  "Beryllium",
			"state": "solid",
		},
		"B": {
			"name":  "Boron",
			"state": "solid",
		},
		"C": {
			"name":  "Carbon",
			"state": "solid",
		},
		"N": {
			"name":  "Nitrogen",
			"state": "gas",
		},
		"O": {
			"name":  "Oxygen",
			"state": "gas",
		},
		"F": {
			"name":  "Fluorine",
			"state": "gas",
		},
		"Ne": {
			"name":  "Neon",
			"state": "gas",
		},
	}

	testMultipleDepthMapFalse2 := map[string]map[string]string{
		"H": {
			"name":  "Hydrogen",
			"state": "gas",
		},
		"He": {
			"name":  "Helium",
			"state": "gas",
		},
		"Li": {
			"name":  "Lithium",
			"state": "solid",
		},
		"Be": {
			"name":  "Beryllium",
			"state": "solid",
		},
		"B": {
			"name":  "Boron",
			"state": "solid",
		},
		"C": {
			"name":  "Carbon",
			"state": "solid",
		},
		"N": {
			"name":  "Nitrogen",
			"state": "gas",
		},
		"O": {
			"name":  "Oxygen1",
			"state": "gas",
		},
		"F": {
			"name":  "Fluorine",
			"state": "gas",
		},
		"Ne": {
			"name1": "Neon",
			"state": "gas",
		},
	}

	retval, _ = strproc.AnyCompare(testMultipleDepthMapFalse1, testMultipleDepthMapFalse2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison.")

	testComplexMap1 := map[string]map[string]map[string]int{
		"F": {
			"name": {
				"first": 1,
				"last":  2,
			},
		},
		"A": {
			"name": {
				"first": 11,
				"last":  21,
			},
		},
	}

	testComplexMap2 := map[string]map[string]map[string]int{
		"F": {
			"name": {
				"first": 11,
				"last":  12222,
			},
		},
		"A": {
			"name": {
				"first": 11,
				"last":  21,
			},
		},
	}

	retval, _ = strproc.AnyCompare(testComplexMap1, testComplexMap2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison.")

	// check : uint in map
	testMDepthUint1 := map[string]map[string]uint{"H": {"name": 1, "state": 2}}
	testMDepthUint2 := map[string]map[string]uint{"H": {"name": 1, "state": 3}}
	retval, _ = strproc.AnyCompare(testMDepthUint1, testMDepthUint2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison.")

	// check : float in map
	testMDepthFloat1 := map[string]map[string]float64{"H": {"name": 1, "state": 2}}
	testMDepthFloat2 := map[string]map[string]float64{"H": {"name": 1, "state": 3}}
	retval, _ = strproc.AnyCompare(testMDepthFloat1, testMDepthFloat2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison.")

	// check : complex in map
	testMDepthComplex1 := map[string]map[string]complex64{"H": {"name": 1, "state": 2}}
	testMDepthComplex2 := map[string]map[string]complex64{"H": {"name": 1, "state": 3}}
	retval, _ = strproc.AnyCompare(testMDepthComplex1, testMDepthComplex2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison.")

	// check : different type
	testDiffType1 := int(32)
	testDiffType2 := int64(32)
	_, err = strproc.AnyCompare(testDiffType1, testDiffType2)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `Different Type`\nError : %v", err)

	// check : different len
	testDiffLen1 := []int{1, 2, 3, 4, 5, 6, 7}
	testDiffLen2 := []int{1, 2, 3}
	_, err = strproc.AnyCompare(testDiffLen1, testDiffLen2)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `Different Len`\nError : %v", err)

	// check : different len
	testDiffMapLen1 := map[int]string{0: "A", 1: "B", 2: "C"}
	testDiffMapLen2 := map[int]string{0: "A", 1: "B"}
	_, err = strproc.AnyCompare(testDiffMapLen1, testDiffMapLen2)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `Different Len`\nError : %v", err)

	// check : not support compre
	testDiffNotSupport1 := paddingTestVal{}
	testDiffNotSupport2 := paddingTestVal{}
	_, err = strproc.AnyCompare(testDiffNotSupport1, testDiffNotSupport2)
	assert.AssertNil(t, err, "Failure : Couldn't check the `Not Support Compre`\nError : %v", err)

	// check : sting != string
	testDiffrentStringMap1 := map[string]map[string]map[string]string{
		"F": {
			"name": {
				"first": "1",
				"last":  "2",
			},
		},
	}

	testDiffrentStringMap2 := map[string]map[string]map[string]int{
		"A": {
			"name": {
				"first": 11,
				"last":  21,
			},
		},
	}

	retval, _ = strproc.AnyCompare(testDiffrentStringMap1, testDiffrentStringMap2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison.")

	// check : uint != uint
	testDiffrentUintMap1 := map[string]map[string]map[string]uint{
		"F": {
			"name": {
				"first": 1,
				"last":  2,
			},
		},
	}

	testDiffrentUintMap2 := map[string]map[string]map[string]int{
		"A": {
			"name": {
				"first": 11,
				"last":  21,
			},
		},
	}

	retval, _ = strproc.AnyCompare(testDiffrentUintMap1, testDiffrentUintMap2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison.")

	// check : float64 != float64
	testDiffrentFloatMap1 := map[string]map[string]map[string]float64{
		"F": {
			"name": {
				"first": 1,
				"last":  2,
			},
		},
	}

	testDiffrentFloatMap2 := map[string]map[string]map[string]int{
		"A": {
			"name": {
				"first": 11,
				"last":  21,
			},
		},
	}

	retval, _ = strproc.AnyCompare(testDiffrentFloatMap1, testDiffrentFloatMap2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison.")

	// check : complex64 != complex64
	testDiffrentComplexMap1 := map[string]map[string]map[string]complex64{
		"F": {
			"name": {
				"first": 1,
				"last":  2,
			},
		},
	}

	testDiffrentComplexMap2 := map[string]map[string]map[string]int{
		"A": {
			"name": {
				"first": 11,
				"last":  21,
			},
		},
	}

	retval, _ = strproc.AnyCompare(testDiffrentComplexMap1, testDiffrentComplexMap2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison.")

	type testStruct1 struct {
		a int
		b int
	}

	testMapStruct1 := map[string]testStruct1{"a": {1, 2}}
	testMapStruct2 := map[string]testStruct1{"a": {1, 2}}
	retval, err = strproc.AnyCompare(testMapStruct1, testMapStruct2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison : %v", err)
}

func Test_strutils_DecodeUnicodeEntities(t *testing.T) {
	t.Parallel()

	str_oneline_ok := `안녕하세요.  방갑습니다.  감사합니다.  おはようございます こんにちは． こんばんは． おやすみなさい． ありがとうございます 你好 再見 谢谢!สวัสดีครับ แล้วเจอกันครับ ขอบคุณครับ Сайн байнауу`
	str_mutipleline_ok := `안녕하세요.
방갑습니다.
감사합니다.
おはようございます
こんにちは．
こんばんは．
おやすみなさい．
ありがとうございます
你好
再見
谢谢!สวัสดีครับ
แล้วเจอกันครับ
ขอบคุณครับ
Сайн байнауу`

	var retval string
	var err error

	tmpStrUnicodeEntityEncodedOneLine := "%uC548%uB155%uD558%uC138%uC694.%20%20%uBC29%uAC11%uC2B5%uB2C8%uB2E4.%20%20%uAC10%uC0AC%uD569%uB2C8%uB2E4.%20%20%u304A%u306F%u3088%u3046%u3054%u3056%u3044%u307E%u3059%20%u3053%u3093%u306B%u3061%u306F%uFF0E%20%u3053%u3093%u3070%u3093%u306F%uFF0E%20%u304A%u3084%u3059%u307F%u306A%u3055%u3044%uFF0E%20%u3042%u308A%u304C%u3068%u3046%u3054%u3056%u3044%u307E%u3059%20%u4F60%u597D%20%u518D%u898B%20%u8C22%u8C22%21%u0E2A%u0E27%u0E31%u0E2A%u0E14%u0E35%u0E04%u0E23%u0E31%u0E1A%20%u0E41%u0E25%u0E49%u0E27%u0E40%u0E08%u0E2D%u0E01%u0E31%u0E19%u0E04%u0E23%u0E31%u0E1A%20%u0E02%u0E2D%u0E1A%u0E04%u0E38%u0E13%u0E04%u0E23%u0E31%u0E1A%20%u0421%u0430%u0439%u043D%20%u0431%u0430%u0439%u043D%u0430%u0443%u0443"
	tmpStrUnicodeEntityEncodedMultipleLine := "%uC548%uB155%uD558%uC138%uC694.%0A%uBC29%uAC11%uC2B5%uB2C8%uB2E4.%0A%uAC10%uC0AC%uD569%uB2C8%uB2E4.%0A%u304A%u306F%u3088%u3046%u3054%u3056%u3044%u307E%u3059%0A%u3053%u3093%u306B%u3061%u306F%uFF0E%0A%u3053%u3093%u3070%u3093%u306F%uFF0E%0A%u304A%u3084%u3059%u307F%u306A%u3055%u3044%uFF0E%0A%u3042%u308A%u304C%u3068%u3046%u3054%u3056%u3044%u307E%u3059%0A%u4F60%u597D%0A%u518D%u898B%0A%u8C22%u8C22%21%u0E2A%u0E27%u0E31%u0E2A%u0E14%u0E35%u0E04%u0E23%u0E31%u0E1A%0A%u0E41%u0E25%u0E49%u0E27%u0E40%u0E08%u0E2D%u0E01%u0E31%u0E19%u0E04%u0E23%u0E31%u0E1A%0A%u0E02%u0E2D%u0E1A%u0E04%u0E38%u0E13%u0E04%u0E23%u0E31%u0E1A%0A%u0421%u0430%u0439%u043D%20%u0431%u0430%u0439%u043D%u0430%u0443%u0443"

	retval, err = strproc.DecodeUnicodeEntities(tmpStrUnicodeEntityEncodedOneLine)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, str_oneline_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_oneline_ok)

	retval, err = strproc.DecodeUnicodeEntities(tmpStrUnicodeEntityEncodedMultipleLine)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, str_mutipleline_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_mutipleline_ok)
}

func Test_strutils_DecodeURLEncoded(t *testing.T) {
	t.Parallel()

	url_ok := "https://www.google.com/search?source=hp&ei=ChM0W462AYbs8wXAwIaQCw&q=대한민국&oq=대한민국&gs_l=psy-ab.3..0i131k1l3j0l7.930.2247.0.2376.12.11.0.0.0.0.116.955.10j1.11.0....0...1.1j4.64.psy-ab..2.10.874.0...0.2DIo94YBWPI"

	tmpJSURLEncode := "https%3A%2F%2Fwww.google.com%2Fsearch%3Fsource%3Dhp%26ei%3DChM0W462AYbs8wXAwIaQCw%26q%3D%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD%26oq%3D%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD%26gs_l%3Dpsy-ab.3..0i131k1l3j0l7.930.2247.0.2376.12.11.0.0.0.0.116.955.10j1.11.0....0...1.1j4.64.psy-ab..2.10.874.0...0.2DIo94YBWPI"
	tmpJSEncodeURIComponent := "https%3A%2F%2Fwww.google.com%2Fsearch%3Fsource%3Dhp%26ei%3DChM0W462AYbs8wXAwIaQCw%26q%3D%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD%26oq%3D%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD%26gs_l%3Dpsy-ab.3..0i131k1l3j0l7.930.2247.0.2376.12.11.0.0.0.0.116.955.10j1.11.0....0...1.1j4.64.psy-ab..2.10.874.0...0.2DIo94YBWPI"
	tmpJSencodeURI := "https://www.google.com/search?source=hp&ei=ChM0W462AYbs8wXAwIaQCw&q=%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD&oq=%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD&gs_l=psy-ab.3..0i131k1l3j0l7.930.2247.0.2376.12.11.0.0.0.0.116.955.10j1.11.0....0...1.1j4.64.psy-ab..2.10.874.0...0.2DIo94YBWPI"
	tmpJSEscape := "https%3A//www.google.com/search%3Fsource%3Dhp%26ei%3DChM0W462AYbs8wXAwIaQCw%26q%3D%uB300%uD55C%uBBFC%uAD6D%26oq%3D%uB300%uD55C%uBBFC%uAD6D%26gs_l%3Dpsy-ab.3..0i131k1l3j0l7.930.2247.0.2376.12.11.0.0.0.0.116.955.10j1.11.0....0...1.1j4.64.psy-ab..2.10.874.0...0.2DIo94YBWPI"
	tmpPHPurlEncode := "https%3A%2F%2Fwww.google.com%2Fsearch%3Fsource%3Dhp%26ei%3DChM0W462AYbs8wXAwIaQCw%26q%3D%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD%26oq%3D%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD%26gs_l%3Dpsy-ab.3..0i131k1l3j0l7.930.2247.0.2376.12.11.0.0.0.0.116.955.10j1.11.0....0...1.1j4.64.psy-ab..2.10.874.0...0.2DIo94YBWPI"
	tmpPHPrawurlEncode := "https%3A%2F%2Fwww.google.com%2Fsearch%3Fsource%3Dhp%26ei%3DChM0W462AYbs8wXAwIaQCw%26q%3D%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD%26oq%3D%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD%26gs_l%3Dpsy-ab.3..0i131k1l3j0l7.930.2247.0.2376.12.11.0.0.0.0.116.955.10j1.11.0....0...1.1j4.64.psy-ab..2.10.874.0...0.2DIo94YBWPI"

	tmpPyURLlibQuote := "https%3A//www.google.com/search%3Fsource%3Dhp%26ei%3DChM0W462AYbs8wXAwIaQCw%26q%3D%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD%26oq%3D%EB%8C%80%ED%95%9C%EB%AF%BC%EA%B5%AD%26gs_l%3Dpsy-ab.3..0i131k1l3j0l7.930.2247.0.2376.12.11.0.0.0.0.116.955.10j1.11.0....0...1.1j4.64.psy-ab..2.10.874.0...0.2DIo94YBWPI"

	url_with_japan_world_ok := "http://hello.世界.com/foo"
	url_with_japna_keyword_ok := "https://www.google.co.kr/search?q=アパルトヘイトで世界せかいから孤立していた南アフリカ共和国には多くの企業が進出し、以前から比較的密接な関係を築いていた。&oq=アパルトヘイトで世界せかいから孤立していた南アフリカ共和国には多くの企業が進出し、以前から比較的密接な関係を築いていた。&aqs=chrome..69i57.417j0j4&sourceid=chrome&ie=UTF-8"

	tmpURLWithJapanWorld := "http://hello.%E4%B8%96%E7%95%8C.com/foo"
	tmpURLWithJapanKeyword := "https://www.google.co.kr/search?q=%E3%82%A2%E3%83%91%E3%83%AB%E3%83%88%E3%83%98%E3%82%A4%E3%83%88%E3%81%A7%E4%B8%96%E7%95%8C%E3%81%9B%E3%81%8B%E3%81%84%E3%81%8B%E3%82%89%E5%AD%A4%E7%AB%8B%E3%81%97%E3%81%A6%E3%81%84%E3%81%9F%E5%8D%97%E3%82%A2%E3%83%95%E3%83%AA%E3%82%AB%E5%85%B1%E5%92%8C%E5%9B%BD%E3%81%AB%E3%81%AF%E5%A4%9A%E3%81%8F%E3%81%AE%E4%BC%81%E6%A5%AD%E3%81%8C%E9%80%B2%E5%87%BA%E3%81%97%E3%80%81%E4%BB%A5%E5%89%8D%E3%81%8B%E3%82%89%E6%AF%94%E8%BC%83%E7%9A%84%E5%AF%86%E6%8E%A5%E3%81%AA%E9%96%A2%E4%BF%82%E3%82%92%E7%AF%89%E3%81%84%E3%81%A6%E3%81%84%E3%81%9F%E3%80%82&oq=%E3%82%A2%E3%83%91%E3%83%AB%E3%83%88%E3%83%98%E3%82%A4%E3%83%88%E3%81%A7%E4%B8%96%E7%95%8C%E3%81%9B%E3%81%8B%E3%81%84%E3%81%8B%E3%82%89%E5%AD%A4%E7%AB%8B%E3%81%97%E3%81%A6%E3%81%84%E3%81%9F%E5%8D%97%E3%82%A2%E3%83%95%E3%83%AA%E3%82%AB%E5%85%B1%E5%92%8C%E5%9B%BD%E3%81%AB%E3%81%AF%E5%A4%9A%E3%81%8F%E3%81%AE%E4%BC%81%E6%A5%AD%E3%81%8C%E9%80%B2%E5%87%BA%E3%81%97%E3%80%81%E4%BB%A5%E5%89%8D%E3%81%8B%E3%82%89%E6%AF%94%E8%BC%83%E7%9A%84%E5%AF%86%E6%8E%A5%E3%81%AA%E9%96%A2%E4%BF%82%E3%82%92%E7%AF%89%E3%81%84%E3%81%A6%E3%81%84%E3%81%9F%E3%80%82&aqs=chrome..69i57.417j0j4&sourceid=chrome&ie=UTF-8"

	str_ok := `안녕하세요.
방갑습니다.
감사합니다.
おはようございます
こんにちは．
こんばんは．
おやすみなさい．
ありがとうございます
你好
再見
谢谢!สวัสดีครับ
แล้วเจอกันครับ
ขอบคุณครับ
Сайн байнауу`

	tmpStrUnicodeEntityEncoded := "%uC548%uB155%uD558%uC138%uC694.%0A%uBC29%uAC11%uC2B5%uB2C8%uB2E4.%0A%uAC10%uC0AC%uD569%uB2C8%uB2E4.%0A%u304A%u306F%u3088%u3046%u3054%u3056%u3044%u307E%u3059%0A%u3053%u3093%u306B%u3061%u306F%uFF0E%0A%u3053%u3093%u3070%u3093%u306F%uFF0E%0A%u304A%u3084%u3059%u307F%u306A%u3055%u3044%uFF0E%0A%u3042%u308A%u304C%u3068%u3046%u3054%u3056%u3044%u307E%u3059%0A%u4F60%u597D%0A%u518D%u898B%0A%u8C22%u8C22%21%u0E2A%u0E27%u0E31%u0E2A%u0E14%u0E35%u0E04%u0E23%u0E31%u0E1A%0A%u0E41%u0E25%u0E49%u0E27%u0E40%u0E08%u0E2D%u0E01%u0E31%u0E19%u0E04%u0E23%u0E31%u0E1A%0A%u0E02%u0E2D%u0E1A%u0E04%u0E38%u0E13%u0E04%u0E23%u0E31%u0E1A%0A%u0421%u0430%u0439%u043D%20%u0431%u0430%u0439%u043D%u0430%u0443%u0443"

	tmpStrEncodedJSURIComponent := "%EC%95%88%EB%85%95%ED%95%98%EC%84%B8%EC%9A%94.%0A%EB%B0%A9%EA%B0%91%EC%8A%B5%EB%8B%88%EB%8B%A4.%0A%EA%B0%90%EC%82%AC%ED%95%A9%EB%8B%88%EB%8B%A4.%0A%E3%81%8A%E3%81%AF%E3%82%88%E3%81%86%E3%81%94%E3%81%96%E3%81%84%E3%81%BE%E3%81%99%0A%E3%81%93%E3%82%93%E3%81%AB%E3%81%A1%E3%81%AF%EF%BC%8E%0A%E3%81%93%E3%82%93%E3%81%B0%E3%82%93%E3%81%AF%EF%BC%8E%0A%E3%81%8A%E3%82%84%E3%81%99%E3%81%BF%E3%81%AA%E3%81%95%E3%81%84%EF%BC%8E%0A%E3%81%82%E3%82%8A%E3%81%8C%E3%81%A8%E3%81%86%E3%81%94%E3%81%96%E3%81%84%E3%81%BE%E3%81%99%0A%E4%BD%A0%E5%A5%BD%0A%E5%86%8D%E8%A6%8B%0A%E8%B0%A2%E8%B0%A2!%E0%B8%AA%E0%B8%A7%E0%B8%B1%E0%B8%AA%E0%B8%94%E0%B8%B5%E0%B8%84%E0%B8%A3%E0%B8%B1%E0%B8%9A%0A%E0%B9%81%E0%B8%A5%E0%B9%89%E0%B8%A7%E0%B9%80%E0%B8%88%E0%B8%AD%E0%B8%81%E0%B8%B1%E0%B8%99%E0%B8%84%E0%B8%A3%E0%B8%B1%E0%B8%9A%0A%E0%B8%82%E0%B8%AD%E0%B8%9A%E0%B8%84%E0%B8%B8%E0%B8%93%E0%B8%84%E0%B8%A3%E0%B8%B1%E0%B8%9A%0A%D0%A1%D0%B0%D0%B9%D0%BD%20%D0%B1%D0%B0%D0%B9%D0%BD%D0%B0%D1%83%D1%83"

	var retval, valid_ok string
	var err error
	retval, err = strproc.DecodeURLEncoded(tmpJSURLEncode)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, url_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, url_ok)

	retval, err = strproc.DecodeURLEncoded(tmpJSEncodeURIComponent)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, url_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, url_ok)

	retval, err = strproc.DecodeURLEncoded(tmpJSencodeURI)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, url_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, url_ok)

	retval, err = strproc.DecodeURLEncoded(tmpJSEscape)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, url_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, url_ok)

	retval, err = strproc.DecodeURLEncoded(tmpPHPurlEncode)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, url_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, url_ok)

	retval, err = strproc.DecodeURLEncoded(tmpPHPrawurlEncode)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, url_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, url_ok)

	retval, err = strproc.DecodeURLEncoded(tmpPyURLlibQuote)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, url_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, url_ok)

	retval, err = strproc.DecodeURLEncoded(tmpURLWithJapanWorld)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, url_with_japan_world_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, url_with_japan_world_ok)

	retval, err = strproc.DecodeURLEncoded(tmpURLWithJapanKeyword)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, url_with_japna_keyword_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, url_with_japna_keyword_ok)

	retval, err = strproc.DecodeURLEncoded(tmpStrUnicodeEntityEncoded)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, str_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_ok)

	retval, err = strproc.DecodeURLEncoded(tmpStrEncodedJSURIComponent)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, str_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_ok)

	tmpW3schoolsAsciiEncodingReferenceFromWindows1252 := "%20 %21 %22 %23 %24 %25 %26 %27 %28 %29 %2A %2B %2C %2D %2E %2F %30 %31 %32 %33 %34 %35 %36 %37 %38 %39 %3A %3B %3C %3D %3E %3F %40 %41 %42 %43 %44 %45 %46 %47 %48 %49 %4A %4B %4C %4D %4E %4F %50 %51 %52 %53 %54 %55 %56 %57 %58 %59 %5A %5B %5C %5D %5E %5F %60 %61 %62 %63 %64 %65 %66 %67 %68 %69 %6A %6B %6C %6D %6E %6F %70 %71 %72 %73 %74 %75 %76 %77 %78 %79 %7A %7B %7C %7D %7E %7F %80 %81 %82 %83 %84 %85 %86 %87 %88 %89 %8A %8B %8C %8D %8E %8F %90 %91 %92 %93 %94 %95 %96 %97 %98 %99 %9A %9B %9C %9D %9E %9F %A0 %A1 %A2 %A3 %A4 %A5 %A6 %A7 %A8 %A9 %AA %AB %AC %AD %AE %AF %B0 %B1 %B2 %B3 %B4 %B5 %B6 %B7 %B8 %B9 %BA %BB %BC %BD %BE %BF %C0 %C1 %C2 %C3 %C4 %C5 %C6 %C7 %C8 %C9 %CA %CB %CC %CD %CE %CF %D0 %D1 %D2 %D3 %D4 %D5 %D6 %D7 %D8 %D9 %DA %DB %DC %DD %DE %DF %E0 %E1 %E2 %E3 %E4 %E5 %E6 %E7 %E8 %E9 %EA %EB %EC %ED %EE %EF %F0 %F1 %F2 %F3 %F4 %F5 %F6 %F7 %F8 %F9 %FA %FB %FC %FD %FE %FF"

	tmpW3schoolsAsciiEncodingReferenceFromUTF8 := "%20 %21 %22 %23 %24 %25 %26 %27 %28 %29 %2A %2B %2C %2D %2E %2F %30 %31 %32 %33 %34 %35 %36 %37 %38 %39 %3A %3B %3C %3D %3E %3F %40 %41 %42 %43 %44 %45 %46 %47 %48 %49 %4A %4B %4C %4D %4E %4F %50 %51 %52 %53 %54 %55 %56 %57 %58 %59 %5A %5B %5C %5D %5E %5F %60 %61 %62 %63 %64 %65 %66 %67 %68 %69 %6A %6B %6C %6D %6E %6F %70 %71 %72 %73 %74 %75 %76 %77 %78 %79 %7A %7B %7C %7D %7E %7F %E2%82%AC %81 %E2%80%9A %C6%92 %E2%80%9E %E2%80%A6 %E2%80%A0 %E2%80%A1 %CB%86 %E2%80%B0 %C5%A0 %E2%80%B9 %C5%92 %C5%8D %C5%BD %8F %C2%90 %E2%80%98 %E2%80%99 %E2%80%9C %E2%80%9D %E2%80%A2 %E2%80%93 %E2%80%94 %CB%9C %E2%84 %C5%A1 %E2%80 %C5%93 %9D %C5%BE %C5%B8 %C2%A0 %C2%A1 %C2%A2 %C2%A3 %C2%A4 %C2%A5 %C2%A6 %C2%A7 %C2%A8 %C2%A9 %C2%AA %C2%AB %C2%AC %C2%AD %C2%AE %C2%AF %C2%B0 %C2%B1 %C2%B2 %C2%B3 %C2%B4 %C2%B5 %C2%B6 %C2%B7 %C2%B8 %C2%B9 %C2%BA %C2%BB %C2%BC %C2%BD %C2%BE %C2%BF %C3%80 %C3%81 %C3%82 %C3%83 %C3%84 %C3%85 %C3%86 %C3%87 %C3%88 %C3%89 %C3%8A %C3%8B %C3%8C %C3%8D %C3%8E %C3%8F %C3%90 %C3%91 %C3%92 %C3%93 %C3%94 %C3%95 %C3%96 %C3%97 %C3%98 %C3%99 %C3%9A %C3%9B %C3%9C %C3%9D %C3%9E %C3%9F %C3%A0 %C3%A1 %C3%A2 %C3%A3 %C3%A4 %C3%A5 %C3%A6 %C3%A7 %C3%A8 %C3%A9 %C3%AA %C3%AB %C3%AC %C3%AD %C3%AE %C3%AF %C3%B0 %C3%B1 %C3%B2 %C3%B3 %C3%B4 %C3%B5 %C3%B6 %C3%B7 %C3%B8 %C3%B9 %C3%BA %C3%BB %C3%BC %C3%BD %C3%BE %C3%BF"

	retval, err = strproc.DecodeURLEncoded(tmpW3schoolsAsciiEncodingReferenceFromWindows1252)
	assert.AssertNil(t, err, "Error : %v", err)
	valid_ok, err = url.QueryUnescape(tmpW3schoolsAsciiEncodingReferenceFromWindows1252)

	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, valid_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, valid_ok)

	retval, err = strproc.DecodeURLEncoded(tmpW3schoolsAsciiEncodingReferenceFromUTF8)
	assert.AssertNil(t, err, "Error : %v", err)
	valid_ok, err = url.QueryUnescape(tmpW3schoolsAsciiEncodingReferenceFromUTF8)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, valid_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, valid_ok)
}

func Test_strutils_StripTags(t *testing.T) {
	t.Parallel()

	str_ok := `
Just! a String Processing Library for Go-lang
Just! a String Processing Library for Go-lang
Just a few methods for helping processing and validation the string
View on GitHub
Just! a String Processing Library for Go-lang
Just a few methods for helping processing the string
README.md haven’t contain all the examples. Please refer to the the XXXtest.go files.
`

	str_original_html := `
<!DOCTYPE html>
<html lang="en-us">
<head>
<meta charset="UTF-8">
<title>                            Just! a String Processing Library for Go-lang</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#157878">
<link href='https://fonts.googleapis.com/css?family=Open+Sans:400,700' rel='stylesheet' type='text/css'>
<link rel="stylesheet" href="/go-strutil/assets/css/style.css?v=dae229423409070462d2ce364eba3b5721930df0">
</head>
<body>
<section class="page-header">
<h1 class="project-name">Just! a String Processing Library for Go-lang</h1>
<h2 class="project-tagline">Just a few methods for helping processing and validation the string</h2>
<a href="https://github.com/torden/go-strutil" class="btn">View on GitHub</a>
</section>
<section class="main-content">
<h1 id="just-a-string-processing-library-for-go-lang">Just! a String Processing Library for Go-lang</h1>
<p>Just a few methods for helping processing the string</p>
<p>README.md haven’t contain all the examples. Please refer to the the XXXtest.go files.</p>
</body>
</html>
`

	str_html_entity_encoded := `
&lt;!DOCTYPE html&gt;
&lt;html lang=&quot;en-us&quot;&gt;
&lt;head&gt;
&lt;meta charset=&quot;UTF-8&quot;&gt;
&lt;title&gt;                            Just! a String Processing Library for Go-lang&lt;/title&gt;
&lt;meta name=&quot;viewport&quot; content=&quot;width=device-width, initial-scale=1&quot;&gt;
&lt;meta name=&quot;theme-color&quot; content=&quot;#157878&quot;&gt;
&lt;link href='https://fonts.googleapis.com/css?family=Open+Sans:400,700' rel='stylesheet' type='text/css'&gt;
&lt;link rel=&quot;stylesheet&quot; href=&quot;/go-strutil/assets/css/style.css?v=dae229423409070462d2ce364eba3b5721930df0&quot;&gt;
&lt;/head&gt;
&lt;body&gt;
&lt;section class=&quot;page-header&quot;&gt;
&lt;h1 class=&quot;project-name&quot;&gt;Just! a String Processing Library for Go-lang&lt;/h1&gt;
&lt;h2 class=&quot;project-tagline&quot;&gt;Just a few methods for helping processing and validation the string&lt;/h2&gt;
&lt;a href=&quot;https://github.com/torden/go-strutil&quot; class=&quot;btn&quot;&gt;View on GitHub&lt;/a&gt;
&lt;/section&gt;
&lt;section class=&quot;main-content&quot;&gt;
&lt;h1 id=&quot;just-a-string-processing-library-for-go-lang&quot;&gt;Just! a String Processing Library for Go-lang&lt;/h1&gt;
&lt;p&gt;Just a few methods for helping processing the string&lt;/p&gt;
&lt;p&gt;README.md haven&rsquo;t contain all the examples. Please refer to the the XXXtest.go files.&lt;/p&gt;
&lt;/body&gt;
&lt;/html&gt
`
	str_html_urlencoded := `
%3C%21DOCTYPE+html%3E%0A%3Chtml+lang%3D%22en-us%22%3E%0A%3Chead%3E%0A%3Cmeta+charset%3D%22UTF-8%22%3E%0A%3Ctitle%3E++++++++++++++++++++++++++++Just%21+a+String+Processing+Library+for+Go-lang%3C%2Ftitle%3E%0A%3Cmeta+name%3D%22viewport%22+content%3D%22width%3Ddevice-width%2C+initial-scale%3D1%22%3E%0A%3Cmeta+name%3D%22theme-color%22+content%3D%22%23157878%22%3E%0A%3Clink+href%3D%27https%3A%2F%2Ffonts.googleapis.com%2Fcss%3Ffamily%3DOpen%2BSans%3A400%2C700%27+rel%3D%27stylesheet%27+type%3D%27text%2Fcss%27%3E%0A%3Clink+rel%3D%22stylesheet%22+href%3D%22%2Fgo-strutil%2Fassets%2Fcss%2Fstyle.css%3Fv%3Ddae229423409070462d2ce364eba3b5721930df0%22%3E%0A%3C%2Fhead%3E%0A%3Cbody%3E%0A%3Csection+class%3D%22page-header%22%3E%0A%3Ch1+class%3D%22project-name%22%3EJust%21+a+String+Processing+Library+for+Go-lang%3C%2Fh1%3E%0A%3Ch2+class%3D%22project-tagline%22%3EJust+a+few+methods+for+helping+processing+and+validation+the+string%3C%2Fh2%3E%0A%3Ca+href%3D%22https%3A%2F%2Fgithub.com%2Ftorden%2Fgo-strutil%22+class%3D%22btn%22%3EView+on+GitHub%3C%2Fa%3E%0A%3C%2Fsection%3E%0A%3Csection+class%3D%22main-content%22%3E%0A%3Ch1+id%3D%22just-a-string-processing-library-for-go-lang%22%3EJust%21+a+String+Processing+Library+for+Go-lang%3C%2Fh1%3E%0A%3Cp%3EJust+a+few+methods+for+helping+processing+the+string%3C%2Fp%3E%0A%3Cp%3EREADME.md+haven%E2%80%99t+contain+all+the+examples.+Please+refer+to+the+the+XXXtest.go+files.%3C%2Fp%3E%0A%3C%2Fbody%3E%0A%3C%2Fhtml%3E
`

	// html entity encoded after url encoded
	str_htmentity_urlencoded := `
%3C%21DOCTYPE+html%3E%0A%3Chtml+lang%3D%22en-us%22%3E%0A%3Chead%3E%0A%3Cmeta+charset%3D%22UTF-8%22%3E%0A%3Ctitle%3E++++++++++++++++++++++++++++Just%21+a+String+Processing+Library+for+Go-lang%3C%2Ftitle%3E%0A%3Cmeta+name%3D%22viewport%22+content%3D%22width%3Ddevice-width%2C+initial-scale%3D1%22%3E%0A%3Cmeta+name%3D%22theme-color%22+content%3D%22%23157878%22%3E%0A%3Clink+href%3D%27https%3A%2F%2Ffonts.googleapis.com%2Fcss%3Ffamily%3DOpen%2BSans%3A400%2C700%27+rel%3D%27stylesheet%27+type%3D%27text%2Fcss%27%3E%0A%3Clink+rel%3D%22stylesheet%22+href%3D%22%2Fgo-strutil%2Fassets%2Fcss%2Fstyle.css%3Fv%3Ddae229423409070462d2ce364eba3b5721930df0%22%3E%0A%3C%2Fhead%3E%0A%3Cbody%3E%0A%3Csection+class%3D%22page-header%22%3E%0A%3Ch1+class%3D%22project-name%22%3EJust%21+a+String+Processing+Library+for+Go-lang%3C%2Fh1%3E%0A%3Ch2+class%3D%22project-tagline%22%3EJust+a+few+methods+for+helping+processing+and+validation+the+string%3C%2Fh2%3E%0A%3Ca+href%3D%22https%3A%2F%2Fgithub.com%2Ftorden%2Fgo-strutil%22+class%3D%22btn%22%3EView+on+GitHub%3C%2Fa%3E%0A%3C%2Fsection%3E%0A%3Csection+class%3D%22main-content%22%3E%0A%3Ch1+id%3D%22just-a-string-processing-library-for-go-lang%22%3EJust%21+a+String+Processing+Library+for+Go-lang%3C%2Fh1%3E%0A%3Cp%3EJust+a+few+methods+for+helping+processing+the+string%3C%2Fp%3E%0A%3Cp%3EREADME.md+haven%E2%80%99t+contain+all+the+examples.+Please+refer+to+the+the+XXXtest.go+files.%3C%2Fp%3E%0A%3C%2Fbody%3E%0A%3C%2Fhtml%3E
`

	// url encoded adter html entity encoded
	str_urlencoded_htmlentity := `
%26lt%3B%21DOCTYPE+html%26gt%3B%0A%26lt%3Bhtml+lang%3D%26quot%3Ben-us%26quot%3B%26gt%3B%0A%26lt%3Bhead%26gt%3B%0A%26lt%3Bmeta+charset%3D%26quot%3BUTF-8%26quot%3B%26gt%3B%0A%26lt%3Btitle%26gt%3B++++++++++++++++++++++++++++Just%21+a+String+Processing+Library+for+Go-lang%26lt%3B%2Ftitle%26gt%3B%0A%26lt%3Bmeta+name%3D%26quot%3Bviewport%26quot%3B+content%3D%26quot%3Bwidth%3Ddevice-width%2C+initial-scale%3D1%26quot%3B%26gt%3B%0A%26lt%3Bmeta+name%3D%26quot%3Btheme-color%26quot%3B+content%3D%26quot%3B%23157878%26quot%3B%26gt%3B%0A%26lt%3Blink+href%3D%27https%3A%2F%2Ffonts.googleapis.com%2Fcss%3Ffamily%3DOpen%2BSans%3A400%2C700%27+rel%3D%27stylesheet%27+type%3D%27text%2Fcss%27%26gt%3B%0A%26lt%3Blink+rel%3D%26quot%3Bstylesheet%26quot%3B+href%3D%26quot%3B%2Fgo-strutil%2Fassets%2Fcss%2Fstyle.css%3Fv%3Ddae229423409070462d2ce364eba3b5721930df0%26quot%3B%26gt%3B%0A%26lt%3B%2Fhead%26gt%3B%0A%26lt%3Bbody%26gt%3B%0A%26lt%3Bsection+class%3D%26quot%3Bpage-header%26quot%3B%26gt%3B%0A%26lt%3Bh1+class%3D%26quot%3Bproject-name%26quot%3B%26gt%3BJust%21+a+String+Processing+Library+for+Go-lang%26lt%3B%2Fh1%26gt%3B%0A%26lt%3Bh2+class%3D%26quot%3Bproject-tagline%26quot%3B%26gt%3BJust+a+few+methods+for+helping+processing+and+validation+the+string%26lt%3B%2Fh2%26gt%3B%0A%26lt%3Ba+href%3D%26quot%3Bhttps%3A%2F%2Fgithub.com%2Ftorden%2Fgo-strutil%26quot%3B+class%3D%26quot%3Bbtn%26quot%3B%26gt%3BView+on+GitHub%26lt%3B%2Fa%26gt%3B%0A%26lt%3B%2Fsection%26gt%3B%0A%26lt%3Bsection+class%3D%26quot%3Bmain-content%26quot%3B%26gt%3B%0A%26lt%3Bh1+id%3D%26quot%3Bjust-a-string-processing-library-for-go-lang%26quot%3B%26gt%3BJust%21+a+String+Processing+Library+for+Go-lang%26lt%3B%2Fh1%26gt%3B%0A%26lt%3Bp%26gt%3BJust+a+few+methods+for+helping+processing+the+string%26lt%3B%2Fp%26gt%3B%0A%26lt%3Bp%26gt%3BREADME.md+haven%26rsquo%3Bt+contain+all+the+examples.+Please+refer+to+the+the+XXXtest.go+files.%26lt%3B%2Fp%26gt%3B%0A%26lt%3B%2Fbody%26gt%3B%0A%26lt%3B%2Fhtml%26gt%3B
`

	// similar html entity tag
	str_infinity_loop_test1 := `ABC PERC%NT DEF`

	// fb url encoded
	str_fb_urlencoded := `https://m.facebook.com/story.php?story_fbid=2566258156763959&id=636551103068017&ref=page_internal&_ft_=mf_story_key.2566258156763959%3Atop_level_post_id.2566258156763959%3Atl_objid.2566258156763959%3Acontent_owner_id_new.636551103068017%3Athrowback_story_fbid.2566258156763959%3Apage_id.636551103068017%3Aphoto_attachments_list.%5B2566257123430729%2C2566257250097383%2C2566257386764036%2C2566257496764025%5D%3Astory_location.4%3Astory_attachment_style.album%3Apage_insights.%7B%22636551103068017%22%3A%7B%22page_id%22%3A636551103068017%2C%22actor_id%`

	str_fb_urlencoded_ok := `https://m.facebook.com/story.php?story_fbid=2566258156763959&id=636551103068017&ref=page_internal&_ft_=mf_story_key.2566258156763959:top_level_post_id.2566258156763959:tl_objid.2566258156763959:content_owner_id_new.636551103068017:throwback_story_fbid.2566258156763959:page_id.636551103068017:photo_attachments_list.[2566257123430729,2566257250097383,2566257386764036,2566257496764025]:story_location.4:story_attachment_style.album:page_insights.{"636551103068017":{"page_id":636551103068017,"actor_id%`

	var retval string
	var err error

	// check : original html
	retval, err = strproc.StripTags(str_original_html)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, str_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_ok)

	// check : html entity encoded html
	retval, err = strproc.StripTags(str_html_entity_encoded)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, str_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_ok)

	// check : url encoded html
	retval, err = strproc.StripTags(str_html_urlencoded)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, str_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_ok)

	// check : html entity encoded after url encoded
	retval, err = strproc.StripTags(str_htmentity_urlencoded)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, str_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_ok)

	// check : url encoded after html entity encoded
	retval, err = strproc.StripTags(str_urlencoded_htmlentity)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, str_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_ok)

	// check : #77 - infinity loop
	retval, err = strproc.StripTags(str_infinity_loop_test1)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, str_infinity_loop_test1, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_infinity_loop_test1)

	// check : singular point
	retval, err = strproc.StripTags(str_fb_urlencoded)
	assert.AssertNil(t, err, "Error : %v", err)
	assert.AssertEquals(t, retval, str_fb_urlencoded_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_fb_urlencoded_ok)
}

func Test_strutils_ConvertToStr(t *testing.T) {
	t.Parallel()

	dataset := map[interface{}]string{
		string("1234567"): "1234567",
		int(1):            "1",
		int8(1):           "1",
		int16(256):        "256",
		int32(256):        "256",
		int64(1234567):    "1234567",
		uint(1):           "1",
		uint8(1):          "1",
		uint16(256):       "256",
		uint32(256):       "256",
		uint64(1234567):   "1234567",
		float32(12):       "12",
		float64(12):       "12",
		complex64(12):     "(12+0i)",
		complex128(12):    "(12+0i)",

		int(-1):         "-1",
		int8(-1):        "-1",
		int16(-256):     "-256",
		int32(-256):     "-256",
		int64(-1234567): "-1234567",
		float32(-12):    "-12",
		float64(-12):    "-12",
		complex64(-12):  "(-12+0i)",
		complex128(-12): "(-12+0i)",

		float32(12.1):    "12.1",
		float64(12.1):    "12.1",
		complex64(12.1):  "(12.1+0i)",
		complex128(12.1): "(12.1+0i)",

		float32(-12.1):    "-12.1",
		float64(-12.1):    "-12.1",
		complex64(-12.1):  "(-12.1+0i)",
		complex128(-12.1): "(-12.1+0i)",

		true:  "true",
		false: "false",
	}

	// check : common
	for k, v := range dataset {
		retval, err := strproc.ConvertToStr(k)
		assert.AssertEquals(t, retval, v, "Return Value mismatch.\nExpected: %v\nActual: %v\nError : %v", retval, v, err)
	}
}

func Test_strutils_ConvertToArByte(t *testing.T) {
	t.Parallel()

	dataset := map[interface{}]string{
		string("1234567"): "1234567",
		int(1):            "1",
		int8(1):           "1",
		int16(256):        "256",
		int32(256):        "256",
		int64(1234567):    "1234567",
		uint(1):           "1",
		uint8(1):          "1",
		uint16(256):       "256",
		uint32(256):       "256",
		uint64(1234567):   "1234567",
		float32(12):       "12",
		float64(12):       "12",
		complex64(12):     "(12+0i)",
		complex128(12):    "(12+0i)",

		int(-1):         "-1",
		int8(-1):        "-1",
		int16(-256):     "-256",
		int32(-256):     "-256",
		int64(-1234567): "-1234567",
		float32(-12):    "-12",
		float64(-12):    "-12",
		complex64(-12):  "(-12+0i)",
		complex128(-12): "(-12+0i)",

		float32(12.1):    "12.1",
		float64(12.1):    "12.1",
		complex64(12.1):  "(12.1+0i)",
		complex128(12.1): "(12.1+0i)",

		float32(-12.1):    "-12.1",
		float64(-12.1):    "-12.1",
		complex64(-12.1):  "(-12.1+0i)",
		complex128(-12.1): "(-12.1+0i)",

		true:  "true",
		false: "false",
	}

	for k, v := range dataset {
		retval, err := strproc.ConvertToArByte(k)
		assert.AssertEquals(t, string(retval), v, "Return Value mismatch.\nExpected: %v\nActual: %v\nError : %v", retval, v, err)
	}
}

func Test_strutils_ReverseStr(t *testing.T) {
	t.Parallel()

	dataset := map[string]string{
		"0123456789": "9876543210",
		"가나다라마바사":    "사바마라다나가",
		"あいうえお":      "おえういあ",
		"天地玄黃宇宙洪荒":   "荒洪宙宇黃玄地天",
	}

	// check : common
	for k, v := range dataset {
		retval := strproc.ReverseStr(k)
		assert.AssertEquals(t, retval, v, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func Test_strutils_ReverseNormalStr(t *testing.T) {
	t.Parallel()

	dataset := map[string]string{
		"0123456789": "9876543210",
		"abcdefg":    "gfedcba",
	}

	// check : common
	for k, v := range dataset {
		retval := strproc.ReverseNormalStr(k)
		assert.AssertEquals(t, retval, v, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func Test_strutils_ReverseReverseUnicode(t *testing.T) {
	t.Parallel()

	dataset := map[string]string{
		"0123456789": "9876543210",
		"abcdefg":    "gfedcba",
		"가나다라마바사":    "사바마라다나가",
		"あいうえお":      "おえういあ",
		"天地玄黃宇宙洪荒":   "荒洪宙宇黃玄地天",
	}

	// check : common
	for k, v := range dataset {
		retval := strproc.ReverseUnicode(k)
		assert.AssertEquals(t, retval, v, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func Test_strutils_FileMD5Hash(t *testing.T) {
	t.Parallel()

	var retval string
	var err error

	// check : common
	retval, err = strproc.FileMD5Hash("./LICENSE")
	assert.AssertNil(t, err, "Error : %v", err)

	str_ok := "64e17a4e1c96bbfce57ab19cd0153e6a"
	assert.AssertEquals(t, retval, str_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_ok)

	// check : os.Open
	_, err = strproc.FileMD5Hash("./HELLO_GOLANG")
	assert.AssertNotNil(t, err, "Couldn't check the `os.Open`\nError : %v", err)
}

func Test_strutils_MD5Hash(t *testing.T) {
	t.Parallel()

	dataset := map[string]string{
		"0123456789": "781e5e245d69b566979b86e28d23f2c7",
		"abcdefg":    "7ac66c0f148de9519b8bd264312c4d64",
		"abcdefgqwdoisef;oijawe;fijq2039jdfs.dnc;oa283hr08uj3o;ijwaef;owhjefo;uhwefwef": "15f764f21d09b11102eb015fc8824d00",
	}

	// check : common
	for k, v := range dataset {
		retval, err := strproc.MD5Hash(k)

		assert.AssertNil(t, err, "Error : %v", err)
		assert.AssertEquals(t, retval, v, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func Test_strutils_RegExpNamedGroups(t *testing.T) {
	t.Parallel()

	var ok bool

	// refer : https://golang.org/doc/devel/release.html#policy
	regexGoVersion := regexp.MustCompile(`go(?P<major>([0-9]{1,3}))\.(?P<minor>([0-9]{1,3}))(\.(?P<rev>([0-9]{1,3})))?`)

	verdic, err := strproc.RegExpNamedGroups(regexGoVersion, runtime.Version())
	assert.AssertNil(t, err, "Error : %v", err)

	_, ok = verdic["major"]
	assert.AssertTrue(t, ok, "Not Exists Major ver. in Return Value")

	_, ok = verdic["minor"]
	assert.AssertTrue(t, ok, "Not Exists Minor ver. in Return Value")
}
