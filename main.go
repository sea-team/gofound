package main

import (
	_ "go.uber.org/automaxprocs"

	"github.com/sea-team/gofound/core"
)

func main() {
	//初始化容器和参数解析
	core.Initialize()
}
