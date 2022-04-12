package sorts

import (
	"gofound/searcher/model"
	"gofound/searcher/utils"
	"log"
	"sort"
	"sync"
)

const (
	DESC = "desc"
	ASC  = "asc"
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

type Uint32Slice []*model.SliceItem

func (x Uint32Slice) Len() int           { return len(x) }
func (x Uint32Slice) Less(i, j int) bool { return x[i].Id < x[j].Id }
func (x Uint32Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type FastSort struct {
	data []*model.SliceItem
	sync.Mutex
}

func (f *FastSort) Add(values []*model.SliceItem) {
	if values == nil {
		return
	}
	f.Lock()
	defer f.Unlock()
	f.data = append(f.data, values...)

}

// Count 获取数量
func (f *FastSort) Count() int {
	return len(f.data)
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

func (f *FastSort) GetAll(order string) []model.SliceItem {

	//声明大小，避免重复合并数组
	var result = make([]model.SliceItem, len(f.data))

	//降序排序
	_tt := utils.ExecTime(func() {

		if order == DESC {
			sort.Sort(sort.Reverse(Uint32Slice(f.data)))
		} else {
			sort.Sort(Uint32Slice(f.data))
		}
	})
	log.Println("排序 time:", _tt)

	k := 0
	_ttt := utils.ExecTime(func() {
		for _, item := range f.data {
			found, index := find(result, item.Id)
			if found {
				//log.Println("重复数据:", item.Id)
				result[index].Score += item.Score
			} else {
				result[k] = model.SliceItem{
					Id:    item.Id,
					Score: item.Score,
				}
				k++
			}
		}
	})
	log.Println("去重耗时", _ttt)

	//对分数进行排序
	sort.Sort(sort.Reverse(ScoreSlice(result)))

	return result[:k]
}
