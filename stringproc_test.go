package strutils_test

import (
	"io/ioutil"
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/torden/go-strutil"
)

func Test_strutils_AddSlashes(t *testing.T) {

	t.Parallel()
	dataset := map[string]string{
		`대한민국만세`:     `대한민국만세`,
		`대한\민국만세`:    `대한\\민국만세`,
		`대한\\민국만세`:   `대한\\민국만세`,
		"abcdefgz":   "abcdefgz",
		`a\bcdefgz`:  `a\\bcdefgz`,
		`a\\bcdefgz`: `a\\bcdefgz`,
	}

	//check : common
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

	//check : common
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

	//check : common
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
		"<a href='http://www.president.go.kr/'>abcde</a><br>fgh": "abcde\nfgh",

		"world peace!!<a href='http://www.president.go.kr/'>대한민국만세</a><br>":   "world peace!!<a href='http://www.president.go.kr/'>대한민국만세</a>\n",
		"world peace!!<a href='http://www.president.go.kr/'>abcde</a><br>fgh": "world peace!!<a href='http://www.president.go.kr/'>abcde</a>\nfgh",

		"world peace!!<a href='http://www.president.go.kr/'><br />대한민국만세</a><br>":   "world peace!!<a href='http://www.president.go.kr/'>\n대한민국만세</a>\n",
		"world peace!!<a href='http://www.president.go.kr/'><br />abcde</a><br>fgh": "world peace!!<a href='http://www.president.go.kr/'>\nabcde</a>\nfgh",
	}

	//check : common
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

	//check : common
	for _, v := range dataset {

		retval, _ := strproc.WordWrapSimple(v.str, v.wd, v.breakstr)
		assert.AssertEquals(t, v.okstr, retval, "Original Value : %v\nReturn Value mismatch.\nExpected: %v\nActual: %v", v.str, retval, v)
	}

	//check : wd = 0
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

	//check : common
	for _, v := range dataset {

		retval, _ := strproc.WordWrapAround(v.str, v.wd, v.breakstr)
		assert.AssertEquals(t, v.okstr, retval, "Original Value : %v\nReturn Value mismatch.\nExpected: %v\nActual: %v", v.str, retval, v.okstr)
	}

	var err error

	//check : wd = 0
	_, err = strproc.WordWrapAround("test", 0, "1234")
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `wd at least 1`")

	//check : lastspc = 1
	_, err = strproc.WordWrapAround("ttttttt tttttttttt", 2, "1111")
	assert.AssertNil(t, err, "Failure : Couldn't check the `lastspc = 1`")

	//check : except
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
		//BUG(r) :
		//int8(math.MinInt8):       "-128",
		//float32(math.SmallestNonzeroFloat32): "1e-45",
		//float64(math.SmallestNonzeroFloat64): "5e-324",

	}

	//check : common
	for k, v := range dataset {
		retval, err := strproc.NumberFmt(k)

		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		assert.AssertNil(t, err, "Return Error : %v", err)
	}

	var err error

	//check : ParseFloat
	_, err = strproc.NumberFmt("12.11111111111111111111111111111111111111111111111111111111111e12e12e1p029ekj12e")
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `Not Support strconv.ParseFloat`")

	//check : not support obj
	_, err = strproc.NumberFmt(complex128(123))
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `not support obj`")

	//check : not support obj
	_, err = strproc.NumberFmt(complex64(123))
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `not support obj`")

	//check : not support obj
	_, err = strproc.NumberFmt(true)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `not support obj`")

	//check : not support numric string
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

	//check : common
	for _, v := range dataset {

		retval := strproc.Padding(v.str, v.fill, v.m, v.mx)
		assert.AssertEquals(t, v.okstr, retval, "Original Value : %v\nReturn Value mismatch.\nExpected: %v\nActual: %v", v.str, retval, v.okstr)
	}

	//check : mx >= byteStrLen
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

	//check : common
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

	//check : common
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

	//check : common
	for k, v := range dataset {
		retval := strproc.SwapCaseFirstWords(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func Test_strutils_HumanByteSize(t *testing.T) {

	t.Parallel()
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
		3.40282346638528859811704183484516925440e+38:                                   "288230358971842560.00Yb",
		"12121212121212121212121212121212121212121211212121212211212121212121.1234e+3": "0.00NaN",
	}

	//check : common
	for k, v := range dataset {
		retval, err := strproc.HumanByteSize(k, 2, strutils.CamelCaseDouble)

		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		assert.AssertNil(t, err, "Error : %v", err)
	}

	var err error

	//check : unit < UpperCaseSingle || unit > CamelCaseLong
	_, err = strproc.HumanByteSize(`1234`, 2, 123)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `retval, err := strproc.HumanByteSize(k, 2, strutils.CamelCaseDouble)`")

	//check : numberToString
	_, err = strproc.HumanByteSize(`abc`, 2, strutils.UpperCaseDouble)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `can't convert number to string`")

	//check : ParseFloat
	_, err = strproc.HumanByteSize("100.7976931348623157e+308", 2, strutils.UpperCaseDouble)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `strconv.ParseFloat(strNum, 64)`")

	//check : Complex64
	_, err = strproc.HumanByteSize(complex64(13), 2, strutils.UpperCaseDouble)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `not support obj.(complex128)`")

	//check : Complex128
	_, err = strproc.HumanByteSize(complex128(2+3i), 2, strutils.UpperCaseDouble)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `not support obj.(complex128)`")
}

func Test_strutils_HumanFileSize(t *testing.T) {

	t.Parallel()

	const tmpFilePath = "./filesizecheck.touch"
	const tmpPath = "./testdir"

	var retval string
	var err error

	//generating a touch file
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

	//check : common
	for k, v := range dataset {
		retval, err = strproc.HumanFileSize(tmpFilePath, 2, k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
		assert.AssertNil(t, err, "Error : %v", err)
	}

	//check : lost file description
	go func() {
		time.Sleep(time.Nanosecond * 10)
		err := os.Remove(tmpFilePath)
		assert.AssertLog(t, err, "lost file : %s but it's OK", tmpFilePath)

	}()
	_, err = strproc.HumanFileSize(tmpFilePath, 2, strutils.CamelCaseDouble)
	assert.AssertLog(t, err, "PASS")

	defer func(t *testing.T) {
		err := os.Remove(tmpFilePath)
		assert.AssertLog(t, err, "lost file : %s but it's OK", tmpFilePath)

	}(t)

	//check : isDir
	err = os.MkdirAll(tmpPath, 0777)
	assert.AssertNil(t, err, "Failure : Couldn't Mkdir %q: %s", tmpPath, err)

	_, err = strproc.HumanFileSize(tmpPath, 2, strutils.CamelCaseDouble)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `stat.IsDir()`")

	//check : os.Open
	_, err = strproc.HumanFileSize("/hello_word_txt", 2, strutils.CamelCaseDouble)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `os.Open()`")

	//check : not support obj.(complex128)
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

	//check : uint in map
	testMDepthUint1 := map[string]map[string]uint{"H": {"name": 1, "state": 2}}
	testMDepthUint2 := map[string]map[string]uint{"H": {"name": 1, "state": 3}}
	retval, _ = strproc.AnyCompare(testMDepthUint1, testMDepthUint2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison.")

	//check : float in map
	testMDepthFloat1 := map[string]map[string]float64{"H": {"name": 1, "state": 2}}
	testMDepthFloat2 := map[string]map[string]float64{"H": {"name": 1, "state": 3}}
	retval, _ = strproc.AnyCompare(testMDepthFloat1, testMDepthFloat2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison.")

	//check : complex in map
	testMDepthComplex1 := map[string]map[string]complex64{"H": {"name": 1, "state": 2}}
	testMDepthComplex2 := map[string]map[string]complex64{"H": {"name": 1, "state": 3}}
	retval, _ = strproc.AnyCompare(testMDepthComplex1, testMDepthComplex2)
	assert.AssertFalse(t, retval, "Couldn't make an accurate comparison.")

	//check : different type
	testDiffType1 := int(32)
	testDiffType2 := int64(32)
	_, err = strproc.AnyCompare(testDiffType1, testDiffType2)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `Different Type`\nError : %v", err)

	//check : different len
	testDiffLen1 := []int{1, 2, 3, 4, 5, 6, 7}
	testDiffLen2 := []int{1, 2, 3}
	_, err = strproc.AnyCompare(testDiffLen1, testDiffLen2)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `Different Len`\nError : %v", err)

	//check : different len
	testDiffMapLen1 := map[int]string{0: "A", 1: "B", 2: "C"}
	testDiffMapLen2 := map[int]string{0: "A", 1: "B"}
	_, err = strproc.AnyCompare(testDiffMapLen1, testDiffMapLen2)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `Different Len`\nError : %v", err)

	//check : not support compre
	testDiffNotSupport1 := paddingTestVal{}
	testDiffNotSupport2 := paddingTestVal{}
	_, err = strproc.AnyCompare(testDiffNotSupport1, testDiffNotSupport2)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `Not Support Compre`\nError : %v", err)

	//check : sting != string
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

	//check : uint != uint
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

	//check : float64 != float64
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

	//check : complex64 != complex64
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

	// check : failure at urldecode
	failTestStr := `html%26gt%3` // clear str is `html%26gt%3B`
	_, err = strproc.StripTags(failTestStr)
	assert.AssertNotNil(t, err, "Failure : Couldn't check the `failure at url decoding`")

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

	//check : common
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

	//check : common
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

	//check : common
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

	//check : common
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

	//check : common
	for k, v := range dataset {
		retval := strproc.ReverseUnicode(k)
		assert.AssertEquals(t, retval, v, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func Test_strutils_FileMD5Hash(t *testing.T) {

	t.Parallel()

	var retval string
	var err error

	//check : common
	retval, err = strproc.FileMD5Hash("./LICENSE")
	assert.AssertNil(t, err, "Error : %v", err)

	str_ok := "64e17a4e1c96bbfce57ab19cd0153e6a"
	assert.AssertEquals(t, retval, str_ok, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, str_ok)

	//check : os.Open
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

	//check : common
	for k, v := range dataset {
		retval, err := strproc.MD5Hash(k)

		assert.AssertNil(t, err, "Error : %v", err)
		assert.AssertEquals(t, retval, v, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}
