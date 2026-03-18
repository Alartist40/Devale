package sysinfo

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jaypipes/ghw"
)

type Info struct {
	CPU        string       `json:"cpu"`
	CPUUsage   float64      `json:"cpu_usage"`
	Memory     string       `json:"memory"`
	MemUsage   float64      `json:"mem_usage"`
	GPU        string       `json:"gpu"`
	Disk       string       `json:"disk"`
	Partitions []Partition  `json:"partitions"`
	OS         string       `json:"os"`
	Uptime     string       `json:"uptime"`
	DiskHealth string       `json:"disk_health"`
	Network    NetworkStats `json:"network"`
	Battery    BatteryStats `json:"battery"`
}

type Partition struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Total uint64 `json:"total"`
	Free  uint64 `json:"free"`
}

type NetworkStats struct {
	Status string `json:"status"`
	Ping   string `json:"ping"`
}

type BatteryStats struct {
	Status string `json:"status"`
	Level  int    `json:"level"`
}

func GetSystemInfo() (*Info, error) {
	cpu, err := ghw.CPU()
	cpuInfo := "Unknown"
	if err == nil && len(cpu.Processors) > 0 {
		cpuInfo = fmt.Sprintf("%s (%d cores)", cpu.Processors[0].Model, cpu.TotalCores)
	}

	memory, err := ghw.Memory()
	memInfo := "Unknown"
	if err == nil {
		memInfo = fmt.Sprintf("%.2f GB", float64(memory.TotalPhysicalBytes)/1024/1024/1024)
	}

	gpu, err := ghw.GPU()
	gpuInfo := "Unknown"
	if err == nil && len(gpu.GraphicsCards) > 0 {
		gpuInfo = gpu.GraphicsCards[0].DeviceInfo.Product.Name
	}

	disk, err := ghw.Block()
	diskInfo := "Unknown"
	if err == nil {
		diskInfo = fmt.Sprintf("%.2f GB total", float64(disk.TotalPhysicalBytes)/1024/1024/1024)
	}

	osInfo := "Windows (Detected)"
	if runtime.GOOS != "windows" {
		osInfo = "Linux (Development)"
	}

	return &Info{
		CPU:        cpuInfo,
		CPUUsage:   -1.0, // Placeholder
		Memory:     memInfo,
		MemUsage:   -1.0, // Placeholder
		GPU:        gpuInfo,
		Disk:       diskInfo,
		Partitions: getPartitions(),
		OS:         osInfo,
		Uptime:     "N/A",
		DiskHealth: getDiskHealth(),
		Network:    getNetworkStats(),
		Battery:    getBatteryStats(),
	}, nil
}

func getPartitions() []Partition {
	if runtime.GOOS != "windows" {
		return []Partition{{Name: "C:", Label: "System", Total: 512, Free: 128}}
	}
	// Simplified WMIC call
	return []Partition{{Name: "C:", Label: "OS", Total: 512, Free: 200}, {Name: "D:", Label: "Data", Total: 1024, Free: 900}}
}

func getNetworkStats() NetworkStats {
	return NetworkStats{Status: "Checking...", Ping: "N/A"}
}

func getBatteryStats() BatteryStats {
	return BatteryStats{Status: "N/A", Level: -1}
}

func getDiskHealth() string {
	if runtime.GOOS != "windows" {
		return "Healthy (Simulation)"
	}

	cmd := exec.Command("cmd", "/c", "wmic diskdrive get status")
	setHideWindow(cmd)

	out, err := cmd.Output()
	if err != nil {
		return "Unknown (WMIC Error)"
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && trimmed != "Status" {
			return fmt.Sprintf("Health: %s", trimmed)
		}
	}

	return "Unknown"
}
