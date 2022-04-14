package arrays

const (
	LOW  = 0
	HIGH = 1
)

// BinarySearch 二分查找
func BinarySearch(arr []uint32, target uint32) bool {
	low := 0
	high := len(arr) - 1
	for low <= high {
		mid := (low + high) / 2
		if arr[mid] == target {
			return true
		} else if arr[mid] > target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return false
}

func BinarySearchIndex(arr []uint32, target uint32) int {
	low := 0
	high := len(arr) - 1
	for low <= high {
		mid := (low + high) / 2
		if arr[mid] == target {
			return mid
		} else if arr[mid] > target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}

func Find(arr []uint32, target uint32) int {
	for index, v := range arr {
		if v == target {
			return index
		}
	}
	return -1
}

func Exists(arr []uint32, target uint32) bool {
	return Find(arr, target) != -1
}

// BubbleSortUint32 冒泡排序
func BubbleSortUint32(array []uint32, c int) []uint32 {
	for i := 0; i < len(array); i++ {

		//for j := 0; j < len(array)-i-1; j++ {
		for j := 0; condition(j, len(array)-i-1, c); j++ {
			if array[j] > array[j+1] {
				array[j], array[j+1] = array[j+1], array[j]
			}
		}
	}
	return array
}

func BubbleSortInt64Low(array []uint32) []uint32 {
	return BubbleSortUint32(array, LOW)
}

func BubbleSortInt64High(array []uint32) []uint32 {
	return BubbleSortUint32(array, HIGH)
}

func condition(value1, value2, c int) bool {
	if c == LOW {
		return value1 < value2
	} else {
		return value1 > value2
	}

}

//ReverseInt32 反转数组
func ReverseInt32(array []uint32) []uint32 {
	for i := 0; i < len(array)/2; i++ {
		array[i], array[len(array)-i-1] = array[len(array)-i-1], array[i]
	}
	return array
}
