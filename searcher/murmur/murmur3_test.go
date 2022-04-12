package murmur

import (
	"fmt"
	"testing"
)

func TestMurmur3(t *testing.T) {
	r := Murmur3([]byte("test"))

	fmt.Println(r)
}
