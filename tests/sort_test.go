package tests

import (
	"fmt"
	"gofound/searcher/utils"
	"math/rand"
	"sort"
	"testing"
)

//排序测试

//冒泡排序测试
func TestSort(t *testing.T) {
	//测试数据
	var data []int

	for i := 0; i < 100000; i++ {
		//随机数
		data = append(data, rand.Intn(100))
	}

	//fmt.Println("原始：", data)

	//排序
	data1 := data
	var data2 = make([]int, len(data))
	copy(data2, data)

	var data3 = make([]int, len(data))
	copy(data3, data)

	var data4 = make([]int, len(data))
	copy(data4, data)

	_time := utils.ExecTime(func() {
		BubbleSort(data1)
	})
	fmt.Println("冒泡排序耗时：", _time)
	//fmt.Println("data1", data1)
	//快速排序

	_time = utils.ExecTime(func() {
		//QuickSortAsc(data2, 0, len(data2)-1)
		utils.QuickSortAsc(data2, 0, len(data2)-1, func(i int, j int) {
			//log.Println(i, j)
		})
	})

	fmt.Println("快速排序耗时：", _time)
	//fmt.Println("data2", data2)

	_time = utils.ExecTime(func() {
		SelectSort(data3)
	})
	fmt.Println("选择排序耗时：", _time)
	//fmt.Println("data3", data3)

	_time = utils.ExecTime(func() {
		InsertSort(data4)
	})
	fmt.Println("插入排序耗时：", _time)

}

//冒泡排序
func BubbleSort(data []int) {
	//排序
	for i := 0; i < len(data)-1; i++ {
		for j := 0; j < len(data)-1-i; j++ {
			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}
}

func SelectSort(arr []int) {
	for j := 0; j < len(arr)-1; j++ {
		max := arr[j]
		maxIndex := j
		for i := j + 1; i < len(arr); i++ {
			if max < arr[i] {
				//记录
				max = arr[i]
				maxIndex = i
			}
		}
		//交换
		if maxIndex != j {
			arr[j], arr[maxIndex] = arr[maxIndex], arr[j]
		}
		//fmt.Printf("数据第 %v 次交换后为:\t%v\n", j+1, arr)
	}
}

//快速排序
func QuickSort(arr []int, start, end int) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		for i <= j {
			for arr[i] < key {
				i++
			}
			for arr[j] > key {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}

		if start < j {
			QuickSort(arr, start, j)
		}
		if end > i {
			QuickSort(arr, i, end)
		}
	}
}

func InsertSort(list []int) {
	n := len(list)
	// 进行 N-1 轮迭代
	for i := 1; i <= n-1; i++ {
		deal := list[i] // 待排序的数
		j := i - 1      // 待排序的数左边的第一个数的位置

		// 如果第一次比较，比左边的已排好序的第一个数小，那么进入处理
		if deal < list[j] {
			// 一直往左边找，比待排序大的数都往后挪，腾空位给待排序插入
			for ; j >= 0 && deal < list[j]; j-- {
				list[j+1] = list[j] // 某数后移，给待排序留空位
			}
			list[j+1] = deal // 结束了，待排序的数插入空位
		}
	}
}

func TestFastSort(t *testing.T) {

	//QuickSortDesc
	//测试数据
	var data []int

	for i := 0; i < 1000; i++ {
		//随机数
		data = append(data, rand.Intn(100))

	}

	_time := utils.ExecTime(func() {
		//utils.QuickSortDesc(data, 0, len(data)-1, func(i int, j int) {

		//})
		//sort.Ints(data)
		sort.Sort(sort.Reverse(sort.IntSlice(data)))
		//sort.Reverse(data)
	})
	fmt.Println("时间", _time)
	fmt.Println(data)

}

//获取数组最大值
func getMaxInArr(arr []int) int {
	max := arr[0]
	for i := 1; i < len(arr); i++ {
		if arr[i] > max {
			max = arr[i]
		}
	}
	return max
}
func sortInBucket(bucket []int) { //此处实现插入排序方式，其实可以用任意其他排序方式
	length := len(bucket)
	if length == 1 {
		return
	}
	for i := 1; i < length; i++ {
		backup := bucket[i]
		j := i - 1
		//将选出的被排数比较后插入左边有序区
		for j >= 0 && backup < bucket[j] { //注意j >= 0必须在前边，否则会数组越界
			bucket[j+1] = bucket[j] //移动有序数组
			j--                     //反向移动下标
		}
		bucket[j+1] = backup //插队插入移动后的空位
	}
}

//桶排序
func BucketSort(arr []int) []int {
	//桶数
	num := len(arr)
	//k（数组最大值）
	max := getMaxInArr(arr)
	//二维切片
	buckets := make([][]int, num)
	//分配入桶
	index := 0
	for i := 0; i < num; i++ {
		index = arr[i] * (num - 1) / max //分配桶index = value * (n-1) /k
		buckets[index] = append(buckets[index], arr[i])
	}
	//桶内排序
	tmpPos := 0
	for i := 0; i < num; i++ {
		bucketLen := len(buckets[i])
		if bucketLen > 0 {
			sortInBucket(buckets[i])
			copy(arr[tmpPos:], buckets[i])
			tmpPos += bucketLen
		}
	}
	return arr
}

func TestFind(t *testing.T) {

	data := make([]int, 0)
	data2 := make([]int, 0)
	for i := 0; i < 100000; i++ {
		val := rand.Intn(100000)
		data = append(data, val)
		data2 = append(data2, val)
	}

	t1 := utils.ExecTime(func() {
		sort.Sort(sort.IntSlice(data))
	})
	fmt.Println("快排用时", t1)

	//fmt.Println(find(data, 1))
	t2 := utils.ExecTime(func() {
		BucketSort(data2)
		for i, j := 0, len(data2)-1; i < j; i, j = i+1, j-1 {
			data2[i], data2[j] = data2[j], data2[i]
		}
	})
	fmt.Println("捅排", t2)
	//fmt.Println("捅排", sort.Reverse(sort.IntSlice(data2)))

	//查找优化，桶排序+map去重

}
func find(data []uint32, target uint32) (bool, int) {
	low := 0
	high := len(data) - 1
	for low <= high {
		mid := (low + high) / 2
		if data[mid] == target {
			return true, mid
		} else if data[mid] < target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return false, -1
}
func TestMerge(t *testing.T) {

	data1 := make([]uint32, 0)
	data2 := make([]uint32, 0)
	for i := 0; i < 10000; i++ {
		v := rand.Intn(10)
		data1 = append(data1, uint32(v))
		data2 = append(data2, uint32(v))
	}

	t1 := utils.ExecTime(func() {
		temp := make([]uint32, 0)
		for _, v := range data1 {
			if found, _ := find(temp, v); found {
				temp = append(temp, v)
			}
		}
		fmt.Println(temp)
	})

	fmt.Println("二分法去重", t1)

	t2 := utils.ExecTime(func() {
		temp := make(map[uint32]bool, len(data2))
		d := make([]uint32, 0)
		for _, val := range data2 {
			if _, ok := temp[val]; !ok {
				temp[val] = true
				d = append(d, val)
			}
		}
		fmt.Println(d)
	})
	fmt.Println("map去重", t2)
}
