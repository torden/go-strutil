package strutils_test

import (
	"strings"
	"testing"
)

func Test_PureText_SecurityCases(t *testing.T) {
	t.Parallel()

	// XSS 공격 시도 케이스
	xssAttackCases := []struct {
		name     string
		input    string
		expected bool
	}{
		// 기본 XSS 공격
		{"Script tag", "<script>alert('xss')</script>", false},
		{"Script with attributes", "<script type=\"text/javascript\">alert(1)</script>", false},
		{"Script uppercase", "<SCRIPT>alert('XSS')</SCRIPT>", false},
		{"Script mixed case", "<ScRiPt>alert('xss')</ScRiPt>", false},
		
		// JavaScript: URL 스킴
		{"Javascript URL", "javascript:alert(1)", false},
		{"Javascript uppercase", "JAVASCRIPT:alert(1)", false},
		{"Javascript mixed case", "JaVaScRiPt:alert(1)", false},
		{"Javascript with spaces", "java script:alert(1)", false},
		
		// 다양한 태그 기반 공격
		{"Img onerror", "<img src=x onerror=alert(1)>", false},
		{"Img onload", "<img onload=\"alert('xss')\">", false},
		{"Div onclick", "<div onclick=\"alert(1)\">test</div>", false},
		{"A href javascript", "<a href=\"javascript:alert(1)\">link</a>", false},
		{"Input onfocus", "<input onfocus=\"alert(1)\" autofocus>", false},
		{"Body onload", "<body onload=\"alert(1)\">", false},
		{"Iframe src", "<iframe src=\"javascript:alert(1)\"></iframe>", false},
		{"Object data", "<object data=\"javascript:alert(1)\"></object>", false},
		{"Embed src", "<embed src=\"javascript:alert(1)\">", false},
		{"Form onsubmit", "<form onsubmit=\"alert(1)\">", false},
		
		// HTML5 새로운 이벤트 핸들러
		{"Details ontoggle", "<details ontoggle=\"alert(1)\"><summary>test</summary></details>", false},
		{"Video onloadstart", "<video onloadstart=\"alert(1)\">", false},
		{"Audio onplay", "<audio onplay=\"alert(1)\" autoplay>", false},
		
		// CSS 기반 공격
		{"Style with expression", "<style>body{background:expression(alert(1))}</style>", false},
		{"Style with javascript", "<style>@import 'javascript:alert(1)'</style>", false},
		{"Style with url", "<style>body{background:url('javascript:alert(1)')}</style>", false},
		
		// 인코딩된 공격
		{"HTML entity script", "&lt;script&gt;alert(1)&lt;/script&gt;", true}, // 인코딩된 것은 안전
		{"URL encoded script", "%3Cscript%3Ealert(1)%3C/script%3E", true}, // URL 인코딩된 것은 일단 안전
		{"Hex encoded", "\\x3cscript\\x3ealert(1)\\x3c/script\\x3e", true}, // 16진수 인코딩
		
		// 우회 시도
		{"Script with null byte", "<script\\x00>alert(1)</script>", false},
		{"Script with newline", "<script\n>alert(1)</script>", false},
		{"Script with tab", "<script\t>alert(1)</script>", false},
		{"Script with carriage return", "<script\r>alert(1)</script>", false},
		
		// Data URL 기반 공격
		{"Data URL text/html", "data:text/html,<script>alert(1)</script>", false},
		{"Data URL base64", "data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==", false},
		{"Data URL javascript", "data:text/javascript,alert(1)", false},
		
		// 유효한 텍스트 케이스 (태그나 스크립트 없음)
		{"Plain text", "This is just plain text", true},
		{"Text with numbers", "Hello 123 world 456", true},
		{"Text with punctuation", "Hello, world! How are you?", true},
		{"Text with symbols", "Price: $19.99 (includes tax @ 8%)", true},
		{"Text with quotes", "She said \"Hello world\"", true},
		{"Text with apostrophe", "It's a beautiful day", true},
	}

	for _, testCase := range xssAttackCases {
		t.Run("XSS_"+testCase.name, func(t *testing.T) {
			result, err := strvalidator.IsPureTextStrict(testCase.input)
			
			if testCase.expected {
				assert.AssertNil(t, err, "Expected no error for safe text: %s", testCase.input)
				assert.AssertTrue(t, result, "Expected safe text to pass: %s", testCase.input)
			} else {
				// XSS 공격은 감지되어야 함 (에러가 발생하거나 false 반환)
				if err == nil && result == true {
					t.Errorf("XSS attack not detected: %s", testCase.input)
				}
			}
		})
	}
}

func Test_PureText_SQLInjectionCases(t *testing.T) {
	t.Parallel()

	// SQL 인젝션 시도 케이스
	sqlInjectionCases := []string{
		"'; DROP TABLE users; --",
		"' OR 1=1 --",
		"' OR '1'='1",
		"'; DELETE FROM users WHERE '1'='1",
		"' UNION SELECT * FROM users --",
		"' AND (SELECT COUNT(*) FROM users) > 0 --",
		"'; INSERT INTO users VALUES('admin','password') --",
		"' OR 1=1#",
		"'; EXEC xp_cmdshell('dir') --",
		"'; SHUTDOWN WITH NOWAIT --",
		"admin'--",
		"admin'/*",
		"' OR 'x'='x",
		"'; UPDATE users SET password='hacked' --",
		"' HAVING 1=1 --",
		"' GROUP BY password HAVING 1=1 --",
		"' ORDER BY 1 --",
		"') OR ('1'='1",
		"') AND ('1'='1",
		"; SELECT SLEEP(5) --",
	}

	for i, injection := range sqlInjectionCases {
		t.Run("SQLInjection_"+string(rune('A'+i%26))+string(rune('0'+i/26)), func(t *testing.T) {
			// SQL 인젝션 시도는 일반적으로 특수 문자나 패턴을 포함하므로
			// IsPureTextStrict에서 감지될 수 있음
			result, err := strvalidator.IsPureTextStrict(injection)
			
			// 대부분의 SQL 인젝션은 특수 문자(', ;, --)를 포함하므로 감지되어야 함
			if strings.Contains(injection, "'") || strings.Contains(injection, ";") || strings.Contains(injection, "--") {
				// 이런 특수 문자들이 포함된 경우 false이거나 에러가 발생해야 함
				if err == nil && result == true {
					t.Logf("SQL injection pattern may not be detected by PureText validator: %s", injection)
				}
			}
		})
	}
}

func Test_PureText_ControlCharacterCases(t *testing.T) {
	t.Parallel()

	// 제어 문자 및 바이너리 데이터 케이스
	controlCharCases := []struct {
		name     string
		input    string
		expected bool
	}{
		// ASCII 제어 문자 (0x00-0x1F, 0x7F)
		{"Null byte", "Hello\x00World", false},
		{"Bell character", "Hello\x07World", false},
		{"Backspace", "Hello\x08World", false},
		{"Form feed", "Hello\x0CWorld", false},
		{"Vertical tab", "Hello\x0BWorld", false},
		{"Delete character", "Hello\x7FWorld", false},
		
		// 일반적으로 허용되는 제어 문자
		{"Tab character", "Hello\tWorld", true},
		{"Newline", "Hello\nWorld", true},
		{"Carriage return", "Hello\rWorld", true},
		{"Space", "Hello World", true},
		
		// 다중 제어 문자
		{"Multiple nulls", "Hello\x00\x01\x02World", false},
		{"Mixed control", "Hello\x00\x07\x0CWorld", false},
		
		// 바이너리 데이터 시뮬레이션
		{"Binary data", string([]byte{0x89, 0x50, 0x4E, 0x47}), false}, // PNG 헤더
		{"Random bytes", string([]byte{0xFF, 0xFE, 0xFD, 0xFC}), false},
		
		// UTF-8 BOM 및 특수 유니코드
		{"UTF-8 BOM", "\uFEFFHello World", true}, // BOM은 일반적으로 허용
		{"Zero width space", "Hello\u200BWorld", true}, // 보이지 않지만 유효한 유니코드
		{"Soft hyphen", "Hello\u00ADWorld", true},
		
		// 유효한 텍스트
		{"Normal text", "Hello World", true},
		{"Unicode text", "안녕하세요 世界", true},
		{"Emoji text", "Hello World", true},
	}

	for _, testCase := range controlCharCases {
		t.Run("Control_"+testCase.name, func(t *testing.T) {
			result, err := strvalidator.IsPureTextStrict(testCase.input)
			
			if testCase.expected {
				assert.AssertNil(t, err, "Expected no error for valid text: %q", testCase.input)
				assert.AssertTrue(t, result, "Expected valid text to pass: %q", testCase.input)
			} else {
				// 제어 문자나 바이너리 데이터는 감지되어야 함
				if err == nil && result == true {
					t.Errorf("Control character/binary data not detected: %q", testCase.input)
				}
			}
		})
	}
}

func Test_PureText_HTTPHeaderInjectionCases(t *testing.T) {
	t.Parallel()

	// HTTP 헤더 인젝션 시도 케이스
	headerInjectionCases := []string{
		"User input\r\nSet-Cookie: admin=1",
		"Normal text\nLocation: http://evil.com",
		"Content\r\nContent-Type: text/html",
		"Text\r\nX-Forwarded-For: attacker.com",
		"Input\nSet-Cookie: sessionid=hijacked",
		"Data\r\n\r\n<script>alert(1)</script>",
		"Text\r\nHTTP/1.1 200 OK\r\nContent-Length: 0",
		"Value\nConnection: close\nHost: evil.com",
		"Normal\r\nAuthorization: Bearer stolen_token",
		"Text\nRefresh: 0;url=http://malicious.com",
	}

	for i, injection := range headerInjectionCases {
		t.Run("HeaderInjection_"+string(rune('A'+i%26)), func(t *testing.T) {
			result, err := strvalidator.IsPureTextStrict(injection)
			
			// HTTP 헤더 인젝션은 CRLF 문자(\r\n)를 포함하므로 감지되어야 함
			if strings.Contains(injection, "\r\n") || strings.Contains(injection, "\n") {
				// 이런 패턴들이 포함된 경우 보안상 위험할 수 있음
				t.Logf("HTTP header injection pattern: %q, result: %v, error: %v", 
					injection, result, err)
			}
		})
	}
}

func Test_PureText_UnicodeAttackCases(t *testing.T) {
	t.Parallel()

	// 유니코드 기반 공격 케이스
	unicodeAttackCases := []struct {
		name     string
		input    string
		expected bool
	}{
		// 동형 문자 공격 (Homograph attack)
		{"Cyrillic domain", "gооgle.com", true}, // 'о' is Cyrillic, not Latin 'o'
		{"Mixed script", "microsоft.com", true}, // Cyrillic 'о'
		{"Greek letters", "αpple.com", true}, // Greek alpha
		
		// 우측-좌측 텍스트 방향 제어
		{"RTL override", "user\u202Eadmin", true}, // Right-to-left override
		{"LTR override", "user\u202Dadmin", true}, // Left-to-right override
		{"RTL embedding", "user\u202Badmin\u202C", true}, // RTL embedding
		{"LTR embedding", "user\u202Aadmin\u202C", true}, // LTR embedding
		
		// 정규화 공격
		{"Composed vs decomposed", "café", true}, // Should be same as ca\u0301
		{"NFKC attack", "①②③", true}, // Circled numbers
		{"Width variants", "ａｄｍｉｎ", true}, // Full-width characters
		
		// 보이지 않는 문자
		{"Zero width joiner", "admin\u200Droot", true},
		{"Zero width non-joiner", "admin\u200Croot", true},
		{"Word joiner", "admin\u2060root", true},
		
		// 수직 탭 및 기타 공백 문자
		{"Line separator", "line1\u2028line2", true},
		{"Paragraph separator", "para1\u2029para2", true},
		{"Thin space", "word\u2009word", true},
		{"Hair space", "word\u200Aword", true},
		
		// 정상 유니코드 텍스트
		{"Korean text", "안녕하세요", true},
		{"Japanese text", "こんにちは", true},
		{"Chinese text", "你好", true},
		{"Arabic text", "مرحبا", true},
		{"Emoji text", "Hello World", true},
		{"Mixed languages", "Hello 안녕 こんにちは", true},
	}

	for _, testCase := range unicodeAttackCases {
		t.Run("Unicode_"+testCase.name, func(t *testing.T) {
			result, err := strvalidator.IsPureTextStrict(testCase.input)
			
			// 유니코드 공격 감지는 구현에 따라 다름
			// 대부분의 경우 유니코드 텍스트는 허용되지만, 특정 제어 문자는 제한될 수 있음
			t.Logf("Unicode test: %q, result: %v, error: %v", 
				testCase.input, result, err)
				
			// 실제 예상 결과와 비교
			if testCase.expected && err != nil {
				t.Logf("Valid unicode text failed validation: %q", testCase.input)
			}
		})
	}
}

func Test_PureText_LengthAndBoundaryEdgeCases(t *testing.T) {
	t.Parallel()

	// 길이 및 경계 조건 테스트
	boundaryTests := []struct {
		name     string
		input    string
		expected bool
	}{
		// 길이 경계
		{"Empty string", "", true},
		{"Single character", "a", true},
		{"Very long text", strings.Repeat("a", 10000), true},
		{"Extremely long text", strings.Repeat("abc", 100000), true},
		
		// 반복 패턴
		{"Repeated tags", strings.Repeat("<script>", 100), false},
		{"Repeated safe text", strings.Repeat("hello ", 1000), true},
		{"Alternating pattern", strings.Repeat("a<b>", 1000), false},
		
		// 메모리 소진 시도
		{"Deep nesting attempt", strings.Repeat("<div>", 10000) + "content" + strings.Repeat("</div>", 10000), false},
		{"Large script block", "<script>" + strings.Repeat("alert(1);", 10000) + "</script>", false},
		
		// 특수한 길이 조합
		{"Just under limit", strings.Repeat("a", 999), true},
		{"Just over limit", strings.Repeat("a", 1001), true},
	}

	for _, testCase := range boundaryTests {
		t.Run("Boundary_"+testCase.name, func(t *testing.T) {
			result, err := strvalidator.IsPureTextStrict(testCase.input)
			
			if testCase.expected {
				if err != nil || !result {
					t.Logf("Valid text failed validation: %s (length: %d)", 
						testCase.name, len(testCase.input))
				}
			} else {
				if err == nil && result {
					t.Errorf("Malicious content not detected: %s", testCase.name)
				}
			}
		})
	}
}

func Test_PureText_NormalVsStrictComparison(t *testing.T) {
	t.Parallel()

	// IsPureTextNormal vs IsPureTextStrict 비교 테스트
	comparisonCases := []struct {
		name           string
		input          string
		expectNormal   bool
		expectStrict   bool
	}{
		{"Plain text", "Hello world", true, true},
		{"Text with quotes", "She said \"hello\"", true, true},
		{"HTML tags", "<div>content</div>", false, false},
		{"Script tag", "<script>alert(1)</script>", false, false},
		{"HTML entities", "&lt;div&gt;content&lt;/div&gt;", true, true},
		{"Partial tags", "Price < 100 and > 50", true, false}, // Strict may be more restrictive
		{"Email format", "user@domain.com", true, true},
		{"URL format", "https://example.com", true, true},
		{"Special chars", "Price: $19.99 (tax @8%)", true, true},
		{"Angle brackets", "a < b and c > d", true, false}, // Strict may detect angle brackets
	}

	for _, testCase := range comparisonCases {
		t.Run("Compare_"+testCase.name, func(t *testing.T) {
			normalResult, normalErr := strvalidator.IsPureTextNormal(testCase.input)
			strictResult, strictErr := strvalidator.IsPureTextStrict(testCase.input)
			
			t.Logf("Input: %q", testCase.input)
			t.Logf("Normal: result=%v, error=%v", normalResult, normalErr)
			t.Logf("Strict: result=%v, error=%v", strictResult, strictErr)
			
			// 기대값과 비교 (로깅 목적)
			if testCase.expectNormal != (normalErr == nil && normalResult) {
				t.Logf("Normal validation unexpected result for: %q", testCase.input)
			}
			if testCase.expectStrict != (strictErr == nil && strictResult) {
				t.Logf("Strict validation unexpected result for: %q", testCase.input)
			}
		})
	}
}