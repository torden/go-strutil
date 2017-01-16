# Simple Go-lang String Util

Just a few methods for helping string handing

`go get github.com/torden/go-strutils`, import it as `"github.com/torden/go-strutils"`, use it as `stringutils`

## AddSlashes
quote string with slashes.

Example:
```go
strutil := strutils.NewStringUrils()
example_str := "a\bcdefgz"
fmt.Println("%v", strutil.AddSlashes(example_str))
```
Run: 
```bash
a\\bcdefgz
```

## StripSlashes
Un-quotes a quoted string. 

Example:
```go
strutil := NewStringUrils()
example_str := "a\\bcdefgz"
fmt.Println("%v", strutil.StripSlashes(example_str))
```
Run: 
```bash
a\bcdefgz
```

## NL2BR
breakstr inserted before looks like space (CRLF , LFCR, SPACE, NL).

Example:
```go
strutil := strutils.NewStringUrils()
example_str := "abc\ndefgh"
fmt.Println("%v", strutil.Nl2Br(example_str))
```
Run: 
```bash
abc<br />defgh
```

## WordWrapSimple , WordWrapAround
Wraps a string to a given number of characters using break characters (TAB, SPACE)

Example:
```go
strutil := strutils.NewStringUrils()
example_str := "The quick brown fox jumped over the lazy dog."
fmt.Printf("%v\n", strutil.WordWrapSimple(example_str, 3, "*"))
fmt.Printf("%v\n", strutil.WordWrapSimple(example_str, 8, "*"))

fmt.Printf("%v\n", strutil.WordWrapAround(example_str, 3, "*"))
fmt.Printf("%v\n", strutil.WordWrapAround(example_str, 8, "*"))
```
Run: 
```bash
The*quick*brown*fox*jumped*over*the*lazy*dog.
The quick*brown fox*jumped over*the lazy*dog.

The*quick*brown*fox*jumped*over*the*lazy*dog.
The quick*brown fox*jumped*over the*lazy*dog.
```

##NumberFmt
format a number with english notation grouped thousands

Example:
```go
strutil := strutils.NewStringUrils()
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
Run: 
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

Example:
```go
strutil := strutils.NewStringUrils()

example_str := "Life isn't always what one like."

fmt.Printf("%v\n", strutil.PaddingBoth(example_str, "*", PAD_BOTH, 38))
fmt.Printf("%v\n", strutil.PaddingLeft(example_str, "*", PAD_BOTH, 38))
fmt.Printf("%v\n", strutil.PaddingRight(example_str, "*", PAD_BOTH, 38))

fmt.Printf("%v\n", strutil.PaddingBoth(example_str, "*-=", PAD_BOTH, 37))
fmt.Printf("%v\n", strutil.PaddingLeft(example_str, "*-=", PAD_BOTH, 37))
fmt.Printf("%v\n", strutil.PaddingRight(example_str, "*-=", PAD_BOTH, 37))
```
Run: 
```bash
***Life isn't always what one like.***
******Life isn't always what one like.
Life isn't always what one like.******
*-Life isn't always what one like.*-=
*-=*-Life isn't always what one like.
Life isn't always what one like.*-=*-
```
