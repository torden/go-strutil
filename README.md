# Just String Util for Go-lang

Just a few methods for helping string handing

[![Build Status](https://travis-ci.org/torden/go-strutil.svg?branch=master)](https://travis-ci.org/torden/go-strutil)

## Installation
`go get github.com/torden/go-strutils`, import it as `"github.com/torden/go-strutils"`, use it as `stringutils`


## AddSlashes
quote string with slashes.
```go
func (s *stringUtils) AddSlashes(str string) string
````

Example:
```go
strutil := strutils.NewStringUtils()
example_str := "a\bcdefgz"
fmt.Println("%v", strutil.AddSlashes(example_str))
```
The above example will output:
```bash
a\\bcdefgz
```

## StripSlashes
Un-quotes a quoted string.
```go
func (s *stringUtils) StripSlashes(str string) string
```

Example:
```go
strutil := NewStringUtils()
example_str := "a\\bcdefgz"
fmt.Println("%v", strutil.StripSlashes(example_str))
```
The above example will output:
```bash
a\bcdefgz
```

## NL2BR
breakstr inserted before looks like space (CRLF , LFCR, SPACE, NL).
```go
func (s *stringUtils) Nl2Br(str string) string
```

Example:
```go
strutil := strutils.NewStringUtils()
example_str := "abc\ndefgh"
fmt.Println("%v", strutil.Nl2Br(example_str))
```
The above example will output:
```bash
abc<br />defgh
```

## WordWrapSimple , WordWrapAround
Wraps a string to a given number of characters using break characters (TAB, SPACE)
```go
func (s *stringUtils) WordWrapSimple(str string, wd int, breakstr string) string
func (s *stringUtils) WordWrapAround(str string, wd int, breakstr string) string
```

Example:
```go
strutil := strutils.NewStringUtils()
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

##NumberFmt
format a number with english notation grouped thousands
```go
func (s *stringUtils) NumberFmt(obj interface{}) (string, error)
````

Example:
```go
strutil := strutils.NewStringUtils()
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

##PaddingBoth , PaddingLeft, PaddingRight
pad a string to a certain length with another string
```go
func (s *stringUtils) PaddingBoth(str string, fill string, mx int) string
func (s *stringUtils) PaddingLeft(str string, fill string, mx int) string
func (s *stringUtils) PaddingRight(str string, fill string, mx int) string
```

Example:
```go
strutil := strutils.NewStringUtils()
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

## LowerCaseFirstWords
Lowercase the first character of each word in a string
```go
// TOKEN : \t \r \n \f \v \s
func (s *stringUtils) LowerCaseFirstWords(str string) string
```

Example:
```go
strutil := strutils.NewStringUtils()
example_str := "LIFE ISN'T ALWAYS WHAT ONE LIKE."
fmt.Printf("%v\n", strutil.LowerCaseFirstWords(example_str))
```
The above example will output:
```bash
lIFE iSN'T aLWAYS wHAT oNE lIKE.
```

## UpperCaseFirstWords
Uppercase the first character of each word in a string

```go
// TOKEN : \t \r \n \f \v \s
func (s *stringUtils) UpperCaseFirstWords(str string) string
```

Example:
```go
strutil := strutils.NewStringUtils()
example_str := "life isn't always what one like."
fmt.Printf("%v\n", strutil.UpperCaseFirstWords(example_str))
```
The above example will output:
```bash
Life Isn't Always What One Like.
```

## SwapCaseFirstWords
Switch the first character case of each word in a string

```go
// TOKEN : \t \r \n \f \v \s
func (s *stringUtils) SwapCaseFirstWords(str string) string
```

Example:
```go
strutil := strutils.NewStringUtils()
example_str := "O SAY, CAN YOU SEE, BY THE DAWN’S EARLY LIGHT,"
fmt.Printf("%v\n", strutil.UpperCaseFirstWords(example_str))
```
The above example will output:
```bash
o sAY, cAN yOU sEE, bY tHE dAWN’S eARLY lIGHT,
```







