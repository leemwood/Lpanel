package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"net/http"
	"time"
)

type SystemStatus struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemTotal    uint64  `json:"mem_total"`
	MemUsed     uint64  `json:"mem_used"`
	MemPercent  float64 `json:"mem_percent"`
	DiskTotal   uint64  `json:"disk_total"`
	DiskUsed    uint64  `json:"disk_used"`
	DiskPercent float64 `json:"disk_percent"`
	Uptime      uint64  `json:"uptime"`
	Platform    string  `json:"platform"`
}

func GetSystemStatus(c *gin.Context) {
	v, _ := mem.VirtualMemory()
	cUsage, _ := cpu.Percent(time.Second, false)
	d, _ := disk.Usage("/")
	h, _ := host.Info()

	cpuPct := 0.0
	if len(cUsage) > 0 {
		cpuPct = cUsage[0]
	}

	status := SystemStatus{
		CPUUsage:    cpuPct,
		MemTotal:    v.Total,
		MemUsed:     v.Used,
		MemPercent:  v.UsedPercent,
		DiskTotal:   d.Total,
		DiskUsed:    d.Used,
		DiskPercent: d.UsedPercent,
		Uptime:      h.Uptime,
		Platform:    h.Platform,
	}

	c.JSON(http.StatusOK, status)
}
