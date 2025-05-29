package utils

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"slices"
	"time"
)

type Result struct {
	Operation  string
	AvgLatency float64
	P95Latency float64
	Throughput float64
	AvgRAM     float64
	GCPause    float64
}

type Operation interface {
	Name() string
	Execute(int) error
}

func Run(op Operation, iterations int) Result {
	times := make([]time.Duration, iterations)
	var memStart, memEnd runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&memStart)
	gcPauseStart := memStart.PauseTotalNs

	for i := 0; i < iterations; i++ {
		start := time.Now()

		if err := op.Execute(i); err != nil {
			log.Fatalf("failed to execute operation \"%s\": %v", op.Name(), err)
		}

		times[i] = time.Since(start)
	}

	runtime.ReadMemStats(&memEnd)
	gcPauseEnd := memEnd.PauseTotalNs

	var totalTime time.Duration
	for _, t := range times {
		totalTime += t
	}

	avgLatency := float64(totalTime.Nanoseconds()) / float64(iterations) / 1e6

	slices.Sort(times)
	idx := max(int(math.Ceil(0.95*float64(iterations)))-1, 0)
	p95Latency := float64(times[idx].Nanoseconds()) / 1e6

	throughput := float64(iterations) / totalTime.Seconds()

	avgRAM := float64(memEnd.TotalAlloc-memStart.TotalAlloc) / float64(iterations) / (1024 * 1024)

	gcPause := float64(gcPauseEnd-gcPauseStart) / 1e6

	return Result{
		Operation:  op.Name(),
		AvgLatency: avgLatency,
		P95Latency: p95Latency,
		Throughput: throughput,
		AvgRAM:     avgRAM,
		GCPause:    gcPause,
	}
}

func PrintResult(operations []Operation, iterations int) {
	fmt.Printf(
		"%-60s %-20s %-20s %-20s %-20s %-20s\n",
		"Operation", "Avg Latency (ms)", "P95 Latency (ms)", "Throughput (ops/s)", "Avg RAM (MB)", "GC Pause (ms)",
	)
	for _, operation := range operations {
		result := Run(operation, iterations)
		fmt.Printf(
			"%-60s %-20.4f %-20.4f %-20.4f %-20.4f %-20.4f\n",
			result.Operation,
			result.AvgLatency,
			result.P95Latency,
			result.Throughput,
			result.AvgRAM,
			result.GCPause,
		)
	}
}
