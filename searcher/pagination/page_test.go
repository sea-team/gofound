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

	pagination.Init(10, data)

	for i := 0; i < 10; i++ {
		r := pagination.GetPage(i)
		fmt.Println(r)
	}
}
