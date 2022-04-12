package dump

import (
	"gofound/searcher/utils"
	"gofound/tree"
)

func Serialize(node *tree.Node) []uint32 {

	d := make([]uint32, 0)
	if node == nil {
		return d
	}
	d = append(d, node.Key.(uint32))

	left := node.Children[0]
	right := node.Children[1]

	d = append(d, Serialize(right)...)
	d = append(d, Serialize(left)...)

	return d

}

func Write(node *tree.Node, filename string) {
	data := Serialize(node)
	utils.Write(&data, filename)
}

func Read(filename string) *tree.Tree {

	data := make([]uint32, 0)
	utils.Read(&data, filename)

	tree := &tree.Tree{Comparator: utils.Uint32Comparator}
	//遍历重新组装成内存树
	for _, id := range data {
		tree.Insert(id)
	}
	return tree
}
