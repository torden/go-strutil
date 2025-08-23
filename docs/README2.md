# Just! a String Processing Library for Go-lang

Just several methods for helping processing/handling the string in Go (go-lang)

README.md haven't contain all the examples. Please refer to the XXXtest.go files.

[![Build Status](https://github.com/torden/go-strutil/actions/workflows/go.yml/badge.svg)](https://github.com/torden/go-strutil/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/torden/go-strutil)](https://goreportcard.com/report/github.com/torden/go-strutil)
[![GoDoc](https://godoc.org/github.com/torden/go-strutil?status.svg)](https://godoc.org/github.com/torden/go-strutil)
[![codecov](https://codecov.io/gh/torden/go-strutil/branch/master/graph/badge.svg?token=tyG0IJBH11)](https://codecov.io/gh/torden/go-strutil)
[![Coverage Status](https://coveralls.io/repos/github/torden/go-strutil/badge.svg?branch=master)](https://coveralls.io/github/torden/go-strutil?branch=master)
[![CodeQL](https://github.com/torden/go-strutil/workflows/CodeQL/badge.svg)](https://github.com/torden/go-strutil/actions/workflows/codeql-analysis.yml)
[![GitHub version](https://img.shields.io/github/v/release/torden/go-strutil)](https://github.com/torden/go-strutil)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Multi-version Go support**: Compatible with Go 1.6 to Go 1.21+
- **Security-focused**: Enhanced validation and XSS/injection protection  
- **Cross-platform**: Automated testing on Ubuntu, Rocky Linux, CentOS, Alpine, macOS
- **High performance**: Optimized algorithms with comprehensive benchmarks
- **Concurrency-safe**: Thread-safe operations with race condition testing
- **Comprehensive**: 50+ string processing and validation methods

## Table of Contents

- [Installation](#installation)
- [Go Version Compatibility](#go-version-compatibility)
- [Testing & Quality](#testing--quality)
- [Examples](#examples)
- Methods
  - [Processing Methods](#processing-methods)
    - [AddSlashes](#addslashes)
    - [StripSlashes](#stripslashes)
    - [NL2BR](#nl2br)
    - [BR2NL](#br2nl)
    - [WordWrapSimple , WordWrapAround](#wordwrapsimple--wordwraparound)
    - [NumberFmt](#numberfmt)
    - [PaddingBoth , PaddingLeft, PaddingRight](#paddingboth--paddingleft-paddingright)
    - [LowerCaseFirstWords](#lowercasefirstwords)
    - [UpperCaseFirstWords](#uppercasefirstwords)
    - [SwapCaseFirstWords](#swapcasefirstwords)
    - [HumanByteSize](#humanbytesize)
    - [HumanFileSize](#humanfilesize)
    - [AnyCompare](#anycompare)
    - [DecodeUnicodeEntities](#decodeunicodeentities)
    - [DecodeURLEncoded](#decodeurlencoded)
    - [StripTags](#striptags)
    - [ConvertToStr](#converttostr)
    - [ReverseStr](#reversestr)
    - [ReverseNormalStr](#reversenormalstr)
    - [ReverseUnicode](#reverseunicode)
    - [FileMD5Hash](#filemd5hash)
    - [MD5Hash](#md5hash)
    - [RegExpNamedGroups](#regexpnamedgroups)
  - [Validation Methods](#validation-methods)
    - [IsValidEmail](#isvalidemail)
    - [IsValidDomain](#isvaliddomain)
    - [IsValidURL](#isvalidurl)
    - [IsValidMACAddr](#isvalidmacaddr)
    - [IsValidIPAddr](#isvalidipaddr)
    - [IsValidFilePath](#isvalidfilepath)
    - [IsValidFilePathWithRelativePath](#isvalidfilepathwithrelativepath)
    - [IsPureTextStrict](#ispuretextstrict)
    - [IsPureTextNormal](#ispuretextnormal)
  - [Assertion Methods](#assertion-methods)
    - [AssertLog](#assertlog)
    - [AssertEquals](#assertequals)
    - [AssertNotEquals](#assertnotequals)
    - [AssertFalse](#assertfalse)
    - [AssertTrue](#asserttrue)
    - [AssertNil](#assertnil)
    - [AssertNotNil](#assertnotnil)
    - [AssertLessThan](#assertlessthan)
    - [AssertLessThanEqualTo](#assertlessthanequalto)
    - [AssertGreaterThan](#assertgreaterthan)
    - [AssertGreaterThanEqualTo](#assertgreaterthanequalto)
    - [AssertLengthOf](#assertlengthof)

## Installation

```bash
go get github.com/torden/go-strutil
```

Import it as `"github.com/torden/go-strutil"`, use it as `strutils.NewStringProc()` or `strutils.NewStringValidator()`

### Requirements

- Go 1.6+ (supports up to Go 1.21+)
- No external dependencies for core functionality

## Go Version Compatibility

This library supports a wide range of Go versions with automatic compatibility detection:

- **Go 1.6-1.7**: Legacy support with interface-based implementations
- **Go 1.18+**: Modern support with generics and enhanced performance
- **Automatic detection**: Uses build tags for optimal performance per Go version

## Testing & Quality

### Comprehensive Test Suite

- **91.1%+ code coverage** with 1000+ test cases
- **9 categories of enhanced tests**:
  - Edge case testing (NumberFmt, file processing)
  - Security validation (XSS, injection, Unicode attacks)
  - IP address security (IPv4/IPv6, CIDR, injection attempts)  
  - Email validation (RFC compliance, internationalization)
  - Concurrent usage and race condition detection
  - Performance and memory benchmarks
  - Cross-platform compatibility

### Cross-Platform Testing

Automated testing on multiple platforms:

```bash
# Run all cross-platform tests
make test-cross-platform

# Test specific platforms
make test-ubuntu      # Ubuntu latest
make test-alpine      # Alpine Linux  
make test-centos      # CentOS Stream
make test-multiarch   # Multiple architectures

# Performance benchmarks across platforms
make bench-cross-platform
```

### Available Test Commands

```bash
# Local development
./run-tests.sh --local          # Local tests only
./run-tests.sh --podman         # Podman-based tests  
./run-tests.sh --all           # All available tests

# Makefile targets
make test                      # Basic tests
make test-race                 # Race condition detection
make test-coverage            # Coverage report
make bench                    # Benchmarks
make ci-extended              # Full CI pipeline
```

## Examples

See the [Example Source](https://github.com/torden/go-strutil/blob/master/example_test.go) for more details

## Processing Methods

### AddSlashes

Quote string with slashes.

```go
func (s *StringProc) AddSlashes(str string) string
```

Example:

```go
strproc := strutils.NewStringProc()
example_str := "a\\bcdefgz"
fmt.Printf("%v", strproc.AddSlashes(example_str))
```

The above example will output:

```bash
a\\\\bcdefgz
```

### StripSlashes

Un-quotes a quoted string.

```go
func (s *StringProc) StripSlashes(str string) string
```

Example:

```go
strproc := strutils.NewStringProc()
example_str := "a\\\\bcdefgz"
fmt.Printf("%v", strproc.StripSlashes(example_str))
```

The above example will output:

```bash
a\bcdefgz
```

### NL2BR

Breakstr inserted before looks like space (CRLF , LFCR, SPACE, NL).

```go
func (s *StringProc) Nl2Br(str string) string
```

Example:

```go
strproc := strutils.NewStringProc()
example_str := "abc\\ndefgh"
fmt.Printf("%v", strproc.Nl2Br(example_str))
```

The above example will output:

```bash
abc<br />defgh
```

### BR2NL

Replaces HTML line breaks to a newline

```go
func (s *StringProc) Br2Nl(str string) string
```

Example:

```go
strproc := strutils.NewStringProc()
example_str1 := "abc<br>defgh"
fmt.Printf("%v", strproc.Br2Nl(example_str1))

example_str2 := "abc<br />defgh"
fmt.Printf("%v", strproc.Br2Nl(example_str2))

example_str3 := "abc<br/>defgh"
fmt.Printf("%v", strproc.Br2Nl(example_str3))
```

The above example will output:

```bash
abc\ndefgh
abc\ndefgh
abc\ndefgh
```

### WordWrapSimple , WordWrapAround

Wraps a string to a given number of characters using break characters (TAB, SPACE)

```go
func (s *StringProc) WordWrapSimple(str string, wd int, breakstr string) (string, error)
func (s *StringProc) WordWrapAround(str string, wd int, breakstr string) (string, error)
```

Example:

```go
strproc := strutils.NewStringProc()
example_str := "The quick brown fox jumped over the lazy dog."
retval1, _ := strproc.WordWrapSimple(example_str, 3, "*")
fmt.Printf("%v\n", retval1)
retval2, _ := strproc.WordWrapSimple(example_str, 8, "*")
fmt.Printf("%v\n", retval2)

retval3, _ := strproc.WordWrapAround(example_str, 3, "*")
fmt.Printf("%v\n", retval3)
retval4, _ := strproc.WordWrapAround(example_str, 8, "*")
fmt.Printf("%v\n", retval4)
```

The above example will output:

```bash
The*quick*brown*fox*jumped*over*the*lazy*dog.
The quick*brown fox*jumped over*the lazy*dog.

The*quick*brown*fox*jumped*over*the*lazy*dog.
The quick*brown fox*jumped*over the*lazy*dog.
```

### NumberFmt

Format a number with english notation grouped thousands

```go
func (s *StringProc) NumberFmt(obj interface{}) (string, error)
```

Example:

```go
strproc := strutils.NewStringProc()
dataset := map[interface{}]string{
    123456789101112: "123,456,789,101,112",
    123456.1234:     "123,456.1234",
    -123456.1234:    "-123,456.1234",
    1.1234561e+06:   "1.1234561e+06",
    1234.1234:       "1,234.1234",
    12345.1234:      "12,345.1234",
    -1.1234561e+06:  "-1.1234561e+06",
    -12345.16:       "-12,345.16",
    12345.16:        "12,345.16",
    1234:            "1,234",
    12.12123098123:  "12.12123098123",
    1.212e+24:       "1.212e+24",
    123456789:       "123,456,789",
}

for k, v := range dataset {
    retval, err := strproc.NumberFmt(k)
    if v != retval {
        fmt.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", v, retval)
    } else if err != nil {
        fmt.Errorf("Return Error : %v", err)
    } else {
        fmt.Printf("%v\n", retval)
    }
}
```

The above example will output:

```bash
123,456,789,101,112
123,456.1234
-123,456.1234
1.1234561e+06
1,234.1234
12,345.1234
-1.1234561e+06
-12,345.16
12,345.16
1,234
12.12123098123
1.212e+24
123,456,789
```

### PaddingBoth , PaddingLeft, PaddingRight

Pad a string to a certain length with another string

```go
func (s *StringProc) PaddingBoth(str string, fill string, mx int) string
func (s *StringProc) PaddingLeft(str string, fill string, mx int) string
func (s *StringProc) PaddingRight(str string, fill string, mx int) string
```

Example:

```go
strproc := strutils.NewStringProc()
example_str := "Life isn't always what one like."

fmt.Printf("%v\n", strproc.PaddingBoth(example_str, "*", 38))
fmt.Printf("%v\n", strproc.PaddingLeft(example_str, "*", 38))
fmt.Printf("%v\n", strproc.PaddingRight(example_str, "*", 38))

fmt.Printf("%v\n", strproc.PaddingBoth(example_str, "*-=", 37))
fmt.Printf("%v\n", strproc.PaddingLeft(example_str, "*-=", 37))
fmt.Printf("%v\n", strproc.PaddingRight(example_str, "*-=", 37))
```

The above example will output:

```bash
***Life isn't always what one like.***
******Life isn't always what one like.
Life isn't always what one like.******
*-Life isn't always what one like.*-=
*-=*-Life isn't always what one like.
Life isn't always what one like.*-=*-
```

### LowerCaseFirstWords

Lowercase the first character of each word in a string

```go
// TOKEN : \t \r \n \f \v \s
func (s *StringProc) LowerCaseFirstWords(str string) string
```

Example:

```go
strproc := strutils.NewStringProc()
example_str := "LIFE ISN'T ALWAYS WHAT ONE LIKE."
fmt.Printf("%v\n", strproc.LowerCaseFirstWords(example_str))
```

The above example will output:

```bash
lIFE iSN'T aLWAYS wHAT oNE lIKE.
```

### UpperCaseFirstWords

Uppercase the first character of each word in a string

```go
// TOKEN : \t \r \n \f \v \s
func (s *StringProc) UpperCaseFirstWords(str string) string
```

Example:

```go
strproc := strutils.NewStringProc()
example_str := "life isn't always what one like."
fmt.Printf("%v\n", strproc.UpperCaseFirstWords(example_str))
```

The above example will output:

```bash
Life Isn't Always What One Like.
```

### SwapCaseFirstWords

Switch the first character case of each word in a string

```go
// TOKEN : \t \r \n \f \v \s
func (s *StringProc) SwapCaseFirstWords(str string) string
```

Example:

```go
strproc := strutils.NewStringProc()
example_str := "O SAY, CAN YOU SEE, BY THE DAWN'S EARLY LIGHT,"
fmt.Printf("%v\n", strproc.SwapCaseFirstWords(example_str))
```

The above example will output:

```bash
o sAY, cAN yOU sEE, bY tHE dAWN'S eARLY lIGHT,
```

### HumanByteSize

Byte Size convert to Easy Readable Size String

```go
func (s *StringProc) HumanByteSize(obj interface{}, decimals int, unit uint8) (string, error)
```

Example:

```go
strproc := strutils.NewStringProc()
example_bytes := 3276537856
retval, err := strproc.HumanByteSize(example_bytes, 2, strutils.CamelCaseLong)
if err != nil {
    fmt.Printf("Error: %v", err)
} else {
    fmt.Printf("%v\n", retval)
}
```

The above example will output:

```bash
3.05GigaByte
```

### HumanFileSize

File Size convert to Easy Readable Size String

```go
func (s *StringProc) HumanFileSize(filepath string, decimals int, unit uint8) (string, error)
```

Example:

```go
strproc := strutils.NewStringProc()
retval, err := strproc.HumanFileSize("./LICENSE", 2, strutils.CamelCaseLong)
if err != nil {
    fmt.Printf("Error: %v", err)
} else {
    fmt.Printf("%v\n", retval)
}
```

The above example will output:

```bash
1.05KiloByte
```

### AnyCompare

AnyCompare compares two same basic type (without ptr) dataset (slice,map,single data).

```go
func (s *StringProc) AnyCompare(obj1 interface{}, obj2 interface{}) (bool, error)
```

Example:

```go
strproc := strutils.NewStringProc()

testComplexMap1 := map[string]map[string]map[string]int{
    "F": map[string]map[string]int{
        "name": map[string]int{
            "first": 1,
            "last":  2,
        },
    },
    "A": map[string]map[string]int{
        "name": map[string]int{
            "first": 11,
            "last":  21,
        },
    },
}

testComplexMap2 := map[string]map[string]map[string]int{
    "F": map[string]map[string]int{
        "name": map[string]int{
            "first": 11,
            "last":  12222,
        },
    },
    "A": map[string]map[string]int{
        "name": map[string]int{
            "first": 11,
            "last":  21,
        },
    },
}

retval, err := strproc.AnyCompare(testComplexMap1, testComplexMap2)

fmt.Println("Return : ", retval)
fmt.Println("Error : ", err)
```

The above example will output:

```bash
Return :  false
Error :  Different Value : (obj1[F][name][first] := 1) != (obj2[F][name][first] := 11)
```

### DecodeUnicodeEntities

DecodeUnicodeEntities Decodes Unicode Entities

```go
func (s *StringProc) DecodeUnicodeEntities(val string) (string, error)
```

Example:

```go
strproc := strutils.NewStringProc()
strUnicodeEntityEncodedMultipleLine := "%uC548%uB155%uD558%uC138%uC694.%0A%uBC29%uAC11%uC2B5%uB2C8%uB2E4.%0A%uAC10%uC0AC%uD569%uB2C8%uB2E4.%0A%u304A%u306F%u3088%u3046%u3054%u3056%u3044%u307E%u3059%0A%u3053%u3093%u306B%u3061%u306F%uFF0E%0A%u3053%u3093%u3070%u3093%u306F%uFF0E%0A%u304A%u3084%u3059%u307F%u306A%u3055%u3044%uFF0E%0A%u3042%u308A%u304C%u3068%u3046%u3054%u3056%u3044%u307E%u3059%0A%u4F60%u597D%0A%u518D%u898B%0A%u8C22%u8C22%21%u0E2A%u0E27%u0E31%u0E2A%u0E14%u0E35%u0E04%u0E23%u0E31%u0E1A%0A%u0E41%u0E25%u0E49%u0E27%u0E40%u0E08%u0E2D%u0E01%u0E31%u0E19%u0E04%u0E23%u0E31%u0E1A%0A%u0E02%u0E2D%u0E1A%u0E04%u0E38%u0E13%u0E04%u0E23%u0E31%u0E1A%0A%u0421%u0430%u0439%u043D%20%u0431%u0430%u0439%u043D%u0430%u0443%u0443"

retval, err := strproc.DecodeUnicodeEntities(strUnicodeEntityEncodedMultipleLine)

fmt.Println("Return : ", retval)
fmt.Println("Error : ", err)
```

The above example will output:

```bash
Return : 안녕하세요.
방갑습니다.
감사합니다.
おはようございます
こんにちは．
こんばんは．
おやすみなさい．
ありがとうございます
你好
再見
谢谢!สวัสดีครับ
แล้วเจอกันครับ
ขอบคุณครับ
Сайн байнауу
Error : <nil>
```

### DecodeURLEncoded

DecodeURLEncoded Decodes URL-encoded string (including unicode entities)

```go
func (s *StringProc) DecodeURLEncoded(val string) (string, error)
```

Example:

```go
strproc := strutils.NewStringProc()
urlWithJapanWorld := "http://hello.%E4%B8%96%E7%95%8C.com/foo"

retval, err := strproc.DecodeURLEncoded(urlWithJapanWorld)

fmt.Println("Return : ", retval)
fmt.Println("Error : ", err)
```

The above example will output:

```bash
Return : http://hello.世界.com/foo
Error : <nil>
```

### StripTags

StipTags remove all tag in string (Pure String or URL Encoded or Html (Unicode) Entities Encoded or Mixed String)

```go
func (s *StringProc) StripTags(str string) (string, error)
```

Example:

```go
strproc := strutils.NewStringProc()
example_str := `
<!DOCTYPE html>
<html lang="en-us">
<head>
<meta charset="UTF-8">
<title>                            Just! a String Processing Library for Go-lang</title>
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="theme-color" content="#157878">
<link href='https://fonts.googleapis.com/css?family=Open+Sans:400,700' rel='stylesheet' type='text/css'>
<link rel="stylesheet" href="/go-strutil/assets/css/style.css?v=dae229423409070462d2ce364eba3b5721930df0">
</head>
<body>
<section class="page-header">
<h1 class="project-name">Just! a String Processing Library for Go-lang</h1>
<h2 class="project-tagline">Just a few methods for helping processing and validation the string</h2>
<a href="https://github.com/torden/go-strutil" class="btn">View on GitHub</a>
</section>
<section class="main-content">
<h1 id="just-a-string-processing-library-for-go-lang">Just! a String Processing Library for Go-lang</h1>
<p>Just a few methods for helping processing the string</p>
<p>README.md haven't contain all the examples. Please refer to the the XXXtest.go files.</p>
</body>
</html>
`
retval, err := strproc.StripTags(example_str)
if err != nil {
    fmt.Println("Error : ", err)
} else {
    fmt.Println(retval)
}
```

The above example will output:

```bash
Just! a String Processing Library for Go-lang
Just! a String Processing Library for Go-lang
Just a few methods for helping processing and validation the string
View on GitHub
Just! a String Processing Library for Go-lang
Just a few methods for helping processing the string
README.md haven't contain all the examples. Please refer to the the XXXtest.go files.
```

### ConvertToStr

ConvertToStr converts basic data type to string

```go
func (s *StringProc) ConvertToStr(obj interface{}) (string, error)
```

Example:

```go
strproc := strutils.NewStringProc()
example_val := uint64(1234567)
retval, err := strproc.ConvertToStr(example_val)
if err != nil {
    fmt.Println("Error : ", err)
} else {
    fmt.Println(retval)
}
```

The above example will output:

```bash
"1234567"
```

### ReverseStr

ReverseStr reverses a String, according to value type between ascii (ReverseNormalStr) or rune (ReverseUnicode)

```go
func (s *StringProc) ReverseStr(str string) string
```

Example:

```go
strproc := strutils.NewStringProc()

dataset := []string{
  "0123456789",
  "가나다라마바사",
  "あいうえお",
  "天地玄黃宇宙洪荒",
}

for _, v := range dataset {
  fmt.Println(strproc.ReverseStr(v))
}
```

The above example will output:

```bash
9876543210
사바마라다나가
おえういあ
荒洪宙宇黃玄地天
```

### ReverseNormalStr

ReverseNormalStr reverses a None-unicode String.
Faster than ReverseUnicode or ReverseStr

```go
func (s *StringProc) ReverseNormalStr(str string) string
```

Example:

```go
strproc := strutils.NewStringProc()

dataset := []string{
  "0123456789",
  "abcdefg",
}

for _, v := range dataset {
  fmt.Println(strproc.ReverseNormalStr(v))
}
```

The above example will output:

```bash
9876543210
gfedcba
```

### ReverseUnicode

ReverseUnicode reverses a Unicode String

```go
func (s *StringProc) ReverseUnicode(str string) string
```

Example:

```go
strproc := strutils.NewStringProc()

dataset := []string{
  "0123456789",
  "가나다라마바사",
  "あいうえお",
  "天地玄黃宇宙洪荒",
}

for _, v := range dataset {
  fmt.Println(strproc.ReverseUnicode(v))
}
```

The above example will output:

```bash
9876543210
사바마라다나가
おえういあ
荒洪宙宇黃玄地天
```

### FileMD5Hash

FileMD5Hash calculates MD5 checksum of the file

```go
func (s *StringProc) FileMD5Hash(filepath string) (string, error)
```

Example:

```go
strproc := strutils.NewStringProc()

retval, err := strproc.FileMD5Hash("./LICENSE")
if err != nil {
    fmt.Printf("Error : %v", err)
} else {
    fmt.Println(retval)
}
```

The above example will output:

```bash
f3f8954bac465686f0bfc2a757c5200b
```

### MD5Hash

MD5Hash calculates MD5 checksum of the string

```go
func (s *StringProc) MD5Hash(str string) (string, error)
```

Example:

```go
strproc := strutils.NewStringProc()

dataset := []string{
    "0123456789",
    "abcdefg",
    "abcdefgqwdoisef;oijawe;fijq2039jdfs.dnc;oa283hr08uj3o;ijwaef;owhjefo;uhwefwef",
}

for _, v := range dataset {
    retval, err := strproc.MD5Hash(v)
    if err != nil {
        fmt.Printf("Error : %v", err)
    } else {
        fmt.Println(retval)
    }
}
```

The above example will output:

```bash
781e5e245d69b566979b86e28d23f2c7
7ac66c0f148de9519b8bd264312c4d64
15f764f21d09b11102eb015fc8824d00
```

### RegExpNamedGroups

RegExpNamedGroups captures the text matched by regex into the group name

```go
func (s *StringProc) RegExpNamedGroups(regex *regexp.Regexp, val string) (map[string]string, error)
```

Example:

```go
strproc := strutils.NewStringProc()

regexGoVersion := regexp.MustCompile(`go(?P<major>([0-9]{1,3}))\.(?P<minor>([0-9]{1,3}))(\.(?P<rev>([0-9]{1,3})))?`)
retval, err := strproc.RegExpNamedGroups(regexGoVersion, runtime.Version())
if err != nil {
    fmt.Printf("Error: %v", err)
} else {
    fmt.Println(retval)
}
```

The above example will output:

```bash
map[major:1 minor:21 rev:0]
```

----

## Validation Methods

### IsValidEmail

IsValidEmail validates whether the value is a valid e-mail address.

```go
func (s *StringValidator) IsValidEmail(str string) bool
```

Example:

```go
strvalidator := strutils.NewStringValidator()
example_str := "a@golang.org"
fmt.Printf("%v\n", strvalidator.IsValidEmail(example_str))
```

The above example will output:

```bash
true
```

### IsValidDomain

IsValidDomain validates whether the value is a valid domain address

```go
func (s *StringValidator) IsValidDomain(str string) bool
```

Example:

```go
strvalidator := strutils.NewStringValidator()
example_str := "golang.org"
fmt.Printf("%v\n", strvalidator.IsValidDomain(example_str))
```

The above example will output:

```bash
true
```

### IsValidURL

IsValidURL validates whether the value is a valid url

```go
func (s *StringValidator) IsValidURL(str string) bool
```

Example:

```go
strvalidator := strutils.NewStringValidator()
example_str := "https://www.google.co.kr/url?sa=t&rct=j&q=&esrc=s&source=web"
fmt.Printf("%v\n", strvalidator.IsValidURL(example_str))
```

The above example will output:

```bash
true
```

### IsValidMACAddr

IsValidMACAddr validates whether the value is a valid h/w mac address

```go
func (s *StringValidator) IsValidMACAddr(str string) bool
```

Example:

```go
strvalidator := strutils.NewStringValidator()
example_str := "02-f3-71-eb-9e-4b"
fmt.Printf("%v\n", strvalidator.IsValidMACAddr(example_str))
```

The above example will output:

```bash
true
```

### IsValidIPAddr

IsValidIPAddr validates whether the value to be exactly a given validation types
(IPv4, IPv6, IPv4MappedIPv6, IPv4CIDR, IPv6CIDR, IPv4MappedIPv6CIDR OR IPAny)

```go
func (s *StringValidator) IsValidIPAddr(str string, cktypes ...int) (bool, error)
```

Example:

```go
strvalidator := strutils.NewStringValidator()
example_str := "2001:470:1f09:495::3:217.126.185.21"
result, err := strvalidator.IsValidIPAddr(example_str, strutils.IPv4MappedIPv6, strutils.IPv4)
if err != nil {
    fmt.Printf("Error: %v\n", err)
} else {
    fmt.Printf("%v\n", result)
}
```

The above example will output:

```bash
true
```

### IsValidFilePath

IsValidFilePath validates whether the value is a valid FilePath without relative path

```go
func (s *StringValidator) IsValidFilePath(str string) bool
```

Example:

```go
strvalidator := strutils.NewStringValidator()
example_str := "a-1-s-d-v-we-wd_+qwd-qwd-qwd.txt"
fmt.Printf("%v\n", strvalidator.IsValidFilePath(example_str))
```

The above example will output:

```bash
false
```

### IsValidFilePathWithRelativePath

IsValidFilePathWithRelativePath validates whether the value is a valid FilePath (allow with relative path)

```go
func (s *StringValidator) IsValidFilePathWithRelativePath(str string) bool
```

Example:

```go
strvalidator := strutils.NewStringValidator()
example_str := "/asdasd/asdasdasd/qwdqwd_qwdqwd/12-12/a-1-e-r-t-_1_21234_d_1234_qwed_1423_.txt"
fmt.Printf("%v\n", strvalidator.IsValidFilePathWithRelativePath(example_str))
```

The above example will output:

```bash
true
```

### IsPureTextStrict

IsPureTextStrict validates whether the value is a pure text, validation uses native methods

```go
func (s *StringValidator) IsPureTextStrict(str string) (bool, error)
```

Example:

```go
strvalidator := strutils.NewStringValidator()
example_str := `abcd/>qwdqwdoijhwer/>qwdojiqwdqwd</a>qwdoijqwdoiqjd`
result, err := strvalidator.IsPureTextStrict(example_str)
if err != nil {
    fmt.Printf("Error: %v\n", err)
} else {
    fmt.Printf("%v\n", result)
}
```

The above example will output:

```bash
Error: Detect Tag (<[!|?]~>)
```

### IsPureTextNormal

IsPureTextNormal validates whether the value is a pure text, validation uses Regular Expressions

```go
func (s *StringValidator) IsPureTextNormal(str string) (bool, error)
```

Example:

```go
strvalidator := strutils.NewStringValidator()
example_str := `Foo<script type="text/javascript">alert(1337)</script>Bar`
result, err := strvalidator.IsPureTextNormal(example_str)
if err != nil {
    fmt.Printf("Error: %v\n", err)
} else {
    fmt.Printf("%v\n", result)
}
```

The above example will output:

```bash
Error: Detect HTML Element
```

## Assertion Methods

### AssertLog

AssertLog formats its arguments using default formatting, analogous to t.Log

```go
func (a *Assert) AssertLog(t *testing.T, err error, msgfmt string, args ...interface{})
```

### AssertEquals

AssertEquals asserts that two objects are equal.

```go
func (a *Assert) AssertEquals(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{})
```

### AssertNotEquals

AssertNotEquals asserts that two objects are not equal.

```go
func (a *Assert) AssertNotEquals(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{})
```

### AssertFalse

AssertFalse asserts that the specified value is false.

```go
func (a *Assert) AssertFalse(t *testing.T, v1 bool, msgfmt string, args ...interface{})
```

### AssertTrue

AssertTrue asserts that the specified value is true.

```go
func (a *Assert) AssertTrue(t *testing.T, v1 bool, msgfmt string, args ...interface{})
```

### AssertNil

AssertNil asserts that the specified value is nil.

```go
func (a *Assert) AssertNil(t *testing.T, v1 interface{}, msgfmt string, args ...interface{})
```

### AssertNotNil

AssertNotNil asserts that the specified value isn't nil.

```go
func (a *Assert) AssertNotNil(t *testing.T, v1 interface{}, msgfmt string, args ...interface{})
```

### AssertLessThan

AssertLessThan asserts that the specified value are v1 less than v2

```go
func (a *Assert) AssertLessThan(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{})
```

### AssertLessThanEqualTo

AssertLessThanEqualTo asserts that the specified value are v1 less than v2 or equal to

```go
func (a *Assert) AssertLessThanEqualTo(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{})
```

### AssertGreaterThan

AssertGreaterThan asserts that the specified value are v1 greater than v2

```go
func (a *Assert) AssertGreaterThan(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{})
```

### AssertGreaterThanEqualTo

AssertGreaterThanEqualTo asserts that the specified value are v1 greater than v2 or equal to

```go
func (a *Assert) AssertGreaterThanEqualTo(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{})
```

### AssertLengthOf

AssertLengthOf asserts that object has a length property with the expected value.

```go
func (a *Assert) AssertLengthOf(t *testing.T, v1 interface{}, v2 interface{}, msgfmt string, args ...interface{})
```

## Security Features

This library includes comprehensive security features:

- **XSS Protection**: Advanced validation against script injection
- **SQL Injection Prevention**: Input sanitization and validation  
- **Unicode Attack Detection**: Protection against homograph and normalization attacks
- **Path Traversal Prevention**: Secure file path validation
- **Control Character Filtering**: Protection against binary and control character injection
- **Header Injection Prevention**: HTTP header injection detection

## Performance

Optimized for performance with:

- **Efficient algorithms**: Optimized string processing implementations
- **Memory management**: Minimal memory allocations and garbage collection pressure
- **Concurrent safety**: Thread-safe operations without performance penalties
- **Benchmarking**: Comprehensive performance testing across platforms

See benchmark results:
```bash
make bench-cross-platform  # Cross-platform performance comparison
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make ci-extended`)
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

### Development Tools

```bash
make setup              # Setup development environment
make lint               # Run linter
make test-coverage      # Generate coverage report  
make bench              # Run benchmarks
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

*Please feel free. I hope it is helpful for you*