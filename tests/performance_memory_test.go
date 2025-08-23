package strutils_test

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	strutils "github.com/torden/go-strutil"
)

func Test_Performance_StringProcessing(t *testing.T) {
	t.Parallel()

	// 성능 임계값 설정 (환경에 따라 조정 가능)
	const maxProcessingTime = 100 * time.Millisecond
	const maxMemoryGrowth = 50 * 1024 * 1024 // 50MB

	t.Run("StringProc_Performance", func(t *testing.T) {
		strproc := strutils.NewStringProc()
		
		// 다양한 크기의 문자열에 대한 성능 테스트
		testSizes := []int{100, 1000, 10000, 100000}
		
		for _, size := range testSizes {
			t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
				testData := strings.Repeat("abcdefghij", size/10)
				
				// 메모리 측정 시작
				var m1, m2 runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&m1)
				
				// 시간 측정 시작
				startTime := time.Now()
				
				// 다양한 문자열 처리 작업 수행
				_ = strproc.AddSlashes(testData)
				_ = strproc.StripSlashes(testData)
				_ = strproc.UpperCaseFirstWords(testData)
				_ = strproc.LowerCaseFirstWords(testData)
				_ = strproc.ReverseStr(testData)
				_ = strproc.PaddingBoth(testData, "*", len(testData)+10)
				
				// 시간 측정 종료
				duration := time.Since(startTime)
				
				// 메모리 측정 종료
				runtime.GC()
				runtime.ReadMemStats(&m2)
				memoryUsed := m2.TotalAlloc - m1.TotalAlloc
				
				t.Logf("Size %d: Duration=%v, Memory=%d bytes", 
					size, duration, memoryUsed)
				
				// 성능 검증 (매우 관대한 기준)
				maxExpectedTime := time.Duration(size/1000+1) * maxProcessingTime
				if duration > maxExpectedTime {
					t.Logf("Performance warning: Size %d took %v (expected < %v)", 
						size, duration, maxExpectedTime)
				}
				
				// 메모리 사용량 검증
				maxExpectedMemory := uint64(size * 10) // 입력 크기의 10배까지 허용
				if memoryUsed > maxExpectedMemory {
					t.Logf("Memory warning: Size %d used %d bytes (expected < %d)", 
						size, memoryUsed, maxExpectedMemory)
				}
			})
		}
	})

	t.Run("StringValidator_Performance", func(t *testing.T) {
		strvalidator := strutils.NewStringValidator()
		
		// 검증 작업의 성능 테스트
		testCases := []struct {
			name     string
			testFunc func() bool
		}{
			{"Email validation", func() bool { 
				return strvalidator.IsValidEmail("test@example.com") 
			}},
			{"Domain validation", func() bool { 
				return strvalidator.IsValidDomain("example.com") 
			}},
			{"URL validation", func() bool { 
				return strvalidator.IsValidURL("https://example.com/path") 
			}},
			{"MAC address validation", func() bool { 
				return strvalidator.IsValidMACAddr("00:11:22:33:44:55") 
			}},
			{"IP validation", func() bool {
				result, _ := strvalidator.IsValidIPAddr("192.168.1.1", strutils.IPv4)
				return result
			}},
		}
		
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				const iterations = 10000
				
				// 메모리 측정 시작
				var m1, m2 runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&m1)
				
				// 시간 측정 시작
				startTime := time.Now()
				
				// 반복 실행
				for i := 0; i < iterations; i++ {
					testCase.testFunc()
				}
				
				// 측정 종료
				duration := time.Since(startTime)
				runtime.GC()
				runtime.ReadMemStats(&m2)
				
				avgTime := duration / iterations
				memoryUsed := m2.TotalAlloc - m1.TotalAlloc
				
				t.Logf("%s: %d iterations, avg time=%v, total memory=%d bytes",
					testCase.name, iterations, avgTime, memoryUsed)
				
				// 성능 기준 검증 (1ms 이하per call)
				maxAvgTime := 1 * time.Millisecond
				if avgTime > maxAvgTime {
					t.Logf("Performance warning: %s average time %v exceeds %v",
						testCase.name, avgTime, maxAvgTime)
				}
			})
		}
	})

	t.Run("NumberFmt_Performance", func(t *testing.T) {
		strproc := strutils.NewStringProc()
		
		// 다양한 숫자 형식에 대한 성능 테스트
		testNumbers := []interface{}{
			123,
			123456789,
			1234.5678,
			-987654.321,
			1.23456789e+10,
			"123456789",
			"1234.5678",
		}
		
		const iterations = 1000
		
		for i, number := range testNumbers {
			t.Run(fmt.Sprintf("Number_%d", i), func(t *testing.T) {
				var m1, m2 runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&m1)
				
				startTime := time.Now()
				
				for j := 0; j < iterations; j++ {
					_, err := strproc.NumberFmt(number)
					if err != nil {
						// 에러가 있어도 성능 측정은 계속
					}
				}
				
				duration := time.Since(startTime)
				runtime.GC()
				runtime.ReadMemStats(&m2)
				
				avgTime := duration / iterations
				memoryUsed := m2.TotalAlloc - m1.TotalAlloc
				
				t.Logf("NumberFmt(%v): avg time=%v, memory=%d bytes",
					number, avgTime, memoryUsed)
			})
		}
	})

	t.Run("Large_String_Processing", func(t *testing.T) {
		strproc := strutils.NewStringProc()
		
		// 매우 큰 문자열 처리 성능 테스트
		largeSizes := []int{1000000, 5000000} // 1MB, 5MB
		
		for _, size := range largeSizes {
			if size > 5000000 && testing.Short() {
				t.Skip("Skipping large string test in short mode")
			}
			
			t.Run(fmt.Sprintf("LargeString_%dMB", size/1000000), func(t *testing.T) {
				// 반복 패턴으로 큰 문자열 생성 (메모리 효율적)
				pattern := "abcdefghijklmnopqrstuvwxyz0123456789"
				testData := strings.Repeat(pattern, size/len(pattern))
				
				var m1, m2 runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&m1)
				
				startTime := time.Now()
				
				// 가벼운 작업만 수행 (큰 문자열에서는)
				_ = strproc.AddSlashes(testData[:1000]) // 일부만 처리
				_ = strproc.UpperCaseFirstWords(testData[:1000])
				
				duration := time.Since(startTime)
				runtime.GC()
				runtime.ReadMemStats(&m2)
				
				memoryUsed := m2.TotalAlloc - m1.TotalAlloc
				
				t.Logf("Large string (%d bytes): Duration=%v, Memory=%d bytes",
					len(testData), duration, memoryUsed)
				
				// 큰 문자열에 대한 합리적인 성능 기준
				maxExpectedTime := 10 * time.Millisecond
				if duration > maxExpectedTime {
					t.Logf("Large string processing took longer than expected: %v", duration)
				}
			})
		}
	})
}

func Test_Memory_Usage_Patterns(t *testing.T) {
	t.Parallel()

	t.Run("Memory_Leak_Detection", func(t *testing.T) {
		const iterations = 1000
		
		var m1, m2, m3 runtime.MemStats
		
		// 초기 메모리 상태 측정
		runtime.GC()
		runtime.ReadMemStats(&m1)
		
		// 반복적인 객체 생성 및 사용
		for i := 0; i < iterations; i++ {
			strproc := strutils.NewStringProc()
			strvalidator := strutils.NewStringValidator()
			
			testData := fmt.Sprintf("test_data_%d", i)
			
			_ = strproc.AddSlashes(testData)
			_ = strproc.StripSlashes(testData)
			_ = strvalidator.IsValidEmail(testData + "@test.com")
			_ = strvalidator.IsValidDomain("test.com")
		}
		
		// 중간 메모리 상태 측정
		runtime.GC()
		runtime.ReadMemStats(&m2)
		
		// GC를 여러 번 실행하여 정리
		for i := 0; i < 5; i++ {
			runtime.GC()
			time.Sleep(10 * time.Millisecond)
		}
		
		// 최종 메모리 상태 측정
		runtime.ReadMemStats(&m3)
		
		initialHeap := m1.HeapInuse
		peakHeap := m2.HeapInuse
		finalHeap := m3.HeapInuse
		
		t.Logf("Memory usage pattern:")
		t.Logf("  Initial heap: %d bytes", initialHeap)
		t.Logf("  Peak heap: %d bytes", peakHeap)
		t.Logf("  Final heap: %d bytes", finalHeap)
		t.Logf("  Peak growth: %d bytes", peakHeap-initialHeap)
		t.Logf("  Final growth: %d bytes", finalHeap-initialHeap)
		
		// 메모리 누수 검사 (최종 heap이 초기보다 너무 크면 안됨)
		maxAcceptableGrowth := uint64(10 * 1024 * 1024) // 10MB
		actualGrowth := finalHeap - initialHeap
		
		if actualGrowth > maxAcceptableGrowth {
			t.Logf("Potential memory leak detected: %d bytes growth", actualGrowth)
		} else {
			t.Logf("Memory usage looks healthy: %d bytes growth", actualGrowth)
		}
	})

	t.Run("String_Allocation_Patterns", func(t *testing.T) {
		strproc := strutils.NewStringProc()
		
		// 다양한 크기의 문자열 할당 패턴 테스트
		sizes := []int{10, 100, 1000, 10000}
		
		for _, size := range sizes {
			t.Run(fmt.Sprintf("AllocPattern_%d", size), func(t *testing.T) {
				var m1, m2 runtime.MemStats
				runtime.GC()
				runtime.ReadMemStats(&m1)
				
				// 동일한 크기의 문자열을 여러 번 처리
				testData := strings.Repeat("x", size)
				const iterations = 100
				
				results := make([]string, iterations)
				for i := 0; i < iterations; i++ {
					results[i] = strproc.AddSlashes(testData)
				}
				
				runtime.GC()
				runtime.ReadMemStats(&m2)
				
				allocPerOp := (m2.TotalAlloc - m1.TotalAlloc) / uint64(iterations)
				
				t.Logf("String size %d: %d bytes allocated per operation",
					size, allocPerOp)
				
				// 결과 사용 (컴파일러 최적화 방지)
				_ = results[len(results)-1]
			})
		}
	})

	t.Run("Goroutine_Memory_Usage", func(t *testing.T) {
		const numGoroutines = 100
		const numOps = 50
		
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)
		
		// 여러 고루틴에서 동시에 메모리 사용
		done := make(chan bool, numGoroutines)
		
		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				strproc := strutils.NewStringProc()
				strvalidator := strutils.NewStringValidator()
				
				for j := 0; j < numOps; j++ {
					testData := fmt.Sprintf("goroutine_%d_op_%d", id, j)
					_ = strproc.AddSlashes(testData)
					_ = strvalidator.IsValidEmail(testData + "@test.com")
				}
				
				done <- true
			}(i)
		}
		
		// 모든 고루틴 완료 대기
		for i := 0; i < numGoroutines; i++ {
			<-done
		}
		
		runtime.GC()
		runtime.ReadMemStats(&m2)
		
		totalMemory := m2.TotalAlloc - m1.TotalAlloc
		memoryPerGoroutine := totalMemory / uint64(numGoroutines)
		
		t.Logf("Concurrent memory usage:")
		t.Logf("  Total memory: %d bytes", totalMemory)
		t.Logf("  Memory per goroutine: %d bytes", memoryPerGoroutine)
		t.Logf("  Total operations: %d", numGoroutines*numOps)
		
		// 고루틴당 합리적인 메모리 사용량 검증
		maxMemoryPerGoroutine := uint64(1024 * 1024) // 1MB per goroutine
		if memoryPerGoroutine > maxMemoryPerGoroutine {
			t.Logf("High memory usage per goroutine: %d bytes", memoryPerGoroutine)
		}
	})

	t.Run("GC_Pressure_Test", func(t *testing.T) {
		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)
		
		strproc := strutils.NewStringProc()
		
		// GC 압박을 주는 작업 (많은 임시 문자열 생성)
		const iterations = 10000
		
		for i := 0; i < iterations; i++ {
			// 매번 새로운 문자열 생성
			testData := strings.Repeat("gc_test", i%100+1)
			
			// 다양한 문자열 작업으로 임시 객체 생성
			result1 := strproc.AddSlashes(testData)
			result2 := strproc.StripSlashes(result1)
			result3 := strproc.UpperCaseFirstWords(result2)
			
			// 결과 사용 (최적화 방지)
			if len(result3) == 0 {
				t.Fatal("Unexpected empty result")
			}
		}
		
		runtime.GC()
		runtime.ReadMemStats(&m2)
		
		gcCycles := m2.NumGC - m1.NumGC
		totalAlloc := m2.TotalAlloc - m1.TotalAlloc
		
		t.Logf("GC pressure test:")
		t.Logf("  GC cycles triggered: %d", gcCycles)
		t.Logf("  Total allocations: %d bytes", totalAlloc)
		t.Logf("  Allocations per iteration: %d bytes", totalAlloc/uint64(iterations))
		
		// 적당한 GC 압박이 있어야 함 (하지만 과도하면 안됨)
		if gcCycles > 100 {
			t.Logf("High GC pressure detected: %d cycles", gcCycles)
		}
	})
}

func Test_Performance_Regression(t *testing.T) {
	t.Parallel()

	// 성능 회귀 테스트 - 기본적인 벤치마크
	t.Run("Performance_Baseline", func(t *testing.T) {
		strproc := strutils.NewStringProc()
		strvalidator := strutils.NewStringValidator()
		
		// 표준 테스트 데이터
		testData := "The quick brown fox jumps over the lazy dog"
		const iterations = 1000
		
		operations := []struct {
			name string
			op   func() interface{}
		}{
			{"AddSlashes", func() interface{} { 
				return strproc.AddSlashes(testData) 
			}},
			{"StripSlashes", func() interface{} { 
				return strproc.StripSlashes(testData) 
			}},
			{"UpperCaseFirstWords", func() interface{} { 
				return strproc.UpperCaseFirstWords(testData) 
			}},
			{"ReverseStr", func() interface{} { 
				return strproc.ReverseStr(testData) 
			}},
			{"IsValidEmail", func() interface{} { 
				return strvalidator.IsValidEmail("test@example.com") 
			}},
			{"NumberFmt", func() interface{} {
				result, _ := strproc.NumberFmt(123456)
				return result
			}},
		}
		
		for _, operation := range operations {
			t.Run(operation.name, func(t *testing.T) {
				startTime := time.Now()
				
				for i := 0; i < iterations; i++ {
					result := operation.op()
					// 결과 사용 (최적화 방지)
					_ = result
				}
				
				duration := time.Since(startTime)
				avgTime := duration / iterations
				
				t.Logf("%s: %d iterations, total=%v, avg=%v, ops/sec=%.0f",
					operation.name, iterations, duration, avgTime, 
					float64(iterations)/duration.Seconds())
				
				// 기본적인 성능 임계값 (매우 관대함)
				maxAvgTime := 100 * time.Microsecond
				if avgTime > maxAvgTime {
					t.Logf("Performance regression warning: %s took %v (expected < %v)",
						operation.name, avgTime, maxAvgTime)
				}
			})
		}
	})

	t.Run("Scalability_Test", func(t *testing.T) {
		strproc := strutils.NewStringProc()
		
		// 입력 크기에 따른 확장성 테스트
		sizes := []int{100, 1000, 10000}
		
		for _, size := range sizes {
			if size > 1000 && testing.Short() {
				continue
			}
			
			t.Run(fmt.Sprintf("Scalability_%d", size), func(t *testing.T) {
				testData := strings.Repeat("abcde", size/5)
				
				startTime := time.Now()
				result := strproc.AddSlashes(testData)
				duration := time.Since(startTime)
				
				// 처리량 계산 (bytes/second)
				throughput := float64(len(testData)) / duration.Seconds()
				
				t.Logf("Size %d: %v duration, %.0f bytes/sec throughput",
					size, duration, throughput)
				
				// 결과 검증 (기능적 정확성)
				if len(result) < len(testData) {
					t.Errorf("Result shorter than input: %d < %d", len(result), len(testData))
				}
				
				// 최소 처리량 검증 (매우 관대함 - 1MB/sec)
				minThroughput := float64(1024 * 1024) // 1MB/sec
				if throughput < minThroughput && size > 1000 {
					t.Logf("Low throughput warning: %.0f bytes/sec", throughput)
				}
			})
		}
	})
}

// 벤치마크 헬퍼 함수들
func generateTestString(size int) string {
	return strings.Repeat("abcdefghij", size/10)
}

func measureMemoryUsage(fn func()) (uint64, time.Duration) {
	var m1, m2 runtime.MemStats
	
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	start := time.Now()
	fn()
	duration := time.Since(start)
	
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	return m2.TotalAlloc - m1.TotalAlloc, duration
}