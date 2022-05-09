package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gofound/assets"
	"gofound/router/api"
	"gofound/searcher"
	"gofound/searcher/words"
	"log"
	"net/http"
	"os"
	"runtime"
)

type Args struct {
	Addr           string
	DataDir        string
	Debug          bool
	DictionaryPath *string
}

func parseArgs() Args {
	var addr string
	flag.StringVar(&addr, "addr", "0.0.0.0:5678", "设置监听地址和端口")

	var dataDir string

	//兼容windows
	dir := fmt.Sprintf(".%sdata", string(os.PathSeparator))

	flag.StringVar(&dataDir, "data", dir, "设置数据存储目录")

	var debug bool
	flag.BoolVar(&debug, "debug", true, "设置是否开启调试模式")

	var dictionaryPath = flag.String("dictionary", "./data/dictionary.txt", "设置词典路径")
	flag.Parse()

	return Args{
		Addr:           addr,
		DataDir:        dataDir,
		Debug:          debug,
		DictionaryPath: dictionaryPath,
	}
}

func initTokenizer(dictionaryPath string) *words.Tokenizer {
	return words.NewTokenizer(dictionaryPath)
}

func initContainer(args Args, tokenizer *words.Tokenizer) *searcher.Container {
	container := &searcher.Container{
		Dir:       args.DataDir,
		Debug:     args.Debug,
		Tokenizer: tokenizer,
	}
	container.Init()
	return container
}

func initGin(args Args, container *searcher.Container) {
	if args.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	//处理异常
	router.Use(api.Recover)
	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}

	router.StaticFS("/static", http.FS(assets.Static))
	//注册api
	api.Register(router)

	api.SetContainer(container)

	log.Println("API url： \t http://" + args.Addr + "/api")

	err = router.Run(args.Addr)
}

func main() {

	defer func() {

		if r := recover(); r != nil {
			fmt.Printf("panic: %s\n", r)
		}
	}()
	//解析参数
	args := parseArgs()

	//线程数=cpu数
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	//初始化分词器
	tokenizer := initTokenizer(*args.DictionaryPath)

	container := initContainer(args, tokenizer)

	//初始化gin
	initGin(args, container)
}
