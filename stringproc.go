package strutils

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var numericPattern = regexp.MustCompile(`^[-+]?[0-9]*\.?[0-9]+([eE][-+]?[0-9]+)?$`)

// StringProc is String processing methods, All operations on this object
type StringProc struct{}

// NewStringProc is Creates and returns a String processing methods's pointer.
func NewStringProc() *StringProc {
	return &StringProc{}
}

// AddSlashes is quote string with slashes
func (s *StringProc) AddSlashes(str string) string {

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

// StripSlashes is Un-quotes a quoted string
func (s *StringProc) StripSlashes(str string) string {

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

// Nl2Br is breakstr inserted before looks like space (CRLF , LFCR, SPACE, NL)
func (s *StringProc) Nl2Br(str string) string {

	// BenchmarkNl2Br-8                   	10000000	      3398 ns/op
	// BenchmarkNl2BrUseStringReplace-8   	10000000	      4535 ns/op
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

// WordWrapSimple is Wraps a string to a given number of characters using break characters (TAB, SPACE)
func (s *StringProc) WordWrapSimple(str string, wd int, breakstr string) string {

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

// WordWrapAround is Wraps a string to a given number of characters using break characters (TAB, SPACE)
func (s *StringProc) WordWrapAround(str string, wd int, breakstr string) string {

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

func numberToString(obj interface{}) (string, error) {

	var strNum string

	switch obj.(type) {

	case string:
		strNum = obj.(string)
		if numericPattern.MatchString(strNum) == false {
			return "", fmt.Errorf("Not Support obj.(%v) := %v ", reflect.TypeOf(obj), strNum)
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
		return "", fmt.Errorf("Not Support obj.(%v)", reflect.TypeOf(obj))
	}

	return strNum, nil
}

// NumberFmt is format a number with english notation grouped thousands
// TODO : support other country notation
func (s *StringProc) NumberFmt(obj interface{}) (string, error) {

	strNum, err := numberToString(obj)
	if err != nil {
		return "", err
	}

	bufbyteStr := []byte(strNum)
	bufbyteStrLen := len(bufbyteStr)

	//subffix after dot
	bufbyteTail := make([]byte, bufbyteStrLen-1)

	//init.
	foundDot := 0
	foundPos := 0
	dotcnt := 0
	bufbyteSize := 0

	//looking for dot
	for i := bufbyteStrLen - 1; i >= 0; i-- {
		if bufbyteStr[i] == 46 {
			copy(bufbyteTail, bufbyteStr[i:])
			foundDot = i
			foundPos = i
			break
		}
	}

	//make bufbyte size
	if foundDot == 0 { //numeric without dot
		bufbyteSize = int(math.Ceil(float64(bufbyteStrLen) + (float64(bufbyteStrLen) / 3)))
		foundDot = bufbyteStrLen
		foundPos = bufbyteSize - 2

		bufbyteSize--

	} else { //with dot

		var calFoundDot int

		if bufbyteStr[0] == 45 { //if startwith "-"(45)
			calFoundDot = foundDot - 1
		} else {
			calFoundDot = foundDot
		}

		bufbyteSize = int(math.Ceil(float64(calFoundDot) + (float64(calFoundDot) / 3) + float64(bufbyteStrLen-calFoundDot) - 1))
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
	if foundDot != bufbyteStrLen {
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
	PadLeft  = 0 //left padding
	PadRight = 1 //right padding
	PadBoth  = 2 //both padding
)

// PaddingBoth is Padding method alias with PadBoth Option
func (s *StringProc) PaddingBoth(str string, fill string, mx int) string {
	return s.Padding(str, fill, PadBoth, mx)
}

// PaddingLeft is Padding method alias with PadRight Option
func (s *StringProc) PaddingLeft(str string, fill string, mx int) string {
	return s.Padding(str, fill, PadLeft, mx)
}

// PaddingRight is Padding method alias with PadRight Option
func (s *StringProc) PaddingRight(str string, fill string, mx int) string {
	return s.Padding(str, fill, PadRight, mx)
}

// Padding is Pad a string to a certain length with another string
// BenchmarkPadding-8                   10000000	       271 ns/op
// BenchmarkPaddingUseStringRepeat-8   	 3000000	       418 ns/op
func (s *StringProc) Padding(str string, fill string, m int, mx int) string {

	byteStr := []byte(str)
	byteStrLen := len(byteStr)
	if byteStrLen >= mx || mx < 1 {
		return str
	}

	var leftsize int
	var rightsize int

	switch m {
	case PadBoth:
		rlsize := float64(mx-byteStrLen) / 2
		leftsize = int(rlsize)
		rightsize = int(rlsize + math.Copysign(0.5, rlsize))

	case PadLeft:
		leftsize = mx - byteStrLen

	case PadRight:
		rightsize = mx - byteStrLen

	}

	buf := make([]byte, 0, mx)

	if m == PadLeft || m == PadBoth {
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

	if m == PadRight || m == PadBoth {
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

// LowerCaseFirstWords is Lowercase the first character of each word in a string
// INFO : (Support Token Are \t(9)\r(13)\n(10)\f(12)\v(11)\s(32))
func (s *StringProc) LowerCaseFirstWords(str string) string {

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

// UpperCaseFirstWords is Uppercase the first character of each word in a string
// INFO : (Support Token Are \t(9)\r(13)\n(10)\f(12)\v(11)\s(32))
func (s *StringProc) UpperCaseFirstWords(str string) string {

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

// SwapCaseFirstWords is Switch the first character case of each word in a string
func (s *StringProc) SwapCaseFirstWords(str string) string {

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

// Unit type control
const (
	_               = uint8(iota)
	LowerCaseSingle // Single Unit character converted to Lower-case
	LowerCaseDouble // Double Unit characters converted to Lower-case

	UpperCaseSingle // Single Unit character converted to Uppper-case
	UpperCaseDouble // Double Unit characters converted to Upper-case

	CamelCaseDouble // Double Unit characters converted to Camel-case
	CamelCaseLong   // Full Unit characters converted to Camel-case
)

var sizeStrLowerCaseSingle = []string{"b", "k", "m", "g", "t", "p", "e", "z", "y"}
var sizeStrLowerCaseDouble = []string{"b", "kb", "mb", "gb", "tb", "pb", "eb", "zb", "yb"}
var sizeStrUpperCaseSingle = []string{"B", "K", "M", "G", "T", "P", "E", "Z", "Y"}
var sizeStrUpperCaseDouble = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
var sizeStrCamelCaseDouble = []string{"B", "Kb", "Mb", "Gb", "Tb", "Eb", "Zb", "Yb"}
var sizeStrCamelCaseLong = []string{"Byte", "KiloByte", "MegaByte", "GigaByte", "TeraByte", "ExaByte", "ZettaByte", "YottaByte"}

//HumanByteSize is Byte Size convert to Easy Readable Size String
func (s *StringProc) HumanByteSize(obj interface{}, decimals int, unit uint8) (string, error) {

	if unit < UpperCaseSingle || unit > CamelCaseLong {
		return "", fmt.Errorf("Not allow unit parameter : %v", unit)
	}

	strNum, err := numberToString(obj)
	if err != nil {
		return "", err
	}

	var bufStrFloat64 float64

	switch obj.(type) {
	case string:
		bufStrFloat64, err = strconv.ParseFloat(strNum, 64)
		if err != nil {
			return "", fmt.Errorf("Not Support %v (obj.(%v))", obj, reflect.TypeOf(obj))
		}

	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32:

		float64Type := reflect.TypeOf(float64(0))
		tmpVal := reflect.Indirect(reflect.ValueOf(obj))

		if tmpVal.Type().ConvertibleTo(float64Type) == false {
			return "", fmt.Errorf("Not Support obj.(%v)", reflect.TypeOf(obj))
		}

		bufStrFloat64 = tmpVal.Convert(float64Type).Float()

	case float64:
		bufStrFloat64 = obj.(float64)

	default:
		return "", fmt.Errorf("Not Support obj.(%v)", reflect.TypeOf(obj))
	}

	var sizeStr []string

	switch unit {
	case LowerCaseSingle:
		sizeStr = sizeStrLowerCaseSingle
	case LowerCaseDouble:
		sizeStr = sizeStrLowerCaseDouble
	case UpperCaseSingle:
		sizeStr = sizeStrUpperCaseSingle
	case UpperCaseDouble:
		sizeStr = sizeStrUpperCaseDouble
	case CamelCaseDouble:
		sizeStr = sizeStrCamelCaseDouble
	case CamelCaseLong:
		sizeStr = sizeStrCamelCaseLong
	}

	strNumLen := len(strNum)

	factor := int(math.Floor(float64(strNumLen)-1) / 3)

	decimalsFmt := `%.` + strconv.Itoa(decimals) + `f%s`
	humanSize := bufStrFloat64 / math.Pow(1024, float64(factor))

	return fmt.Sprintf(decimalsFmt, humanSize, sizeStr[factor]), nil
}

//HumanFileSize is File Size convert to Easy Readable Size String
func (s *StringProc) HumanFileSize(filepath string, decimals int, unit uint8) (string, error) {

	fd, err := os.Open(filepath)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	stat, err := fd.Stat()
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	if stat.IsDir() == true {
		return "", fmt.Errorf("%v isn't file", filepath)
	}

	return s.HumanByteSize(stat.Size(), decimals, unit)
}

// compare with map
var recursiveDepth = 0
var recursiveDepthKeypList []string

func compareMap(compObj1 reflect.Value, compObj2 reflect.Value) (bool, error) {

	recursiveDepth++
	var valueCompareErr bool

	for _, k := range compObj1.MapKeys() {

		recursiveDepthKeypList = append(recursiveDepthKeypList, k.String())

		//check : Type
		if compObj1.MapIndex(k).Kind() != compObj2.MapIndex(k).Kind() {
			return false, fmt.Errorf("Different Type : (obj1[%v] is  %v) != (obj2[%v] is  %v)", k, compObj1.MapIndex(k).Kind(), k, compObj1.MapIndex(k).Kind())
		}

		switch compObj1.MapIndex(k).Kind() {

		//String
		case reflect.String:
			if compObj1.MapIndex(k).String() != compObj2.MapIndex(k).String() {
				valueCompareErr = true
			}

		//Integer
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if compObj1.MapIndex(k).Int() != compObj2.MapIndex(k).Int() {
				valueCompareErr = true
			}

		//Un-signed Integer
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if compObj1.MapIndex(k).Uint() != compObj2.MapIndex(k).Uint() {
				valueCompareErr = true
			}

		//Float
		case reflect.Float32, reflect.Float64:
			if compObj1.MapIndex(k).Float() != compObj2.MapIndex(k).Float() {
				valueCompareErr = true
			}

		//Boolean
		case reflect.Bool:
			if compObj1.MapIndex(k).Bool() != compObj2.MapIndex(k).Bool() {
				valueCompareErr = true
			}

		//Complex
		case reflect.Complex64, reflect.Complex128:
			if compObj1.MapIndex(k).Complex() != compObj2.MapIndex(k).Complex() {
				valueCompareErr = true
			}

		//Map : recursive loop
		case reflect.Map:
			retval, err := compareMap(compObj1.MapIndex(k), compObj2.MapIndex(k))
			if retval == false {
				return retval, err
			}

		default:
			return false, fmt.Errorf("Not Support Compare : (obj1[%v] := %v) != (obj2[%v] := %v)", k, compObj1.MapIndex(k), k, compObj2.MapIndex(k))
		}

		if valueCompareErr == true {
			if recursiveDepth == 1 {
				return false, fmt.Errorf("Different Value : (obj1[%v] := %v) != (obj2[%v] := %v)", k, compObj1.MapIndex(k), k, compObj2.MapIndex(k))
			}

			depthStr := strings.Join(recursiveDepthKeypList, "][")
			return false, fmt.Errorf("Different Value : (obj1[%v] := %v) != (obj2[%v] := %v)", depthStr, compObj1.MapIndex(k).Interface(), depthStr, compObj2.MapIndex(k))

		}
	}

	return true, nil
}

// AnyCompare is compares two same basic type (without prt) dataset (slice,map,single data).
// TODO : support interface, struct ...
// NOTE : Not safe , Not Test Complete. Require more test data based on the complex dataset.
func (s *StringProc) AnyCompare(obj1 interface{}, obj2 interface{}) (bool, error) {

	if reflect.TypeOf(obj1) != reflect.TypeOf(obj2) {
		return false, fmt.Errorf("Not Compare type, obj1.(%v) != obj2.(%v)", obj1, obj2)
	}

	recursiveDepthKeypList = make([]string, 0)

	switch obj1.(type) {

	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, complex64, complex128, bool:
		if reflect.TypeOf(obj1).Comparable() == true && reflect.TypeOf(obj2).Comparable() == true {
			return (obj1 == obj2), nil
		}

	default:

		compObj1 := reflect.ValueOf(obj1)
		compObj2 := reflect.ValueOf(obj2)

		if compObj1.Len() != compObj2.Len() {
			return false, fmt.Errorf("Different Size : obj1(%d) != obj2(%d)", compObj1.Len(), compObj2.Len())
		}

		switch {

		case compObj1.Kind() == reflect.Slice:

			for i := 0; i < compObj1.Len(); i++ {
				if compObj1.Index(i).Interface() != compObj2.Index(i).Interface() {
					return false, fmt.Errorf("Different Value : (obj1[%d] := %v) != (obj2[%d] := %v)", i, compObj1.Index(i).Interface(), i, compObj2.Index(i).Interface())
				}
			}

		case compObj1.Kind() == reflect.Map:
			recursiveDepth = 0
			retval, err := compareMap(compObj1, compObj2)
			if retval == false {
				return retval, err
			}

		default:
			return false, fmt.Errorf("Not Support Compare : (obj1[%v]) , (obj2[%v])", compObj1.Kind(), compObj2.Kind())

		}
	}
	return true, nil
}