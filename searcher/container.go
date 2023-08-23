package searcher

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/sea-team/gofound/searcher/words"
)

type Container struct {
	Dir       string             //文件夹
	engines   map[string]*Engine //引擎
	Debug     bool               //调试
	Tokenizer *words.Tokenizer   //分词器
	Shard     int                //分片
	Timeout   int64              //超时关闭数据库
	BufferNum int                //分片缓冲数
	rmu       sync.RWMutex
}

func (c *Container) Init() error {

	c.engines = make(map[string]*Engine)

	//读取当前路径下的所有目录，就是数据库名称
	dirs, err := os.ReadDir(c.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			//创建
			return os.MkdirAll(c.Dir, os.ModePerm)
		}
		return err
	}
	//初始化数据库
	for _, dir := range dirs {
		if dir.IsDir() {
			c.engines[dir.Name()] = c.GetDataBase(dir.Name())
			log.Println("db:", dir.Name())
		}
	}

	return nil
}

// NewEngine 创建一个引擎
func (c *Container) NewEngine(name string) *Engine {
	var engine = &Engine{
		IndexPath:    fmt.Sprintf("%s%c%s", c.Dir, os.PathSeparator, name),
		DatabaseName: name,
		Tokenizer:    c.Tokenizer,
		Shard:        c.Shard,
		Timeout:      c.Timeout,
		BufferNum:    c.BufferNum,
	}
	option := engine.GetOptions()

	engine.InitOption(option)
	engine.IsDebug = c.Debug

	return engine
}

// GetDataBase 获取或创建引擎
func (c *Container) GetDataBase(name string) *Engine {
	//默认数据库名为default
	if name == "" {
		name = "default"
	}

	c.rmu.Lock()
	defer c.rmu.Unlock()

	engine, ok := c.engines[name]
	if !ok {
		engine = c.NewEngine(name)
		c.engines[name] = engine
	}

	return engine
}

// GetDataBases 获取数据库列表
func (c *Container) GetDataBases() map[string]*Engine {
	c.rmu.RLock()
	defer c.rmu.RUnlock()

	out := make(map[string]*Engine, len(c.engines))
	for name, engine := range c.engines {
		out[name] = engine
	}
	return out
}

func (c *Container) GetDataBaseNumber() int {
	c.rmu.RLock()
	defer c.rmu.RUnlock()

	return len(c.engines)
}

func (c *Container) GetIndexCount() int64 {
	c.rmu.RLock()
	defer c.rmu.RUnlock()

	var count int64
	for _, engine := range c.engines {
		count += engine.GetIndexCount()
	}
	return count
}

func (c *Container) GetDocumentCount() int64 {
	c.rmu.RLock()
	defer c.rmu.RUnlock()

	var count int64
	for _, engine := range c.engines {
		count += engine.GetDocumentCount()
	}
	return count
}

// DropDataBase 删除数据库
func (c *Container) DropDataBase(name string) error {
	c.rmu.Lock()
	defer c.rmu.Unlock()

	if _, ok := c.engines[name]; !ok {
		return errors.New("数据库不存在")
	}
	err := c.engines[name].Drop()
	if err != nil {
		return err
	}

	delete(c.engines, name)
	//释放资源
	runtime.GC()

	return nil
}
