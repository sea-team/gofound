package initialize

import (
	"gofound/global"
	"gofound/searcher"
	"gofound/searcher/words"
)

// Container
func Container(tokenizer *words.Tokenizer) *searcher.Container {
	container := &searcher.Container{
		Dir:       global.CONFIG.Engine.DataDir,
		Debug:     global.CONFIG.System.Debug,
		Tokenizer: tokenizer,
		Shard:     global.CONFIG.Engine.Shard,
	}
	go container.Init()

	return container
}
