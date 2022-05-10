package system

import (
	"fmt"
	"strconv"
)

func GetFloat64MB(size int64) float64 {
	val, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(size)/1024/1024), 64)
	return val
}
func GetUint64GB(size uint64) float64 {
	val, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(size)/1024/1024/1024), 64)
	return val
}

func GetPercent(val float64) float64 {
	v, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", val), 64)
	return v
}
