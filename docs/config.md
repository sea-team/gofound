# 配置

在编译好[gofound](./compile.md)之后，就可以启动了。

```shell
./gofound
```

## 参数

```shell
./gofound -h

  -addr string
    	设置监听地址和端口 (default "127.0.0.1:5678")
  -auth string
    	开启认证，例如: admin:123456
  -data string
    	设置数据存储目录 (default "./data")
  -debug
    	设置是否开启调试模式 (default true)
  -dictionary string
    	设置词典路径 (default "./data/dictionary.txt")
  -enableAdmin
    	设置是否开启后台管理 (default true)
  -enableGzip
    	是否开启gzip压缩 (default true)
  -gomaxprocs int
    	设置GOMAXPROCS (default 20)
  -shard int
    	文件分片数量 (default 10)

```

### addr

指定要监听的地址和端口。默认为`127.0.0.1:5678` 监听本地地址。

```shell
./gofound --addr=127.0.0.1:5678  
./gofound --addr=:5678  
./gofound --addr=0.0.0.0:5678  
./gofound --addr=192.168.1.1:5678  
```

### auth

设置admin和api接口的用户名密码，采用basic auth

```shell
./gofound --auth=admin:123456
```

### data

指定索引数据存储的目录，可以是相对路径，也可以是绝对路径。

相对路径是存在`gofound`所在目录下的。

```shell

```shell
./gofound --data=./data
./gofound --data=/www/data
```

### debug

设置是否开启调试模式。默认为`true`。

```shell
./gofound --debug=false
```

### dictionary

设置自定义词典路径。默认为`./data/dictionary.txt`。

```shell
./gofound --dictionary=./data/dictionary.txt
```

### enableAdmin

设置是否开启后台管理。默认为`true`。

```shell
./gofound --enableAdmin=false
```

### enableGzip

设置是否开启gzip压缩。默认为`true`。

```shell
./gofound --enableGzip=false
```

### gomaxprocs

设置GOMAXPROCS。默认为CPU数量X2。

```shell
./gofound --gomaxprocs=10
```

### shard

设置文件分片数量。默认为`10`。分片越多查询会越快，相反的磁盘IO和CPU会越多。

```shell
./gofound --shard=10
```
