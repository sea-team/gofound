package storage

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
)

func TestLeveldb(t *testing.T) {
	db, err := leveldb.OpenFile("/Users/panjing/GolandProjects/gofound/cache/doc_6.db", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	//_time := utils.ExecTime(func() {
	//
	//	for i := 0; i < 10000; i++ {
	//		db.Put([]byte(strconv.Itoa(i)), []byte(strconv.Itoa(i)), nil)
	//	}
	//})
	//fmt.Println("leveldb put 1000:", _time)
	db.Put([]byte("1"), []byte("1"), nil)
	value, err := db.Get([]byte("1"), nil)
	fmt.Println(string(value), err)
}
