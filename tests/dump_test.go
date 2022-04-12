package tests

import (
	"fmt"
	"gofound/searcher"
	"gofound/searcher/dump"
	"testing"
)

func Test(t *testing.T) {

	tree := searcher.NewUInt32ComparatorTree()

	for i := 0; i < 10; i++ {
		tree.Insert(uint32(i))
	}

	data := dump.Serialize(tree.Root)
	fmt.Println(data)

}
