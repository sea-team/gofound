package sorts

import (
	"gofound/searcher/model"
	"sort"
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

type Uint32Slice []uint32

func (x Uint32Slice) Len() int           { return len(x) }
func (x Uint32Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Uint32Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type FastSort struct {
	sync.Mutex

	IsDebug bool

	data []model.SliceItem
}

func (f *FastSort) Add(ids []uint32, frequency int) {
	f.Lock()
	defer f.Unlock()

	for _, id := range ids {

		found, index := find(f.data, id)
		if found {
			f.data[index].Score += 1
		} else {

			f.data = append(f.data, model.SliceItem{
				Id:    id,
				Score: 1,
			})
		}
	}
}

// 二分法查找
func find(data []model.SliceItem, target uint32) (bool, int) {
	low := 0
	high := len(data) - 1
	for low <= high {
		mid := (low + high) / 2
		if data[mid].Id == target {
			return true, mid
		} else if data[mid].Id < target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return false, -1
}

// Count 获取数量
func (f *FastSort) Count() int {
	return len(f.data)
}

func (f *FastSort) GetAll(order string) []model.SliceItem {

	//对分数进行排序
	sort.Sort(sort.Reverse(ScoreSlice(f.data)))

	return f.data
}
