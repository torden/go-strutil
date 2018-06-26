package strutils_test

import (
	"sync"

	"github.com/torden/go-strutil"
)

var mutx sync.RWMutex
var assert = strutils.NewAssert()
var strvalidator = strutils.NewStringValidator()
var strproc = strutils.NewStringProc()

//var randnum = rand.New(rand.NewSource(time.Now().UnixNano()))
