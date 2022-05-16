package global

// Config 服务器设置
type Config struct {
	Addr          string // 监听地址
	DataDir       string // 数据目录
	Debug         bool   // 调试模式
	DictionaryDir string // 字典路径
	EnableAdmin   bool   //启用admin
	Gomaxprocs    int    //GOMAXPROCS
	Shard         int    //分片数
	Auth          string //认证
	EnableGzip    bool   //是否开启gzip压缩
}
