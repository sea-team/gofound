# 持久化

持久化采用golang版本的leveldb

+ 关键词与ID映射

二叉树的每个关键词都与ID相关联，这样在搜索的时候，可以先找到索引的key，然后在通过key找到对应的id数组。

映射文件采用的是`leveldb`存储，编码格式为`gob`

[查看源码](../searcher/storage/leveldb_storage.go)


+ 文档

文档是指在索引时传入的数据，在搜索的时候会原样返回。

存储文件采用的是leveldb存储，编码格式为gob

[查看源码](../searcher/storage/leveldb_storage.go)