package tests

import (
	"fmt"
	"testing"
	"time"
)

func TestExecTime(t *testing.T) {
	startT := time.Now()
	time.Sleep(time.Millisecond * 10)

	tc := time.Since(startT)
	fmt.Println(tc)
}
