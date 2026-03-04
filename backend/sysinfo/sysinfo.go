package sysinfo

import (
	"fmt"
	"github.com/jaypipes/ghw"
)

type Info struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
	GPU    string `json:"gpu"`
	Disk   string `json:"disk"`
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

	return &Info{
		CPU:    cpuInfo,
		Memory: memInfo,
		GPU:    gpuInfo,
		Disk:   diskInfo,
	}, nil
}
