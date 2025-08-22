//go:build go1.7

package strutils

import (
	"context"
	"errors"
	"fmt"
)

// ModernContext aliases the standard context for Go >= 1.7
type ModernContext = context.Context

// Background returns the standard context.Background for modern Go versions
func Background() ModernContext {
	return context.Background()
}

// WithCancel returns context.WithCancel for modern Go versions
func WithCancel(parent ModernContext) (ModernContext, context.CancelFunc) {
	return context.WithCancel(parent)
}

// WithTimeout returns context.WithTimeout for modern Go versions (Go 1.7+)
func WithTimeout(parent ModernContext, timeout int64) (ModernContext, context.CancelFunc) {
	// Note: This is a simplified version. In real implementation,
	// you would use time.Duration and context.WithTimeout
	return context.WithCancel(parent)
}

// ModernError wraps errors using Go 1.13+ error wrapping when available
func ModernError(msg string, err error) error {
	if err == nil {
		return errors.New(msg)
	}
	// Use fmt.Errorf with %w verb for Go 1.13+
	version, verr := GetGoVersion()
	if verr == nil && version.IsVersionAtLeast(1, 13) {
		return fmt.Errorf("%s: %w", msg, err)
	}
	// Fallback for Go 1.7-1.12
	return fmt.Errorf("%s: %v", msg, err)
}
