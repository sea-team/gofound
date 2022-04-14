package tests

import (
	"fmt"
	"testing"
)

func DeleteArray(array []uint32, index int) []uint32 {
	return append(array[:index], array[index+1:]...)
}

func TestArray(t *testing.T) {
	array := []uint32{1}
	fmt.Println(DeleteArray(array, 0))
}
