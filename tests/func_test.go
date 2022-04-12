package tests

import (
	"fmt"
	"testing"
)

type FuncTest struct {
	name string
}

func aa(a int, b *int, d *FuncTest) int {
	a = 111
	*b = 3
	fmt.Printf("b=%p\n", &b)
	fmt.Printf("d=%p\n", d)
	d.name = "aa"
	return a + *b
}

func TestA(t *testing.T) {

	var a int = 1
	var b int = 2
	d := &FuncTest{name: "test"}

	fmt.Printf("b=%p\n", &b)
	fmt.Printf("d=%p\n", d)

	fmt.Println(aa(a, &b, d))
	fmt.Println(d)
}
