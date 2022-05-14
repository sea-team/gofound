package config

// Auth 权限设置
type Auth struct {
	Enable   bool   `mapstructure:"enable" json:"enable" yaml:"enable"`       // 是否开启
	Username string `mapstructure:"username" json:"username" yaml:"username"` // 用户名
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 用户密码
}
