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

	test := c.GetOrCreate("test")

	fmt.Println(test.GetIndexSize())

	all := c.GetEngines()
	for name, engine := range all {
		fmt.Println(name)
		fmt.Println(engine)
	}
}
