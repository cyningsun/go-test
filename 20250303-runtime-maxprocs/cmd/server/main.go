package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/cyningsun/go-test/20250303-runtime-maxprocs/pkg/compute"
	"github.com/cyningsun/go-test/20250303-runtime-maxprocs/pkg/metrics"
	"github.com/cyningsun/go-test/20250303-runtime-maxprocs/pkg/proc"
)

var (
	port            = flag.Int("port", 8080, "HTTP服务器端口")
	autoAdjustProcs = flag.Bool("auto-adjust", false, "是否自动调整GOMAXPROCS")
	requestLatency  = metrics.NewRequestLatency()
)

type ServerStats struct {
	GOMAXPROCS    int     `json:"gomaxprocs"`
	NumCPU        int     `json:"num_cpu"`
	AutoAdjusted  bool    `json:"auto_adjusted"`
	Uptime        string  `json:"uptime"`
	TotalRequests int64   `json:"total_requests"`
	QPS           float64 `json:"qps"`
	P50Latency    string  `json:"p50_latency"`
	P90Latency    string  `json:"p90_latency"`
	P99Latency    string  `json:"p99_latency"`
	NrThrottled   int     `json:"nr_throttled"`
	ThrottledTime int64   `json:"throttled_time_ns"`
}

func main() {
	flag.Parse()

	// 设置GOMAXPROCS
	maxProcs := proc.SetMaxProcs(*autoAdjustProcs)
	log.Printf("服务启动，GOMAXPROCS=%d, 自动调整=%v, 系统CPU数=%d\n",
		maxProcs, *autoAdjustProcs, runtime.NumCPU())

	// 记录启动时间以计算运行时长
	startTime := time.Now()

	// 计算密集型任务处理接口
	http.HandleFunc("/fibonacci", fibonacciHandler)

	// 矩阵乘法计算接口
	http.HandleFunc("/matrix", matrixHandler)

	// 健康检查和状态接口
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		statsHandler(w, r, startTime)
	})

	// 请求统计重置接口
	http.HandleFunc("/reset", resetHandler)

	// 创建服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: http.DefaultServeMux,
	}

	// 启动HTTP服务器
	go func() {
		log.Printf("HTTP服务器启动，监听端口: %d\n", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP服务器启动失败: %v\n", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 输出统计信息
	log.Println("服务退出前统计信息:")
	log.Println(requestLatency.Summary())

	// 输出CPU限流统计
	throttled, throttledTime, err := proc.GetCPUThrottlingStats()
	if err == nil {
		log.Printf("CPU限流统计: 限流次数=%d, 限流时间=%d纳秒\n",
			throttled, throttledTime)
	} else {
		log.Printf("无法获取CPU限流统计: %v\n", err)
	}
}

func fibonacciHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		requestLatency.RecordLatency(time.Since(start))
	}()

	// 获取请求参数，默认计算第30个斐波那契数
	n := 30
	if nParam := r.URL.Query().Get("n"); nParam != "" {
		if val, err := strconv.Atoi(nParam); err == nil && val > 0 {
			n = val
		}
	}

	// 计算斐波那契数
	result := compute.FibonacciMatrix(n)

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"n":      n,
		"result": result,
	})
}

func matrixHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		requestLatency.RecordLatency(time.Since(start))
	}()

	// 默认矩阵大小
	size := 50
	if sizeParam := r.URL.Query().Get("size"); sizeParam != "" {
		if val, err := strconv.Atoi(sizeParam); err == nil && val > 0 {
			size = val
		}
	}

	// 创建两个随机矩阵
	matrixA := make([][]int, size)
	matrixB := make([][]int, size)
	for i := 0; i < size; i++ {
		matrixA[i] = make([]int, size)
		matrixB[i] = make([]int, size)
		for j := 0; j < size; j++ {
			matrixA[i][j] = i + j
			matrixB[i][j] = i * j
		}
	}

	// 执行矩阵乘法
	result := compute.MatrixMultiply(matrixA, matrixB)

	// 返回结果摘要
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"size":           size,
		"sum_of_results": sumMatrix(result),
	})
}

func sumMatrix(matrix [][]int) int {
	sum := 0
	for i := range matrix {
		for j := range matrix[i] {
			sum += matrix[i][j]
		}
	}
	return sum
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	requestLatency.Reset()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("statistics reset"))
}

func statsHandler(w http.ResponseWriter, r *http.Request, startTime time.Time) {
	throttled, throttledTime, _ := proc.GetCPUThrottlingStats()

	stats := ServerStats{
		GOMAXPROCS:    runtime.GOMAXPROCS(0),
		NumCPU:        runtime.NumCPU(),
		AutoAdjusted:  *autoAdjustProcs,
		Uptime:        time.Since(startTime).String(),
		TotalRequests: requestLatency.GetCount(),
		QPS:           requestLatency.GetQPS(),
		P50Latency:    requestLatency.GetLatencyPercentile(50).String(),
		P90Latency:    requestLatency.GetLatencyPercentile(90).String(),
		P99Latency:    requestLatency.GetLatencyPercentile(99).String(),
		NrThrottled:   throttled,
		ThrottledTime: throttledTime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
