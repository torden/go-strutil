package strutils_test

import (
	"fmt"

	"github.com/torden/go-strutil"
)

func ExampleAddSlashes() {

	strproc := strutils.NewStringProc()
	example_str := `a\bcdefgz`
	fmt.Println(strproc.AddSlashes(example_str))
	// Output: a\\bcdefgz
}

func ExmapleStripSlashes() {

	strproc := strutils.NewStringProc()
	example_str := "a\\bcdefgz"
	fmt.Println(strproc.StripSlashes(example_str))
	// Output: abcdefgz
}

func ExampleNl2Br() {

	strproc := strutils.NewStringProc()
	example_str := "abc\ndefgh"
	fmt.Println(strproc.Nl2Br(example_str))
	// Output: abc<br />defgh
}

func ExampleWordWrapSimple() {

	strproc := strutils.NewStringProc()
	example_str := "The quick brown fox jumped over the lazy dog."
	fmt.Printf("%v\n", strproc.WordWrapSimple(example_str, 3, "*"))
	fmt.Printf("%v\n", strproc.WordWrapSimple(example_str, 8, "*"))
	// Output: The*quick*brown*fox*jumped*over*the*lazy*dog.
	// The quick*brown fox*jumped over*the lazy*dog.
}

func ExampleWordWrapAround() {

	strproc := strutils.NewStringProc()
	example_str := "The quick brown fox jumped over the lazy dog."
	fmt.Printf("%v\n", strproc.WordWrapAround(example_str, 3, "*"))
	fmt.Printf("%v\n", strproc.WordWrapAround(example_str, 8, "*"))
	// Output: The*quick*brown*fox*jumped*over*the*lazy*dog.
	// The quick*brown fox*jumped*over the*lazy*dog.
}

func ExampleNumberFmt() {

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

func ExamplePaddingBoth() {

	strproc := strutils.NewStringProc()
	example_str := "Life isn't always what one like."

	fmt.Printf("%v\n", strproc.PaddingBoth(example_str, "*", 38))
	fmt.Printf("%v\n", strproc.PaddingBoth(example_str, "*-=", 37))

	// Output: ***Life isn't always what one like.***
	// *-Life isn't always what one like.*-=
}

func ExamplePaddingLeft() {

	strproc := strutils.NewStringProc()
	example_str := "Life isn't always what one like."

	fmt.Printf("%v\n", strproc.PaddingLeft(example_str, "*", 38))
	fmt.Printf("%v\n", strproc.PaddingLeft(example_str, "*-=", 37))

	// Output: ******Life isn't always what one like.
	// *-=*-Life isn't always what one like.
}

func ExamplePaddingRight() {

	strproc := strutils.NewStringProc()
	example_str := "Life isn't always what one like."

	fmt.Printf("%v\n", strproc.PaddingRight(example_str, "*", 38))
	fmt.Printf("%v\n", strproc.PaddingRight(example_str, "*-=", 37))

	// Output: Life isn't always what one like.******
	// Life isn't always what one like.*-=*-
}

func ExampleLowerCaseFirstWords() {

	strproc := strutils.NewStringProc()
	example_str := "LIFE ISN'T ALWAYS WHAT ONE LIKE."
	fmt.Printf("%v\n", strproc.LowerCaseFirstWords(example_str))
	// Output: lIFE iSN'T aLWAYS wHAT oNE lIKE.
}

func ExampleUpperCaseFirstWords() {

	strproc := strutils.NewStringProc()
	example_str := "life isn't always what one like."
	fmt.Printf("%v\n", strproc.UpperCaseFirstWords(example_str))
	// Output: Life Isn't Always What One Like.
}

func ExampleSwapCaseFirstWords() {

	strproc := strutils.NewStringProc()
	example_str := "O SAY, CAN YOU SEE, BY THE DAWN’S EARLY LIGHT,"
	fmt.Printf("%v\n", strproc.UpperCaseFirstWords(example_str))
	// Output: O SAY, CAN YOU SEE, BY THE DAWN’S EARLY LIGHT,
}

func ExampleHumanByteSize() {

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

func ExampleHumanFileSize() {

	strproc := strutils.NewStringProc()
	retval, err := strproc.HumanFileSize("/etc/ssh/sshd_config", 2, strutils.CamelCaseLong)
	if err != nil {
		fmt.Println("Error : ", err)
	} else {
		fmt.Println(retval)
	}
	// Output: 2.44KiloByte
}

func ExampleAnyCompare() {

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

func ExampleIsValidEmail() {

	strvalidator := strutils.NewStringValidator()
	example_str := "a@golang.org"
	fmt.Printf("%v\n", strvalidator.IsValidEmail(example_str))
	// Output: true
}

func ExampleIsValidDomain() {

	strvalidator := strutils.NewStringValidator()
	example_str := "golang.org"
	fmt.Printf("%v\n", strvalidator.IsValidDomain(example_str))
	// Output: true

}

func ExampleIsValidURL() {

	strvalidator := strutils.NewStringValidator()
	example_str := "https://www.google.co.kr/url?sa=t&rct=j&q=&esrc=s&source=web"
	fmt.Printf("%v\n", strvalidator.IsValidURL(example_str))
	// Output: true

}

func ExampleIsValidMACAddr() {

	strvalidator := strutils.NewStringValidator()
	example_str := "02-f3-71-eb-9e-4b"
	fmt.Printf("%v\n", strvalidator.IsValidMACAddr(example_str))
	// Output: true

}

func ExampleIsValidIPAddr() {

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

func ExampleIsValidFilePath() {

	strvalidator := strutils.NewStringValidator()
	example_str := "a-1-s-d-v-we-wd_+qwd-qwd-qwd.txt"
	fmt.Printf("%v\n", strvalidator.IsValidFilePath(example_str))
	// Output: false
}

func ExampleIsValidFilePathWithRelativePath() {

	strvalidator := strutils.NewStringValidator()
	example_str := "/asdasd/asdasdasd/qwdqwd_qwdqwd/12-12/a-1-e-r-t-_1_21234_d_1234_qwed_1423_.txt"
	fmt.Printf("%v\n", strvalidator.IsValidFilePathWithRelativePath(example_str))
	// Output: true
}

func ExampleIsPureTextStrict() {

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

func ExampleIsPureTextNormal() {

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
