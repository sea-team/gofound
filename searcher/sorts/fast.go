package sorts

import (
	"gofound/searcher/model"
	"sort"
	"strings"
	"sync"
)

const (
	DESC = "desc"
)

type ScoreSlice []model.SliceItem

func (x ScoreSlice) Len() int {
	return len(x)
}
func (x ScoreSlice) Less(i, j int) bool {
	return x[i].Score < x[j].Score
}
func (x ScoreSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

type SortSlice []uint32

func (x SortSlice) Len() int {
	return len(x)
}
func (x SortSlice) Less(i, j int) bool {
	return x[i] < x[j]
}
func (x SortSlice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]

}

type Uint32Slice []uint32

func (x Uint32Slice) Len() int           { return len(x) }
func (x Uint32Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Uint32Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type FastSort struct {
	sync.Mutex

	IsDebug bool

	data []model.SliceItem

	temps []uint32

	count int //总数

	Order string //排序方式
}

func (f *FastSort) Add(ids *[]uint32) {
	//f.Lock()
	//defer f.Unlock()

	//for _, id := range *ids {
	//
	//	found, index := f.find(&id)
	//	if found {
	//		f.data[index].Score += 1
	//	} else {
	//
	//		f.data = append(f.data, model.SliceItem{
	//			Id:    id,
	//			Score: 1,
	//		})
	//		f.Sort()
	//	}
	//}
	//f.count = len(f.data)
	f.temps = append(f.temps, *ids...)
}

// 二分法查找
func (f *FastSort) find(target *uint32) (bool, int) {

	low := 0
	high := f.count - 1
	for low <= high {
		mid := (low + high) / 2
		if f.data[mid].Id == *target {
			return true, mid
		} else if f.data[mid].Id < *target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return false, -1
	//for index, item := range f.data {
	//	if item.Id == *target {
	//		return true, index
	//	}
	//}
	//return false, -1
}

// Count 获取数量
func (f *FastSort) Count() int {
	return f.count
}

// Sort 排序
func (f *FastSort) Sort() {
	if strings.ToLower(f.Order) == DESC {
		sort.Sort(sort.Reverse(SortSlice(f.temps)))
	} else {
		sort.Sort(SortSlice(f.temps))
	}
}

// Process 处理数据
func (f *FastSort) Process() {
	//计算重复
	f.Sort()

	for _, temp := range f.temps {
		if found, index := f.find(&temp); found {
			f.data[index].Score += 1
		} else {
			f.data = append(f.data, model.SliceItem{
				Id:    temp,
				Score: 1,
			})
			f.count++
		}
	}
	//对分数进行排序
	sort.Sort(sort.Reverse(ScoreSlice(f.data)))
}
func (f *FastSort) GetAll(result *[]model.SliceItem, start int, end int) {

	*result = f.data[start:end]
}
