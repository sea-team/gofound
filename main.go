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

func main() {

	var addr string
	flag.StringVar(&addr, "addr", "127.0.0.1:5678", "设置监听地址和端口")

	var dataDir string

	//兼容windows
	dir := fmt.Sprintf(".%sdata", string(os.PathSeparator))

	flag.StringVar(&dataDir, "data", dir, "设置数据存储目录")

	var debug bool
	flag.BoolVar(&debug, "debug", true, "设置是否开启调试模式")

	flag.Parse()

	var engine = &searcher.Engine{
		IndexPath: dataDir,
	}
	option := engine.GetOptions()

	go engine.InitOption(option)

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine.IsDebug = debug

	router := gin.Default()
	//处理异常
	router.Use(api.Recover)
	err := router.SetTrustedProxies(nil)
	if err != nil {
		return
	}

	//注册api
	api.Register(router)

	//保存索引到磁盘
	defer engine.Close()
	api.SetEngine(engine)

	log.Println("API url： \t http://" + addr + "/api")

	err = router.Run(addr)
	defer func() {

		if r := recover(); r != nil {

			fmt.Printf("panic: %s\n", r)

		}

		fmt.Println("-- 2 --")

	}()
	fmt.Println("-- 1 --")
	if err != nil {
		fmt.Println("错误", err)
		return
	}
}
