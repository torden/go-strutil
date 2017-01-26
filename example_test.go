/*
The MIT License (MIT)

Copyright (C) 2016-2017 Torden Cho <https://github.com/torden>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package strutils_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/torden/go-strutil"
)

func Example_strutils_AddSlashes() {

	strproc := strutils.NewStringProc()
	example_str := `a\bcdefgz`
	fmt.Println(strproc.AddSlashes(example_str))
	// Output: a\\bcdefgz
}

func Example_strutils_StripSlashes() {

	strproc := strutils.NewStringProc()
	example_str := "a\\bcdefgz"
	fmt.Println(strproc.StripSlashes(example_str))
	// Output: abcdefgz
}

func Example_strutils_Nl2Br() {

	strproc := strutils.NewStringProc()
	example_str := "abc\ndefgh"
	fmt.Println(strproc.Nl2Br(example_str))
	// Output: abc<br />defgh
}

func Example_strutils_WordWrapSimple() {

	strproc := strutils.NewStringProc()
	example_str := "The quick brown fox jumped over the lazy dog."

	var retval string
	retval, _ = strproc.WordWrapSimple(example_str, 3, "*")
	fmt.Printf("%v\n", retval)

	retval, _ = strproc.WordWrapSimple(example_str, 8, "*")
	fmt.Printf("%v\n", retval)

	// Output: The*quick*brown*fox*jumped*over*the*lazy*dog.
	// The quick*brown fox*jumped over*the lazy*dog.
}

func Example_strutils_WordWrapAround() {

	strproc := strutils.NewStringProc()
	example_str := "The quick brown fox jumped over the lazy dog."

	var retval string
	var err error
	retval, err = strproc.WordWrapAround(example_str, 3, "*")
	if err != nil {
		fmt.Println("Error : ", err)
	} else {
		fmt.Printf("%v\n", retval)
	}

	retval, _ = strproc.WordWrapAround(example_str, 8, "*")
	if err != nil {
		fmt.Println("Error : ", err)
	} else {
		fmt.Printf("%v\n", retval)
	}

	// Output: The*quick*brown*fox*jumped*over*the*lazy*dog.
	// The quick*brown fox*jumped*over the*lazy*dog.
}

func Example_strutils_NumberFmt() {

	strproc := strutils.NewStringProc()

	var retval string

	retval, _ = strproc.NumberFmt(123456789101112)
	fmt.Println(retval)
	// Output: 123,456,789,101,112

	retval, _ = strproc.NumberFmt(123456.1234)
	fmt.Println(retval)
	//123,456.1234

	retval, _ = strproc.NumberFmt(-123456.1234)
	fmt.Println(retval)
	//-123,456.1234

	retval, _ = strproc.NumberFmt(1.1234561e+06)
	fmt.Println(retval)
	//1.1234561e+06

	retval, _ = strproc.NumberFmt(1234.1234)
	fmt.Println(retval)
	//1,234.1234

	retval, _ = strproc.NumberFmt(12345.1234)
	fmt.Println(retval)
	//12,345.1234

	retval, _ = strproc.NumberFmt(-1.1234561e+06)
	fmt.Println(retval)
	//-1.1234561e+06

	retval, _ = strproc.NumberFmt(-12345.16)
	fmt.Println(retval)
	//-12,345.16

	retval, _ = strproc.NumberFmt(12345.16)
	fmt.Println(retval)
	//12,345.16

	retval, _ = strproc.NumberFmt(1234)
	fmt.Println(retval)
	//1,234

	retval, _ = strproc.NumberFmt(12.12123098123)
	fmt.Println(retval)
	//12.12123098123

	retval, _ = strproc.NumberFmt(1.212e+24)
	fmt.Println(retval)
	//1.212e+24

	retval, _ = strproc.NumberFmt(123456789)
	fmt.Println(retval)
	//123,456,789

	retval, _ = strproc.NumberFmt("123456789101112")
	fmt.Println(retval)
	//123,456,789,101,112

	retval, _ = strproc.NumberFmt("123456.1234")
	fmt.Println(retval)
	//123,456.1234

	retval, _ = strproc.NumberFmt("-123456.1234")
	fmt.Println(retval)
	//-123,456.1234

	retval, _ = strproc.NumberFmt("1.1234561e+06")
	fmt.Println(retval)
	//1.1234561e+06

	retval, _ = strproc.NumberFmt("1234.1234")
	fmt.Println(retval)
	//1,234.1234

	retval, _ = strproc.NumberFmt("12345.1234")
	fmt.Println(retval)
	//12,345.1234

	retval, _ = strproc.NumberFmt("-1.1234561e+06")
	fmt.Println(retval)
	//-1.1234561e+06

	retval, _ = strproc.NumberFmt("-12345.16")
	fmt.Println(retval)
	//-12,345.16

	retval, _ = strproc.NumberFmt("12345.16")
	fmt.Println(retval)
	//12,345.16

	retval, _ = strproc.NumberFmt("1234")
	fmt.Println(retval)
	//1,234

	retval, _ = strproc.NumberFmt("12.12123098123")
	fmt.Println(retval)
	//12.12123098123

	retval, _ = strproc.NumberFmt("1.212e+24")
	fmt.Println(retval)
	//1.212e+24

	retval, _ = strproc.NumberFmt("123456789")
	fmt.Println(retval)
	//123,456,789

}

func Example_strutils_PaddingBoth() {

	strproc := strutils.NewStringProc()
	example_str := "Life isn't always what one like."

	fmt.Printf("%v\n", strproc.PaddingBoth(example_str, "*", 38))
	fmt.Printf("%v\n", strproc.PaddingBoth(example_str, "*-=", 37))

	// Output: ***Life isn't always what one like.***
	// *-Life isn't always what one like.*-=
}

func Example_strutils_PaddingLeft() {

	strproc := strutils.NewStringProc()
	example_str := "Life isn't always what one like."

	fmt.Printf("%v\n", strproc.PaddingLeft(example_str, "*", 38))
	fmt.Printf("%v\n", strproc.PaddingLeft(example_str, "*-=", 37))

	// Output: ******Life isn't always what one like.
	// *-=*-Life isn't always what one like.
}

func Example_strutils_PaddingRight() {

	strproc := strutils.NewStringProc()
	example_str := "Life isn't always what one like."

	fmt.Printf("%v\n", strproc.PaddingRight(example_str, "*", 38))
	fmt.Printf("%v\n", strproc.PaddingRight(example_str, "*-=", 37))

	// Output: Life isn't always what one like.******
	// Life isn't always what one like.*-=*-
}

func Example_strutils_LowerCaseFirstWords() {

	strproc := strutils.NewStringProc()
	example_str := "LIFE ISN'T ALWAYS WHAT ONE LIKE."
	fmt.Printf("%v\n", strproc.LowerCaseFirstWords(example_str))
	// Output: lIFE iSN'T aLWAYS wHAT oNE lIKE.
}

func Example_strutils_UpperCaseFirstWords() {

	strproc := strutils.NewStringProc()
	example_str := "life isn't always what one like."
	fmt.Printf("%v\n", strproc.UpperCaseFirstWords(example_str))
	// Output: Life Isn't Always What One Like.
}

func Example_strutils_SwapCaseFirstWords() {

	strproc := strutils.NewStringProc()
	example_str := "O SAY, CAN YOU SEE, BY THE DAWN’S EARLY LIGHT,"
	fmt.Printf("%v\n", strproc.UpperCaseFirstWords(example_str))
	// Output: O SAY, CAN YOU SEE, BY THE DAWN’S EARLY LIGHT,
}

func Example_strutils_HumanByteSize() {

	strproc := strutils.NewStringProc()
	example_str := 3276537856
	retval, err := strproc.HumanByteSize(example_str, 2, strutils.CamelCaseLong)
	if err != nil {
		fmt.Println("Error : ", err)
	} else {
		fmt.Println(retval)
	}
	// Output: 3.05GigaByte
}

func Example_strutils_HumanFileSize() {

	const tmpFilePath = "./filesizecheck.touch"
	var retval string
	var err error

	//generating a touch file
	tmpdata := []byte("123456789")

	ioutil.WriteFile(tmpFilePath, tmpdata, 0750)

	strproc := strutils.NewStringProc()
	retval, err = strproc.HumanFileSize(tmpFilePath, 2, strutils.CamelCaseLong)
	if err != nil {
		fmt.Println("Error : ", err)
	} else {
		fmt.Println(retval)
	}

	retval, err = strproc.HumanFileSize(tmpFilePath, 2, strutils.CamelCaseDouble)
	if err != nil {
		fmt.Println("Error : ", err)
	} else {
		fmt.Println(retval)
	}

	os.Remove(tmpFilePath)

	// Output: 9.00Byte
	// 9.00B

}

func Example_strutils_AnyCompare() {

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

	var retval bool
	var err error

	retval, err = strproc.AnyCompare(testComplexMap1, testComplexMap2)
	fmt.Println("Return : ", retval)
	fmt.Println("Error : ", err)

	// Output: Return :  false
	// Error :  Different Value : (obj1[F][name][first] := 1) != (obj2[F][name][first] := 11)

	testSliceInt1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	testSliceInt2 := []int{1, 2, 3, 4, 5, 6, 7, 8, 8}

	retval, err = strproc.AnyCompare(testSliceInt1, testSliceInt2)
	fmt.Println("Return : ", retval)
	fmt.Println("Error : ", err)
	// Return :  false
	// Error :  Different Value : (obj1[8] := 9) != (obj2[8] := 8)

	testSliceStr1 := []string{"a", "b", "c"}
	testSliceNotStr1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 8}

	retval, err = strproc.AnyCompare(testSliceStr1, testSliceNotStr1)
	fmt.Println("Return : ", retval)
	fmt.Println("Error : ", err)

	// Return :  false
	// Error :  Not Compare type, obj1.([a b c]) != obj2.([1 2 3 4 5 6 7 8 8])
}

func Example_strutils_IsValidEmail() {

	strvalidator := strutils.NewStringValidator()
	example_str := "a@golang.org"
	fmt.Printf("%v\n", strvalidator.IsValidEmail(example_str))
	// Output: true
}

func Example_strutils_IsValidDomain() {

	strvalidator := strutils.NewStringValidator()
	example_str := "golang.org"
	fmt.Printf("%v\n", strvalidator.IsValidDomain(example_str))
	// Output: true

}

func Example_strutils_IsValidURL() {

	strvalidator := strutils.NewStringValidator()
	example_str := "https://www.google.co.kr/url?sa=t&rct=j&q=&esrc=s&source=web"
	fmt.Printf("%v\n", strvalidator.IsValidURL(example_str))
	// Output: true

}

func Example_strutils_IsValidMACAddr() {

	strvalidator := strutils.NewStringValidator()
	example_str := "02-f3-71-eb-9e-4b"
	fmt.Printf("%v\n", strvalidator.IsValidMACAddr(example_str))
	// Output: true

}

func Example_strutils_IsValidIPAddr() {

	strvalidator := strutils.NewStringValidator()
	example_str := "2001:470:1f09:495::3:217.126.185.21"
	retval, err := strvalidator.IsValidIPAddr(example_str, strutils.IPv4MappedIPv6, strutils.IPv4)
	if err != nil {
		fmt.Println("Error : ", err)
	} else {
		fmt.Println(retval)
	}
	// Output: true
}

func Example_strutils_IsValidFilePath() {

	strvalidator := strutils.NewStringValidator()
	example_str := "a-1-s-d-v-we-wd_+qwd-qwd-qwd.txt"
	fmt.Printf("%v\n", strvalidator.IsValidFilePath(example_str))
	// Output: false
}

func Example_strutils_IsValidFilePathWithRelativePath() {

	strvalidator := strutils.NewStringValidator()
	example_str := "/asdasd/asdasdasd/qwdqwd_qwdqwd/12-12/a-1-e-r-t-_1_21234_d_1234_qwed_1423_.txt"
	fmt.Printf("%v\n", strvalidator.IsValidFilePathWithRelativePath(example_str))
	// Output: true
}

func Example_strutils_IsPureTextStrict() {

	strvalidator := strutils.NewStringValidator()
	example_str := `abcd/>qwdqwdoijhwer/>qwdojiqwdqwd</a>qwdoijqwdoiqjd`
	retval, err := strvalidator.IsPureTextStrict(example_str)
	if err != nil {
		fmt.Println("Error : ", err)
	} else {
		fmt.Println(retval)
	}
	// Output: Error :  Detect Tag (<[!|?]~>)
}

func Example_strutils_IsPureTextNormal() {

	strvalidator := strutils.NewStringValidator()
	example_str := `Foo<script type="text/javascript">alert(1337)</script>Bar`
	retval, err := strvalidator.IsPureTextNormal(example_str)

	if err != nil {
		fmt.Println("Error : ", err)
	} else {
		fmt.Println(retval)
	}
	// Output: Error :  Detect HTML Element
}

func Example_strutils_StripTags() {

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
<p>README.md haven’t contain all the examples. Please refer to the the XXXtest.go files.</p>
</body>
</html>
`
	retval, err := strproc.StripTags(example_str)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	fmt.Println(retval)

	// 	Output :Just! a String Processing Library for Go-lang
	//Just! a String Processing Library for Go-lang
	//Just a few methods for helping processing and validation the string
	//View on GitHub
	//Just! a String Processing Library for Go-lang
	//Just a few methods for helping processing the string
	//README.md haven’t contain all the examples. Please refer to the the XXXtest.go files.
}

func Example_strutils_ConvertToStr() {

	strproc := strutils.NewStringProc()
	example_val := uint64(1234567)
	retval, err := strproc.ConvertToStr(example_val)
	if err != nil {
		fmt.Println("Error : ", err)
	}
	fmt.Println(retval)

	// Output : "1234567"
}

func Example_strutils_ReverseStr() {

	dataset := []string{
		"0123456789",
		"가나다라마바사",
		"あいうえお",
		"天地玄黃宇宙洪荒",
	}

	strproc := strutils.NewStringProc()
	for _, v := range dataset {
		fmt.Println(strproc.ReverseStr(v))
	}

	// Output : 9876543210
	//사바마라다나가
	//おえういあ
	//荒洪宙宇黃玄地天
}

func Example_strutils_ReverseNormalStr() {

	dataset := []string{
		"0123456789",
		"abcdefg",
	}

	strproc := strutils.NewStringProc()
	for _, v := range dataset {
		fmt.Println(strproc.ReverseNormalStr(v))
	}

	// Output : 9876543210
	//gfedcba
}

func Example_strutils_ReverseReverseUnicode() {

	dataset := []string{
		"0123456789",
		"가나다라마바사",
		"あいうえお",
		"天地玄黃宇宙洪荒",
	}

	strproc := strutils.NewStringProc()
	for _, v := range dataset {
		fmt.Println(strproc.ReverseUnicode(v))
	}

	// Output : 9876543210
	//사바마라다나가
	//おえういあ
	//荒洪宙宇黃玄地天
}

func Example_strutils_FileMD5Hash() {
	strproc := strutils.NewStringProc()

	retval, err := strproc.FileMD5Hash("./LICENSE")
	if err != nil {
		fmt.Println("Error : %v", err)
	}

	fmt.Println(retval)

	// Output: f3f8954bac465686f0bfc2a757c5200b
}

func Example_strutils_MD5Hash() {

	dataset := []string{
		"0123456789",
		"abcdefg",
		"abcdefgqwdoisef;oijawe;fijq2039jdfs.dnc;oa283hr08uj3o;ijwaef;owhjefo;uhwefwef",
	}

	strproc := strutils.NewStringProc()

	//check : common
	for _, v := range dataset {
		retval, err := strproc.MD5Hash(v)
		if err != nil {
			fmt.Println("Error : %v", err)
		} else {
			fmt.Println(retval)
		}
	}

	// Output : 781e5e245d69b566979b86e28d23f2c7
	// 7ac66c0f148de9519b8bd264312c4d64
	// 15f764f21d09b11102eb015fc8824d00

}
