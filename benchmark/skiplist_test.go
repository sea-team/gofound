package benchmark

import (
	"gofound/searcher/arrays"
	"math/rand"
	"testing"
)
import "github.com/ryszard/goskiplist/skiplist"

func BenchmarkSkipList(b *testing.B) {

	//产生1万个随机数
	var nums []int
	for i := 0; i < 10000; i++ {
		num := rand.Intn(100000)
		nums = append(nums, num)
	}

	b.ResetTimer()

	b.Run("skip", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sl := skiplist.NewIntSet()

			for _, num := range nums {
				if !sl.Contains(num) {
					sl.Add(num)
				}
			}
		}
	})

	b.Run("binary", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			temps := make([]uint32, len(nums))
			for index, num := range nums {
				if !arrays.BinarySearch(temps, uint32(num)) {
					temps[index] = uint32(num)
				}
			}
		}
	})

}
