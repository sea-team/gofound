package initialize

import "gofound/searcher/words"

// Tokenizer
func Tokenizer(dictionaryPath string) *words.Tokenizer {
	return words.NewTokenizer(dictionaryPath)
}
