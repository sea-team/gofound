package tests

import (
	"fmt"
	"github.com/wangbin/jiebago"
	"gofound/searcher/utils"
	"testing"
)

func TestWord(t *testing.T) {
	var seg jiebago.Segmenter

	seg.LoadDictionary("/Users/panjing/GolandProjects/gofound/data/dictionary.txt")
	r := seg.CutForSearch("深圳是中国的深圳，上海不是中国的上海", true)
	for {
		w, ok := <-r
		if !ok {
			break
		}
		fmt.Println(utils.ExecTime(func() {

			seg.SuggestFrequency(w)
		}))
		f := seg.SuggestFrequency(w)
		fmt.Printf("%s\t%f\n", w, f)
	}
}
