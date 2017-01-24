# Just! a String Processing Library for Go-lang

Just a few methods for helping processing the string

README.md haven't contain all the examples. Please refer to the the XXXtest.go files.
[[Referrer to Example Code](https://github.com/torden/go-strutil/blob/master/example_test.go)]

[![Build Status](https://travis-ci.org/torden/go-strutil.svg?branch=master)](https://travis-ci.org/torden/go-strutil)
[![Go Report Card](https://goreportcard.com/badge/github.com/torden/go-strutil)](https://goreportcard.com/report/github.com/torden/go-strutil)
[![GoDoc](https://godoc.org/github.com/torden/go-strutil?status.svg)](https://godoc.org/github.com/torden/go-strutil)
[![Coverage Status](https://coveralls.io/repos/github/torden/go-strutil/badge.svg?branch=master)](https://coveralls.io/github/torden/go-strutil?branch=master)
[![Go Walker](http://gowalker.org/api/v1/badge)](https://gowalker.org/github.com/torden/go-strutil)

## Installation

`go get github.com/torden/go-strutils`, import it as `"github.com/torden/go-strutils"`, use it as `StringProc or StringValidator`

## Processing Methods

### AddSlashes

quote string with slashes.

```go
func (s *StringProc) AddSlashes(str string) string
```

Example:

```go
strutil := strutils.NewStringProc()
example_str := "a\bcdefgz"
fmt.Println("%v", strutil.AddSlashes(example_str))
```

The above example will output:

```bash
a\\bcdefgz
```

### StripSlashes
Un-quotes a quoted string.
```go
func (s *StringProc) StripSlashes(str string) string
```

Example:
```go
strutil := NewStringProc()
example_str := "a\\bcdefgz"
fmt.Println("%v", strutil.StripSlashes(example_str))
```
The above example will output:
```bash
a\bcdefgz
```

### NL2BR
breakstr inserted before looks like space (CRLF , LFCR, SPACE, NL).
```go
func (s *StringProc) Nl2Br(str string) string
```

Example:
```go
strutil := strutils.NewStringProc()
example_str := "abc\ndefgh"
fmt.Println("%v", strutil.Nl2Br(example_str))
```
The above example will output:
```bash
abc<br />defgh
```

### WordWrapSimple , WordWrapAround
Wraps a string to a given number of characters using break characters (TAB, SPACE)
```go
func (s *StringProc) WordWrapSimple(str string, wd int, breakstr string) string
func (s *StringProc) WordWrapAround(str string, wd int, breakstr string) string
```

Example:
```go
strutil := strutils.NewStringProc()
example_str := "The quick brown fox jumped over the lazy dog."
fmt.Printf("%v\n", strutil.WordWrapSimple(example_str, 3, "*"))
fmt.Printf("%v\n", strutil.WordWrapSimple(example_str, 8, "*"))

fmt.Printf("%v\n", strutil.WordWrapAround(example_str, 3, "*"))
fmt.Printf("%v\n", strutil.WordWrapAround(example_str, 8, "*"))
```
The above example will output:
```bash
The*quick*brown*fox*jumped*over*the*lazy*dog.
The quick*brown fox*jumped over*the lazy*dog.

The*quick*brown*fox*jumped*over*the*lazy*dog.
The quick*brown fox*jumped*over the*lazy*dog.
```

## NumberFmt
format a number with english notation grouped thousands
```go
func (s *StringProc) NumberFmt(obj interface{}) (string, error)
```

Example:
```go
strutil := strutils.NewStringProc()
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
    retval, err := strutil.NumberFmt(k)
    if v != retval {
        fmt.Errorf("Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
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

## PaddingBoth , PaddingLeft, PaddingRight
pad a string to a certain length with another string
```go
func (s *StringProc) PaddingBoth(str string, fill string, mx int) string
func (s *StringProc) PaddingLeft(str string, fill string, mx int) string
func (s *StringProc) PaddingRight(str string, fill string, mx int) string
```

Example:
```go
strutil := strutils.NewStringProc()
example_str := "Life isn't always what one like."

fmt.Printf("%v\n", strutil.PaddingBoth(example_str, "*", 38))
fmt.Printf("%v\n", strutil.PaddingLeft(example_str, "*", 38))
fmt.Printf("%v\n", strutil.PaddingRight(example_str, "*", 38))

fmt.Printf("%v\n", strutil.PaddingBoth(example_str, "*-=", 37))
fmt.Printf("%v\n", strutil.PaddingLeft(example_str, "*-=", 37))
fmt.Printf("%v\n", strutil.PaddingRight(example_str, "*-=", 37))
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
strutil := strutils.NewStringProc()
example_str := "LIFE ISN'T ALWAYS WHAT ONE LIKE."
fmt.Printf("%v\n", strutil.LowerCaseFirstWords(example_str))
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
strutil := strutils.NewStringProc()
example_str := "life isn't always what one like."
fmt.Printf("%v\n", strutil.UpperCaseFirstWords(example_str))
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
strutil := strutils.NewStringProc()
example_str := "O SAY, CAN YOU SEE, BY THE DAWN’S EARLY LIGHT,"
fmt.Printf("%v\n", strutil.UpperCaseFirstWords(example_str))
```
The above example will output:
```bash
o sAY, cAN yOU sEE, bY tHE dAWN’S eARLY lIGHT,
```

### HumanByteSize
Byte Size convert to Easy Readable Size String

```go
func (s *StringProc) HumanByteSize(obj interface{}, decimals int, unit uint8) (string, error)
```

Example:
```go
strutil := strutils.NewStringProc()
example_str := 3276537856
fmt.Printf("%v\n", strutil.HumanByteSize(k, 2, CamelCaseDouble)
```
The above example will output:
```bash
3.05Gb
```

### HumanFileSize
File Size convert to Easy Readable Size String

```go
func (s *StringProc) HumanFileSize(filepath string, decimals int, unit uint8) (string, error)
```

Example:
```go
strutil := strutils.NewStringProc()
example_str := 3276537856
fmt.Printf("%v\n", strutil.HumanFileSize("/tmp/java.tomcat.core", 2, CamelCaseDouble)
```
The above example will output:
```bash
3.05Gb
```

### AnyCompare

AnyCompare is compares two same basic type (without prt) dataset (slice,map,single data).

```go
func (s *StringProc) AnyCompare(obj1 interface{}, obj2 interface{}) (bool, error)
```

Example:
```go
strutil := strutils.NewStringProc()

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

retval, err = strproc.AnyCompare(testComplexMap1, testComplexMap2)

fmt.Println("Return : ", retval)
fmt.Println("Error : ", err)


```
The above example will output:
```bash
Return :  false
Error :  different value : (obj1[A][name][first][last][F][name][first] := 1) != (obj2[A][name][first][last][F][name][first] := 11)
```

----

## Validation Methods
### IsValidEmail
IsValidEmail is Validates whether the value is a valid e-mail address.

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
IsValidDomain is Validates whether the value is a valid domain address
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
IsValidURL is Validates whether the value is a valid url
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
IsValidMACAddr is Validates whether the value is a valid h/w mac address
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
IsValidIPAddr is Validates whether the value to be exactly a given validation type (IPv4, IPv6, IPv4MappedIPv6, IPv4CIDR, IPv6CIDR, IPv4MappedIPv6CIDR OR IPAny)
```go
func (s *StringValidator) IsValidIPAddr(str string, cktypes ...int) (bool, error)
```

Example:
```go
strvalidator := strutils.NewStringValidator()
example_str := "2001:470:1f09:495::3:217.126.185.21"
fmt.Printf("%v\n", strvalidator.IsValidIPAddr(example_str,strutils.IPv4MappedIPv6,strutils.IPv4))
```
The above example will output:
```bash
true
```

### IsValidFilePath
IsValidFilePath is Validates whether the value is a valid FilePath without relative path
```go
func (s *StringValidator) IsValidFilePath(str string) bool
```

Example:
```go
strvalidator := strutils.NewStringValidator()
example_str := "a-1-s-d-v-we-wd_+qwd-qwd-qwd.txt
fmt.Printf("%v\n", strvalidator.IsValidFilePath(example_str))
```
The above example will output:
```bash
true
```

### IsValidFilePathWithRelativePath
IsValidFilePathWithRelativePath is Validates whether the value is a valid FilePath (allow with relative path)
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

```

### IsPureTextStrict
IsPureTextStrict is Validates whether the value is a pure text, Validation use native
```go
func (s *StringValidator) IsPureTextStrict(str string) (bool, error)
```

Example:
```go
strvalidator := strutils.NewStringValidator()
example_str := `abcd/>qwdqwdoijhwer/>qwdojiqwdqwd</a>qwdoijqwdoiqjd`
fmt.Printf("%v\n", strvalidator.IsPureTextStrict(example_str))
```
The above example will output:
```bash
false
```

### IsPureTextNormal
IsPureTextNormal is Validates whether the value is a pure text, Validation use Regular Expressions
```go
func (s *StringValidator) IsPureTextNormal(str string) (bool, error)
```

Example:
```go
strvalidator := strutils.NewStringValidator()
example_str := `Foo<script type="text/javascript">alert(1337)</script>Bar`
fmt.Printf("%v\n", strvalidator.IsPureTextNormal(example_str))
```
The above example will output:
```bash
false
```

----
Please feel free
