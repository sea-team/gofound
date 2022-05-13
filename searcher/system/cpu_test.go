package system

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"testing"
)

func TestCPU(t *testing.T) {
	fmt.Println(GetCPUStatus())
	c, _ := cpu.Info()
	fmt.Println(c)
}
