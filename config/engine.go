package config

// Engine 搜索引擎
type Engine struct {
	DataDir       string `mapstructure:"dataDir" json:"dataDir" yaml:"dataDir"`                   // 数据存放路径
	DictionaryDir string `mapstructure:"dictionaryDir" json:"dictionaryDir" yaml:"dictionaryDir"` // 词库路径
	Shard         int    `mapstructure:"shard" json:"shard" yaml:"shard"`                         // 文件分块
	QueueMax      int    `mapstructure:"queue_max" json:"queueMax" yaml:"queueMax"`               // 索引队列最大值
	GcInterval    int    `mapstructure:"gcInterval" json:"gcInterval" yaml:"gcInterval"`          // GC时间间隔
}
