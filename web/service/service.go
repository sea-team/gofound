package service

import (
	"gofound/global"
	"gofound/searcher/system"
	"gofound/searcher/utils"
	"os"
	"runtime"
)

func Callback() map[string]interface{} {
	return map[string]interface{}{
		"os":             runtime.GOOS,
		"arch":           runtime.GOARCH,
		"cores":          runtime.NumCPU(),
		"version":        runtime.Version(),
		"goroutines":     runtime.NumGoroutine(),
		"dataPath":       global.CONFIG.Data,
		"dictionaryPath": global.CONFIG.Dictionary,
		"gomaxprocs":     runtime.NumCPU() * 2,
		"debug":          global.CONFIG.Debug,
		"shard":          global.CONFIG.Shard,
		"dataSize":       system.GetFloat64MB(utils.DirSizeB(global.CONFIG.Data)),
		"executable":     os.Args[0],
		"dbs":            global.Container.GetDataBaseNumber(),
		//"indexCount":     global.Container.GetIndexCount(),
		//"documentCount":  global.Container.GetDocumentCount(),
		"pid":        os.Getpid(),
		"enableAuth": global.CONFIG.Auth != "",
		"enableGzip": global.CONFIG.EnableGzip,
	}
}
