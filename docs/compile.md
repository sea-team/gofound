# 编译

`gofound` 基于`golang-1.18`，编译之前需要安装对于的golang版本。 

推荐使用编译好的[二进制文件](https://github.com/newpanjing/gofound/releases)

## Admin
> 如果需要Admin部分，请先构建admin，admin基于vue+element-ui+vite，而这些也需要安装nodejs

构建命令：

```shell
cd ./web/admin/assets/web/

npm install

npm run build
```

完成以上步骤之后，才能使用admin

## 编译

```shell
go get
go build -o gofound
```

## 依赖

```shell
go 1.18

require (
	github.com/emirpasic/gods v1.12.0
	github.com/gin-gonic/gin v1.7.7
	github.com/yanyiwu/gojieba v1.1.2
)

require (
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang/protobuf v1.3.3 // indirect
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/syndtr/goleveldb v1.0.0 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

```