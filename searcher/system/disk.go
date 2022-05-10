package system

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v3/disk"
)

type DiskStatus struct {
	Total       float64 `json:"total"`
	Used        float64 `json:"used"`
	Free        float64 `json:"free"`
	FsType      string  `json:"fsType"`
	UsedPercent float64 `json:"usedPercent"`
	Path        string  `json:"path"`
}

func (d *DiskStatus) String() string {
	buf, _ := json.Marshal(d)
	return string(buf)
}

func GetDiskStat() DiskStatus {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)

	d := DiskStatus{
		Path:        diskInfo.Path,
		Total:       GetUint64GB(diskInfo.Total),
		Free:        GetUint64GB(diskInfo.Free),
		Used:        GetUint64GB(diskInfo.Used),
		UsedPercent: GetPercent(diskInfo.UsedPercent),
		FsType:      diskInfo.Fstype,
	}
	return d
}
