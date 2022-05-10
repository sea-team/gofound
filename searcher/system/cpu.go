package system

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"runtime"
	"time"
)

type CPUStatus struct {
	Cores       int     `json:"cores"`
	UsedPercent float64 `json:"usedPercent"`
	ModelName   string  `json:"modelName"`
}

func GetCPUStatus() CPUStatus {
	percent, _ := cpu.Percent(time.Second, false)
	info, _ := cpu.Info()
	c := CPUStatus{
		UsedPercent: GetPercent(percent[0]),
		Cores:       runtime.NumCPU(),
		ModelName:   info[0].ModelName,
	}

	return c
}
