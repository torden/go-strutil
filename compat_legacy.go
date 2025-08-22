//go:build !go1.7

package strutils

import (
	"errors"
	"fmt"
	"reflect"
)

// LegacyContext provides a minimal context interface for Go < 1.7
type LegacyContext interface {
	Deadline() (deadline int64, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

// legacyContext is a basic implementation for older Go versions
type legacyContext struct {
	done   chan struct{}
	err    error
	values map[interface{}]interface{}
}

// Background returns a basic context for legacy Go versions
func Background() LegacyContext {
	return &legacyContext{
		done:   make(chan struct{}),
		values: make(map[interface{}]interface{}),
	}
}

func (c *legacyContext) Deadline() (deadline int64, ok bool) {
	return 0, false
}

func (c *legacyContext) Done() <-chan struct{} {
	return c.done
}

func (c *legacyContext) Err() error {
	return c.err
}

func (c *legacyContext) Value(key interface{}) interface{} {
	return c.values[key]
}

// DeepEqual provides reflect.DeepEqual for compatibility
func DeepEqual(x, y interface{}) bool {
	return reflect.DeepEqual(x, y)
}

// LegacyError wraps errors for older Go versions without error wrapping
func LegacyError(msg string, err error) error {
	if err == nil {
		return errors.New(msg)
	}
	return fmt.Errorf("%s: %v", msg, err)
}
