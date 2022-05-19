package tests

import (
	"bufio"
	"fmt"
	"gofound/searcher"
	"gofound/searcher/model"
	"gofound/searcher/utils"
	"gofound/searcher/words"
	"os"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {

	tokenizer := words.NewTokenizer("../searcher/words/data/dictionary.txt")

	var engine = &searcher.Engine{
		IndexPath: "./index/db2",
		Tokenizer: tokenizer,
	}
	option := engine.GetOptions()

	engine.InitOption(option)

	f, err := os.Open("index/toutiao_cat_data.txt")
	if err != nil {
		t.Errorf("open file: %v", err)
	}

	id := uint32(0)
	rd := bufio.NewReader(f)
	index := 0
	for {
		line, isPrefix, err := rd.ReadLine()
		if err != nil {
			return
		}
		if isPrefix {
			t.Errorf("A long line has been cut, %s", line)
		}

		if len(line) == 0 {
			break
		}

		lineString := string(line)
		//fmt.Println(lineString)
		array := strings.Split(lineString, "_!_")
		if index%1000 == 0 {
			fmt.Println(index)
		}
		index++
		//if index == 6000 {
		//	break
		//}
		data := make(map[string]interface{})
		id++

		data["id"] = id
		data["title"] = array[3]
		data["category"] = array[2]
		data["cid"] = array[1]

		doc := model.IndexDoc{
			Id:       id,
			Text:     array[3],
			Document: data,
		}
		engine.IndexDocument(&doc)
	}
	for engine.GetQueue() > 0 {
	}
	fmt.Println("index finish")
}

func TestRepeat(t *testing.T) {
	//判断是否重复

	tokenizer := words.NewTokenizer("../searcher/words/data/dictionary.txt")
	var engine = &searcher.Engine{
		IndexPath: "./index",
		Tokenizer: tokenizer,
	}
	option := engine.GetOptions()

	engine.InitOption(option)

	f, err := os.Open("index/toutiao_cat_data.txt")
	if err != nil {
		t.Errorf("open file: %v", err)
	}

	container := make(map[uint32][]string)

	rd := bufio.NewReader(f)
	index := 0
	for {

		line, _, err := rd.ReadLine()
		if err != nil {
			break
		}

		lineString := string(line)
		array := strings.Split(lineString, "_!_")
		if index%10000 == 0 {
			fmt.Println(index)
		}
		index++

		data := struct {
			Id       string
			Title    string
			Category string
			Cid      string
		}{
			Id:       array[0],
			Title:    array[3],
			Category: array[2],
			Cid:      array[1],
		}

		//分词
		words := engine.Tokenizer.Cut(data.Title)
		for _, word := range words {
			key := Murmur3([]byte(word))
			val := container[key]
			if val == nil {
				val = make([]string, 0)
			}
			if !exists(val, word) {
				val = append(val, word)
			}
			container[key] = val
		}
	}

	//输出 value大于2的key
	for key, val := range container {
		if len(val) > 1 {
			fmt.Println("key:", key, "value:", val)
		}
	}

	fmt.Println("index finish")

}

func exists(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false

}

func TestStringToInt(t *testing.T) {
	/*
		key: 3756240089 value: [现场版 58.6]
		key: 2832448212 value: [树下 初展]
	*/

	fmt.Println(utils.StringToInt("现场版"))
	fmt.Println(utils.StringToInt("58.6"))

	fmt.Println(utils.StringToInt("树下"))
	fmt.Println(utils.StringToInt("初展"))
}

const (
	c1 = 0xcc9e2d51
	c2 = 0x1b873593
	c3 = 0x85ebca6b
	c4 = 0xc2b2ae35
	r1 = 15
	r2 = 13
	m  = 5
	n  = 0xe6546b64
)

var (
	Seed = uint32(1)
)

func Murmur3(key []byte) (hash uint32) {
	hash = Seed
	iByte := 0
	for ; iByte+4 <= len(key); iByte += 4 {
		k := uint32(key[iByte]) | uint32(key[iByte+1])<<8 | uint32(key[iByte+2])<<16 | uint32(key[iByte+3])<<24
		k *= c1
		k = (k << r1) | (k >> (32 - r1))
		k *= c2
		hash ^= k
		hash = (hash << r2) | (hash >> (32 - r2))
		hash = hash*m + n
	}

	var remainingBytes uint32
	switch len(key) - iByte {
	case 3:
		remainingBytes += uint32(key[iByte+2]) << 16
		fallthrough
	case 2:
		remainingBytes += uint32(key[iByte+1]) << 8
		fallthrough
	case 1:
		remainingBytes += uint32(key[iByte])
		remainingBytes *= c1
		remainingBytes = (remainingBytes << r1) | (remainingBytes >> (32 - r1))
		remainingBytes = remainingBytes * c2
		hash ^= remainingBytes
	}

	hash ^= uint32(len(key))
	hash ^= hash >> 16
	hash *= c3
	hash ^= hash >> 13
	hash *= c4
	hash ^= hash >> 16

	// 出发吧，狗嬷嬷！
	return
}
