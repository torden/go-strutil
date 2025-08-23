//go:build !go1.18
// +build !go1.18

package strutils

import (
	"bytes"
	"fmt"
	"reflect"
)

// LegacyStringBuilder provides string building functionality for older Go versions
type LegacyStringBuilder struct {
	buffer bytes.Buffer
}

// NewLegacyStringBuilder creates a new legacy string builder
func NewLegacyStringBuilder() *LegacyStringBuilder {
	return &LegacyStringBuilder{}
}

// Append appends values to the builder
func (sb *LegacyStringBuilder) Append(values ...interface{}) *LegacyStringBuilder {
	for _, v := range values {
		switch val := v.(type) {
		case string:
			sb.buffer.WriteString(val)
		case fmt.Stringer:
			sb.buffer.WriteString(val.String())
		default:
			sb.buffer.WriteString(fmt.Sprintf("%v", val))
		}
	}
	return sb
}

// AppendLine appends values with a newline
func (sb *LegacyStringBuilder) AppendLine(values ...interface{}) *LegacyStringBuilder {
	sb.Append(values...)
	sb.buffer.WriteString("\n")
	return sb
}

// String returns the built string
func (sb *LegacyStringBuilder) String() string {
	return sb.buffer.String()
}

// Reset clears the builder
func (sb *LegacyStringBuilder) Reset() {
	sb.buffer.Reset()
}

// Len returns the length of the current content
func (sb *LegacyStringBuilder) Len() int {
	return sb.buffer.Len()
}

// LegacySliceProcessor provides slice processing for older Go versions
type LegacySliceProcessor struct{}

// NewLegacySliceProcessor creates a new legacy slice processor
func NewLegacySliceProcessor() *LegacySliceProcessor {
	return &LegacySliceProcessor{}
}

// ReverseSlice reverses a slice using reflection
func (sp *LegacySliceProcessor) ReverseSlice(slice interface{}) interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return slice
	}

	length := v.Len()
	if length <= 1 {
		return slice
	}

	// Create new slice of the same type
	result := reflect.MakeSlice(v.Type(), length, length)

	for i := 0; i < length; i++ {
		result.Index(i).Set(v.Index(length - 1 - i))
	}

	return result.Interface()
}

// UniqueSlice removes duplicates using reflection
func (sp *LegacySliceProcessor) UniqueSlice(slice interface{}) interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return slice
	}

	seen := make(map[interface{}]bool)
	var result []reflect.Value

	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		key := elem.Interface()

		if !seen[key] {
			seen[key] = true
			result = append(result, elem)
		}
	}

	// Create result slice
	resultSlice := reflect.MakeSlice(v.Type(), len(result), len(result))
	for i, val := range result {
		resultSlice.Index(i).Set(val)
	}

	return resultSlice.Interface()
}

// CompareSlices compares two slices using reflection
func (sp *LegacySliceProcessor) CompareSlices(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// NumberFmtInterface provides interface-based number formatting for legacy versions
func (s *StringProc) NumberFmtInterface(obj interface{}) (string, error) {
	// This delegates to the existing NumberFmt method
	return s.NumberFmt(obj)
}
