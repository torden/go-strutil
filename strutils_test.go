package strutils_test

import (
	"github.com/torden/go-strutil"
)

var (
	assert       = strutils.NewAssert()
	strvalidator = strutils.NewStringValidator()
	strproc      = strutils.NewStringProc()
)

// var randnum = rand.New(rand.NewSource(time.Now().UnixNano()))
