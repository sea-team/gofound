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
		"dataPath":       global.CONFIG.Engine.DataDir,
		"dictionaryPath": global.CONFIG.Engine.DictionaryDir,
		"gomaxprocs":     runtime.NumCPU() * 2,
		"debug":          global.CONFIG.System.Debug,
		"shard":          global.CONFIG.Engine.Shard,
		"dataSize":       system.GetFloat64MB(utils.DirSizeB(global.CONFIG.Engine.DataDir)),
		"executable":     os.Args[0],
		"dbs":            global.Container.GetDataBaseNumber(),
		"indexCount":     global.Container.GetIndexCount(),
		"documentCount":  global.Container.GetDocumentCount(),
		"pid":            os.Getpid(),
		"enableAuth":     global.CONFIG.Auth.Enable,
		"enableGzip":     global.CONFIG.System.EnableGzip,
	}
}
