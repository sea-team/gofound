package tests

import (
	"bufio"
	"fmt"
	"gofound/searcher"
	"gofound/searcher/model"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {

	var Engine = &searcher.Engine{
		IndexPath: "./index",
	}
	option := Engine.GetOptions()

	go Engine.InitOption(option)

	f, err := os.Open("./txt/toutiao_cat_data.txt")
	if err != nil {
		return
	}

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
		if index%10000 == 0 {
			fmt.Println(index)
		}
		index++

		data := make(map[string]interface{})
		data["id"] = array[0]
		data["title"] = array[3]
		data["category"] = array[2]
		data["cid"] = array[1]

		id, _ := strconv.Atoi(array[0])
		model := &model.IndexDoc{
			Id:       uint32(id),
			Text:     array[3],
			Document: data,
		}
		Engine.AddDocument(model)
	}

	for {
		queue := len(Engine.AddDocumentWorkerChan)
		if queue == 0 {
			break
		}
	}
	fmt.Println("index finish")
}
