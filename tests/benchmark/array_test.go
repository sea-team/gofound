package benchmark

import (
	"github.com/emirpasic/gods/sets/hashset"
	"gofound/searcher/arrays"
	"math/rand"
	"sort"
	"testing"
)

const dir = "abcdefghijklmnopqrstuvwxyzABCDEFGXIJKLMNOPQRSTUVWXYZ1234567890"

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

/*
*
string array length = 10000, single string length = 15
BenchmarkArrayStringExists
BenchmarkArrayStringExists-6    	1000000000	         0.0000568 ns/op
BenchmarkArrayStringExists1
BenchmarkArrayStringExists1-6   	1000000000	         0.002870 ns/op

string array length = 1000000, single string length = 15
BenchmarkArrayStringExists
BenchmarkArrayStringExists-6    	1000000000	         0.0007200 ns/op
BenchmarkArrayStringExists1
BenchmarkArrayStringExists1-6   	       1	1029444051 ns/op
*/
func BenchmarkArrayStringExists(b *testing.B) {
	path := hashset.New()
	s := make([]string, 0)
	for i := 0; i < 1000000; i++ {
		p := ""
		for j := 0; j < 15; j++ {
			p += string(dir[rand.Intn(len(dir))])
		}
		for path.Contains(p) {
			for j := 0; j < 15; j++ {
				p += string(dir[rand.Intn(len(dir))])
			}
		}
		path.Add(p)
		s = append(s, p)
	}
	b.ResetTimer()
	arrays.ArrayStringExists(s, s[rand.Intn(len(s))])
}

func BenchmarkArrayStringExists1(b *testing.B) {
	path := hashset.New()
	s := make([]string, 0)
	for i := 0; i < 1000000; i++ {
		p := ""
		for j := 0; j < 15; j++ {
			p += string(dir[rand.Intn(len(dir))])
		}
		for path.Contains(p) {
			for j := 0; j < 15; j++ {
				p += string(dir[rand.Intn(len(dir))])
			}
		}
		path.Add(p)
		s = append(s, p)
	}
	b.ResetTimer()
	arrayStringExists2(s, s[rand.Intn(len(s))])
}

func arrayStringExists2(arr []string, str string) bool {
	sort.Strings(arr)
	index := sort.SearchStrings(arr, str)
	if index < len(arr) && arr[index] == str {
		return true
	}
	return false
}
