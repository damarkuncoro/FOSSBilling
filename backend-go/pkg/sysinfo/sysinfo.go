package sysinfo

import (
	"fmt"
	"runtime"
	"time"
)

type RuntimeStats struct {
	GoVersion, MemoryAlloc, MemoryTotal, MemorySys, Uptime string; Goroutines, NumCPU int; UptimeSeconds int64
}

var start = time.Now()

func GetRuntimeStats() RuntimeStats {
	var m runtime.MemStats; runtime.ReadMemStats(&m); up := time.Since(start)
	f := func(v uint64) string { return fmt.Sprintf("%.2f MB", float64(v)/1024/1024) }
	return RuntimeStats{GoVersion: runtime.Version(), Goroutines: runtime.NumGoroutine(), MemoryAlloc: f(m.Alloc), MemoryTotal: f(m.TotalAlloc), MemorySys: f(m.Sys), NumCPU: runtime.NumCPU(), Uptime: up.Round(time.Second).String(), UptimeSeconds: int64(up.Seconds())}
}
