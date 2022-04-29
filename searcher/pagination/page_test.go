package pagination

import (
	"fmt"
	"testing"
)

func TestPagination_GetPage(t *testing.T) {
	pagination := new(Pagination)

	var data []int64
	for i := 0; i < 100; i++ {
		data = append(data, int64(i))
	}

	pagination.Init(10, 100)

	for i := 1; i <= 10; i++ {
		start, end := pagination.GetPage(i)
		fmt.Println(start, end)
	}
}
