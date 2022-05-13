package main

import (
	"context"
	"flag"
	"fmt"
	"gofound/controller"
	"gofound/core"
	"gofound/global"
	"gofound/initialize"
	"gofound/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
	tokenizer := initialize.Tokenizer(global.CONFIG.Engine.DictionaryDir)
	global.Container = initialize.Container(tokenizer)

	// 初始化业务逻辑
	controller.NewServices()

	// 注册路由
	r := router.SetupRouter()
	// 启动服务
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", global.CONFIG.System.Addr),
		Handler: r,
	}
	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("listen:", err)
		}
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shtdown:", err)
	}

	log.Println("Server exiting")

}
