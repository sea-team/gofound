package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
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
	"strings"
)

type Args struct {
	Addr           string
	DataDir        string
	Debug          bool
	DictionaryPath string
	EnableAdmin    bool
	GOMAXPROCS     int
	Shard          int
	Auth           string //认证
}

func parseArgs() Args {

	var addr = flag.String("addr", "127.0.0.1:5678", "设置监听地址和端口")
	//兼容windows
	dir := fmt.Sprintf(".%sdata", string(os.PathSeparator))

	var dataDir = flag.String("data", dir, "设置数据存储目录")

	var debug = flag.Bool("debug", true, "设置是否开启调试模式")

	var dictionaryPath = flag.String("dictionary", "./data/dictionary.txt", "设置词典路径")

	var enableAdmin = flag.Bool("enableAdmin", true, "设置是否开启后台管理")

	var gomaxprocs = flag.Int("gomaxprocs", runtime.NumCPU()*2, "设置GOMAXPROCS")

	var shard = flag.Int("shard", 5, "文件分片数量")

	var auth = flag.String("auth", "", "开启认证，例如: admin:123456")

	flag.Parse()

	return Args{
		Addr:           *addr,
		DataDir:        *dataDir,
		Debug:          *debug,
		DictionaryPath: *dictionaryPath,
		EnableAdmin:    *enableAdmin,
		GOMAXPROCS:     *gomaxprocs,
		Shard:          *shard,
		Auth:           *auth,
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
		Shard:     args.Shard,
	}
	err := container.Init()
	if err != nil {
		panic(err)
	}
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
				"dataPath":       args.DataDir,
				"dictionaryPath": args.DictionaryPath,
				"gomaxprocs":     args.GOMAXPROCS,
				"debug":          args.Debug,
				"shard":          args.Shard,
				"dataSize":       system.GetFloat64MB(utils.DirSizeB(args.DataDir)),
				"executable":     os.Args[0],
				"dbs":            container.GetDataBaseNumber(),
				"indexCount":     container.GetIndexCount(),
				"documentCount":  container.GetDocumentCount(),
				"pid":            os.Getpid(),
				"enableAuth":     args.Auth != "",
			}

			return *s
		},
	}
	var handlers []gin.HandlerFunc

	if args.Auth != "" {
		splits := strings.Split(args.Auth, ":")
		if len(splits) != 2 {
			panic("auth format error")
		}

		handlers = append(handlers, gin.BasicAuth(gin.Accounts{splits[0]: splits[1]}))
		log.Println("Enable Auth:", args.Auth)
	}

	//注册api
	api.Register(router, handlers...)

	log.Println("API Url： \t http://" + args.Addr + "/api")

	//注册admin
	if args.EnableAdmin {
		admin.Register(router, handlers...)
		log.Println("Admin Url： \t http://" + args.Addr + "/admin")
	}

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
	runtime.GOMAXPROCS(args.GOMAXPROCS)

	//初始化分词器
	tokenizer := initTokenizer(args.DictionaryPath)

	container := initContainer(args, tokenizer)

	//初始化gin
	initGin(args, container)

	fmt.Printf("Done!")
}
