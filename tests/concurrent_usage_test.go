package strutils_test

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	strutils "github.com/torden/go-strutil"
)

func Test_Concurrent_StringProc_Usage(t *testing.T) {
	t.Parallel()

	const numGoroutines = 50
	const numOperations = 100

	// 여러 StringProc 인스턴스를 동시에 사용
	t.Run("Multiple_StringProc_Instances", func(t *testing.T) {
		var wg sync.WaitGroup
		results := make(chan string, numGoroutines*numOperations)
		errors := make(chan error, numGoroutines*numOperations)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				
				// 각 고루틴에서 새로운 StringProc 인스턴스 생성
				localStrproc := strutils.NewStringProc()
				
				for j := 0; j < numOperations; j++ {
					testData := fmt.Sprintf("test_data_%d_%d", goroutineID, j)
					
					// 다양한 StringProc 메소드 테스트
					switch j % 5 {
					case 0:
						result := localStrproc.AddSlashes(testData)
						results <- fmt.Sprintf("AddSlashes_%d_%d: %s", goroutineID, j, result)
					case 1:
						result := localStrproc.StripSlashes(testData)
						results <- fmt.Sprintf("StripSlashes_%d_%d: %s", goroutineID, j, result)
					case 2:
						result := localStrproc.UpperCaseFirstWords(testData)
						results <- fmt.Sprintf("UpperCase_%d_%d: %s", goroutineID, j, result)
					case 3:
						result, err := localStrproc.NumberFmt(12345 + j)
						if err != nil {
							errors <- err
						} else {
							results <- fmt.Sprintf("NumberFmt_%d_%d: %s", goroutineID, j, result)
						}
					case 4:
						result := localStrproc.PaddingLeft(testData, "*", 20)
						results <- fmt.Sprintf("PaddingLeft_%d_%d: %s", goroutineID, j, result)
					}
				}
			}(i)
		}

		wg.Wait()
		close(results)
		close(errors)

		// 결과 검증
		resultCount := 0
		for range results {
			resultCount++
		}

		errorCount := 0
		for range errors {
			errorCount++
		}

		expectedResults := numGoroutines * numOperations
		assert.AssertEquals(t, expectedResults, resultCount, 
			"Expected %d results, got %d", expectedResults, resultCount)
		assert.AssertEquals(t, 0, errorCount, 
			"Expected 0 errors, got %d", errorCount)
	})

	// 단일 StringProc 인스턴스를 여러 고루틴에서 공유 사용
	t.Run("Shared_StringProc_Instance", func(t *testing.T) {
		var wg sync.WaitGroup
		sharedStrproc := strutils.NewStringProc()
		results := make([]string, numGoroutines*numOperations)
		var mu sync.Mutex

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				
				for j := 0; j < numOperations; j++ {
					testData := fmt.Sprintf("shared_test_%d_%d", goroutineID, j)
					
					// 공유된 StringProc 인스턴스 사용
					result := sharedStrproc.ReverseStr(testData)
					
					mu.Lock()
					results[goroutineID*numOperations+j] = result
					mu.Unlock()
				}
			}(i)
		}

		wg.Wait()

		// 결과 검증 - 모든 결과가 올바르게 저장되었는지 확인
		for i, result := range results {
			assert.AssertNotEquals(t, "", result, 
				"Result at index %d should not be empty", i)
		}
	})
}

func Test_Concurrent_StringValidator_Usage(t *testing.T) {
	t.Parallel()

	const numGoroutines = 30
	const numValidations = 50

	// 여러 StringValidator 인스턴스를 동시에 사용
	t.Run("Multiple_StringValidator_Instances", func(t *testing.T) {
		var wg sync.WaitGroup
		results := make(chan bool, numGoroutines*numValidations)
		
		testEmails := []string{
			"test@example.com",
			"user@domain.org",
			"admin@site.net",
			"invalid@",
			"@invalid.com",
			"not_an_email",
			"valid.email@test.com",
		}

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				
				// 각 고루틴에서 새로운 StringValidator 인스턴스 생성
				localValidator := strutils.NewStringValidator()
				
				for j := 0; j < numValidations; j++ {
					email := testEmails[j%len(testEmails)]
					result := localValidator.IsValidEmail(email)
					results <- result
				}
			}(i)
		}

		wg.Wait()
		close(results)

		// 결과 검증
		resultCount := 0
		for range results {
			resultCount++
		}

		expectedResults := numGoroutines * numValidations
		assert.AssertEquals(t, expectedResults, resultCount, 
			"Expected %d validation results, got %d", expectedResults, resultCount)
	})

	// 다양한 검증 메소드 동시 실행
	t.Run("Mixed_Validation_Methods", func(t *testing.T) {
		var wg sync.WaitGroup
		validator := strutils.NewStringValidator()
		var results sync.Map

		testData := map[string]interface{}{
			"emails": []string{"test@example.com", "invalid@", "user@domain.org"},
			"domains": []string{"example.com", "invalid_domain", "google.com"},
			"urls": []string{"https://example.com", "invalid_url", "http://test.org"},
			"ips": []string{"192.168.1.1", "256.256.256.256", "2001:db8::1"},
		}

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				
				// 이메일 검증
				for _, email := range testData["emails"].([]string) {
					result := validator.IsValidEmail(email)
					key := fmt.Sprintf("email_%d_%s", goroutineID, email)
					results.Store(key, result)
				}
				
				// 도메인 검증
				for _, domain := range testData["domains"].([]string) {
					result := validator.IsValidDomain(domain)
					key := fmt.Sprintf("domain_%d_%s", goroutineID, domain)
					results.Store(key, result)
				}
				
				// URL 검증
				for _, url := range testData["urls"].([]string) {
					result := validator.IsValidURL(url)
					key := fmt.Sprintf("url_%d_%s", goroutineID, url)
					results.Store(key, result)
				}
				
				// IP 주소 검증
				for _, ip := range testData["ips"].([]string) {
					result, _ := validator.IsValidIPAddr(ip, strutils.IPv4)
					key := fmt.Sprintf("ip_%d_%s", goroutineID, ip)
					results.Store(key, result)
				}
			}(i)
		}

		wg.Wait()

		// 결과 카운팅
		count := 0
		results.Range(func(key, value interface{}) bool {
			count++
			return true
		})

		expectedCount := numGoroutines * (len(testData["emails"].([]string)) + 
			len(testData["domains"].([]string)) + 
			len(testData["urls"].([]string)) + 
			len(testData["ips"].([]string)))
		
		assert.AssertEquals(t, expectedCount, count, 
			"Expected %d results, got %d", expectedCount, count)
	})
}

func Test_Concurrent_StressTest(t *testing.T) {
	t.Parallel()

	// 시스템 리소스 기반 동적 고루틴 수 결정
	numCPU := runtime.NumCPU()
	numGoroutines := numCPU * 10 // CPU 수의 10배
	numOperations := 200

	t.Run("High_Load_Stress_Test", func(t *testing.T) {
		var wg sync.WaitGroup
		startTime := time.Now()
		
		// 에러 카운터
		var errorCount int64
		var successCount int64
		var mu sync.Mutex

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				
				localStrproc := strutils.NewStringProc()
				localValidator := strutils.NewStringValidator()
				
				for j := 0; j < numOperations; j++ {
					// 랜덤 데이터 생성
					dataSize := rand.Intn(1000) + 10
					testData := generateRandomString(dataSize)
					
					// 랜덤 작업 선택
					operation := rand.Intn(10)
					
					success := true
					switch operation {
					case 0, 1:
						localStrproc.AddSlashes(testData)
					case 2, 3:
						localStrproc.StripSlashes(testData)
					case 4:
						localStrproc.UpperCaseFirstWords(testData)
					case 5:
						localStrproc.ReverseStr(testData)
					case 6:
						_, err := localStrproc.NumberFmt(rand.Intn(1000000))
						if err != nil {
							success = false
						}
					case 7:
						localValidator.IsValidEmail(testData + "@test.com")
					case 8:
						localValidator.IsValidDomain("test.com")
					case 9:
						_, err := localValidator.IsPureTextStrict(testData)
						if err != nil {
							// PureText 검증에서 에러는 예상된 동작일 수 있음
						}
					}
					
					mu.Lock()
					if success {
						successCount++
					} else {
						errorCount++
					}
					mu.Unlock()
				}
			}(i)
		}

		wg.Wait()
		duration := time.Since(startTime)

		totalOperations := int64(numGoroutines * numOperations)
		operationsPerSecond := float64(totalOperations) / duration.Seconds()

		t.Logf("Stress test completed:")
		t.Logf("  Total operations: %d", totalOperations)
		t.Logf("  Success operations: %d", successCount)
		t.Logf("  Error operations: %d", errorCount)
		t.Logf("  Duration: %v", duration)
		t.Logf("  Operations per second: %.2f", operationsPerSecond)
		t.Logf("  Goroutines: %d", numGoroutines)

		// 기본 성공률 검증 (최소 90% 성공)
		successRate := float64(successCount) / float64(totalOperations) * 100
		assert.AssertTrue(t, successRate >= 90.0, 
			"Success rate should be at least 90%%, got %.2f%%", successRate)
	})
}

func Test_Concurrent_RaceConditionDetection(t *testing.T) {
	t.Parallel()

	// race condition 감지를 위한 테스트
	// go test -race로 실행할 때 race condition 감지 가능

	const numGoroutines = 100
	const numIterations = 100

	t.Run("Race_Condition_Detection", func(t *testing.T) {
		var wg sync.WaitGroup
		sharedStrproc := strutils.NewStringProc()
		sharedValidator := strutils.NewStringValidator()
		
		// 공유 데이터 (의도적으로 race condition 유발 가능성 테스트)
		sharedData := "shared_test_data"

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(goroutineID int) {
				defer wg.Done()
				
				for j := 0; j < numIterations; j++ {
					// 동일한 데이터로 동시에 여러 작업 수행
					_ = sharedStrproc.AddSlashes(sharedData)
					_ = sharedStrproc.StripSlashes(sharedData)
					_ = sharedStrproc.UpperCaseFirstWords(sharedData)
					_ = sharedValidator.IsValidEmail(sharedData + "@test.com")
					_ = sharedValidator.IsValidDomain("test.com")
				}
			}(i)
		}

		wg.Wait()
	})
}

func Test_Concurrent_MemoryUsage(t *testing.T) {
	t.Parallel()

	// 메모리 사용량 모니터링 테스트
	t.Run("Memory_Usage_Monitoring", func(t *testing.T) {
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		const numGoroutines = 50
		const numOperations = 200

		var wg sync.WaitGroup
		
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				
				strproc := strutils.NewStringProc()
				validator := strutils.NewStringValidator()
				
				for j := 0; j < numOperations; j++ {
					largeData := generateRandomString(1000)
					
					_ = strproc.AddSlashes(largeData)
					_ = strproc.StripSlashes(largeData)
					_ = validator.IsValidEmail("test@example.com")
					_, _ = strproc.NumberFmt(j)
				}
			}()
		}

		wg.Wait()
		
		runtime.GC()
		runtime.ReadMemStats(&m2)

		allocatedBytes := m2.TotalAlloc - m1.TotalAlloc
		heapGrowth := m2.HeapInuse - m1.HeapInuse

		t.Logf("Memory usage after concurrent operations:")
		t.Logf("  Total allocated: %d bytes", allocatedBytes)
		t.Logf("  Heap growth: %d bytes", heapGrowth)
		t.Logf("  Number of GC cycles: %d", m2.NumGC-m1.NumGC)

		// 메모리 누수 기본 검사 (힙 증가가 너무 크지 않은지 확인)
		maxExpectedHeapGrowth := uint64(50 * 1024 * 1024) // 50MB
		assert.AssertTrue(t, heapGrowth < maxExpectedHeapGrowth, 
			"Heap growth %d bytes exceeds expected maximum %d bytes", 
			heapGrowth, maxExpectedHeapGrowth)
	})
}

// 헬퍼 함수: 랜덤 문자열 생성
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// 벤치마크: 동시성 성능 측정
func BenchmarkConcurrentStringProc(b *testing.B) {
	strproc := strutils.NewStringProc()
	testData := "benchmark_test_data_with_some_length"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = strproc.AddSlashes(testData)
			_ = strproc.StripSlashes(testData)
			_ = strproc.UpperCaseFirstWords(testData)
		}
	})
}

func BenchmarkConcurrentStringValidator(b *testing.B) {
	validator := strutils.NewStringValidator()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = validator.IsValidEmail("test@example.com")
			_ = validator.IsValidDomain("example.com")
			_ = validator.IsValidURL("https://example.com")
		}
	})
}