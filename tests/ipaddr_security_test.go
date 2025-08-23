package strutils_test

import (
	"testing"

	strutils "github.com/torden/go-strutil"
)

func Test_IPAddr_SecurityCases(t *testing.T) {
	t.Parallel()

	// IPv4 보안 관련 엣지 케이스
	ipv4SecurityCases := []struct {
		name     string
		ip       string
		ipType   int
		expected bool
	}{
		// 네트워크 주소와 브로드캐스트
		{"Network address", "0.0.0.0", strutils.IPv4, true},
		{"Broadcast address", "255.255.255.255", strutils.IPv4, true},
		{"Loopback", "127.0.0.1", strutils.IPv4, true},
		{"Private range A", "10.0.0.1", strutils.IPv4, true},
		{"Private range B", "172.16.0.1", strutils.IPv4, true},
		{"Private range C", "192.168.1.1", strutils.IPv4, true},
		
		// 범위 초과 (공격에 자주 사용)
		{"Octet overflow 1", "256.1.1.1", strutils.IPv4, false},
		{"Octet overflow 2", "1.256.1.1", strutils.IPv4, false},
		{"Octet overflow 3", "1.1.256.1", strutils.IPv4, false},
		{"Octet overflow 4", "1.1.1.256", strutils.IPv4, false},
		{"Max overflow", "999.999.999.999", strutils.IPv4, false},
		
		// 선행 0 (일부 시스템에서 8진수로 해석)
		{"Leading zero 1", "01.1.1.1", strutils.IPv4, false},
		{"Leading zero 2", "1.01.1.1", strutils.IPv4, false},
		{"Leading zero 3", "1.1.01.1", strutils.IPv4, false},
		{"Leading zero 4", "1.1.1.01", strutils.IPv4, false},
		{"Octal format", "010.010.010.010", strutils.IPv4, false},
		
		// 구조적 문제
		{"Missing octet", "1.1.1", strutils.IPv4, false},
		{"Extra octet", "1.1.1.1.1", strutils.IPv4, false},
		{"Double dot", "1..1.1", strutils.IPv4, false},
		{"Trailing dot", "1.1.1.1.", strutils.IPv4, false},
		{"Leading dot", ".1.1.1.1", strutils.IPv4, false},
		
		// 빈 옥텟
		{"Empty octet 1", ".1.1.1", strutils.IPv4, false},
		{"Empty octet 2", "1..1.1", strutils.IPv4, false},
		{"Empty octet 3", "1.1..1", strutils.IPv4, false},
		{"Empty octet 4", "1.1.1.", strutils.IPv4, false},
		
		// 문자 포함 (인젝션 시도)
		{"Alpha character", "1.1.1.a", strutils.IPv4, false},
		{"Hex notation", "0x1.0x1.0x1.0x1", strutils.IPv4, false},
		{"Mixed format", "192.168.0x1.1", strutils.IPv4, false},
	}

	for _, testCase := range ipv4SecurityCases {
		t.Run("IPv4_"+testCase.name, func(t *testing.T) {
			result, _ := strvalidator.IsValidIPAddr(testCase.ip, testCase.ipType)
			assert.AssertEquals(t, testCase.expected, result,
				"IPv4 security test failed for %s (%s).\nExpected: %v\nActual: %v",
				testCase.name, testCase.ip, testCase.expected, result)
		})
	}
}

func Test_IPAddr_IPv6SecurityCases(t *testing.T) {
	t.Parallel()

	// IPv6 보안 관련 엣지 케이스
	ipv6SecurityCases := []struct {
		name     string
		ip       string
		ipType   int
		expected bool
	}{
		// 기본 유효 케이스
		{"All zeros short", "::", strutils.IPv6, true},
		{"All zeros full", "0000:0000:0000:0000:0000:0000:0000:0000", strutils.IPv6, true},
		{"Loopback", "::1", strutils.IPv6, true},
		{"Documentation", "2001:db8::1", strutils.IPv6, true},
		{"Link-local", "fe80::1", strutils.IPv6, true},
		
		// 압축 규칙 위반
		{"Double compression", "2001:db8::1::2", strutils.IPv6, false},
		{"Invalid compression", "2001:db8:::1", strutils.IPv6, false},
		{"Triple colon", "2001:::db8::1", strutils.IPv6, false},
		
		// 잘못된 16진수 문자
		{"Invalid hex g", "2001:db8:85a3::8a2e:370g:7334", strutils.IPv6, false},
		{"Invalid hex z", "2001:db8:85a3::8a2e:370z:7334", strutils.IPv6, false},
		{"Non-hex character", "2001:db8:85a3::8a2e:370!:7334", strutils.IPv6, false},
		
		// 그룹 길이 초과
		{"Group too long 1", "12345:db8::1", strutils.IPv6, false},
		{"Group too long 2", "2001:123456::1", strutils.IPv6, false},
		{"Group too long 3", "2001:db8::12345", strutils.IPv6, false},
		
		// 그룹 수 초과
		{"Too many groups", "2001:db8:85a3:8d3:1319:8a2e:370:7344:extra", strutils.IPv6, false},
		{"Nine groups", "1:2:3:4:5:6:7:8:9", strutils.IPv6, false},
		
		// 잘못된 구조
		{"Leading colon single", ":2001:db8::1", strutils.IPv6, false},
		{"Trailing colon single", "2001:db8::1:", strutils.IPv6, false},
		{"Empty group middle", "2001::db8", strutils.IPv6, true}, // 이것은 유효
		{"Four colons", "2001::::db8", strutils.IPv6, false},
		
		// Zone ID 포함 (일부 구현에서 지원하지 않음)
		{"Zone ID percent", "fe80::1%eth0", strutils.IPv6, false},
		{"Zone ID number", "fe80::1%1", strutils.IPv6, false},
		
		// 임베디드 IPv4 관련
		{"Invalid embedded IPv4", "::ffff:256.1.1.1", strutils.IPv6, false},
		{"Malformed embedded", "::ffff:192.168.1", strutils.IPv6, false},
		{"Extra embedded", "::ffff:192.168.1.1.1", strutils.IPv6, false},
	}

	for _, testCase := range ipv6SecurityCases {
		t.Run("IPv6_"+testCase.name, func(t *testing.T) {
			result, _ := strvalidator.IsValidIPAddr(testCase.ip, testCase.ipType)
			assert.AssertEquals(t, testCase.expected, result,
				"IPv6 security test failed for %s (%s).\nExpected: %v\nActual: %v",
				testCase.name, testCase.ip, testCase.expected, result)
		})
	}
}

func Test_IPAddr_CIDRSecurityCases(t *testing.T) {
	t.Parallel()

	// CIDR 표기법 보안 테스트
	cidrSecurityCases := []struct {
		name     string
		ip       string
		ipType   int
		expected bool
	}{
		// IPv4 CIDR 엣지 케이스
		{"Valid IPv4 CIDR min", "192.168.1.0/0", strutils.IPv4CIDR, true},
		{"Valid IPv4 CIDR max", "192.168.1.0/32", strutils.IPv4CIDR, true},
		{"Invalid IPv4 CIDR over", "192.168.1.0/33", strutils.IPv4CIDR, false},
		{"Invalid IPv4 CIDR negative", "192.168.1.0/-1", strutils.IPv4CIDR, false},
		{"Invalid IPv4 CIDR string", "192.168.1.0/abc", strutils.IPv4CIDR, false},
		{"Missing CIDR value", "192.168.1.0/", strutils.IPv4CIDR, false},
		{"Double slash", "192.168.1.0//24", strutils.IPv4CIDR, false},
		{"Leading zero CIDR", "192.168.1.0/08", strutils.IPv4CIDR, false},
		
		// IPv6 CIDR 엣지 케이스
		{"Valid IPv6 CIDR min", "2001:db8::/0", strutils.IPv6CIDR, true},
		{"Valid IPv6 CIDR max", "2001:db8::/128", strutils.IPv6CIDR, true},
		{"Invalid IPv6 CIDR over", "2001:db8::/129", strutils.IPv6CIDR, false},
		{"Invalid IPv6 CIDR negative", "2001:db8::/-1", strutils.IPv6CIDR, false},
		{"IPv6 missing CIDR", "2001:db8::/", strutils.IPv6CIDR, false},
		{"IPv6 string CIDR", "2001:db8::/xyz", strutils.IPv6CIDR, false},
		
		// 혼합 공격 케이스
		{"IPv4 in IPv6 CIDR test", "192.168.1.1/24", strutils.IPv6CIDR, false},
		{"IPv6 in IPv4 CIDR test", "2001:db8::/64", strutils.IPv4CIDR, false},
	}

	for _, testCase := range cidrSecurityCases {
		t.Run("CIDR_"+testCase.name, func(t *testing.T) {
			result, _ := strvalidator.IsValidIPAddr(testCase.ip, testCase.ipType)
			assert.AssertEquals(t, testCase.expected, result,
				"CIDR security test failed for %s (%s).\nExpected: %v\nActual: %v",
				testCase.name, testCase.ip, testCase.expected, result)
		})
	}
}

func Test_IPAddr_MixedTypeSecurityCases(t *testing.T) {
	t.Parallel()

	// 여러 IP 유형을 동시에 검증하는 보안 테스트
	mixedSecurityCases := []struct {
		name     string
		ip       string
		ipTypes  []int
		expected []bool
	}{
		{
			"IPv4-mapped IPv6",
			"2001:470:1f09:495::3:217.126.185.21",
			[]int{strutils.IPv4, strutils.IPv6, strutils.IPv4MappedIPv6},
			[]bool{false, false, true},
		},
		{
			"Invalid IPv4-mapped",
			"2001:470:1f09:495::3:256.126.185.21",
			[]int{strutils.IPv4, strutils.IPv6, strutils.IPv4MappedIPv6},
			[]bool{false, false, false},
		},
		{
			"Pure IPv4",
			"192.168.1.1",
			[]int{strutils.IPv4, strutils.IPv6, strutils.IPv4MappedIPv6},
			[]bool{true, false, false},
		},
		{
			"Pure IPv6",
			"2001:db8::1",
			[]int{strutils.IPv4, strutils.IPv6, strutils.IPv4MappedIPv6},
			[]bool{false, true, false},
		},
		{
			"Malformed mixed",
			"not.an.ip.address",
			[]int{strutils.IPv4, strutils.IPv6, strutils.IPv4MappedIPv6},
			[]bool{false, false, false},
		},
	}

	for _, testCase := range mixedSecurityCases {
		t.Run("Mixed_"+testCase.name, func(t *testing.T) {
			for i, ipType := range testCase.ipTypes {
				result, _ := strvalidator.IsValidIPAddr(testCase.ip, ipType)
				assert.AssertEquals(t, testCase.expected[i], result,
					"Mixed type security test failed for %s (%s) with type %d.\nExpected: %v\nActual: %v",
					testCase.name, testCase.ip, ipType, testCase.expected[i], result)
			}
		})
	}
}

func Test_IPAddr_InjectionAttempts(t *testing.T) {
	t.Parallel()

	// SQL 인젝션, 명령어 인젝션 등을 시도하는 케이스
	injectionCases := []string{
		"192.168.1.1'; DROP TABLE users; --",
		"192.168.1.1 OR 1=1",
		"192.168.1.1; cat /etc/passwd",
		"192.168.1.1 && rm -rf /",
		"192.168.1.1 | nc -l 1234",
		"192.168.1.1`whoami`",
		"192.168.1.1$(id)",
		"192.168.1.1{cat,/etc/passwd}",
		"../../../etc/passwd",
		"\\x41\\x42\\x43\\x44",
		"<script>alert('xss')</script>",
		"javascript:alert(1)",
		"data:text/html,<h1>test</h1>",
		"\x00\x01\x02\x03",
		"192.168.1.1\r\nSet-Cookie: admin=1",
		"192.168.1.1\n\nHTTP/1.1 200 OK",
	}

	for i, maliciousInput := range injectionCases {
		t.Run("Injection_"+string(rune('A'+i%26)), func(t *testing.T) {
			// 모든 IP 타입에 대해 false여야 함
			ipTypes := []int{strutils.IPv4, strutils.IPv6, strutils.IPv4CIDR, strutils.IPv6CIDR, strutils.IPv4MappedIPv6}
			
			for _, ipType := range ipTypes {
				result, _ := strvalidator.IsValidIPAddr(maliciousInput, ipType)
				assert.AssertFalse(t, result,
					"Injection attempt should be rejected: %q with type %d", maliciousInput, ipType)
			}
		})
	}
}

func Test_IPAddr_BufferOverflowAttempts(t *testing.T) {
	t.Parallel()

	// 버퍼 오버플로우 시도
	longStrings := []struct {
		name   string
		input  string
		length int
	}{
		{"Long IPv4", "192.168.1." + string(make([]byte, 1000)), 1011},
		{"Long IPv6", "2001:db8::" + string(make([]byte, 1000)), 1012},
		{"Very long input", string(make([]byte, 10000)), 10000},
		{"Long CIDR", "192.168.1.1/" + string(make([]byte, 100)), 112},
	}

	for _, testCase := range longStrings {
		t.Run("BufferOverflow_"+testCase.name, func(t *testing.T) {
			ipTypes := []int{strutils.IPv4, strutils.IPv6, strutils.IPv4CIDR, strutils.IPv6CIDR}
			
			for _, ipType := range ipTypes {
				result, _ := strvalidator.IsValidIPAddr(testCase.input, ipType)
				assert.AssertFalse(t, result,
					"Buffer overflow attempt should be rejected: %s (length: %d) with type %d",
					testCase.name, testCase.length, ipType)
			}
		})
	}
}