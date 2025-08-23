//go:build go1.18
// +build go1.18

package strutils

import (
	"fmt"
	"strings"
)

// Generic type constraints for Go 1.18+
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type Stringable interface {
	~string
}

// NumberFmtGeneric is a generic version of NumberFmt for Go 1.18+
func NumberFmtGeneric[T Numeric](s *StringProc, obj T) (string, error) {
	// Convert to string using fmt.Sprintf for simplicity and type safety
	strNum := fmt.Sprintf("%v", obj)

	// Use existing NumberFmt logic
	return s.NumberFmt(strNum)
}

// PaddingGeneric provides generic padding with type safety
func PaddingGeneric[T Stringable](s *StringProc, str T, fill string, mode int, length int) string {
	var input string
	switch v := any(str).(type) {
	case string:
		input = v
	case fmt.Stringer:
		input = v.String()
	default:
		input = fmt.Sprintf("%v", str)
	}

	return s.Padding(input, fill, mode, length)
}

// SliceProcessor provides generic slice processing for Go 1.18+
type SliceProcessor[T comparable] struct{}

// NewSliceProcessor creates a new generic slice processor
func NewSliceProcessor[T comparable]() *SliceProcessor[T] {
	return &SliceProcessor[T]{}
}

// ReverseSlice reverses a slice of any comparable type
func (sp *SliceProcessor[T]) ReverseSlice(slice []T) []T {
	length := len(slice)
	if length <= 1 {
		return slice
	}

	result := make([]T, length)
	for i := 0; i < length; i++ {
		result[i] = slice[length-1-i]
	}
	return result
}

// UniqueSlice removes duplicates from a slice
func (sp *SliceProcessor[T]) UniqueSlice(slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

// CompareSlices compares two slices for equality
func (sp *SliceProcessor[T]) CompareSlices(a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// StringBuilder is an enhanced string builder with generic methods
type StringBuilder struct {
	builder strings.Builder
}

// NewStringBuilder creates a new enhanced string builder
func NewStringBuilder() *StringBuilder {
	return &StringBuilder{}
}

// Append appends any stringable type
func (sb *StringBuilder) Append(values ...interface{}) *StringBuilder {
	for _, v := range values {
		switch val := v.(type) {
		case string:
			sb.builder.WriteString(val)
		case fmt.Stringer:
			sb.builder.WriteString(val.String())
		default:
			sb.builder.WriteString(fmt.Sprintf("%v", val))
		}
	}
	return sb
}

// AppendLine appends values with a newline
func (sb *StringBuilder) AppendLine(values ...interface{}) *StringBuilder {
	sb.Append(values...)
	sb.builder.WriteString("\n")
	return sb
}

// String returns the built string
func (sb *StringBuilder) String() string {
	return sb.builder.String()
}

// Reset clears the builder
func (sb *StringBuilder) Reset() {
	sb.builder.Reset()
}

// Len returns the length of the current content
func (sb *StringBuilder) Len() int {
	return sb.builder.Len()
}
