package tests

import (
	"fmt"
	"github.com/wangbin/jiebago"
	"strings"
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
func contains(s *[]string, e string, skipIndex int) bool {
	for index, a := range *s {
		if index != skipIndex && strings.Contains(a, e) {
			return true
		}
	}
	return false
}
func getLongWords(words *[]string) []string {

	var newWords = make([]string, 0)
	for index, w := range *words {
		if !contains(words, w, index) {
			newWords = append(newWords, w)
		}
	}
	return newWords
}

func TestLongWord(t *testing.T) {
	words := []string{"博物", "博物馆", "深圳北", "深圳", "深圳东"}
	r := getLongWords(&words)
	fmt.Println(r)
}

func BenchmarkTest(b *testing.B) {
	var r []string
	for i := 0; i < b.N; i++ {
		words := []string{"博物", "博物馆", "深圳北", "深圳", "深圳东"}
		r = getLongWords(&words)
	}
	fmt.Println(r)
}
