// metrics/metrics.go
package metrics

import (
	"log"
	"time"

	"github.com/distatus/battery"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/net"
)

type MetricPoint struct {
	TS      int64   `json:"ts"`
	CPU     float64 `json:"cpu"`
	Net     float64 `json:"net"`
	Battery float64 `json:"battery"`
}

// GetMetrics collects CPU %, network usage (scaled), and battery %.
func GetMetrics() MetricPoint {
	// CPU %
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil || len(cpuPercent) == 0 {
		log.Println("Error getting CPU:", err)
	}

	// Network — total bytes sent+received in last second
	netBefore, _ := net.IOCounters(false)
	time.Sleep(time.Second)
	netAfter, _ := net.IOCounters(false)
	var netKB float64
	if len(netBefore) > 0 && len(netAfter) > 0 {
		bytes := (netAfter[0].BytesRecv + netAfter[0].BytesSent) - (netBefore[0].BytesRecv + netBefore[0].BytesSent)
		netKB = float64(bytes) / 1024.0
	}

	// Convert net KB to scale: 10KB = 1, 10MB = 100 (linear)
	netScaled := (netKB / 10.0)
	if netScaled > 100 {
		netScaled = 100
	}

	// Battery %
	batteryPercent := 0.0
	batteries, err := battery.GetAll()
	if err == nil && len(batteries) > 0 {
		batteryPercent = batteries[0].Current / batteries[0].Full * 100
	}

	return MetricPoint{
		TS:      time.Now().UnixMilli(),
		CPU:     cpuPercent[0],
		Net:     netScaled,
		Battery: batteryPercent,
	}
}
