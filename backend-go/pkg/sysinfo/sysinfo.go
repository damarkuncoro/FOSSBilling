package sysinfo

import (
	"fmt"
	"runtime"
	"time"
)

type RuntimeStats struct {
	GoVersion      string `json:"go_version"`
	Goroutines     int    `json:"goroutines"`
	MemoryAlloc    string `json:"memory_alloc"`
	MemoryTotal    string `json:"memory_total"`
	MemorySys      string `json:"memory_sys"`
	NumCPU         int    `json:"num_cpu"`
	Uptime         string `json:"uptime"`
	UptimeSeconds  int64  `json:"uptime_seconds"`
}

var startTime = time.Now()

func GetRuntimeStats() RuntimeStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	uptime := time.Since(startTime)

	return RuntimeStats{
		GoVersion:      runtime.Version(),
		Goroutines:     runtime.NumGoroutine(),
		MemoryAlloc:    fmt.Sprintf("%.2f MB", float64(m.Alloc)/1024/1024),
		MemoryTotal:    fmt.Sprintf("%.2f MB", float64(m.TotalAlloc)/1024/1024),
		MemorySys:      fmt.Sprintf("%.2f MB", float64(m.Sys)/1024/1024),
		NumCPU:         runtime.NumCPU(),
		Uptime:         fmt.Sprintf("%s", uptime.Round(time.Second)),
		UptimeSeconds:  int64(uptime.Seconds()),
	}
}
