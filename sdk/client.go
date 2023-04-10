package gofound

import (
	"fmt"
	"github.com/sea-team/gofound/core"
	"github.com/sea-team/gofound/global"
	"github.com/sea-team/gofound/searcher"
	"os"
	"runtime"
	"sync"
)

var once sync.Once

// Client 应该对外部屏蔽细节
// 尽量少的提供接口，但是又要保证功能性
type Client struct {
	config    *global.Config      //服务配置
	container *searcher.Container //运行实体
}

func newDefaultConfig() *global.Config {
	return &global.Config{
		Addr:        "127.0.0.1:5678",
		Data:        fmt.Sprintf(".%sdata", string(os.PathSeparator)),
		Debug:       true,
		Dictionary:  "./data/dictionary.txt",
		EnableAdmin: true,
		Gomaxprocs:  runtime.NumCPU() * 2,
		Shard:       0,
		Auth:        "",
		EnableGzip:  true,
		Timeout:     10 * 60,
	}
}
func newTokenizerAndContainer(config *global.Config) *searcher.Container {
	tokenizer := core.NewTokenizer(global.CONFIG.Dictionary)
	return core.NewContainer(tokenizer)
}

// NewClient 通过参数进行配置,必须指定全部参数
func NewClient(config *global.Config) *Client {
	global.CONFIG = config
	//初始化分词器
	container := newTokenizerAndContainer(config)
	global.Container = container
	return &Client{
		config:    config,
		container: container,
	}
}

// Default 使用默认参数创建服务
func Default() *Client {
	global.CONFIG = newDefaultConfig()
	container := newTokenizerAndContainer(global.CONFIG)
	global.Container = container
	return &Client{
		config:    global.CONFIG,
		container: container,
	}
}

// SetAddr 设置Web服务地址
func (c *Client) SetAddr(addr string) *Client {
	if addr == "" {
		return c
	}
	c.config.Addr = addr
	return c
}

// SetData 设置数据存放地址
func (c *Client) SetData(dir string) *Client {
	if dir == "" {
		return c
	}
	c.config.Data = dir
	return c
}

//TODO 其他配置项
