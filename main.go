package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gofound/router/api"
	"gofound/searcher"
	"log"
	"os"
)

type Args struct {
	Addr    string
	DataDir string
	Debug   bool
}

func parseArgs() Args {
	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1:5678", "设置监听地址和端口")

	var dataDir string

	//兼容windows
	dir := fmt.Sprintf(".%sdata", string(os.PathSeparator))

	flag.StringVar(&dataDir, "data", dir, "设置数据存储目录")

	var debug bool
	flag.BoolVar(&debug, "debug", true, "设置是否开启调试模式")

	flag.Parse()

	return Args{
		Addr:    addr,
		DataDir: dataDir,
		Debug:   debug,
	}
}

func initEngine(args Args) *searcher.Engine {
	var engine = &searcher.Engine{
		IndexPath: args.DataDir,
	}
	option := engine.GetOptions()

	go engine.InitOption(option)
	engine.IsDebug = args.Debug

	return engine
}

func initGin(args Args, engine *searcher.Engine) {
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

	//注册api
	api.Register(router)

	api.SetEngine(engine)

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

	//初始化引擎
	engine := initEngine(args)
	//保存索引到磁盘
	defer engine.Close()

	//初始化gin
	initGin(args, engine)
}
