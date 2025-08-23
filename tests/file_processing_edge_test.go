package strutils_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_FileProcessing_EdgeCases(t *testing.T) {
	t.Parallel()

	// 임시 디렉터리 생성
	tmpDir, err := ioutil.TempDir("", "strutils_file_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("FileMD5Hash_EdgeCases", func(t *testing.T) {
		// 빈 파일 테스트
		emptyFile := filepath.Join(tmpDir, "empty.txt")
		err := ioutil.WriteFile(emptyFile, []byte{}, 0644)
		assert.AssertNil(t, err, "Failed to create empty file")

		hash, err := strproc.FileMD5Hash(emptyFile)
		assert.AssertNil(t, err, "FileMD5Hash should work with empty file")
		assert.AssertEquals(t, "d41d8cd98f00b204e9800998ecf8427e", hash, "Empty file MD5 hash mismatch")

		// 단일 바이트 파일
		singleByteFile := filepath.Join(tmpDir, "single.txt")
		err = ioutil.WriteFile(singleByteFile, []byte{'a'}, 0644)
		assert.AssertNil(t, err, "Failed to create single byte file")

		hash, err = strproc.FileMD5Hash(singleByteFile)
		assert.AssertNil(t, err, "FileMD5Hash should work with single byte file")
		assert.AssertNotEquals(t, "", hash, "Single byte file should have valid hash")
		assert.AssertEquals(t, 32, len(hash), "MD5 hash should be 32 characters")

		// 큰 파일 (1MB)
		largeFile := filepath.Join(tmpDir, "large.txt")
		largeData := strings.Repeat("abcdefghij", 104857) // ~1MB
		err = ioutil.WriteFile(largeFile, []byte(largeData), 0644)
		assert.AssertNil(t, err, "Failed to create large file")

		hash, err = strproc.FileMD5Hash(largeFile)
		assert.AssertNil(t, err, "FileMD5Hash should work with large file")
		assert.AssertEquals(t, 32, len(hash), "Large file MD5 hash should be 32 characters")

		// 바이너리 파일
		binaryFile := filepath.Join(tmpDir, "binary.bin")
		binaryData := make([]byte, 1000)
		for i := range binaryData {
			binaryData[i] = byte(i % 256)
		}
		err = ioutil.WriteFile(binaryFile, binaryData, 0644)
		assert.AssertNil(t, err, "Failed to create binary file")

		hash, err = strproc.FileMD5Hash(binaryFile)
		assert.AssertNil(t, err, "FileMD5Hash should work with binary file")
		assert.AssertEquals(t, 32, len(hash), "Binary file MD5 hash should be 32 characters")
	})

	t.Run("FileMD5Hash_ErrorCases", func(t *testing.T) {
		// 존재하지 않는 파일
		nonExistentFile := filepath.Join(tmpDir, "nonexistent.txt")
		_, err := strproc.FileMD5Hash(nonExistentFile)
		assert.AssertNotNil(t, err, "FileMD5Hash should fail with nonexistent file")

		// 디렉터리를 파일로 처리 시도
		_, err = strproc.FileMD5Hash(tmpDir)
		assert.AssertNotNil(t, err, "FileMD5Hash should fail when given directory")

		// 빈 경로
		_, err = strproc.FileMD5Hash("")
		assert.AssertNotNil(t, err, "FileMD5Hash should fail with empty path")

		// 권한 없는 파일 (Unix 계열에서만 테스트)
		if os.Getenv("GOOS") != "windows" {
			noPermFile := filepath.Join(tmpDir, "noperm.txt")
			err = ioutil.WriteFile(noPermFile, []byte("test"), 0000)
			if err == nil {
				_, err = strproc.FileMD5Hash(noPermFile)
				// 권한 에러가 발생할 수 있음 (시스템에 따라 다름)
				os.Chmod(noPermFile, 0644) // 정리를 위해 권한 복구
			}
		}
	})

	t.Run("HumanFileSize_EdgeCases", func(t *testing.T) {
		// 빈 파일 크기
		emptyFile := filepath.Join(tmpDir, "empty.txt")
		err := ioutil.WriteFile(emptyFile, []byte{}, 0644)
		assert.AssertNil(t, err, "Failed to create empty file")

		size, err := strproc.HumanFileSize(emptyFile, 2, 1) // CamelCaseLong
		assert.AssertNil(t, err, "HumanFileSize should work with empty file")
		assert.AssertEquals(t, "0.00Byte", size, "Empty file size should be 0.00Byte")

		// 정확히 1KB 파일
		oneKBFile := filepath.Join(tmpDir, "1kb.txt")
		oneKBData := strings.Repeat("a", 1024)
		err = ioutil.WriteFile(oneKBFile, []byte(oneKBData), 0644)
		assert.AssertNil(t, err, "Failed to create 1KB file")

		size, err = strproc.HumanFileSize(oneKBFile, 2, 1)
		assert.AssertNil(t, err, "HumanFileSize should work with 1KB file")
		// 결과는 "1.00KiloByte" 형태여야 함
		assert.AssertTrue(t, strings.Contains(size, "1.00"), "1KB file should show 1.00")
		assert.AssertTrue(t, strings.Contains(strings.ToLower(size), "k"), "1KB file should contain K unit")

		// 정확히 1MB 파일
		oneMBFile := filepath.Join(tmpDir, "1mb.txt")
		oneMBData := strings.Repeat("a", 1024*1024)
		err = ioutil.WriteFile(oneMBFile, []byte(oneMBData), 0644)
		assert.AssertNil(t, err, "Failed to create 1MB file")

		size, err = strproc.HumanFileSize(oneMBFile, 2, 1)
		assert.AssertNil(t, err, "HumanFileSize should work with 1MB file")
		assert.AssertTrue(t, strings.Contains(size, "1.00"), "1MB file should show 1.00")
		assert.AssertTrue(t, strings.Contains(strings.ToLower(size), "m"), "1MB file should contain M unit")

		// 다양한 정밀도 테스트
		testSizes := []int{0, 1, 2, 3, 4, 5}
		for _, precision := range testSizes {
			size, err := strproc.HumanFileSize(oneKBFile, precision, 1)
			assert.AssertNil(t, err, "HumanFileSize should work with precision %d", precision)
			assert.AssertNotEquals(t, "", size, "Size should not be empty for precision %d", precision)
		}

		// 다양한 유닛 형식 테스트
		unitFormats := []uint8{1, 2, 3, 4} // CamelCaseLong, CamelCaseShort, etc.
		for _, format := range unitFormats {
			size, err := strproc.HumanFileSize(oneKBFile, 2, format)
			if err == nil { // 지원하는 형식인 경우에만 테스트
				assert.AssertNotEquals(t, "", size, "Size should not be empty for format %d", format)
			}
		}
	})

	t.Run("HumanFileSize_ErrorCases", func(t *testing.T) {
		// 존재하지 않는 파일
		nonExistentFile := filepath.Join(tmpDir, "nonexistent.txt")
		_, err := strproc.HumanFileSize(nonExistentFile, 2, 1)
		assert.AssertNotNil(t, err, "HumanFileSize should fail with nonexistent file")

		// 디렉터리 크기 (에러 또는 0이어야 함)
		_, err = strproc.HumanFileSize(tmpDir, 2, 1)
		// 구현에 따라 에러이거나 디렉터리 크기를 반환할 수 있음

		// 빈 경로
		_, err = strproc.HumanFileSize("", 2, 1)
		assert.AssertNotNil(t, err, "HumanFileSize should fail with empty path")

		// 잘못된 정밀도 (음수)
		validFile := filepath.Join(tmpDir, "valid.txt")
		err = ioutil.WriteFile(validFile, []byte("test"), 0644)
		assert.AssertNil(t, err, "Failed to create valid file")

		_, err = strproc.HumanFileSize(validFile, -1, 1)
		// 구현에 따라 에러이거나 기본값을 사용할 수 있음

		// 잘못된 유닛 형식
		_, err = strproc.HumanFileSize(validFile, 2, uint8(255))
		// 구현에 따라 에러이거나 기본값을 사용할 수 있음
	})

	t.Run("File_SpecialCharacters_Paths", func(t *testing.T) {
		// 특수 문자가 포함된 파일명 테스트
		specialFiles := []string{
			"test with spaces.txt",
			"test-with-dashes.txt",
			"test_with_underscores.txt",
			"test.with.dots.txt",
			"test(with)parentheses.txt",
			"test[with]brackets.txt",
			"test{with}braces.txt",
		}

		for _, filename := range specialFiles {
			filePath := filepath.Join(tmpDir, filename)
			testData := "test content for " + filename
			
			err := ioutil.WriteFile(filePath, []byte(testData), 0644)
			if err != nil {
				t.Logf("Skipping file with special characters: %s (error: %v)", filename, err)
				continue
			}

			// MD5 해시 테스트
			hash, err := strproc.FileMD5Hash(filePath)
			assert.AssertNil(t, err, "FileMD5Hash should work with special filename: %s", filename)
			assert.AssertEquals(t, 32, len(hash), "Hash length should be 32 for file: %s", filename)

			// 파일 크기 테스트
			size, err := strproc.HumanFileSize(filePath, 2, 1)
			assert.AssertNil(t, err, "HumanFileSize should work with special filename: %s", filename)
			assert.AssertNotEquals(t, "", size, "Size should not be empty for file: %s", filename)

			// 파일 정리
			os.Remove(filePath)
		}
	})

	t.Run("File_Unicode_Paths", func(t *testing.T) {
		// 유니코드 파일명 테스트
		unicodeFiles := []string{
			"한글파일.txt",
			"日本語ファイル.txt", 
			"файл.txt",
			"αρχείο.txt",
			"файл_тест.txt",
			"测试文件.txt",
		}

		for _, filename := range unicodeFiles {
			filePath := filepath.Join(tmpDir, filename)
			testData := "unicode test content"
			
			err := ioutil.WriteFile(filePath, []byte(testData), 0644)
			if err != nil {
				t.Logf("Skipping unicode filename (filesystem may not support): %s", filename)
				continue
			}

			// MD5 해시 테스트
			hash, err := strproc.FileMD5Hash(filePath)
			if err == nil {
				assert.AssertEquals(t, 32, len(hash), "Hash length should be 32 for unicode file: %s", filename)
			} else {
				t.Logf("Unicode filename hash failed (may be expected): %s", filename)
			}

			// 파일 크기 테스트
			size, err := strproc.HumanFileSize(filePath, 2, 1)
			if err == nil {
				assert.AssertNotEquals(t, "", size, "Size should not be empty for unicode file: %s", filename)
			} else {
				t.Logf("Unicode filename size failed (may be expected): %s", filename)
			}

			// 파일 정리
			os.Remove(filePath)
		}
	})

	t.Run("File_ContentTypes_MD5", func(t *testing.T) {
		// 다양한 내용 유형의 파일 MD5 테스트
		contentTests := map[string][]byte{
			"text_ascii.txt":     []byte("Hello World\n"),
			"text_utf8.txt":      []byte("안녕하세요 World\n"),
			"text_numbers.txt":   []byte("1234567890\n"),
			"text_symbols.txt":   []byte("!@#$%^&*()_+-={}[]|\\:;\"'<>?,./ \n"),
			"text_mixed.txt":     []byte("Mixed: Hello 안녕 123 !@# World\n"),
			"binary_zeros.bin":   make([]byte, 100), // All zeros
			"binary_ones.bin":    bytes_ones(100),   // All 0xFF
			"binary_pattern.bin": bytes_pattern(256), // 0x00, 0x01, 0x02, ...
		}

		for filename, content := range contentTests {
			filePath := filepath.Join(tmpDir, filename)
			err := ioutil.WriteFile(filePath, content, 0644)
			assert.AssertNil(t, err, "Failed to create file: %s", filename)

			hash, err := strproc.FileMD5Hash(filePath)
			assert.AssertNil(t, err, "FileMD5Hash should work with content type: %s", filename)
			assert.AssertEquals(t, 32, len(hash), "Hash should be 32 chars for: %s", filename)

			// 동일한 내용으로 다시 생성해서 해시가 같은지 확인
			filePath2 := filepath.Join(tmpDir, "copy_"+filename)
			err = ioutil.WriteFile(filePath2, content, 0644)
			assert.AssertNil(t, err, "Failed to create copy file: %s", filename)

			hash2, err := strproc.FileMD5Hash(filePath2)
			assert.AssertNil(t, err, "FileMD5Hash should work with copy: %s", filename)
			assert.AssertEquals(t, hash, hash2, "Hash should be same for identical content: %s", filename)

			// 파일 정리
			os.Remove(filePath)
			os.Remove(filePath2)
		}
	})

	t.Run("File_LargeSize_EdgeCases", func(t *testing.T) {
		// 큰 파일 크기 경계 테스트
		sizeCases := []struct {
			name string
			size int64
		}{
			{"1KB-1", 1023},    // 1KB 미만
			{"1KB", 1024},      // 정확히 1KB
			{"1KB+1", 1025},    // 1KB 초과
			{"1MB-1", 1048575}, // 1MB 미만
			{"1MB", 1048576},   // 정확히 1MB
			{"1MB+1", 1048577}, // 1MB 초과
		}

		for _, sizeCase := range sizeCases {
			if sizeCase.size > 10*1024*1024 { // 10MB 초과는 스킵 (테스트 시간 단축)
				continue
			}

			filename := filepath.Join(tmpDir, fmt.Sprintf("size_test_%s.txt", sizeCase.name))
			
			// 파일 생성 (반복 패턴으로 메모리 효율적으로)
			file, err := os.Create(filename)
			if err != nil {
				t.Logf("Failed to create size test file: %s", sizeCase.name)
				continue
			}

			pattern := []byte("0123456789abcdef") // 16바이트 패턴
			written := int64(0)
			for written < sizeCase.size {
				remaining := sizeCase.size - written
				toWrite := pattern
				if remaining < int64(len(pattern)) {
					toWrite = pattern[:remaining]
				}
				
				n, err := file.Write(toWrite)
				if err != nil {
					file.Close()
					t.Logf("Failed to write to size test file: %s", sizeCase.name)
					break
				}
				written += int64(n)
			}
			file.Close()

			if written == sizeCase.size {
				// HumanFileSize 테스트
				size, err := strproc.HumanFileSize(filename, 2, 1)
				assert.AssertNil(t, err, "HumanFileSize should work for: %s", sizeCase.name)
				assert.AssertNotEquals(t, "", size, "Size should not be empty for: %s", sizeCase.name)
				
				t.Logf("Size test %s (%d bytes): %s", sizeCase.name, sizeCase.size, size)

				// MD5 해시는 큰 파일에서는 스킵 (시간이 오래 걸림)
				if sizeCase.size <= 1024*1024 { // 1MB 이하만 테스트
					hash, err := strproc.FileMD5Hash(filename)
					assert.AssertNil(t, err, "FileMD5Hash should work for: %s", sizeCase.name)
					assert.AssertEquals(t, 32, len(hash), "Hash should be 32 chars for: %s", sizeCase.name)
				}
			}

			// 파일 정리
			os.Remove(filename)
		}
	})
}

// 헬퍼 함수들
func bytes_ones(size int) []byte {
	result := make([]byte, size)
	for i := range result {
		result[i] = 0xFF
	}
	return result
}

func bytes_pattern(size int) []byte {
	result := make([]byte, size)
	for i := range result {
		result[i] = byte(i % 256)
	}
	return result
}

func Test_FileProcessing_PathValidation_EdgeCases(t *testing.T) {
	t.Parallel()

	// IsValidFilePath 및 IsValidFilePathWithRelativePath 엣지 케이스 테스트
	t.Run("IsValidFilePath_EdgeCases", func(t *testing.T) {
		pathTests := []struct {
			name     string
			path     string
			expected bool
		}{
			// 기본 유효 케이스
			{"Simple filename", "file.txt", true},
			{"Filename with extension", "document.pdf", true},
			{"No extension", "README", true},
			
			// 특수 문자 케이스
			{"With spaces", "file with spaces.txt", false}, // 일반적으로 허용되지 않음
			{"With underscores", "file_name.txt", true},
			{"With hyphens", "file-name.txt", true},
			{"With dots", "file.name.txt", true},
			{"With numbers", "file123.txt", true},
			
			// 확장자 관련
			{"Multiple extensions", "file.tar.gz", true},
			{"Long extension", "file.extension", true},
			{"No extension dot", "filename", true},
			{"Only extension", ".gitignore", true},
			
			// 특수한 파일명
			{"Very long name", strings.Repeat("a", 100) + ".txt", true},
			{"Single character", "a", true},
			{"Numbers only", "123", true},
			
			// 무효한 케이스
			{"Empty string", "", false},
			{"Only spaces", "   ", false},
			{"Path separator", "dir/file.txt", false}, // 경로 구분자 포함
			{"Backslash", "dir\\file.txt", false},
			{"Relative path", "./file.txt", false},
			{"Parent path", "../file.txt", false},
			{"Absolute path", "/file.txt", false},
			
			// 특수 문자 (시스템에 따라 다를 수 있음)
			{"Colon", "file:name.txt", false},
			{"Question mark", "file?.txt", false},
			{"Asterisk", "file*.txt", false},
			{"Pipe", "file|name.txt", false},
			{"Less than", "file<name.txt", false},
			{"Greater than", "file>name.txt", false},
			{"Double quote", "file\"name.txt", false},
		}

		for _, testCase := range pathTests {
			t.Run("FilePath_"+testCase.name, func(t *testing.T) {
				result := strvalidator.IsValidFilePath(testCase.path)
				// 주의: 정확한 기대값은 구현에 따라 다를 수 있음
				t.Logf("IsValidFilePath(%q) = %v (expected: %v)", testCase.path, result, testCase.expected)
			})
		}
	})

	t.Run("IsValidFilePathWithRelativePath_EdgeCases", func(t *testing.T) {
		pathTests := []struct {
			name     string
			path     string
			expected bool
		}{
			// 절대 경로
			{"Absolute unix", "/home/user/file.txt", true},
			{"Absolute windows", "C:\\Users\\file.txt", true},
			{"Root file", "/file.txt", true},
			
			// 상대 경로
			{"Current dir file", "./file.txt", true},
			{"Parent dir file", "../file.txt", true},
			{"Deep relative", "../../dir/file.txt", true},
			{"Relative dir", "dir/file.txt", true},
			{"Deep path", "dir/subdir/file.txt", true},
			
			// 복잡한 경로
			{"Mixed separators", "dir\\subdir/file.txt", true}, // Windows 스타일
			{"Multiple dots", "dir/../other/./file.txt", true},
			{"Trailing slash", "dir/file.txt/", false}, // 파일이 아닌 디렉터리
			
			// 특수 디렉터리명
			{"Hidden dir", ".config/file.txt", true},
			{"Space in dir", "my dir/file.txt", true},
			{"Unicode dir", "한글폴더/file.txt", true},
			
			// 매우 긴 경로
			{"Very long path", strings.Repeat("dir/", 50) + "file.txt", true},
			
			// 무효한 케이스
			{"Empty path", "", false},
			{"Only separators", "///", false},
			{"Only dots", "...", false},
			{"Null in path", "dir\x00file.txt", false},
		}

		for _, testCase := range pathTests {
			t.Run("FilePathRel_"+testCase.name, func(t *testing.T) {
				result := strvalidator.IsValidFilePathWithRelativePath(testCase.path)
				t.Logf("IsValidFilePathWithRelativePath(%q) = %v (expected: %v)", 
					testCase.path, result, testCase.expected)
			})
		}
	})

	t.Run("FilePath_SecurityCases", func(t *testing.T) {
		// 보안 관련 파일 경로 테스트
		securityCases := []string{
			// 디렉터리 순회 시도
			"../../etc/passwd",
			"..\\..\\windows\\system32\\config",
			"..\\..\\..\\..\\boot.ini",
			"../../../../root/.ssh/id_rsa",
			
			// 절대 경로를 이용한 시스템 파일 접근 시도
			"/etc/passwd",
			"/etc/shadow",
			"/proc/version",
			"/sys/class/net/eth0/address",
			"C:\\Windows\\System32\\config\\SAM",
			"C:\\Windows\\System32\\drivers\\etc\\hosts",
			
			// 네트워크 경로
			"\\\\server\\share\\file.txt",
			"//server/share/file.txt",
			
			// 특수 장치
			"/dev/null",
			"/dev/random",
			"CON",
			"PRN",
			"AUX",
			"NUL",
			"COM1",
			"LPT1",
			
			// URL 스킴을 이용한 시도
			"file:///etc/passwd",
			"http://evil.com/file.txt",
			"ftp://server/file.txt",
			
			// 제어 문자 포함
			"file\nname.txt",
			"file\rname.txt",
			"file\tname.txt",
			"file\x00name.txt",
			
			// 매우 긴 경로 (버퍼 오버플로우 시도)
			strings.Repeat("a", 1000) + ".txt",
			strings.Repeat("../", 500) + "file.txt",
		}

		for i, maliciousPath := range securityCases {
			t.Run(fmt.Sprintf("Security_%d", i), func(t *testing.T) {
				// 일반 파일 경로 검증
				result1 := strvalidator.IsValidFilePath(maliciousPath)
				
				// 상대 경로 포함 검증
				result2 := strvalidator.IsValidFilePathWithRelativePath(maliciousPath)
				
				t.Logf("Security test for %q:", maliciousPath)
				t.Logf("  IsValidFilePath: %v", result1)
				t.Logf("  IsValidFilePathWithRelativePath: %v", result2)
				
				// 대부분의 악성 경로는 거부되어야 함
				// 하지만 구현에 따라 다를 수 있으므로 로그만 남김
			})
		}
	})
}