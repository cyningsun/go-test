package metrics

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// RequestLatency 记录请求延迟的时间
type RequestLatency struct {
	mu        sync.Mutex
	latencies []time.Duration
	startTime time.Time
	count     int64
}

// NewRequestLatency 创建一个新的请求延迟指标收集器
func NewRequestLatency() *RequestLatency {
	return &RequestLatency{
		latencies: make([]time.Duration, 0, 10000),
		startTime: time.Now(),
	}
}

// RecordLatency 记录一个请求的延迟
func (r *RequestLatency) RecordLatency(d time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.latencies = append(r.latencies, d)
	r.count++
}

// GetQPS 获取每秒请求数
func (r *RequestLatency) GetQPS() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	duration := time.Since(r.startTime).Seconds()
	if duration == 0 {
		return 0
	}
	return float64(r.count) / duration
}

// GetCount 获取请求总数
func (r *RequestLatency) GetCount() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}

// GetLatencyPercentile 获取特定百分位的延迟
func (r *RequestLatency) GetLatencyPercentile(p float64) time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.latencies) == 0 {
		return 0
	}

	// 创建一个副本进行排序，避免修改原始数据
	sortedLatencies := make([]time.Duration, len(r.latencies))
	copy(sortedLatencies, r.latencies)
	sort.Slice(sortedLatencies, func(i, j int) bool {
		return sortedLatencies[i] < sortedLatencies[j]
	})

	// 计算百分位索引
	idx := int(math.Ceil(float64(len(sortedLatencies))*p/100)) - 1
	if idx < 0 {
		idx = 0
	}
	return sortedLatencies[idx]
}

// Summary 返回延迟摘要信息
func (r *RequestLatency) Summary() string {
	return fmt.Sprintf("总请求数: %d, QPS: %.2f, P50: %v, P90: %v, P99: %v",
		r.GetCount(),
		r.GetQPS(),
		r.GetLatencyPercentile(50),
		r.GetLatencyPercentile(90),
		r.GetLatencyPercentile(99),
	)
}

// Reset 重置指标
func (r *RequestLatency) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.latencies = r.latencies[:0]
	r.count = 0
	r.startTime = time.Now()
}
