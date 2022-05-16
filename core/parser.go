package core

import (
	"flag"
	"fmt"
	"gofound/global"
	"os"
	"runtime"
	//"github.com/spf13/viper"
)

// Parser 解析器
func Parser() *global.Config {
	//v := viper.New()
	//v.SetConfigFile(config)
	//err := v.ReadInConfig()
	//if err != nil {
	//	panic(err)
	//}
	//
	//if err := v.Unmarshal(&global.CONFIG); err != nil {
	//	panic(err)
	//}
	//
	//return v
	var configPath = flag.String("config", "", "配置文件路径，配置此项其他参数忽略")
	if *configPath != "" {
		//解析配置文件
		return nil
	}

	var addr = flag.String("addr", "0.0.0.0:5678", "设置监听地址和端口")
	//兼容windows
	dir := fmt.Sprintf(".%sdata", string(os.PathSeparator))

	var dataDir = flag.String("data", dir, "设置数据存储目录")

	var debug = flag.Bool("debug", true, "设置是否开启调试模式")

	var dictionaryPath = flag.String("dictionary", "./data/dictionary.txt", "设置词典路径")

	var enableAdmin = flag.Bool("enableAdmin", true, "设置是否开启后台管理")

	var gomaxprocs = flag.Int("gomaxprocs", runtime.NumCPU()*2, "设置GOMAXPROCS")

	var auth = flag.String("auth", "", "开启认证，例如: admin:123456")

	var enableGzip = flag.Bool("enableGzip", true, "是否开启gzip压缩")

	flag.Parse()

	config := &global.Config{
		Addr:          *addr,
		DataDir:       *dataDir,
		Debug:         *debug,
		DictionaryDir: *dictionaryPath,
		EnableAdmin:   *enableAdmin,
		Gomaxprocs:    *gomaxprocs,
		Auth:          *auth,
		EnableGzip:    *enableGzip,
	}
	fmt.Println(config)

	return config
}
