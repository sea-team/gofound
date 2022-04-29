package benchmark

import (
	"gofound/searcher/arrays"
	"testing"
)

func Benchmark(b *testing.B) {

	//测试两种方法的性能
	size := 100
	arrayList := make([][]uint32, size)
	for i := 0; i < size; i++ {
		arrayList[i] = GetRandomUint32(1000)
	}

	b.Run("array", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var temp []uint32
			for _, nums := range arrayList {

				for _, num := range nums {
					if !arrays.BinarySearch(temp, num) {
						temp = append(temp, num)
					}
				}
			}
		}
	})

	b.Run("sort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var temp []uint32
			for _, v := range arrayList {
				temp = append(temp, v...)
			}
			//去重
			var as []uint32
			for _, v := range temp {
				if !arrays.BinarySearch(as, v) {
					as = append(as, v)
				}
			}
		}
	})
}
