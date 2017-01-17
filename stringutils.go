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
var emailPattern = regexp.MustCompile("^[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?$")
var domainPattern = regexp.MustCompile(`^(([a-zA-Z0-9-\p{L}]{1,63}\.)?(xn--)?[a-zA-Z0-9\p{L}]+(-[a-zA-Z0-9\p{L}]+)*\.)+[a-zA-Z\p{L}]{2,63}$`)
var urlPattern = regexp.MustCompile(`^((((https?|ftps?|gopher|telnet|nntp)://)|(mailto:|news:))(%[0-9A-Fa-f]{2}|[-()_.!~*';/?:@#&=+$,A-Za-z0-9\p{L}])+)([).!';/?:,][[:blank:]])?$`)

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
	before_bp := 0
	width := wd

	for _, v := range lastspc {

		if before_bp != v {
			before_bp = v
		}

		// DEBUG: fmt.Printf("V(%v) (%d <= %d || %d <= %d || %d <= %d) && %d <= %d : ", v, width, before_bp, width, before_bp+1, width, before_bp-1, width, v)
		if (width <= before_bp || width <= before_bp+1 || width <= before_bp-1) && width <= v {
			inject = append(inject, before_bp)
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

	var str_num string

	switch obj.(type) {

	case string:
		str_num = obj.(string)
		if numericPattern.MatchString(str_num) == false {
			return "", fmt.Errorf("not support obj.(%v) := %v ", reflect.TypeOf(obj), str_num)
		}
	case int:
		str_num = strconv.FormatInt(int64(obj.(int)), 10)
	case int8:
		str_num = strconv.FormatInt(int64(obj.(int8)), 10)
	case int16:
		str_num = strconv.FormatInt(int64(obj.(int16)), 10)
	case int32:
		str_num = strconv.FormatInt(int64(obj.(int32)), 10)
	case int64:
		str_num = strconv.FormatInt(int64(obj.(int64)), 10)
	case uint:
		str_num = strconv.FormatUint(uint64(obj.(uint)), 10)
	case uint8:
		str_num = strconv.FormatUint(uint64(obj.(uint8)), 10)
	case uint16:
		str_num = strconv.FormatUint(uint64(obj.(uint16)), 10)
	case uint32:
		str_num = strconv.FormatUint(uint64(obj.(uint32)), 10)
	case uint64:
		str_num = strconv.FormatUint(uint64(obj.(uint64)), 10)
	case float32:
		str_num = fmt.Sprintf("%g", obj.(float32))
	case float64:
		str_num = fmt.Sprintf("%g", obj.(float64))
	default:
		return "", fmt.Errorf("not support obj.(%v)", reflect.TypeOf(obj))
	}

	bufbyte_str := []byte(str_num)
	bufbyte_str_len := len(bufbyte_str)

	//subffix after dot
	bufbyte_tail := make([]byte, bufbyte_str_len-1)

	//init.
	found_dot := 0
	found_pos := 0
	dotcnt := 0
	bufbyte_size := 0

	//looking for dot
	for i := bufbyte_str_len - 1; i >= 0; i-- {
		if bufbyte_str[i] == 46 {
			copy(bufbyte_tail, bufbyte_str[i:])
			found_dot = i
			found_pos = i
			break
		}
	}

	//make bufbyte size
	if found_dot == 0 { //numeric without dot
		bufbyte_size = int(math.Ceil(float64(bufbyte_str_len) + (float64(bufbyte_str_len) / 3)))
		found_dot = bufbyte_str_len
		found_pos = bufbyte_size - 2

		bufbyte_size -= 1

	} else { //with dot

		var cal_found_dot int

		if bufbyte_str[0] == 45 { //if startwith "-"(45)
			cal_found_dot = found_dot - 1
		} else {
			cal_found_dot = found_dot
		}

		bufbyte_size = int(math.Ceil(float64(cal_found_dot) + (float64(cal_found_dot) / 3) + float64(bufbyte_str_len-cal_found_dot) - 1))
	}

	//make a buffer byte
	bufbyte := make([]byte, bufbyte_size)

	//skip : need to dot injection
	if 4 > found_dot {
		return str_num, nil
	}

	//injection
	into_pos := found_pos
	for i := found_dot - 1; i >= 0; i-- {
		if dotcnt >= 3 && ((bufbyte_str[i] >= 48 && bufbyte_str[i] <= 57) || bufbyte_str[i] == 69 || bufbyte_str[i] == 101 || bufbyte_str[i] == 43) {
			bufbyte[into_pos] = 44
			into_pos--
			dotcnt = 0
		}
		bufbyte[into_pos] = bufbyte_str[i]
		into_pos--
		dotcnt++
	}

	//into dot to tail
	into_pos = found_pos + 1
	if found_dot != bufbyte_str_len {
		for _, v := range bufbyte_tail {
			if v == 0 { //NULL
				break
			}

			bufbyte[into_pos] = v
			into_pos++
		}
	}

	return string(bufbyte), nil
}

const (
	PAD_LEFT  = 0
	PAD_RIGHT = 1
	PAD_BOTH  = 2
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

	byte_str := []byte(str)
	byte_str_len := len(byte_str)
	if byte_str_len >= mx || mx < 1 {
		return str
	}

	var leftsize int
	var rightsize int

	switch m {
	case PAD_BOTH:
		rlsize := float64(mx-byte_str_len) / 2
		leftsize = int(rlsize)
		rightsize = int(rlsize + math.Copysign(0.5, rlsize))

	case PAD_LEFT:
		leftsize = mx - byte_str_len

	case PAD_RIGHT:
		rightsize = mx - byte_str_len

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

	for _, v := range byte_str {
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
