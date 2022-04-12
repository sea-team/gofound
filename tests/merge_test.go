package tests

import (
	"fmt"
	"gofound/searcher/utils"
	"math/rand"
	"sort"
	"testing"
)

func mergeSort(nums []int) []int {
	if len(nums) < 2 {
		// 分治，两两拆分，一直拆到基础元素才向上递归。
		return nums
	}
	i := len(nums) / 2
	left := mergeSort(nums[0:i])
	// 左侧数据递归拆分
	right := mergeSort(nums[i:])
	// 右侧数据递归拆分
	result := merge(left, right)
	// 排序 & 合并
	return result
}

func merge(left, right []int) []int {
	result := make([]int, 0)
	i, j := 0, 0
	l, r := len(left), len(right)
	for i < l && j < r {
		if left[i] > right[j] {
			result = append(result, right[j])
			j++
		} else {
			result = append(result, left[i])
			i++
		}
	}
	result = append(result, right[j:]...)
	result = append(result, left[i:]...)
	return result
}

//合并排序
func TestMergeSort(t *testing.T) {
	array1 := make([]int, 0)
	array2 := make([]int, 0)

	for i := 0; i < 10000; i++ {
		array1 = append(array1, rand.Intn(1000))
		array2 = append(array2, rand.Intn(1000))
	}

	//fmt.Println(array1, array2)
	t1 := utils.ExecTime(func() {
		mergeSort(array1)
	})
	fmt.Println("归并 time:", t1)

	t2 := utils.ExecTime(func() {
		sort.Sort(sort.IntSlice(array1))
		sort.Sort(sort.IntSlice(array2))
	})
	fmt.Println("快排 time:", t2)

	temp := make([]int, 0)
	temp = append(array1, array2...)

	t3 := utils.ExecTime(func() {
		//sort.Sort(sort.IntSlice(array1))
		//sort.Sort(sort.IntSlice(array2))
		sort.Sort(sort.IntSlice(temp))

	})

	fmt.Println("快排有序 time:", t3)
}
