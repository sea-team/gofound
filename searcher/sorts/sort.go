package sorts

import (
	"github.com/emirpasic/gods/trees/avltree"
	"gofound/searcher/utils"
	"log"
	"sync"
)

// IdSort 二叉树对id 进行打分和排序
type IdSort struct {
	Tree *avltree.Tree
	sync.Mutex
}

func NewIdSortTree() *IdSort {
	return &IdSort{
		Tree: &avltree.Tree{Comparator: utils.Uint32Comparator},
	}

}
func (e *IdSort) Add(key uint32) {
	count, found := e.Tree.Get(key)
	val := 1
	if found {
		val = count.(int) + 1
	}
	e.Lock()
	defer e.Unlock()
	e.Tree.Put(key, val)
}

func (e *IdSort) Size() int {
	return e.Tree.Size()
}

// GetAll 正序获取
func (e *IdSort) GetAll(order string) []uint32 {
	scores := make([]int, 0)
	ids := make([]uint32, 0)
	it := e.Tree.Iterator()
	_tt := utils.ExecTime(func() {
		for it.Next() {
			scores = append(scores, it.Value().(int))
			ids = append(ids, it.Key().(uint32))
		}
	})
	log.Println("迭代耗时:", _tt)

	_t := utils.ExecTime(func() {
		//ids 降序
		if order == "desc" {
			for i, j := 0, len(ids)-1; i < j; i, j = i+1, j-1 {
				ids[i], ids[j] = ids[j], ids[i]
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	})
	log.Println("id排序耗时:", _t)

	_t = utils.ExecTime(func() {
		// 排序，得分越高 排越前
		for i := 0; i < len(scores); i++ {
			for j := i + 1; j < len(scores); j++ {
				if scores[i] < scores[j] {
					scores[i], scores[j] = scores[j], scores[i]
					ids[i], ids[j] = ids[j], ids[i]
				}
			}
		}
	})

	log.Println("得分排序耗时:", _t)

	return ids
}
