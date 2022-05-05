# 持久化

+ 关键词索引

关键词是存在内存中的二叉查找树，每个关键词都有一个唯一的索引，这个索引可以通过 `key` 来获取。

`gofound`启动了一个协程，每隔10s检测一次数据是否有变动，有变动的情况就存入磁盘中。

这里我们没有使用`leveldb`，是因为`leveldb`是一个key-value的数据库，而我们的数据只有key没有value，用`leveldb`存储势必会造成存储空间的浪费。
而且频繁取和存，会造成较高的时延和IO。

存储格式：

格式较为简单，由于二叉查找树使用的`uint32`类型，一个key占用4个字节，所以在存储的时候是直接以二进制的方式写入到缓存，最后在压缩进行存储。
[点击查看源码](../searcher/dump/dump.go)

```go
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


```

+ 关键词与ID映射

二叉树的每个关键词都与ID相关联，这样在搜索的时候，可以先找到索引的key，然后在通过key找到对应的id数组。

映射文件采用的是`leveldb`存储，编码格式为`gob`

[查看源码](../searcher/storage/leveldb_storage.go)

+ 文档

文档是指在索引时传入的数据，在搜索的时候会原样返回。

存储文件采用的是leveldb存储，编码格式为gob

[查看源码](../searcher/storage/leveldb_storage.go)
