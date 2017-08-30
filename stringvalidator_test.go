package strutils_test

import (
	"testing"

	"github.com/torden/go-strutil"
)

func ipaddrTest(t *testing.T, cktype int, dataset map[string]bool, errfmt string) {

	//check : common
	for k, v := range dataset {
		retval, _ := strvalidator.IsValidIPAddr(k, cktype)

		assert.AssertEquals(t, v, retval, errfmt, retval)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func TestIPAddrFalse(t *testing.T) {

	t.Parallel()
	var err error

	//check : wrong IP Addr
	_, err = strvalidator.IsValidIPAddr("A.B.C.D", strutils.IPv4)
	assert.AssertNotNil(t, err, "Failured : Couldn't check the `wrong IP Addr`")

	//check : wrong options
	_, err = strvalidator.IsValidIPAddr("127.0.0.1", 7)
	assert.AssertNotNil(t, err, "Failured : Couldn't check the `wrong option`")

	//check : getIPType
	_, err = strvalidator.IsValidIPAddr(":F", strutils.IPv6CIDR, strutils.IPv6)
	assert.AssertNotNil(t, err, "Failured : Couldn't check the `wrong IP Addr`")

	//check : getIPType
	_, err = strvalidator.IsValidIPAddr("127127127127", strutils.IPv6CIDR, strutils.IPv6)
	assert.AssertNotNil(t, err, "Failured : Couldn't check the `wrong IP Addr`")

}

func TestIPAddr(t *testing.T) {

	//IPv4
	testIpv4Ipaddrs := map[string]bool{
		"192.168.1.1":     true,
		"127.0.0.1":       true,
		"10.10.90.12":     true,
		"8.8.8.8":         true,
		"4.4.4.4":         true,
		"912.456.123.123": false,
		"999.999.999.999": false,
		"192.192.19.999":  false,
	}

	//IPv4 with CIDR
	testIpv4WithCidrIpaddrs := map[string]bool{
		"192.168.1.1/24": true,
		"127.0.0.1/32":   true,
		"10.10.90.12/8":  true,
		"8.8.8.8/1":      true,
		"4.4.4.4/7":      true,
		"192.168.1.1/99": false,
		"127.0.0.1/9999": false,
		"10.10.90.12/33": false,
		"8.8.8.8/128":    false,
		"4.4.4.4/256":    false,
	}

	//IPv6
	testIpv6Ipaddrs := map[string]bool{
		"2607:f0d0:1002:51::4":                    true,
		"2607:f0d0:1002:0051:0000:0000:0000:0004": true,
		"ff05::1:3":                               true,
		"FE80:0000:0000:0000:0202:B3FF:FE1E:8329": true,
		"FE80::0202:B3FF:FE1E:8329":               true,
		"fe80::202:b3ff:fe1e:8329":                true,
		"fe80:0000:0000:0000:0202:b3ff:fe1e:8329": true,
		"2001:470:1f09:495::3":                    true,
		"2001:470:1f1d:275::1":                    true,
		"2600:9000:5304:200::1":                   true,
		"2600:9000:5306:d500::1":                  true,
		"2600:9000:5301:b600::1":                  true,
		"2600:9000:5303:900::1":                   true,
		"127:12:12:12:12:12:!2:!2":                false,
		"127.0.0.1":                               false,
		"234:23:23:23:23:23:23":                   false,
	}

	//IPv6 with CIDR
	testIpv6WithCidrIpaddrs := map[string]bool{
		"2000::/5":      true,
		"2000::/15":     true,
		"2001:db8::/33": true,
		"2001:db8::/48": true,
		"fc00::/7":      true,
	}

	//IPv4-Mapped Embedded IPv6 Address
	testIpv4MappedIpv6Ipaddrs := map[string]bool{
		"2001:470:1f09:495::3:217.126.185.215":         true,
		"2001:470:1f1d:275::1:213.0.69.132":            true,
		"2600:9000:5304:200::1:205.251.196.2":          true,
		"2600:9000:5306:d500::1:205.251.198.213":       true,
		"2600:9000:5301:b600::1:205.251.193.182":       true,
		"2600:9000:5303:900::1:205.251.195.9":          true,
		"0:0:0:0:0:FFFF:222.1.41.90":                   true,
		"::FFFF:222.1.41.90":                           true,
		"0000:0000:0000:0000:0000:FFFF:12.155.166.101": true,
		"12.155.166.101":                               false,
		"12.12/12":                                     false,
	}

	//IPv4-Mapped Embedded IPv6 Address with CIDR
	testIpv4MappedIpv6IpaddrsCIDR := map[string]bool{
		"2600:9000:5301:b600::1:205.251.193.182/32": true,
		"2600:9000:5303:900::1:205.251.195.9/32":    true,
		"0:0:0:0:0:FFFF:222.1.41.90/32":             true,
	}

	//check : common
	ipaddrTest(t, strutils.IPv4, testIpv4Ipaddrs, "invalid (%s) IPv4 address")
	ipaddrTest(t, strutils.IPv4CIDR, testIpv4WithCidrIpaddrs, "invalid (%s) IPv4 with CIDR address")

	ipaddrTest(t, strutils.IPv6, testIpv6Ipaddrs, "invalid (%s) IPv6 address")
	ipaddrTest(t, strutils.IPv6CIDR, testIpv6WithCidrIpaddrs, "invalid (%s) IPv6 with CIDR address")
	ipaddrTest(t, strutils.IPv4MappedIPv6, testIpv4MappedIpv6Ipaddrs, "invalid (%s) IPv4-Mapped Embedded IPv6 address")

	ipaddrTest(t, strutils.IPv4MappedIPv6CIDR, testIpv4MappedIpv6IpaddrsCIDR, "invalid (%s) IPv4-Mapped Embedded IPv6 address With CIDR")
}

func TestMacAddr(t *testing.T) {

	macaddrList := map[string]bool{
		"02:f3:71:eb:9e:4b": true,
		"02-f3-71-eb-9e-4b": true,
		"02f3.71eb.9e4b":    true,
		"87:78:6e:3e:90:40": true,
		"87-78-6e-3e-90-40": true,
		"8778.6e3e.9040":    true,
		"e7:28:b9:57:ab:36": true,
		"e7-28-b9-57-ab-36": true,
		"e728.b957.ab36":    true,
		"eb:f8:2b:d7:e9:62": true,
		"eb-f8-2b-d7-e9-62": true,
		"ebf8.2bd7.e962":    true,
	}

	t.Parallel()

	//check : common
	for k, v := range macaddrList {
		retval := strvalidator.IsValidMACAddr(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}

	//check : return FALSE
	retval := strvalidator.IsValidMACAddr("127.0.0.1")
	assert.AssertFalse(t, retval, "Failured : Couldn't check the `return false`")
}

func TestDomain(t *testing.T) {

	testDomains := map[string]bool{
		"대한민국.xn-korea.co.kr":           true,
		"google.com":                    true,
		"masełkowski.pl":                true,
		"maselkowski.pl":                true,
		"m.maselkowski.pl":              true,
		"www.masełkowski.pl.com":        true,
		"xn--masekowski-d0b.pl":         true,
		"中国互联网络信息中心.中国":                 true,
		"masełkowski.pl.":               false,
		"中国互联网络信息中心.xn--masekowski-d0b": false,
		"a.jp":                     true,
		"a.co":                     true,
		"a.co.jp":                  true,
		"a.co.or":                  true,
		"a.or.kr":                  true,
		"qwd-qwdqwd.com":           true,
		"qwd-qwdqwd.co_m":          false,
		"qwd-qwdqwd.c":             false,
		"qwd-qwdqwd.-12":           false,
		"qwd-qwdqwd.1212":          false,
		"qwd-qwdqwd.org":           true,
		"qwd-qwdqwd.ac.kr":         true,
		"qwd-qwdqwd.gov":           true,
		"chicken.beer":             true,
		"aa.xyz":                   true,
		"google.asn.au":            true,
		"google.com.au":            true,
		"google.net.au":            true,
		"google.priv.at":           true,
		"google.ac.at":             true,
		"google.gv.at":             true,
		"google.avocat.fr":         true,
		"google.geek.nz":           true,
		"google.gen.nz":            true,
		"google.kiwi.nz":           true,
		"google.org.il":            true,
		"google.net.il":            true,
		"www.google.edu.au":        true,
		"www.google.gov.au":        true,
		"www.google.csiro.au":      true,
		"www.google.act.au":        true,
		"www.google.avocat.fr":     true,
		"www.google.aeroport.fr":   true,
		"www.google.co.nz":         true,
		"www.google.geek.nz":       true,
		"www.google.gen.nz":        true,
		"www.google.kiwi.nz":       true,
		"www.google.parliament.nz": true,
		"www.google.muni.il":       true,
		"www.google.idf.il":        true,
	}

	t.Parallel()

	//check : common
	for k, v := range testDomains {
		retval := strvalidator.IsValidDomain(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func TestURL(t *testing.T) {

	testUrls := map[string]bool{
		"https://www.google.co.kr/url?sa=t&rct=j&q=&esrc=s&source=web":                                      true,
		"http://stackoverflow.com/questions/27812164/can-i-import-3rd-party-package-into-golang-playground": true,
		"https://tour.golang.org/welcome/4":                                                                 true,
		"https://revel.github.io/":                                                                          true,
		"https://github.com/revel/revel/commit/bd1d083ee4345e919b3bca1e4c42ca682525e395#diff-972a2b2141d27e9d7a8a4149a7e28eef": true,
		"https://github.com/ndevilla/iniparser/pull/82#issuecomment-261817064":                                                 true,
		"http://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&wd=golang":                                            true,
		"http://www.baidu.com/link?url=DrWkM_beo2M5kB5sLYnItKSQ0Ib3oDhKcPprdtLzAWNfFt_VN5oyD3KwnAKT6Xsk":                       true,
		"http://www.yahoo.com/search?p=대한민국":                                                                                   true,
		"http://www.baidu.com/s?ie=utf-8&f=8&rsv_bp=0&rsv_idx=1&tn=baidu&wd=关于百度&rsv_sug3=2":                                   true,
		"https://search.yahoo.co.jp/search?p=日本&aq=-1&oq=&ei=UTF-8&fr=top_ga1_sa&x=wrt":                                        true,
		"http://ftp.yz.yamagata-u.ac.jp/pub/linux/centos/7/isos/x86_64/CentOS-7-x86_64-DVD-1611.iso":                           true,
	}

	t.Parallel()

	//check : common
	for k, v := range testUrls {
		retval := strvalidator.IsValidURL(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nExpected: %v\nActual: %v", retval, v)
	}
}

func TestPureTextNormal(t *testing.T) {

	testTxts := map[string]bool{
		`<script ?>qwdpijqwd</script>qd08j123lneqw\t\nqwedojiqwd\rqwdoihjqwd1d[08jaedl;jkqwd\r\nqdolijqdwqwd`:       false,
		`a\r\nb<script ?>qwdpijqwd</script>qd08j123lneqw\t\nqwedojiqwd\rqwdoihjqwd1d[08jaedl;jkqwd\r\nqdolijqdwqwd`: false,
		`Foo<script type="text/javascript">alert(1337)</script>Bar`:                                                 false,
		`Foo<12>Bar`:              true,
		`Foo<>Bar`:                true,
		`Foo</br>Bar`:             false,
		`Foo <!-- Bar --> Baz`:    false,
		`I <3 Ponies!`:            true,
		`I &#32; like Golang\t\n`: true,
		`I &amp; like Golang\t\n`: true,
		`<?xml version="1.0" encoding="UTF-8" ?> <!DOCTYPE log4j:configuration SYSTEM "log4j.dtd"> <log4j:configuration debug="true" xmlns:log4j='http://jakarta.apache.org/log4j/'> <appender name="console" class="org.apache.log4j.ConsoleAppender"> <layout class="org.apache.log4j.PatternLayout"> <param name="ConversionPattern" value="%d{yyyy-MM-dd HH:mm:ss} %-5p %c{1}:%L - %m%n" /> </layout> </appender> <root> <level value="DEBUG" /> <appender-ref ref="console" /> </root> </log4j:configuration>`: false,
		`I like Golang\r\n`:       true,
		`I like Golang\r\na`:      true,
		"I &#32; like Golang\t\n": true,
		"I &amp; like Golang\t\n": true,
		`ハイレゾ対応ウォークマン®、ヘッドホン、スピーカー「Winter Gift Collection ～Presented by JUJU～」をソニーストアにて販売開始`:                                                                      true,
		`VAIOパーソナルコンピューター type T TZシリーズ 無償点検・修理のお知らせとお詫び（2009年10月15日更新）`:                                                                                          true,
		`把百度设为主页关于百度About  Baidu百度推广`:                                                                                                                             true,
		`%E6%8A%8A%E7%99%BE%E5%BA%A6%E8%AE%BE%E4%B8%BA%E4%B8%BB%E9%A1%B5%E5%85%B3%E4%BA%8E%E7%99%BE%E5%BA%A6About++Baidu%E7%99%BE%E5%BA%A6%E6%8E%A8%E5%B9%BF`:     true,
		`%E6%8A%8A%E7%99%BE%E5%BA%A6%E8%AE%BE%E4%B8%BA%E4%B8%BB%E9%A1%B5%E5%85%B3%E4%BA%8E%E7%99%BE%E5%BA%A6About%20%20Baidu%E7%99%BE%E5%BA%A6%E6%8E%A8%E5%B9%BF`: true,
		`abcd/>qwdqwdoijhwer/>qwdojiqwdqwd</>qwdoijqwdoiqjd`:                                                                                                      true,
		`abcd/>qwdqwdoijhwer/>qwdojiqwdqwd</a>qwdoijqwdoiqjd`:                                                                                                     false,
		"\tq\tq\t\nq": false,
		"": false,
		"aaaaaaqwdqwdqwwdqwdqw	qwdqwdqwqdw": false,
		"AbcEd-=qwdoijqdwoijaaaaaaqwdqwdqwwdqwdqw	qwdqwdqwqdw": false,
	}

	t.Parallel()

	//check : common
	for k, v := range testTxts {
		retval, _ := strvalidator.IsPureTextNormal(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nTest Value : %#v\nExpected: %v\nActual: %v", k, retval, v)
	}
}

func TestPureTextStrict(t *testing.T) {

	testTxts := map[string]bool{
		`<script ?>qwdpijqwd</script>qd08j123lneqw\t\nqwedojiqwd\rqwdoihjqwd1d[08jaedl;jkqwd\r\nqdolijqdwqwd`:       false,
		`a\r\nb<script ?>qwdpijqwd</script>qd08j123lneqw\t\nqwedojiqwd\rqwdoihjqwd1d[08jaedl;jkqwd\r\nqdolijqdwqwd`: false,
		`Foo<script type="text/javascript">alert(1337)</script>Bar`:                                                 false,
		`Foo<12>Bar`:              true,
		`Foo<>Bar`:                true,
		`Foo</br>Bar`:             false,
		`Foo <!-- Bar --> Baz`:    false,
		`I <3 Ponies!`:            true,
		`I &#32; like Golang\t\n`: true,
		`I &amp; like Golang\t\n`: false,
		`<?xml version="1.0" encoding="UTF-8" ?> <!DOCTYPE log4j:configuration SYSTEM "log4j.dtd"> <log4j:configuration debug="true" xmlns:log4j='http://jakarta.apache.org/log4j/'> <ender name="console" class="org.apache.log4j.ConsoleAppender"> <layout class="org.apache.log4j.PatternLayout"> <param name="ConversionPattern" value="%d{yyyy-MM-dd HH:mm:ss} %-5p 1}:%L - %m%n" /> </layout> </appender> <root> <level value="DEBUG" /> <appender-ref ref="console" /> </root> </log4j:configuration>`: false,
		`I like Golang\r\n`:       true,
		`I like Golang\r\na`:      true,
		"I &#32; like Golang\t\n": true,
		"I &amp; like Golang\t\n": false,
		`ハイレゾ対応ウォークマン®、ヘッドホン、スピーカー「Winter Gift Collection ～Presented by JUJU～」をソニーストアにて販売開始`:                                                                      true,
		`VAIOパーソナルコンピューター type T TZシリーズ 無償点検・修理のお知らせとお詫び（2009年10月15日更新）`:                                                                                          true,
		`把百度设为主页关于百度About  Baidu百度推广`:                                                                                                                             true,
		`%E6%8A%8A%E7%99%BE%E5%BA%A6%E8%AE%BE%E4%B8%BA%E4%B8%BB%E9%A1%B5%E5%85%B3%E4%BA%8E%E7%99%BE%E5%BA%A6About++Baidu%E7%99%BE%E5%BA%A6%E6%8E%A8%E5%B9%BF`:     true,
		`%E6%8A%8A%E7%99%BE%E5%BA%A6%E8%AE%BE%E4%B8%BA%E4%B8%BB%E9%A1%B5%E5%85%B3%E4%BA%8E%E7%99%BE%E5%BA%A6About%20%20Baidu%E7%99%BE%E5%BA%A6%E6%8E%A8%E5%B9%BF`: true,
		`abcd/>qwdqwdoijhwer/>qwdojiqwdqwd</>qwdoijqwdoiqjd`:                                                                                                      true,
		`abcd/>qwdqwdoijhwer/>qwdojiqwdqwd</a>qwdoijqwdoiqjd`:                                                                                                     false,
		"\tq\tq\t\nq": false,
		"": false,
		"aaaaaaqwdqwdqwwdqwdqw	qwdqwdqwqdw": false,
		"AbcEd-=qwdoijqdwoijaaaaaaqwdqwdqwwdqwdqw	qwdqwdqwqdw": false,
	}

	t.Parallel()

	//check : common
	for k, v := range testTxts {
		retval, _ := strvalidator.IsPureTextStrict(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nTest Value : %#v\nExpected: %v\nActual: %v", k, retval, v)
	}

}

func TestFilePathOnlyFilePath(t *testing.T) {

	testFilepaths := map[string]bool{
		"../../qwdqwdqwd/../qwdqwdqwd.txt": false,
		`../../qwdqwdqwd/..
				        /qwdqwdqwd.txt`: false,
		"\t../../qwdqwdqwd/../qwdqwdqwd.txt": false,
		`../../qwdqwdqwd/../qwdqwdqwd.txt`: false,
		`../../qwdqwdqwd/../qwdqwdqwd.txt`: false,
		"../../etc/passwd":                 false,
		"a.txt;rm -rf /":                   false,
		"sudo rm -rf ../":                  false,
		"a-1-s-d-v-we-wd_+qwd-qwd-qwd.txt": false,
		"a-qwdqwd_qwdqwdqwd-123.txt":       true,
		"a.txt": true,
		"a-1-e-r-t-_1_21234_d_1234_qwed_1423_.txt": true,
	}

	t.Parallel()

	//check : common
	for k, v := range testFilepaths {
		retval := strvalidator.IsValidFilePath(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nTest Value : %#v\nExpected: %v\nActual: %v", k, retval, v)
	}
}

func TestFilePathAllowRelativePath(t *testing.T) {

	testFilepaths := map[string]bool{
		"../../qwdqwdqwd/../qwdqwdqwd.txt": true,
		`../../qwdqwdqwd/..
				        /qwdqwdqwd.txt`: false,
		"\t../../qwdqwdqwd/../qwdqwdqwd.txt": false,
		`../../qwdqwdqwd/../qwdqwdqwd.txt`: false,
		`../../qwdqwdqwd/../qwdqwdqwd.txt`: false,
		"../../etc/passwd":                 true,
		"a.txt;rm -rf /":                   false,
		"sudo rm -rf ../":                  true,
		"a-1-s-d-v-we-wd_+qwd-qwd-qwd.txt": false,
		"a-qwdqwd_qwdqwdqwd-123.txt":       true,
		"a.txt": true,
		"a-1-e-r-t-_1_21234_d_1234_qwed_1423_.txt":                                       true,
		"/asdasd/asdasdasd/qwdqwd_qwdqwd/12-12/a-1-e-r-t-_1_21234_d_1234_qwed_1423_.txt": true,
	}

	t.Parallel()

	//check : common
	for k, v := range testFilepaths {
		retval := strvalidator.IsValidFilePathWithRelativePath(k)
		assert.AssertEquals(t, v, retval, "Return Value mismatch.\nTest Value : %#v\nExpected: %v\nActual: %v", k, retval, v)
	}
}
