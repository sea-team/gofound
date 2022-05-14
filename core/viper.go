package core

import (
	"gofound/global"

	"github.com/spf13/viper"
)

// Viper 解析器
func Viper(config string) *viper.Viper {
	v := viper.New()
	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&global.CONFIG); err != nil {
		panic(err)
	}

	return v
}
