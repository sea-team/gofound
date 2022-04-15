# 配置

在编译好[gofound](./compile.md)之后，就可以启动了。

```shell
./gofound
```

## 参数

```shell
./gofound -h
  -addr string
        set server addr and port. (default "127.0.0.1:5678")
  -data string
        set data dataDir. (default "./data")

```

### addr

指定要监听的地址和端口。默认为`127.0.0.1:5678` 监听本地地址。

```shell
./gofound --addr=127.0.0.1:5678  
./gofound --addr=:5678  
./gofound --addr=0.0.0.0:5678  
./gofound --addr=192.168.1.1:5678  
```

### data

指定索引数据存储的目录，可以是相对路径，也可以是绝对路径。

相对路径是存在`gofound`所在目录下的。

```shell

```shell
./gofound --data=./data
./gofound --data=/www/data
```

## 生产模式
在生产模式下，不会输出一些不必要的日志，只输出错误日志。设置生产模式为`export GIN_MODE=release`
与GIN框架保持一致。

```shell
export GIN_MODE=release && ./gofound
```