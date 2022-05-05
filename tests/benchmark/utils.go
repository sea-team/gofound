package benchmark

import "math/rand"

func GetRandomUint32(n int) []uint32 {
	var array = make([]uint32, n)
	for i := 0; i < n; i++ {
		array[i] = rand.Uint32()
	}
	return array
}
