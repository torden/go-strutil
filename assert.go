package strutils

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// Assert is Methods for helping testing the strutils pkg.
type Assert struct {
	plib *StringProc
}

// NewAssert Creates and returns a String processing methods's pointer.
func NewAssert() *Assert {

	obj := &Assert{}
	obj.plib = NewStringProc()

	return obj
}

//a.printMsg is equivalent to t.Errorf include other information for easy debug
func (a *Assert) printMsg(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{}) {

	funcn, file, line, _ := runtime.Caller(2)
	t.Errorf(strings.Repeat("-", 120))
	t.Errorf("+ %v:%v\n", file, line)
	t.Errorf("+ %+v\n", runtime.FuncForPC(funcn).Name())
	t.Errorf(strings.Repeat("-", 120))
	if len(args) > 1 {
		t.Errorf(msgfmt, args)
	} else {
		t.Error(msgfmt)
	}

	t.Errorf("- value1 : %+v\n", v1)
	if v2 != nil {
		t.Errorf("- value2 : %+v\n", v2)
	}

	t.Errorf(strings.Repeat("-", 120))
}

//isCompareableNum asserts the specified objects are can compareble
func (a *Assert) isComparableNum(t *testing.T, v1 interface{}, v2 interface{}) bool {

	if !reflect.ValueOf(v1).IsValid() || !reflect.ValueOf(v2).IsValid() {
		a.printMsg(t, v1, v2, "Invalid Value")
		return false
	}

	refv1 := reflect.TypeOf(v1)
	refv2 := reflect.TypeOf(v1)

	if refv1.Comparable() != refv2.Comparable() {
		a.printMsg(t, v1, v2, "Not Comparable")
		return false
	}

	refv1k := refv1.Kind()
	refv2k := refv2.Kind()

	//int ~ int64 (0x2 ~ 0x6)
	//uint ~ uint64 (0x7 ~ 0xc)
	//float ~ float64 (0xd ~ 0xe)
	if (refv1k >= 0x2 && refv1k <= 0xe) && (refv2k >= 0x2 && refv2k <= 0xe) {
		return true
	}

	a.printMsg(t, v1, v2, "Different Type v1.(%d) != v2(%v)", refv1k.String(), refv2k.String())
	return false
}

//numericTypeUpCase converts the any numeric type to upsize type of that
func (a *Assert) numericTypeUpCase(val interface{}) (int64, uint64, float64, bool) {

	var tmpint int64
	var tmpuint uint64
	var tmpfloat float64

	switch val.(type) {
	case int:
		tmpint = int64(val.(int))
	case int8:
		tmpint = int64(val.(int8))
	case int16:
		tmpint = int64(val.(int16))
	case int32:
		tmpint = int64(val.(int32))
	case int64:
		tmpint = val.(int64)
	case uint:
		tmpuint = uint64(val.(uint))
	case uint8:
		tmpuint = uint64(val.(uint8))
	case uint16:
		tmpuint = uint64(val.(uint16))
	case uint32:
		tmpuint = uint64(val.(uint32))
	case uint64:
		tmpuint = val.(uint64)
	case float32:
		tmpfloat = float64(val.(float32))
	case float64:
		tmpfloat = val.(float64)
	default:
		return 0, 0, 0, false
	}

	return tmpint, tmpuint, tmpfloat, true
}

//AssertNilLog formats its arguments using default formatting, analogous to t.Log
func (a *Assert) AssertNilLog(t *testing.T, v1 interface{}, msgfmt string, args ...interface{}) {

	if v1 != nil {
		if len(args) > 0 {
			t.Logf(msgfmt, args)
		} else {
			t.Log(msgfmt)
		}
	}
}

//AssertEquals asserts that two objects are equal.
func (a *Assert) AssertEquals(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{}) {

	_, err := a.plib.AnyCompare(v1, v2)
	if err != nil {
		a.printMsg(t, v1, v2, err.Error())
	}
}

//AssertFalse asserts that the specified value is false.
func (a *Assert) AssertFalse(t *testing.T, v1 bool, msgfmt string, args ...interface{}) {

	if v1 {
		a.printMsg(t, v1, nil, msgfmt, args)
	}
}

//AssertTrue asserts that the specified value is true.
func (a *Assert) AssertTrue(t *testing.T, v1 bool, msgfmt string, args ...interface{}) {

	if !v1 {
		a.printMsg(t, v1, nil, msgfmt, args)
	}
}

//AssertNil asserts that the specified value is nil.
func (a *Assert) AssertNil(t *testing.T, v1 interface{}, msgfmt string, args ...interface{}) {

	if v1 != nil {
		a.printMsg(t, v1, nil, msgfmt, args)
	}
}

//AssertNotNil asserts that the specified value isn't nil.
func (a *Assert) AssertNotNil(t *testing.T, v1 interface{}, msgfmt string, args ...interface{}) {

	if v1 == nil {
		a.printMsg(t, v1, nil, msgfmt, args)
	}
}

//AssertLessThan asserts that the specified value are v1 less than v2
func (a *Assert) AssertLessThan(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{}) {

	if !a.isComparableNum(t, v1, v2) {
		return
	}

	tmpv1int, tmpv1uint, tmpv1float, ok := a.numericTypeUpCase(v1)
	if !ok {
		a.printMsg(t, v1, v2, "Required Numeric (int,uint,float with bit (8~64)")
		return
	}

	tmpv2int, tmpv2uint, tmpv2float, ok := a.numericTypeUpCase(v2)
	if !ok {
		a.printMsg(t, v1, v2, "Required Numeric (int,uint,float with bit (8~64)")
		return
	}

	var retval bool

	switch v1.(type) {
	case int, int8, int16, int32, int64:
		retval = (tmpv1int < tmpv2int)
	case uint, uint8, uint16, uint32, uint64:
		retval = (tmpv1uint < tmpv2uint)
	case float32, float64:
		retval = (tmpv1float < tmpv2float)
	}

	if !retval {
		a.printMsg(t, v1, v2, msgfmt, args)
	}
}

//AssertLessThanEqualTo asserts that the specified value are v1 less than v2 or equal to
func (a *Assert) AssertLessThanEqualTo(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{}) {

	if !a.isComparableNum(t, v1, v2) {
		return
	}

	tmpv1int, tmpv1uint, tmpv1float, ok := a.numericTypeUpCase(v1)
	if !ok {
		a.printMsg(t, v1, v2, "Required Numeric (int,uint,float with bit (8~64)")
		return
	}

	tmpv2int, tmpv2uint, tmpv2float, ok := a.numericTypeUpCase(v2)
	if !ok {
		a.printMsg(t, v1, v2, "Required Numeric (int,uint,float with bit (8~64)")
		return
	}

	var retval bool

	switch v1.(type) {
	case int, int8, int16, int32, int64:
		retval = (tmpv1int <= tmpv2int)
	case uint, uint8, uint16, uint32, uint64:
		retval = (tmpv1uint <= tmpv2uint)
	case float32, float64:
		retval = (tmpv1float <= tmpv2float)
	}

	if !retval {
		a.printMsg(t, v1, v2, msgfmt, args)
	}
}

//AssertGreaterThan nsserts that the specified value are v1 greater than v2
func (a *Assert) AssertGreaterThan(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{}) {

	if !a.isComparableNum(t, v1, v2) {
		return
	}

	tmpv1int, tmpv1uint, tmpv1float, ok := a.numericTypeUpCase(v1)
	if !ok {
		return
	}

	tmpv2int, tmpv2uint, tmpv2float, ok := a.numericTypeUpCase(v2)
	if !ok {
		return
	}

	retval := false

	switch v1.(type) {
	case int, int8, int16, int32, int64:
		retval = (tmpv1int > tmpv2int)
	case uint, uint8, uint16, uint32, uint64:
		retval = (tmpv1uint > tmpv2uint)
	case float32, float64:
		retval = (tmpv1float > tmpv2float)
	}

	if !retval {
		a.printMsg(t, v1, v2, msgfmt, args)
	}
}

//AssertGreaterThanEqualTo asserts that the specified value are v1 greater than v2 or equal to
func (a *Assert) AssertGreaterThanEqualTo(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{}) {

	if a.isComparableNum(t, v1, v2) {
		return
	}

	tmpv1int, tmpv1uint, tmpv1float, ok := a.numericTypeUpCase(v1)
	if !ok {
		return
	}

	tmpv2int, tmpv2uint, tmpv2float, ok := a.numericTypeUpCase(v2)
	if !ok {
		return
	}

	retval := false

	switch v1.(type) {
	case int, int8, int16, int32, int64:
		retval = (tmpv1int >= tmpv2int)
	case uint, uint8, uint16, uint32, uint64:
		retval = (tmpv1uint >= tmpv2uint)
	case float32, float64:
		retval = (tmpv1float >= tmpv2float)
	}

	if !retval {
		a.printMsg(t, v1, v2, msgfmt, args)
	}
}