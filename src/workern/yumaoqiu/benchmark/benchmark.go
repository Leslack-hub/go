package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// BenchmarkResult 基准测试结果
type BenchmarkResult struct {
	Method      string
	TotalTime   time.Duration
	AvgTime     time.Duration
	MinTime     time.Duration
	MaxTime     time.Duration
	RequestNum  int
	SuccessNum  int
	FailNum     int
}

// 测试 URL（使用一个稳定的公共 API）
const testURL = "https://httpbin.org/get"

// 创建高性能 HTTP 客户端
func createBenchHTTPClient() *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  true,
		ForceAttemptHTTP2:   true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		TLSHandshakeTimeout: 5 * time.Second,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
}

// 使用 curl 发送请求
func benchmarkCurl(numRequests int) BenchmarkResult {
	result := BenchmarkResult{
		Method:     "curl",
		RequestNum: numRequests,
		MinTime:    time.Hour, // 初始化为一个很大的值
	}

	start := time.Now()
	for i := 0; i < numRequests; i++ {
		reqStart := time.Now()
		
		cmd := exec.Command("curl", "-s", testURL)
		output, err := cmd.Output()
		
		reqDuration := time.Since(reqStart)
		
		if err != nil || len(output) == 0 {
			result.FailNum++
		} else {
			result.SuccessNum++
		}
		
		if reqDuration < result.MinTime {
			result.MinTime = reqDuration
		}
		if reqDuration > result.MaxTime {
			result.MaxTime = reqDuration
		}
	}
	result.TotalTime = time.Since(start)
	result.AvgTime = result.TotalTime / time.Duration(numRequests)
	
	return result
}

// 使用原生 HTTP 客户端发送请求（无连接复用）
func benchmarkHTTPNoReuse(numRequests int) BenchmarkResult {
	result := BenchmarkResult{
		Method:     "http (no reuse)",
		RequestNum: numRequests,
		MinTime:    time.Hour,
	}

	start := time.Now()
	for i := 0; i < numRequests; i++ {
		reqStart := time.Now()
		
		// 每次创建新客户端，模拟无连接复用
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(testURL)
		
		reqDuration := time.Since(reqStart)
		
		if err != nil {
			result.FailNum++
		} else {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			result.SuccessNum++
		}
		
		if reqDuration < result.MinTime {
			result.MinTime = reqDuration
		}
		if reqDuration > result.MaxTime {
			result.MaxTime = reqDuration
		}
	}
	result.TotalTime = time.Since(start)
	result.AvgTime = result.TotalTime / time.Duration(numRequests)
	
	return result
}

// 使用原生 HTTP 客户端发送请求（启用连接复用）
func benchmarkHTTPWithReuse(numRequests int) BenchmarkResult {
	client := createBenchHTTPClient()
	
	result := BenchmarkResult{
		Method:     "http (with reuse)",
		RequestNum: numRequests,
		MinTime:    time.Hour,
	}

	// 预热连接
	resp, _ := client.Get(testURL)
	if resp != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

	start := time.Now()
	for i := 0; i < numRequests; i++ {
		reqStart := time.Now()
		
		resp, err := client.Get(testURL)
		
		reqDuration := time.Since(reqStart)
		
		if err != nil {
			result.FailNum++
		} else {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			result.SuccessNum++
		}
		
		if reqDuration < result.MinTime {
			result.MinTime = reqDuration
		}
		if reqDuration > result.MaxTime {
			result.MaxTime = reqDuration
		}
	}
	result.TotalTime = time.Since(start)
	result.AvgTime = result.TotalTime / time.Duration(numRequests)
	
	return result
}

// 并发测试 - curl
func benchmarkCurlConcurrent(numRequests int, concurrency int) BenchmarkResult {
	result := BenchmarkResult{
		Method:     fmt.Sprintf("curl (concurrent %d)", concurrency),
		RequestNum: numRequests,
		MinTime:    time.Hour,
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan struct{}, concurrency)

	start := time.Now()
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		semaphore <- struct{}{}
		
		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()
			
			reqStart := time.Now()
			cmd := exec.Command("curl", "-s", testURL)
			output, err := cmd.Output()
			reqDuration := time.Since(reqStart)
			
			mu.Lock()
			if err != nil || len(output) == 0 {
				result.FailNum++
			} else {
				result.SuccessNum++
			}
			if reqDuration < result.MinTime {
				result.MinTime = reqDuration
			}
			if reqDuration > result.MaxTime {
				result.MaxTime = reqDuration
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	
	result.TotalTime = time.Since(start)
	result.AvgTime = result.TotalTime / time.Duration(numRequests)
	
	return result
}

// 并发测试 - HTTP with reuse
func benchmarkHTTPConcurrent(numRequests int, concurrency int) BenchmarkResult {
	client := createBenchHTTPClient()
	
	result := BenchmarkResult{
		Method:     fmt.Sprintf("http reuse (concurrent %d)", concurrency),
		RequestNum: numRequests,
		MinTime:    time.Hour,
	}

	// 预热连接
	resp, _ := client.Get(testURL)
	if resp != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	semaphore := make(chan struct{}, concurrency)

	start := time.Now()
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		semaphore <- struct{}{}
		
		go func() {
			defer wg.Done()
			defer func() { <-semaphore }()
			
			reqStart := time.Now()
			resp, err := client.Get(testURL)
			reqDuration := time.Since(reqStart)
			
			mu.Lock()
			if err != nil {
				result.FailNum++
			} else {
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
				result.SuccessNum++
			}
			if reqDuration < result.MinTime {
				result.MinTime = reqDuration
			}
			if reqDuration > result.MaxTime {
				result.MaxTime = reqDuration
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	
	result.TotalTime = time.Since(start)
	result.AvgTime = result.TotalTime / time.Duration(numRequests)
	
	return result
}

func printResult(r BenchmarkResult) {
	fmt.Printf("\n📊 %s\n", r.Method)
	fmt.Printf("   请求数: %d (成功: %d, 失败: %d)\n", r.RequestNum, r.SuccessNum, r.FailNum)
	fmt.Printf("   总耗时: %v\n", r.TotalTime.Round(time.Millisecond))
	fmt.Printf("   平均耗时: %v\n", r.AvgTime.Round(time.Millisecond))
	fmt.Printf("   最小耗时: %v\n", r.MinTime.Round(time.Millisecond))
	fmt.Printf("   最大耗时: %v\n", r.MaxTime.Round(time.Millisecond))
}

func main() {
	runBenchmark()
}

func runBenchmark() {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("🔥 HTTP 请求性能基准测试")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("测试 URL: %s\n", testURL)
	
	// 串行测试
	numSerial := 5
	fmt.Printf("\n【串行测试】每种方法发送 %d 个请求\n", numSerial)
	
	r1 := benchmarkCurl(numSerial)
	printResult(r1)
	
	r2 := benchmarkHTTPNoReuse(numSerial)
	printResult(r2)
	
	r3 := benchmarkHTTPWithReuse(numSerial)
	printResult(r3)
	
	// 计算提升
	if r1.AvgTime > 0 && r3.AvgTime > 0 {
		speedup := float64(r1.AvgTime) / float64(r3.AvgTime)
		fmt.Printf("\n⚡ HTTP(连接复用) 比 curl 快 %.1fx\n", speedup)
	}
	
	// 并发测试
	numConcurrent := 20
	concurrency := 10
	fmt.Printf("\n【并发测试】%d 个请求，并发数 %d\n", numConcurrent, concurrency)
	
	r4 := benchmarkCurlConcurrent(numConcurrent, concurrency)
	printResult(r4)
	
	r5 := benchmarkHTTPConcurrent(numConcurrent, concurrency)
	printResult(r5)
	
	// 计算并发提升
	if r4.TotalTime > 0 && r5.TotalTime > 0 {
		speedup := float64(r4.TotalTime) / float64(r5.TotalTime)
		fmt.Printf("\n⚡ 并发场景下，HTTP(连接复用) 比 curl 快 %.1fx\n", speedup)
	}
	
	fmt.Println("\n" + strings.Repeat("=", 50))
}
