package global

import (
	"gofound/config"
	"gofound/searcher"

	"github.com/spf13/viper"
)

var (
	VP        *viper.Viper   // 解析器
	CONFIG    *config.Server // 服务器设置
	Container *searcher.Container
)
