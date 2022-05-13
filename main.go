package main

import (
	"flag"
	"fmt"
	"gofound/core"
	"gofound/global"
	"gofound/searcher"
	"gofound/searcher/system"
	"gofound/searcher/utils"
	"gofound/searcher/words"
	"gofound/web"
	"gofound/web/admin"
	"log"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func initTokenizer(dictionaryPath string) *words.Tokenizer {
	return words.NewTokenizer(dictionaryPath)
}

func initContainer(tokenizer *words.Tokenizer) *searcher.Container {
	container := &searcher.Container{
		Dir:       global.CONFIG.Engine.DataDir,
		Debug:     global.CONFIG.System.Debug,
		Tokenizer: tokenizer,
		Shard:     global.CONFIG.Engine.Shard,
	}
	go container.Init()

	return container
}

func initGin(container *searcher.Container) {
	if global.CONFIG.System.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// 启用GZIP压缩
	if global.CONFIG.System.EnableGzip {
		router.Use(gzip.Gzip(gzip.DefaultCompression))
	}
	//处理异常
	router.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				debug.PrintStack()
				c.JSON(200, web.Error(err.(error).Error()))
			}
			c.Abort()
		}()
		c.Next()
	})

	router.Use(web.Cors())

	err := router.SetTrustedProxies(nil)
	if err != nil {
		panic(err)
	}
	api := web.Api{
		Container: container,
		Callback: func() map[string]interface{} {

			s := &map[string]any{
				"os":             runtime.GOOS,
				"arch":           runtime.GOARCH,
				"cores":          runtime.NumCPU(),
				"version":        runtime.Version(),
				"goroutines":     runtime.NumGoroutine(),
				"dataPath":       global.CONFIG.Engine.DataDir,
				"dictionaryPath": global.CONFIG.Engine.DictionaryDir,
				"gomaxprocs":     runtime.NumCPU() * 2,
				"debug":          global.CONFIG.System.Debug,
				"shard":          global.CONFIG.Engine.Shard,
				"dataSize":       system.GetFloat64MB(utils.DirSizeB(global.CONFIG.Engine.DataDir)),
				"executable":     os.Args[0],
				"dbs":            container.GetDataBaseNumber(),
				"indexCount":     container.GetIndexCount(),
				"documentCount":  container.GetDocumentCount(),
				"pid":            os.Getpid(),
				"enableAuth":     global.CONFIG.Auth.Enable,
				"enableGzip":     global.CONFIG.System.EnableGzip,
			}

			return *s
		},
	}
	var handlers []gin.HandlerFunc

	if global.CONFIG.Auth.Enable {
		handlers = append(handlers, gin.BasicAuth(gin.Accounts{global.CONFIG.Auth.Username: global.CONFIG.Auth.Password}))
		log.Println("Enable Auth:", global.CONFIG.Auth.Enable)
	}

	//注册api
	api.Register(router, handlers...)

	log.Println("API Url： \t http://" + global.CONFIG.System.Addr + "/api")

	//注册admin
	if global.CONFIG.Auth.Enable {
		admin.Register(router, handlers...)
		log.Println("Admin Url： \t http://" + global.CONFIG.System.Addr + "/admin")
	}

	err = router.Run(global.CONFIG.System.Addr)
}

func main() {
	var config string
	// 默认使用./config.yaml文件，用户可使用一下命令使用自己的配置文件
	// go run main.go -c xx/xx.config
	flag.StringVar(&config, "c", "./config.yaml", "choose config file.")
	flag.Parse()

	global.VP = core.Viper(config)

	defer func() {

		if r := recover(); r != nil {
			fmt.Printf("panic: %s\n", r)
		}
	}()

	//初始化分词器
	tokenizer := initTokenizer(global.CONFIG.Engine.DictionaryDir)

	container := initContainer(tokenizer)

	//初始化gin
	initGin(container)

	fmt.Printf("Done!")
}
