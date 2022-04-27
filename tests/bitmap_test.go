package tests

import (
	"bytes"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/RoaringBitmap/roaring/roaring64"
	"gofound/searcher/utils"
	"gofound/tree"
	"math/rand"
	"testing"
	_ "unsafe"
)

func TestBitmap(t *testing.T) {
	fmt.Println("==roaring64==")

	rb3 := roaring64.New()

	var i uint64
	for i = 0; i < 10000000000; i++ {
		rb3.Add(i)
	}
	//fmt.Println(rb3.String())
	fmt.Println("ok")
}

//测试bitmap和二叉树效率
func TestBitmapAndTree(t *testing.T) {
	fmt.Println("==roaring64==")

	rb3 := roaring64.New()

	var size = 100000000

	var data = make([]uint64, size)

	var found = make([]uint64, 10000)

	for i := 0; i < size; i++ {
		//产生随机数
		data[i] = uint64(rand.Intn(size))
		//查找随机数
		if i < len(found) {
			found[i] = uint64(rand.Intn(size))
		}
	}
	//放入数据
	time := utils.ExecTime(func() {
		for _, v := range data {
			rb3.Add(v)
		}
	})
	fmt.Println("bitmap add time:", time)

	time = utils.ExecTime(func() {
		for _, v := range found {
			rb3.Contains(v)
		}
	})
	fmt.Println("bitmap select time:", time)

	//二叉树
	binaryTree := &tree.Tree{Comparator: utils.Uint32Comparator}
	time = utils.ExecTime(func() {
		for _, v := range data {
			binaryTree.Insert(uint32(v))
		}
	})
	fmt.Println("binaryTree add time:", time)

	time = utils.ExecTime(func() {
		for _, v := range found {
			binaryTree.Exists(uint32(v))
		}
	})
	fmt.Println("binaryTree select time:", time)

	//fmt.Println(rb3.String())
	fmt.Println("ok")
}

func TestBitmapSkip(t *testing.T) {
	rb3 := roaring.New()
	//序列化
	//rb3.WriteTo()
	var size = 10000

	for i := 0; i < size; i++ {
		//产生随机数
		num := rand.Intn(size)
		rb3.Add(uint32(num))
	}
	//roaring.And()
	fmt.Println(rb3.GetCardinality())

	it := rb3.ReverseIterator()
	for it.HasNext() {
		fmt.Print(it.Next())
		fmt.Print(" ")
	}
	//fmt.Println(len(rb3.ToArray()))

	//fmt.Println(rb3.String())
}

func TestEncode(t *testing.T) {
	var size = 100000
	data := make([]uint32, size)
	for i := 0; i < size; i++ {
		//产生随机数
		num := rand.Intn(size)
		data[i] = uint32(num)
	}

	r := roaring.New()
	r.AddMany(data)

	buffer := new(bytes.Buffer)
	r.WriteTo(buffer)

	nr := roaring.NewBitmap()
	_time := utils.ExecTime(func() {
		nr.FromBuffer(buffer.Bytes())
	})
	fmt.Println("roaring decode time:", _time)

	_time = utils.ExecTime(func() {
		r.ToArray()
	})
	fmt.Println("roaring to array time:", _time)
}
