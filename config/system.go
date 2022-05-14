package config

// System 系统设置
type System struct {
	Addr       int  `mapstructure:"addr" json:"addr" yaml:"addr"`                   // 服务监听地址
	Debug      bool `mapstructure:"debug" json:"debug" yaml:"debug"`                // 是否开启debug模式
	EnableGzip bool `mapstructure:"enableGzip" json:"enableGzip" yaml:"enableGzip"` // 是否开启gzip压缩
}
