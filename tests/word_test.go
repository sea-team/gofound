package tests

import (
	"fmt"
	"github.com/wangbin/jiebago"
	"testing"
)

func TestWord(t *testing.T) {
	var seg jiebago.Segmenter

	seg.LoadDictionary("/Users/panjing/GolandProjects/gofound/data/dictionary.txt")
	r := seg.CutForSearch("想在西安买房投资，哪个区域比较好，最好有具体楼盘？", true)
	words := make([]string, 0)
	for {
		w, ok := <-r
		if !ok {
			break
		}
		words = append(words, w)
	}
	for _, w := range words {
		f := int(seg.SuggestFrequency(w))
		if len([]rune(w)) <= 1 {
			f = 0
		} else {
			f = f % len(words)
		}

		fmt.Printf("%s\t%d\n", w, f)
	}
}
