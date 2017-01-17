package strutils

import (
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
)

type stringUtils struct{}

var numericPattern = regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?$`)

// init.
func NewStringUtils() stringUtils {
	return stringUtils{}
}

// quote string with slashes
func (s *stringUtils) AddSlashes(str string) string {

	l := len(str)

	buf := make([]byte, 0, l*2) //prealloca

	for i := 0; i < l; i++ {

		buf = append(buf, byte(str[i]))

		switch str[i] {

		case 92: //Dec : /

			if l >= i+1 {
				buf = append(buf, 92)

				if l > i+1 && str[i+1] == 92 {
					i++
				}
			}
		}
	}

	return string(buf)
}

// Un-quotes a quoted string
func (s *stringUtils) StripSlashes(str string) string {

	l := len(str)
	buf := make([]byte, 0, l) //prealloca

	for i := 0; i < l; i++ {

		buf = append(buf, byte(str[i]))
		if l > i+1 && str[i+1] == 92 {
			i++
		}
	}

	return string(buf)

}

// breakstr inserted before looks like space (CRLF , LFCR, SPACE, NL)
// BenchmarkNl2Br-8                   	10000000	      3398 ns/op
// BenchmarkNl2BrUseStringReplace-8   	10000000	      4535 ns/op
func (s *stringUtils) Nl2Br(str string) string {

	brtag := []byte("<br />")
	l := len(str)
	buf := make([]byte, 0, l) //prealloca

	for i := 0; i < l; i++ {

		switch str[i] {

		case 10, 13: //NL or CR

			for _, v := range brtag {
				buf = append(buf, v)
			}

			if l >= i+1 {
				if l > i+1 && (str[i+1] == 10 || str[i+1] == 13) { //NL+CR or CR+NL
					i++
				}
			}
		default:
			buf = append(buf, str[i])
		}
	}

	return string(buf)
}

// Wraps a string to a given number of characters using break characters (TAB, SPACE)
func (s *stringUtils) WordWrapSimple(str string, wd int, breakstr string) string {

	if wd < 1 {
		return str
	}

	strl := len(str)
	breakstrl := len(breakstr)

	buf := make([]byte, 0, (strl+breakstrl)*2)
	bufstr := []byte(str)

	brpos := 0
	for _, v := range bufstr {

		if (v == 9 || v == 32) && brpos >= wd {
			for _, vbk := range []byte(breakstr) {
				buf = append(buf, vbk)
			}
			brpos = -1

		} else {
			buf = append(buf, v)
		}
		brpos++
	}

	return string(buf)
}

// Wraps a string to a given number of characters using break characters (TAB, SPACE)
func (s *stringUtils) WordWrapAround(str string, wd int, breakstr string) string {

	if wd < 1 {
		return str
	}

	strl := len(str)
	breakstrl := len(breakstr)

	buf := make([]byte, 0, (strl+breakstrl)*2)
	bufstr := []byte(str)

	lastspc := make([]int, 0, strl)
	strpos := 0

	//looking for space or tab
	for _, v := range bufstr {

		if v == 9 || v == 32 {
			lastspc = append(lastspc, strpos)
		}
		strpos++
	}

	inject := make([]int, 0, strl)

	//looking for break point
	beforeBp := 0
	width := wd

	for _, v := range lastspc {

		if beforeBp != v {
			beforeBp = v
		}

		// DEBUG: fmt.Printf("V(%v) (%d <= %d || %d <= %d || %d <= %d) && %d <= %d : ", v, width, beforeBp, width, beforeBp+1, width, beforeBp-1, width, v)
		if (width <= beforeBp || width <= beforeBp+1 || width <= beforeBp-1) && width <= v {
			inject = append(inject, beforeBp)
			width += wd
			//fmt.Print("OK")
		} else if width < v && len(lastspc) == 1 {
			inject = append(inject, v)
		}
		//fmt.Println()
	}

	//injection
	breakno := 0
	loopcnt := 0
	injectcnt := len(inject)
	for _, v := range bufstr {

		//fmt.Printf("(%v) %d > %d && %d <= %d\n", v, injectcnt, breakno, inject[breakno], loopcnt)
		if injectcnt > breakno && inject[breakno] == loopcnt {
			for _, vbk := range []byte(breakstr) {
				buf = append(buf, vbk)
			}

			if injectcnt > breakno+1 {
				breakno++
			}
		} else {
			buf = append(buf, v)
		}

		loopcnt++
	}

	return string(buf)
}

// format a number with english notation grouped thousands
// TODO : improve bytebuffer efficiently
func (s *stringUtils) NumberFmt(obj interface{}) (string, error) {

	var strNum string

	switch obj.(type) {

	case string:
		strNum = obj.(string)
		if numericPattern.MatchString(strNum) == false {
			return "", fmt.Errorf("not support obj.(%v) := %v ", reflect.TypeOf(obj), strNum)
		}
	case int:
		strNum = strconv.FormatInt(int64(obj.(int)), 10)
	case int8:
		strNum = strconv.FormatInt(int64(obj.(int8)), 10)
	case int16:
		strNum = strconv.FormatInt(int64(obj.(int16)), 10)
	case int32:
		strNum = strconv.FormatInt(int64(obj.(int32)), 10)
	case int64:
		strNum = strconv.FormatInt(int64(obj.(int64)), 10)
	case uint:
		strNum = strconv.FormatUint(uint64(obj.(uint)), 10)
	case uint8:
		strNum = strconv.FormatUint(uint64(obj.(uint8)), 10)
	case uint16:
		strNum = strconv.FormatUint(uint64(obj.(uint16)), 10)
	case uint32:
		strNum = strconv.FormatUint(uint64(obj.(uint32)), 10)
	case uint64:
		strNum = strconv.FormatUint(uint64(obj.(uint64)), 10)
	case float32:
		strNum = fmt.Sprintf("%g", obj.(float32))
	case float64:
		strNum = fmt.Sprintf("%g", obj.(float64))
	default:
		return "", fmt.Errorf("not support obj.(%v)", reflect.TypeOf(obj))
	}

	bufbyteStr := []byte(strNum)
	bufbyteStr_len := len(bufbyteStr)

	//subffix after dot
	bufbyteTail := make([]byte, bufbyteStr_len-1)

	//init.
	foundDot := 0
	foundPos := 0
	dotcnt := 0
	bufbyteSize := 0

	//looking for dot
	for i := bufbyteStr_len - 1; i >= 0; i-- {
		if bufbyteStr[i] == 46 {
			copy(bufbyteTail, bufbyteStr[i:])
			foundDot = i
			foundPos = i
			break
		}
	}

	//make bufbyte size
	if foundDot == 0 { //numeric without dot
		bufbyteSize = int(math.Ceil(float64(bufbyteStr_len) + (float64(bufbyteStr_len) / 3)))
		foundDot = bufbyteStr_len
		foundPos = bufbyteSize - 2

		bufbyteSize -= 1

	} else { //with dot

		var cal_foundDot int

		if bufbyteStr[0] == 45 { //if startwith "-"(45)
			cal_foundDot = foundDot - 1
		} else {
			cal_foundDot = foundDot
		}

		bufbyteSize = int(math.Ceil(float64(cal_foundDot) + (float64(cal_foundDot) / 3) + float64(bufbyteStr_len-cal_foundDot) - 1))
	}

	//make a buffer byte
	bufbyte := make([]byte, bufbyteSize)

	//skip : need to dot injection
	if 4 > foundDot {
		return strNum, nil
	}

	//injection
	intoPos := foundPos
	for i := foundDot - 1; i >= 0; i-- {
		if dotcnt >= 3 && ((bufbyteStr[i] >= 48 && bufbyteStr[i] <= 57) || bufbyteStr[i] == 69 || bufbyteStr[i] == 101 || bufbyteStr[i] == 43) {
			bufbyte[intoPos] = 44
			intoPos--
			dotcnt = 0
		}
		bufbyte[intoPos] = bufbyteStr[i]
		intoPos--
		dotcnt++
	}

	//into dot to tail
	intoPos = foundPos + 1
	if foundDot != bufbyteStr_len {
		for _, v := range bufbyteTail {
			if v == 0 { //NULL
				break
			}

			bufbyte[intoPos] = v
			intoPos++
		}
	}

	return string(bufbyte), nil
}

// padding contol const
const (
	PAD_LEFT  = 0 //left padding
	PAD_RIGHT = 1 //right padding
	PAD_BOTH  = 2 //both padding
)

// pad a string to a certain length with another string
func (s *stringUtils) PaddingBoth(str string, fill string, mx int) string {
	return s.padding(str, fill, PAD_BOTH, mx)
}

// pad a string to a certain length with another string
func (s *stringUtils) PaddingLeft(str string, fill string, mx int) string {
	return s.padding(str, fill, PAD_LEFT, mx)
}

// pad a string to a certain length with another string
func (s *stringUtils) PaddingRight(str string, fill string, mx int) string {
	return s.padding(str, fill, PAD_RIGHT, mx)
}

// BenchmarkPadding-8                   10000000	       271 ns/op
// BenchmarkPaddingUseStringRepeat-8   	 3000000	       418 ns/op
func (s *stringUtils) padding(str string, fill string, m int, mx int) string {

	byteStr := []byte(str)
	byteStr_len := len(byteStr)
	if byteStr_len >= mx || mx < 1 {
		return str
	}

	var leftsize int
	var rightsize int

	switch m {
	case PAD_BOTH:
		rlsize := float64(mx-byteStr_len) / 2
		leftsize = int(rlsize)
		rightsize = int(rlsize + math.Copysign(0.5, rlsize))

	case PAD_LEFT:
		leftsize = mx - byteStr_len

	case PAD_RIGHT:
		rightsize = mx - byteStr_len

	}

	buf := make([]byte, 0, mx)

	if m == PAD_LEFT || m == PAD_BOTH {
		for i := 0; i < leftsize; {
			for _, v := range []byte(fill) {
				buf = append(buf, v)
				if i >= leftsize-1 {
					i = leftsize
					break
				} else {
					i++
				}
			}
		}
	}

	for _, v := range byteStr {
		buf = append(buf, v)
	}

	if m == PAD_RIGHT || m == PAD_BOTH {
		for i := 0; i < rightsize; {
			for _, v := range []byte(fill) {
				buf = append(buf, v)
				if i >= rightsize-1 {
					i = rightsize
					break
				} else {
					i++
				}
			}
		}
	}

	return string(buf)
}

// Lowercase the first character of each word in a string
// TOKEN : \t(9)\r(13)\n(10)\f(12)\v(11)\s(32)
func (s *stringUtils) LowerCaseFirstWords(str string) string {

	upper := 1
	bufbyteStr := []byte(str)
	retval := make([]byte, len(bufbyteStr))
	for k, v := range bufbyteStr {

		if upper == 1 && v >= 65 && v <= 90 {
			v = v + 32
		}

		upper = 0

		if (v >= 9 && v <= 13) || v == 32 {
			upper = 1
		}
		retval[k] = v
	}

	return string(retval)
}

// Uppercase the first character of each word in a string
// TOKEN : \t(9)\r(13)\n(10)\f(12)\v(11)\s(32)
func (s *stringUtils) UpperCaseFirstWords(str string) string {

	upper := 1
	bufbyteStr := []byte(str)
	retval := make([]byte, len(bufbyteStr))
	for k, v := range bufbyteStr {

		if upper == 1 && v >= 97 && v <= 122 {
			v = v - 32
		}

		upper = 0

		if (v >= 9 && v <= 13) || v == 32 {
			upper = 1
		}
		retval[k] = v
	}

	return string(retval)
}

// Switch the first character case of each word in a string
func (s *stringUtils) SwapCaseFirstWords(str string) string {

	upper := 1
	bufbyteStr := []byte(str)
	retval := make([]byte, len(bufbyteStr))
	for k, v := range bufbyteStr {

		switch {
		case upper == 1 && v >= 65 && v <= 90:
			v = v + 32

		case upper == 1 && v >= 97 && v <= 122:
			v = v - 32
		}

		upper = 0

		if (v >= 9 && v <= 13) || v == 32 {
			upper = 1
		}
		retval[k] = v
	}

	return string(retval)
}
