# API

`gofound`启动之后，会监听一个TCP端口，接收来自客户端的搜索请求。处理http请求部分使用`gin`框架。

## 多数据库支持

从1.1版本开始，我们支持了多数据库，API接口中通过get参数来指定数据库。

如果不指定，默认数据库为`default`。

如：`api/index?database=db1` 其他post参数不变

如果指定的数据库名没有存在，将会自动创建一个新的数据库。如果需要删除，直接删除改数据库目录，然后重启gofound即可。


## 增加/修改索引

| 接口地址 | /api/index       |
|------|------------------|
| 请求方式 | POST             |
| 请求类型 | application/json |

### 请求

| 字段       | 类型     | 必选  | 描述                                |
|----------|--------|-----|-----------------------------------|
| id       | uint32 | 是   | 文档的主键id，需要保持唯一性，如果id重复，将会覆盖直接的文档。 |
| text     | string | 是   | 需要索引的文本块                          |
| document | object | 是   | 附带的文档数据，json格式，搜索的时候原样返回          |

+ POST /api/index

```json
{
  "id": 88888,
  "text": "深圳北站",
  "document": {
    "title": "阿森松岛所445",
    "number": 223
  }
}
```

+ 命令行

```bash
curl -H "Content-Type:application/json" -X POST --data '{"id":88888,"text":"深圳北站","document":{"title":"阿森松岛所445","number":223}}' http://127.0.0.1:5678/api/index
```

### 响应

```json
{
  "state": true,
  "message": "success"
}
```

## 批量增加/修改索引


| 接口地址 | /api/index/batch       |
|------|------------------|
| 请求方式 | POST             |
| 请求类型 | application/json |

参数与单个一致，只是需要用数组包裹多个json对象，例如：

```json
[{
  "id": 88888,
  "text": "深圳北站",
  "document": {
    "title": "阿森松岛所445",
    "number": 223
  }
},{
  "id": 22222,
  "text": "北京东站",
  "document": {
    "title": "123123123",
    "number": 123123
  }
}]
```


## 删除索引

| 接口地址 | /api/remove      |
|------|------------------|
| 请求方式 | POST             |
| 请求类型 | application/json |

### 请求

| 字段  | 类型     | 必选  | 描述      |
|-----|--------|-----|---------|
| id  | uint32 | 是   | 文档的主键id |

+ POST /api/remove

```json
{
  "id": 88888
}
```

+ 命令行

```bash
curl -H "Content-Type:application/json" -X POST --data '{"id":88888}' http://127.0.0.1:5678/api/remove
```

### 响应

```json
{
  "state": true,
  "message": "success"
}
```

## 查询索引

`GoFound`提供了一种查询方式，按照文本查询。与其他Nosql数据库不同，`GoFound`不支持按照文档的其他查询。

| 接口地址 | /api/query       |
|------|------------------|
| 请求方式 | POST             |
| 请求类型 | application/json |

### 请求

| 字段        | 类型     | 必选  | 描述                                                   |
|-----------|--------|-----|------------------------------------------------------|
| query     | string | 是   | 查询的关键词，都是or匹配                                        |
| page      | int    | 否   | 页码，默认为1                                              |
| limit     | int    | 否   | 返回的文档数量，默认为100，没有最大限制，最好不要超过1000，超过之后速度会比较慢，内存占用会比较多 |
| order     | string | 否   | 排序方式，取值`asc`和`desc`，默认为`desc`，按id排序，然后根据结果得分排序       |
| highlight | object | 否   | 关键字高亮，相对text字段中的文本                                   |

### highlight

> 配置以后，符合条件的关键词将会被preTag和postTag包裹

| 字段      | 描述    |
|---------|-------|
| preTag  | 关键词前缀 |
| postTag | 关键词后缀 |

+ 示例

```json
{
  "query": "上海哪里好玩",
  "page": 1,
  "limit": 10,
  "order": "desc",
  "highlight": {
    "preTag": "<span style='color:red'>",
    "postTag": "</span>"
  }
}
```

+ POST /api/query

```json
{
  "query": "深圳北站",
  "page": 1,
  "limit": 10,
  "order": "desc"
}
```

+ 命令行

```bash
curl -H "Content-Type:application/json" -X POST --data '{"query":"深圳北站","page":1,"limit":10,"order":"desc"}' http://127.0.0.1:5678/api/query
```

### 响应

| 字段        | 类型      | 描述                      |
|-----------|---------|-------------------------|
| time      | float32 | 搜索文档用时                  |
| total     | int     | 符合条件的数量                 |
| pageCount | int     | 页总数                     |
| page      | int     | 当前页码                    |
| limit     | int     | 每页数量                    |
| documents | array   | 文档列表，[参考索引文档](#增加/修改索引) |

```json
{
  "state": true,
  "message": "success",
  "data": {
    "time": 2.75375,
    "total": 13487,
    "pageCount": 1340,
    "page": 1,
    "limit": 10,
    "documents": [
      {
        "id": 1675269553,
        "text": "【深圳消费卡/购物券转让/求购信息】- 深圳赶集网",
        "document": {
          "id": "8c68e948de7c7eb4362de15434a3ace7",
          "title": "【深圳消费卡/购物券转让/求购信息】- 深圳赶集网"
        },
        "score": 3
      },
      {
        "id": 88888,
        "text": "深圳北站",
        "document": {
          "number": 223,
          "title": "阿森松岛所445"
        },
        "score": 2
      },
      {
        "id": 212645608,
        "text": "【深圳美容美发卡转让/深圳美容美发卡求购信息】- 深圳赶集网",
        "document": {
          "id": "d3ce16b68a90833cbc20b8a49e93b9cd",
          "title": "【深圳美容美发卡转让/深圳美容美发卡求购信息】- 深圳赶集网"
        },
        "score": 1.5
      },
      {
        "id": 1191140208,
        "text": "【深圳赶集网】-免费发布信息-深圳分类信息门户",
        "document": {
          "id": "44be60a1d8b54c431e5511804062ae62",
          "title": "【深圳赶集网】-免费发布信息-深圳分类信息门户"
        },
        "score": 1.5
      },
      {
        "id": 4133884907,
        "text": "【深圳购物卡转让/深圳购物卡求购信息】- 深圳赶集网",
        "document": {
          "id": "f25bb8136e8c2b02e3fcd65627a9ddbc",
          "title": "【深圳购物卡转让/深圳购物卡求购信息】- 深圳赶集网"
        },
        "score": 1
      },
      {
        "id": 206909132,
        "text": "【沙嘴门票/电影票转让/求购信息】- 深圳赶集网",
        "document": {
          "id": "63ca3ea4ffd254454e738a0957efedc2",
          "title": "【沙嘴门票/电影票转让/求购信息】- 深圳赶集网"
        },
        "score": 1
      },
      {
        "id": 220071473,
        "text": "【深圳健身卡转让/深圳健身卡求购信息】- 深圳赶集网",
        "document": {
          "id": "72d3d650c8a8a4e73b89b406f6dc76ef",
          "title": "【深圳健身卡转让/深圳健身卡求购信息】- 深圳赶集网"
        },
        "score": 1
      },
      {
        "id": 461974720,
        "text": "铁路_论坛_深圳热线",
        "document": {
          "id": "73c96ac2c23bc0cb4fb12ce7660c8b35",
          "title": "铁路_论坛_深圳热线"
        },
        "score": 1
      },
      {
        "id": 490922879,
        "text": "【深圳购物卡转让/深圳购物卡求购信息】- 深圳赶集网",
        "document": {
          "id": "93be0f35c484ddcd8c83602e27535d96",
          "title": "【深圳购物卡转让/深圳购物卡求购信息】- 深圳赶集网"
        },
        "score": 1
      },
      {
        "id": 525810194,
        "text": "【深圳购物卡转让/深圳购物卡求购信息】- 深圳赶集网",
        "document": {
          "id": "e489dd19dce0de2c9f4e59c969ec9ec0",
          "title": "【深圳购物卡转让/深圳购物卡求购信息】- 深圳赶集网"
        },
        "score": 1
      }
    ],
    "words": [
      "深圳",
      "北站"
    ]
  }
}
```

## 查询状态

| 接口地址 | /api/status      |
|------|------------------|
| 请求方式 | GET              |

### 请求

```bash
curl http://127.0.0.1:5678/api/status
```

### 响应

```json
{
  "state": true,
  "message": "success",
  "data": {
    "index": {
      "queue": 0,
      "shard": 10,
      "size": 531971
    },
    "memory": {
      "alloc": 1824664656,
      "heap": 1824664656,
      "heap_idle": 10008625152,
      "heap_inuse": 2100068352,
      "heap_objects": 3188213,
      "heap_released": 9252003840,
      "heap_sys": 12108693504,
      "sys": 12700504512,
      "total": 11225144273040
    },
    "status": "ok",
    "system": {
      "arch": "arm64",
      "cores": 10,
      "os": "darwin",
      "version": "go1.18"
    }
  }
}
```

## 删除数据库

| 接口地址 | /api/drop |
|------|-----------|
| 请求方式 | GET       |

### 请求

```bash
curl http://127.0.0.1:5678/api/drop?database=db_name
```

### 响应

```json
{
  "state": true,
  "message": "success",
}
```

## 在线分词

| 接口地址 | /api/word/cut   |
|------|-----------------|
| 请求方式 | GET             |

### 请求参数

| 字段  | 类型     | 必选  | 描述  |
|-----|--------|-----|-----|
| q   | string | 关键词 |

### 请求

```bash
curl http://127.0.0.1:5678/api/word/cut?q=上海和深圳哪个城市幸福指数高
```

### 响应

```json
{
  "state": true,
  "message": "success",
  "data": [
    "上海",
    "深圳",
    "哪个",
    "城市",
    "幸福",
    "指数"
  ]
}
```