package strutils_test

import (
	"strings"
	"testing"
)

func Test_Email_ComprehensiveCases(t *testing.T) {
	t.Parallel()

	// 복합 유효 이메일 케이스
	validComplexCases := []string{
		// 기본 유효 케이스
		"user@domain.com",
		"test@example.org",
		"admin@site.net",
		
		// 복합 로컬 파트
		"user.name@domain.com",
		"user_name@domain.com", 
		"user-name@domain.com",
		"user+tag@domain.com",
		"first.last@domain.com",
		"user123@domain.com",
		"123user@domain.com",
		"user.123@domain.com",
		
		// 복합 도메인
		"user@sub.domain.com",
		"user@sub-domain.com",
		"user@domain-name.co.kr",
		"user@my-site.example.org",
		"user@test.museum",
		"user@site.travel",
		
		// 긴 도메인명
		"user@very-long-domain-name-that-is-still-valid.com",
		"test@sub1.sub2.sub3.domain.com",
		
		// 다양한 TLD
		"user@domain.info",
		"user@domain.biz",  
		"user@domain.name",
		"user@domain.pro",
		"user@domain.aero",
		"user@domain.coop",
		
		// 국제 도메인 (ASCII로 변환된 형태)
		"user@xn--nxasmq6b.xn--j6w193g",  // 한국어 도메인의 punycode
		"user@xn--fsq.xn--0zwm56d",       // 중국어 도메인의 punycode
	}

	for _, email := range validComplexCases {
		t.Run("Valid_"+strings.ReplaceAll(email, "@", "_at_"), func(t *testing.T) {
			result := strvalidator.IsValidEmail(email)
			assert.AssertTrue(t, result, "Valid email should pass: %s", email)
		})
	}
}

func Test_Email_SecurityAndAttackCases(t *testing.T) {
	t.Parallel()

	// 보안 공격 및 예외 케이스
	invalidSecurityCases := []string{
		// 기본 형식 오류  
		"@domain.com",                    // 로컬 파트 없음
		"user@",                          // 도메인 없음
		"user@@domain.com",               // 이중 @
		"user@domain",                    // TLD 없음
		"user.domain.com",                // @ 없음
		"user@domain.com@extra.com",      // 이중 @
		
		// 로컬 파트 문제
		".user@domain.com",               // 선행 점
		"user.@domain.com",               // 후행 점  
		"us..er@domain.com",              // 연속 점
		"user with spaces@domain.com",    // 공백
		"user@domain.com.",               // 도메인 후행 점
		
		// 도메인 문제
		"user@.domain.com",               // 선행 점
		"user@domain..com",               // 연속 점
		"user@domain.com.",               // 후행 점
		"user@-domain.com",               // 선행 하이픈
		"user@domain-.com",               // 후행 하이픈
		"user@domain.com-",               // 도메인 끝 하이픈
		
		// 길이 제한 초과
		strings.Repeat("a", 65) + "@domain.com",                    // 로컬 파트 너무 김 (64자 초과)
		"user@" + strings.Repeat("a", 250) + ".com",               // 도메인 너무 김
		"user@" + strings.Repeat("subdomain.", 30) + "com",        // 도메인 라벨 수 초과
		
		// 특수 문자 공격
		"user@domain.com!",                // 도메인에 특수문자
		"user@domain.com#",                // 해시
		"user@domain.com$",                // 달러
		"user@domain.com%",                // 퍼센트
		"user@domain.com^",                // 캐럿
		"user@domain.com&",                // 앰퍼샌드
		"user@domain.com*",                // 별표
		"user@domain.com(",                // 괄호
		"user@domain.com)",                // 괄호
		
		// SQL 인젝션 시도
		"user'; DROP TABLE users; --@domain.com",
		"user' OR 1=1 --@domain.com", 
		"user@domain.com'; DROP TABLE users; --",
		"user@domain.com' OR 1=1 --",
		
		// XSS 시도
		"<script>alert('xss')</script>@domain.com",
		"user@<script>alert('xss')</script>.com",
		"user@domain.com<script>alert(1)</script>",
		"javascript:alert(1)@domain.com",
		"user@javascript:alert(1).com",
		
		// 명령어 인젝션 시도
		"user; cat /etc/passwd@domain.com",
		"user@domain.com; rm -rf /",
		"user`whoami`@domain.com",
		"user@domain.com`whoami`",
		"user$(id)@domain.com",
		"user@domain.com$(id)",
		
		// 제어 문자 및 바이너리
		"user\x00@domain.com",             // 널 바이트
		"user\n@domain.com",               // 개행
		"user\r@domain.com",               // 캐리지 리턴
		"user\t@domain.com",               // 탭
		"user@domain.com\x00",             // 도메인에 널 바이트
		"user@domain.com\n",               // 도메인에 개행
		
		// HTTP 헤더 인젝션
		"user@domain.com\r\nSet-Cookie: admin=1",
		"user@domain.com\n\nHTTP/1.1 200 OK",
		"user\r\nBcc: attacker@evil.com\r\n@domain.com",
		
		// 경로 조작
		"user@../../../etc/passwd",
		"user@..\\..\\..\\windows\\system32",
		"user@domain.com/../admin",
		
		// 유니코드 공격
		"user@domain.com\u202e", // 우측-좌측 오버라이드
		"user@domain.co\u006d",  // 'm'의 유니코드
		"user@domain.c\u043em",  // 키릴 문자 'о'
		
		// 길이 0 또는 빈 문자열
		"",                       // 빈 문자열
		" ",                      // 공백만
		"@",                      // @ 만
		"@@",                     // @@ 만
		
		// 프로토콜 시도
		"ftp://user@domain.com",
		"http://user@domain.com", 
		"https://user@domain.com",
		"file://user@domain.com",
		"data:user@domain.com",
		
		// IP 주소 (일반적으로 유효하지 않음)
		"user@192.168.1.1",
		"user@[192.168.1.1]",
		"user@2001:db8::1",
		"user@[2001:db8::1]",
	}

	for i, email := range invalidSecurityCases {
		t.Run("Invalid_Security_"+string(rune('A'+i%26))+string(rune('0'+i/26)), func(t *testing.T) {
			result := strvalidator.IsValidEmail(email)
			assert.AssertFalse(t, result, "Invalid/malicious email should fail: %s", email)
		})
	}
}

func Test_Email_EdgeCaseBoundaries(t *testing.T) {
	t.Parallel()

	// 경계값 테스트
	boundaryCases := []struct {
		name     string
		email    string
		expected bool
	}{
		// 로컬 파트 길이 경계
		{"Local part 64 chars", strings.Repeat("a", 64) + "@domain.com", true},
		{"Local part 65 chars", strings.Repeat("a", 65) + "@domain.com", false},
		{"Local part 1 char", "a@domain.com", true},
		
		// 도메인 길이 경계  
		{"Domain 253 chars", "user@" + strings.Repeat("a", 249) + ".com", true},
		{"Domain 254 chars", "user@" + strings.Repeat("a", 250) + ".com", false},
		
		// 도메인 라벨 길이 경계
		{"Label 63 chars", "user@" + strings.Repeat("a", 63) + ".com", true},
		{"Label 64 chars", "user@" + strings.Repeat("a", 64) + ".com", false},
		
		// 전체 이메일 길이 경계 (일반적으로 320자 제한)
		{"Email 320 chars", strings.Repeat("a", 64) + "@" + strings.Repeat("b", 251) + ".co", true},
		{"Email over 320", strings.Repeat("a", 64) + "@" + strings.Repeat("b", 252) + ".com", false},
		
		// TLD 길이 경계
		{"TLD 2 chars", "user@domain.co", true},
		{"TLD 1 char", "user@domain.a", false},
		{"TLD 63 chars", "user@domain." + strings.Repeat("a", 63), true},
		{"TLD 64 chars", "user@domain." + strings.Repeat("a", 64), false},
	}

	for _, testCase := range boundaryCases {
		t.Run("Boundary_"+testCase.name, func(t *testing.T) {
			result := strvalidator.IsValidEmail(testCase.email)
			assert.AssertEquals(t, testCase.expected, result,
				"Boundary test failed for %s (%s).\nExpected: %v\nActual: %v",
				testCase.name, testCase.email, testCase.expected, result)
		})
	}
}

func Test_Email_InternationalDomains(t *testing.T) {
	t.Parallel()

	// 국제화 도메인 테스트 (실제로는 punycode로 변환되어야 함)
	internationalCases := []struct {
		name     string
		email    string
		expected bool
	}{
		// 실제 국제 도메인들 (현재 구현에서는 ASCII가 아니므로 실패할 수 있음)
		{"Korean domain", "사용자@한국.kr", false},           // 한글
		{"Chinese domain", "用户@中国.cn", false},            // 중국어
		{"Japanese domain", "ユーザー@日本.jp", false},        // 일본어  
		{"Arabic domain", "مستخدم@مصر.eg", false},          // 아랍어
		{"German domain", "benutzer@münchen.de", false},    // 독일어 (움라우트)
		{"French domain", "utilisateur@café.fr", false},   // 프랑스어 (악센트)
		
		// Punycode로 변환된 형태 (유효해야 함)
		{"Punycode Korean", "user@xn--3e0b707e.kr", true},
		{"Punycode Chinese", "user@xn--fiqs8s.cn", true},
		{"Punycode German", "user@xn--mnchen-3ya.de", true},
		
		// 혼합 문자 (일반적으로 무효)
		{"Mixed script", "user@도메인domain.com", false},
		{"Mixed Unicode", "user@test测试.com", false},
	}

	for _, testCase := range internationalCases {
		t.Run("International_"+testCase.name, func(t *testing.T) {
			result := strvalidator.IsValidEmail(testCase.email)
			assert.AssertEquals(t, testCase.expected, result,
				"International domain test failed for %s (%s).\nExpected: %v\nActual: %v",
				testCase.name, testCase.email, testCase.expected, result)
		})
	}
}

func Test_Email_SpecialCharacterHandling(t *testing.T) {
	t.Parallel()

	// 특수 문자 처리 테스트
	specialCharCases := []struct {
		name     string
		email    string
		expected bool
	}{
		// 로컬 파트에서 허용되는 특수 문자
		{"Plus sign", "user+tag@domain.com", true},
		{"Hyphen", "user-name@domain.com", true},
		{"Underscore", "user_name@domain.com", true},
		{"Dot", "user.name@domain.com", true},
		{"Number", "user123@domain.com", true},
		
		// 로컬 파트에서 허용되지 않는 특수 문자
		{"Exclamation", "user!@domain.com", false},
		{"Hash", "user#@domain.com", false},
		{"Dollar", "user$@domain.com", false},
		{"Percent", "user%@domain.com", false},
		{"Ampersand", "user&@domain.com", false},
		{"Asterisk", "user*@domain.com", false},
		{"Question", "user?@domain.com", false},
		{"Slash", "user/@domain.com", false},
		{"Backslash", "user\\@domain.com", false},
		{"Pipe", "user|@domain.com", false},
		{"Caret", "user^@domain.com", false},
		{"Tilde", "user~@domain.com", false},
		{"Backtick", "user`@domain.com", false},
		
		// 따옴표 관련
		{"Single quote", "user'@domain.com", false},
		{"Double quote", "user\"@domain.com", false},
		
		// 괄호 관련
		{"Parentheses", "user(@domain.com", false},
		{"Parentheses close", "user)@domain.com", false},
		{"Square bracket", "user[@domain.com", false},
		{"Square bracket close", "user]@domain.com", false},
		{"Curly brace", "user{@domain.com", false},
		{"Curly brace close", "user}@domain.com", false},
		
		// 도메인에서 허용되는 특수 문자
		{"Domain hyphen", "user@sub-domain.com", true},
		{"Domain dot", "user@sub.domain.com", true},
		
		// 도메인에서 허용되지 않는 특수 문자
		{"Domain underscore", "user@domain_name.com", false},
		{"Domain plus", "user@domain+name.com", false},
	}

	for _, testCase := range specialCharCases {
		t.Run("SpecialChar_"+testCase.name, func(t *testing.T) {
			result := strvalidator.IsValidEmail(testCase.email)
			assert.AssertEquals(t, testCase.expected, result,
				"Special character test failed for %s (%s).\nExpected: %v\nActual: %v",
				testCase.name, testCase.email, testCase.expected, result)
		})
	}
}

func Test_Email_CaseSensitivity(t *testing.T) {
	t.Parallel()

	// 대소문자 처리 테스트 (이메일은 기술적으로 대소문자를 구분하지만 실용적으로는 무시됨)
	caseCases := []string{
		"USER@DOMAIN.COM",
		"User@Domain.Com",
		"user@DOMAIN.com",
		"USER@domain.COM",
		"uSeR@DoMaIn.CoM",
		"Test.User@Example.Org",
		"ADMIN@SITE.NET",
		"Admin@Site.Net",
	}

	for _, email := range caseCases {
		t.Run("Case_"+strings.ReplaceAll(email, "@", "_at_"), func(t *testing.T) {
			result := strvalidator.IsValidEmail(email)
			assert.AssertTrue(t, result, "Email should be valid regardless of case: %s", email)
		})
	}
}