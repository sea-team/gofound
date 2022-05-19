package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestChan(t *testing.T) {

	data := make(chan int)

	go func() {
		for {
			time.Sleep(time.Second * 1)
			data <- rand.Intn(100)
			break
		}
	}()

	r := <-data
	fmt.Println(r)

}
