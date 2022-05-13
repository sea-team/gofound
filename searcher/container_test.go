package searcher

import (
	"fmt"
	"testing"
)

func TestContainer_Init(t *testing.T) {
	c := &Container{
		Dir:   "/Users/panjing/GolandProjects/gofound/dbs",
		Debug: true,
	}
	err := c.Init()
	if err != nil {
		panic(err)
	}

	test := c.GetDataBase("test")

	fmt.Println(test.GetIndexCount())

	all := c.GetDataBases()
	for name, engine := range all {
		fmt.Println(name)
		fmt.Println(engine)
	}
}
