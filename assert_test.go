package strutils_test

import (
	"math"
	"testing"
)

func Test_strutils_AssertLog(t *testing.T) {

	assert.AssertLog(t, nil, "test", "hello")
	assert.AssertLog(t, nil, "test", "hello %s", "word")
}

func Test_strutils_AssertEquals(t *testing.T) {

	assert.AssertEquals(t, 1, 1, "hello")
	assert.AssertEquals(t, 1, 1, "hello %s", "word")

	assert.AssertEquals(t, "a", "a", "hello")
	assert.AssertEquals(t, "a", "a", "hello %s", "word")

	assert.AssertEquals(t, math.Log(-1.0), math.Log(-1.0), "hello")
}

func Test_strutils_AssertNotEquals(t *testing.T) {

	assert.AssertNotEquals(t, 1, 2, "hello")
	assert.AssertNotEquals(t, 1, 2, "hello %s", "word")

	assert.AssertNotEquals(t, "a", "b", "hello")
	assert.AssertNotEquals(t, "a", "b", "hello %s", "word")

	assert.AssertNotEquals(t, math.Log(-1.0), math.Log(-1.1), "hello")
}

func Test_strutils_AssertFalse(t *testing.T) {

	assert.AssertFalse(t, false, "hello")
	assert.AssertFalse(t, true, "hello %s", "word")
}

func Test_strutils_AssertTrue(t *testing.T) {

	assert.AssertTrue(t, false, "hello")
	assert.AssertTrue(t, true, "hello %s", "word")
}

func Test_strutils_AssertNil(t *testing.T) {

	assert.AssertNil(t, nil, "hello")
	assert.AssertNil(t, true, "hello %s", "word")
}

func Test_strutils_AssertNotNil(t *testing.T) {

	assert.AssertNotNil(t, nil, "hello")
	assert.AssertNotNil(t, true, "hello %s", "word")
}

func Test_strutils_AssertLessThan(t *testing.T) {

	assert.AssertLessThan(t, "a", "a", "hello")
	assert.AssertLessThan(t, "a", "a", "hello %s", "word")

	assert.AssertLessThan(t, 1, 1, "hello")
	assert.AssertLessThan(t, 1, 1, "hello %s", "word")

	assert.AssertLessThan(t, 1, 2, "hello")
	assert.AssertLessThan(t, 1, 2, "hello %s", "word")

	assert.AssertLessThan(t, 3, 2, "hello")
	assert.AssertLessThan(t, 3, 2, "hello %s", "word")

	assert.AssertLessThan(t, nil, 2, "hello")
	assert.AssertLessThan(t, nil, 2, "hello %s", "word")

	assert.AssertLessThan(t, nil, nil, "hello")
	assert.AssertLessThan(t, nil, nil, "hello %s", "word")

	assert.AssertLessThan(t, int(1), int(1), "hello")
	assert.AssertLessThan(t, int8(1), int8(1), "hello")
	assert.AssertLessThan(t, int16(1), int16(1), "hello")
	assert.AssertLessThan(t, int32(1), int32(1), "hello")
	assert.AssertLessThan(t, int64(1), int64(1), "hello")
	assert.AssertLessThan(t, uint(1), uint(1), "hello")
	assert.AssertLessThan(t, uint8(1), uint8(1), "hello")
	assert.AssertLessThan(t, uint16(1), uint16(1), "hello")
	assert.AssertLessThan(t, uint32(1), uint32(1), "hello")
	assert.AssertLessThan(t, uint64(1), uint64(1), "hello")
	assert.AssertLessThan(t, float32(1), float32(1), "hello")
	assert.AssertLessThan(t, float64(1), float64(1), "hello")

	assert.AssertLessThan(t, math.Log(-1.0), math.Log(-1.0), "hello")
	assert.AssertLessThan(t, math.Log(-1.0), nil, "hello")
	assert.AssertLessThan(t, math.Log(1.0), nil, "hello")
}

func Test_strutils_AssertLessThanEqualTo(t *testing.T) {

	assert.AssertLessThanEqualTo(t, "a", "a", "hello")
	assert.AssertLessThanEqualTo(t, "a", "a", "hello %s", "word")

	assert.AssertLessThanEqualTo(t, 1, 1, "hello")
	assert.AssertLessThanEqualTo(t, 1, 1, "hello %s", "word")

	assert.AssertLessThanEqualTo(t, 1, 2, "hello")
	assert.AssertLessThanEqualTo(t, 1, 2, "hello %s", "word")

	assert.AssertLessThanEqualTo(t, 3, 2, "hello")
	assert.AssertLessThanEqualTo(t, 3, 2, "hello %s", "word")

	assert.AssertLessThanEqualTo(t, nil, 2, "hello")
	assert.AssertLessThanEqualTo(t, nil, 2, "hello %s", "word")

	assert.AssertLessThanEqualTo(t, nil, nil, "hello")
	assert.AssertLessThanEqualTo(t, nil, nil, "hello %s", "word")

	assert.AssertLessThanEqualTo(t, int(1), int(1), "hello")
	assert.AssertLessThanEqualTo(t, int8(1), int8(1), "hello")
	assert.AssertLessThanEqualTo(t, int16(1), int16(1), "hello")
	assert.AssertLessThanEqualTo(t, int32(1), int32(1), "hello")
	assert.AssertLessThanEqualTo(t, int64(1), int64(1), "hello")
	assert.AssertLessThanEqualTo(t, uint(1), uint(1), "hello")
	assert.AssertLessThanEqualTo(t, uint8(1), uint8(1), "hello")
	assert.AssertLessThanEqualTo(t, uint16(1), uint16(1), "hello")
	assert.AssertLessThanEqualTo(t, uint32(1), uint32(1), "hello")
	assert.AssertLessThanEqualTo(t, uint64(1), uint64(1), "hello")
	assert.AssertLessThanEqualTo(t, float32(1), float32(1), "hello")
	assert.AssertLessThanEqualTo(t, float64(1), float64(1), "hello")

	assert.AssertLessThanEqualTo(t, math.Log(-1.0), math.Log(-1.0), "hello")
	assert.AssertLessThanEqualTo(t, math.Log(-1.0), nil, "hello")
	assert.AssertLessThanEqualTo(t, math.Log(1.0), nil, "hello")
}

func Test_strutils_AssertGreaterThan(t *testing.T) {

	assert.AssertGreaterThan(t, "a", "a", "hello")
	assert.AssertGreaterThan(t, "a", "a", "hello %s", "word")

	assert.AssertGreaterThan(t, 1, 1, "hello")
	assert.AssertGreaterThan(t, 1, 1, "hello %s", "word")

	assert.AssertGreaterThan(t, 1, 2, "hello")
	assert.AssertGreaterThan(t, 1, 2, "hello %s", "word")

	assert.AssertGreaterThan(t, 3, 2, "hello")
	assert.AssertGreaterThan(t, 3, 2, "hello %s", "word")

	assert.AssertGreaterThan(t, nil, 2, "hello")
	assert.AssertGreaterThan(t, nil, 2, "hello %s", "word")

	assert.AssertGreaterThan(t, nil, nil, "hello")
	assert.AssertGreaterThan(t, nil, nil, "hello %s", "word")

	assert.AssertGreaterThan(t, int(1), int(1), "hello")
	assert.AssertGreaterThan(t, int8(1), int8(1), "hello")
	assert.AssertGreaterThan(t, int16(1), int16(1), "hello")
	assert.AssertGreaterThan(t, int32(1), int32(1), "hello")
	assert.AssertGreaterThan(t, int64(1), int64(1), "hello")
	assert.AssertGreaterThan(t, uint(1), uint(1), "hello")
	assert.AssertGreaterThan(t, uint8(1), uint8(1), "hello")
	assert.AssertGreaterThan(t, uint16(1), uint16(1), "hello")
	assert.AssertGreaterThan(t, uint32(1), uint32(1), "hello")
	assert.AssertGreaterThan(t, uint64(1), uint64(1), "hello")
	assert.AssertGreaterThan(t, float32(1), float32(1), "hello")
	assert.AssertGreaterThan(t, float64(1), float64(1), "hello")

	assert.AssertGreaterThan(t, math.Log(-1.0), math.Log(-1.0), "hello")
	assert.AssertGreaterThan(t, math.Log(-1.0), nil, "hello")
	assert.AssertGreaterThan(t, math.Log(1.0), nil, "hello")

}

func Test_strutils_AssertGreaterThanEqualTo(t *testing.T) {

	assert.AssertGreaterThanEqualTo(t, 1, 2, "hello")
	assert.AssertGreaterThanEqualTo(t, 1, 2, "hello %s", "word")

	assert.AssertGreaterThanEqualTo(t, 3, 2, "hello")
	assert.AssertGreaterThanEqualTo(t, 3, 2, "hello %s", "word")

	assert.AssertGreaterThanEqualTo(t, nil, 2, "hello")
	assert.AssertGreaterThanEqualTo(t, nil, 2, "hello %s", "word")

	assert.AssertGreaterThanEqualTo(t, nil, nil, "hello")
	assert.AssertGreaterThanEqualTo(t, nil, nil, "hello %s", "word")

	assert.AssertGreaterThanEqualTo(t, int(1), int(1), "hello")
	assert.AssertGreaterThanEqualTo(t, int8(1), int8(1), "hello")
	assert.AssertGreaterThanEqualTo(t, int16(1), int16(1), "hello")
	assert.AssertGreaterThanEqualTo(t, int32(1), int32(1), "hello")
	assert.AssertGreaterThanEqualTo(t, int64(1), int64(1), "hello")
	assert.AssertGreaterThanEqualTo(t, uint(1), uint(1), "hello")
	assert.AssertGreaterThanEqualTo(t, uint8(1), uint8(1), "hello")
	assert.AssertGreaterThanEqualTo(t, uint16(1), uint16(1), "hello")
	assert.AssertGreaterThanEqualTo(t, uint32(1), uint32(1), "hello")
	assert.AssertGreaterThanEqualTo(t, uint64(1), uint64(1), "hello")
	assert.AssertGreaterThanEqualTo(t, float32(1), float32(1), "hello")
	assert.AssertGreaterThanEqualTo(t, float64(1), float64(1), "hello")
	assert.AssertGreaterThanEqualTo(t, complex128(1), float64(1), "hello")
	assert.AssertGreaterThanEqualTo(t, float64(1), complex128(1), "hello")

	assert.AssertGreaterThanEqualTo(t, math.Log(-1.0), math.Log(-1.0), "hello")
	assert.AssertGreaterThanEqualTo(t, math.Log(-1.0), nil, "hello")
	assert.AssertGreaterThanEqualTo(t, math.Log(1.0), nil, "hello")
}

func Test_strutils_AssertLengthOf(t *testing.T) {

	assert.AssertLengthOf(t, "asdfg456", 8, "hello")
	assert.AssertLengthOf(t, []int{1, 2, 3, 4, 5, 6, 7}, 7, "hello")
	assert.AssertLengthOf(t, []string{"a", "b", "c"}, 3, "hello")
	assert.AssertLengthOf(t, map[int]string{1: "a", 2: "b", 3: "c"}, 3, "hello")

	var chans [9]chan int
	for i := range chans {
		chans[i] = make(chan int)
	}

	assert.AssertLengthOf(t, chans, 9, "hello")

}
