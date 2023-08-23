package service

import (
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/sea-team/gofound/global"
	"github.com/sea-team/gofound/searcher/system"
	"github.com/sea-team/gofound/searcher/utils"
)

func Callback() map[string]interface{} {
	var threadProfile = pprof.Lookup("threadcreate")
	return map[string]interface{}{
		"os":             runtime.GOOS,
		"arch":           runtime.GOARCH,
		"cores":          runtime.NumCPU(),
		"version":        runtime.Version(),
		"goroutines":     runtime.NumGoroutine(),
		"threads":        threadProfile.Count(),
		"dataPath":       global.CONFIG.Data,
		"dictionaryPath": global.CONFIG.Dictionary,
		"gomaxprocs":     runtime.NumCPU() * 2,
		"debug":          global.CONFIG.Debug,
		"shard":          global.CONFIG.Shard,
		"dataSize":       system.GetFloat64MB(utils.DirSizeB(global.CONFIG.Data)),
		"executable":     os.Args[0],
		"dbs":            global.Container.GetDataBaseNumber(),
		//"indexCount":     global.container.GetIndexCount(),
		//"documentCount":  global.container.GetDocumentCount(),
		"pid":        os.Getpid(),
		"enableAuth": global.CONFIG.Auth != "",
		"enableGzip": global.CONFIG.EnableGzip,
		"bufferNum":  global.CONFIG.BufferNum,
	}
}
