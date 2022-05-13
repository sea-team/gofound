package system

import (
	"fmt"
	"testing"
)

func TestMem(t *testing.T) {

	m := GetMemStat()
	fmt.Println(m)
}
