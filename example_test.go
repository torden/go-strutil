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

func Exmaple_strutils_StripSlashes() {

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
