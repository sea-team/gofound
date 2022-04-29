package tests

import (
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"gofound/searcher/arrays"
	"gofound/searcher/utils"
	"math/rand"
	"testing"
	"unsafe"
)

func TestTime(t *testing.T) {

	//产生1万个随机数
	var nums []uint32
	for i := 0; i < 10000; i++ {
		num := rand.Uint32()
		nums = append(nums, num)
	}

	tt := utils.ExecTime(func() {
		for i := 0; i < 1000; i++ {
			bitmap := roaring.BitmapOf(nums...)
			if i == 0 {
				fmt.Println("bitmap内存占用", unsafe.Sizeof(bitmap))
			}
		}

	})
	fmt.Println("bitmap", tt)

	tt = utils.ExecTime(func() {

		for i := 0; i < 1000; i++ {
			temps := make([]uint32, len(nums))
			for index, num := range nums {
				if !arrays.BinarySearch(temps, num) {
					temps[index] = num
				}
			}

			if i == 0 {
				fmt.Println("array内存占用", unsafe.Sizeof(temps))
			}
		}
	})
	fmt.Println("binary", tt)

}
