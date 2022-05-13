package system

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v3/mem"
	"runtime"
)

type MemStatus struct {
	Total       float64 `json:"total"`
	Used        float64 `json:"used"`
	Free        float64 `json:"free"`
	Self        float64 `json:"self"`
	UsedPercent float64 `json:"usedPercent"`
}

func (m *MemStatus) String() string {
	buf, _ := json.Marshal(m)
	return string(buf)
}

func GetMemStat() MemStatus {

	//内存信息
	info, _ := mem.VirtualMemory()
	m := MemStatus{
		Total:       GetUint64GB(info.Total),
		Used:        GetUint64GB(info.Used),
		Free:        GetUint64GB(info.Free),
		UsedPercent: GetPercent(info.UsedPercent),
	}

	//自身占用
	memStat := new(runtime.MemStats)
	runtime.ReadMemStats(memStat)
	m.Self = GetUint64GB(memStat.Alloc)

	return m
}
