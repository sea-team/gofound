package arrays

import (
	"testing"
)

func TestArray(t *testing.T) {
	tests := []struct {
		array  []uint32
		target uint32
		exist  bool
	}{
		{
			array: []uint32{1, 2, 3, 4}, target: 2, exist: true,
		},
		{
			array: []uint32{1, 2, 3, 4}, target: 1, exist: true,
		},
		{
			array: []uint32{1, 2, 3, 4}, target: 4, exist: true,
		},
		{
			array: []uint32{1, 2, 3, 4}, target: 0, exist: false,
		},
		{
			array: []uint32{1, 2, 3, 4}, target: 5, exist: false,
		},
		{
			array: []uint32{1, 2, 5, 6}, target: 4, exist: false,
		},
		{
			array: []uint32{1, 2, 3, 3, 3, 4}, target: 3, exist: true,
		},
		{
			array: []uint32{1, 2, 3, 3, 3, 4, 4}, target: 4, exist: true,
		},
	}
	for i, test := range tests {
		if got := BinarySearch(test.array, test.target); got != test.exist {
			t.Fatalf("case %v: got %v, want %v", i, got, test.exist)
		}
	}
}
