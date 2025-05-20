package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Stats stores benchmarking statistics
type Stats struct {
	TotalRequests      int64
	SuccessfulRequests int64
	FailedRequests     int64
	ConnectionErrors   int64
	InvalidResponses   int64
	TotalLatency       int64 // in microseconds
	MaxLatency         int64 // in microseconds
	MinLatency         int64 // in microseconds
}

// Config stores benchmarking configuration
type Config struct {
	NumConcurrent int
	Duration      time.Duration
	BaseURL       string
}

var testExpressions = []string{
	"2+2",
	"sin(3.14)",
	"max(2,3,4,5)",
	"1+2*3/4",
	"sin(max(2,3))",
	"2^10",
	"(1+2)*(3+4)",
	"1+2+3+4+5+6+7+8+9+10",
}

func runBenchmark(config Config) *Stats {
	stats := &Stats{
		MinLatency: int64(^uint64(0) >> 1),
	}

	var wg sync.WaitGroup
	start := time.Now()
	deadline := start.Add(config.Duration)

	// Launch worker goroutines
	for i := 0; i < config.NumConcurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{
				Timeout: 5 * time.Second,
			}

			for time.Now().Before(deadline) {
				// add a delay between requests to prevent overwhelming the server
				time.Sleep(5 * time.Millisecond)

				// pick an expression randomly by using the request count as a simple distribution
				expr := testExpressions[atomic.LoadInt64(&stats.TotalRequests)%int64(len(testExpressions))]

				reqStart := time.Now()

				// prepare request
				payload := map[string]string{"expression": expr}
				jsonData, _ := json.Marshal(payload)

				req, err := http.NewRequest("POST", config.BaseURL+"/api/v1/eval", bytes.NewBuffer(jsonData))
				if err != nil {
					fmt.Printf("Error creating request for expr '%s': %v\n", expr, err)
					atomic.AddInt64(&stats.FailedRequests, 1)
					continue
				}

				req.Header.Set("Content-Type", "application/json")

				// make request
				resp, err := client.Do(req)
				if err != nil {
					if strings.Contains(err.Error(), "connection refused") {
						atomic.AddInt64(&stats.ConnectionErrors, 1)
					} else {
						fmt.Printf("Error making request for expr '%s': %v\n", expr, err)
					}
					atomic.AddInt64(&stats.FailedRequests, 1)
					continue
				}

				// read and validate response
				body, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					fmt.Printf("Error reading response for expr '%s': %v\n", expr, err)
					atomic.AddInt64(&stats.FailedRequests, 1)
					continue
				}

				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					fmt.Printf("Invalid JSON response for expr '%s': %v\n", expr, err)
					atomic.AddInt64(&stats.InvalidResponses, 1)
					atomic.AddInt64(&stats.FailedRequests, 1)
					continue
				}

				latency := time.Since(reqStart).Microseconds()

				// update stats atomically
				atomic.AddInt64(&stats.TotalRequests, 1)
				if resp.StatusCode == http.StatusOK {
					atomic.AddInt64(&stats.SuccessfulRequests, 1)
				} else {
					fmt.Printf("Non-OK status code %d for expr '%s': %s\n", resp.StatusCode, expr, body)
					atomic.AddInt64(&stats.FailedRequests, 1)
				}

				atomic.AddInt64(&stats.TotalLatency, latency)

				// update min/max latency
				for {
					currentMax := atomic.LoadInt64(&stats.MaxLatency)
					if latency <= currentMax {
						break
					}
					if atomic.CompareAndSwapInt64(&stats.MaxLatency, currentMax, latency) {
						break
					}
				}

				for {
					currentMin := atomic.LoadInt64(&stats.MinLatency)
					if latency >= currentMin {
						break
					}
					if atomic.CompareAndSwapInt64(&stats.MinLatency, currentMin, latency) {
						break
					}
				}

				resp.Body.Close()
			}
		}()
	}

	wg.Wait()
	return stats
}

func main() {
	numConcurrent := flag.Int("c", 10, "number of concurrent users")
	duration := flag.Duration("d", 20*time.Second, "test duration")
	baseURL := flag.String("url", "http://localhost:8084", "base URL of the parser service")
	flag.Parse()

	config := Config{
		NumConcurrent: *numConcurrent,
		Duration:      *duration,
		BaseURL:       *baseURL,
	}

	fmt.Printf("running benchmark for %v\n", config.Duration)

	stats := runBenchmark(config)

	fmt.Println("\n📊 Benchmark Results")

	// Create tables for results
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Metric", "Value"})

	// Add request statistics
	table.Bulk([][]string{
		{"Total Requests", fmt.Sprintf("%d", stats.TotalRequests)},
		{"Successful Requests", fmt.Sprintf("%d", stats.SuccessfulRequests)},
		{"Failed Requests", fmt.Sprintf("%d", stats.FailedRequests)},
		{"Connection Errors", fmt.Sprintf("%d", stats.ConnectionErrors)},
		{"Invalid Responses", fmt.Sprintf("%d", stats.InvalidResponses)},
	})

	if stats.TotalRequests > 0 {
		avgLatency := float64(stats.TotalLatency) / float64(stats.TotalRequests)
		rps := float64(stats.TotalRequests) / config.Duration.Seconds()
		successRate := float64(stats.SuccessfulRequests) / float64(stats.TotalRequests) * 100

		// Add performance metrics
		table.Bulk([][]string{
			{"Average Latency", fmt.Sprintf("%.2f ms", avgLatency/1000)},
			{"Min Latency", fmt.Sprintf("%.2f ms", float64(stats.MinLatency)/1000)},
			{"Max Latency", fmt.Sprintf("%.2f ms", float64(stats.MaxLatency)/1000)},
			{"Requests/Second", fmt.Sprintf("%.2f", rps)},
			{"Success Rate", fmt.Sprintf("%.2f%%", successRate)},
		})
	}

	table.Render()
}
