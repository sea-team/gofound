package config

// Server 服务器设置
type Server struct {
	System *System `mapstructure:"system" json:"system" yaml:"system"` // 系统设置
	Auth   *Auth   `mapstructure:"auth" json:"auth" yaml:"auth"`       // 认证设置
	Engine *Engine `mapstructure:"engine" json:"engine" yaml:"engine"` // 搜索引擎设置
}
