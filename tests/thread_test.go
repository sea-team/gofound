package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type ThreadTest struct {
	sync.Mutex
}

var wg sync.WaitGroup

func (t *ThreadTest) Test(name int) {
	defer t.Unlock()
	t.Lock()
	time.Sleep(time.Second * 1)
	fmt.Println("我是线程", name, "执行结束")
	wg.Done()
}

func TestThread(t *testing.T) {

	//sync.Mutex
	test := new(ThreadTest)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go test.Test(i)
	}
	wg.Wait()
	fmt.Println("完成了")

}
