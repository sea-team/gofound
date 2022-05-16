package searcher

import (
	"fmt"
	"gofound/searcher/words"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"unsafe"
)

type Container struct {
	Dir       string             //文件夹
	engines   map[string]*Engine //引擎
	Debug     bool               //调试
	Tokenizer *words.Tokenizer   //分词器
	Shard     int                //分片
	Timeout   int64              //超时关闭数据库
}

func (c *Container) Init() error {

	c.engines = make(map[string]*Engine)

	//读取当前路径下的所有目录，就是数据库名称
	dirs, err := ioutil.ReadDir(c.Dir)
	if err != nil {
		if os.IsNotExist(err) {
			//创建
			err := os.MkdirAll(c.Dir, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
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

	//log.Println("Get DataBase:", name)
	engine, ok := c.engines[name]
	if !ok {
		//创建引擎
		engine = c.NewEngine(name)
		c.engines[name] = engine
		//释放引擎

	}

	return engine
}

// GetDataBases 获取数据库列表
func (c *Container) GetDataBases() map[string]*Engine {
	for _, engine := range c.engines {
		size := unsafe.Sizeof(&engine)
		fmt.Printf("%s:%d\n", engine.DatabaseName, size)
	}
	return c.engines
}

func (c *Container) GetDataBaseNumber() int {
	return len(c.engines)
}

func (c *Container) GetIndexCount() int64 {
	var count int64
	for _, engine := range c.engines {
		count += engine.GetIndexCount()
	}
	return count
}

func (c *Container) GetDocumentCount() int64 {
	var count int64
	for _, engine := range c.engines {
		count += engine.GetDocumentCount()
	}
	return count
}

// DropDataBase 删除数据库
func (c *Container) DropDataBase(name string) error {
	err := c.engines[name].Drop()
	if err != nil {
		return err
	}

	delete(c.engines, name)
	//释放资源
	runtime.GC()

	return nil
}
