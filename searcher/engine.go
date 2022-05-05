package searcher

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/wangbin/jiebago"

	"gofound/searcher/arrays"
	"gofound/searcher/dump"
	"gofound/searcher/model"
	"gofound/searcher/pagination"
	"gofound/searcher/sorts"
	"gofound/searcher/storage"
	"gofound/searcher/utils"
	"gofound/tree"
)

type Engine struct {
	IndexPath string

	//关键词索引树
	Indexes []*tree.Tree

	Option *Option

	//关键字和Id映射
	KeyMapperStorages []*storage.LeveldbStorage

	//ID和key映射，用于计算相关度，一个id 对应多个key
	IdKeyMapperStorages []*storage.LeveldbStorage

	//文档仓
	DocStorages []*storage.LeveldbStorage

	//锁
	sync.Mutex
	//等待
	sync.WaitGroup

	//文件分片
	Shard int

	//添加索引的通道
	AddDocumentWorkerChan chan model.IndexDoc

	//是否调试模式
	isDebug bool
}

type Option struct {
	KeyIndexName string
	KeyIdName    string
	IdKeyName    string
	DocIndexName string

	// 搜索结果最大数量，百度也才返回很少一部分
	MaxResultSize int
}

var seg jiebago.Segmenter

func NewUInt32ComparatorTree() *tree.Tree {
	return &tree.Tree{Comparator: utils.Uint32Comparator}
}

func (e *Engine) Init() {
	e.Add(1)
	defer e.Done()
	//线程数=cpu数
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	//保持和gin一致
	e.isDebug = os.Getenv("GIN_MODE") != "release"
	if e.Option == nil {
		e.Option = e.GetOptions()
	}
	log.Println("数据存储目录：", e.IndexPath)

	seg.LoadDictionary("./data/dictionary.txt")

	if e.Shard == 0 {
		e.Shard = 10
	}

	//初始化chan
	e.AddDocumentWorkerChan = make(chan model.IndexDoc, 1000)

	//初始化文件存储
	for shard := 0; shard < e.Shard; shard++ {
		//初始化chan
		go e.DocumentWorkerExec()
		//初始化树
		index := dump.Read(e.getFilePath(fmt.Sprintf("%s.%d", e.Option.KeyIndexName, shard)))
		e.Indexes = append(e.Indexes, index)

		s, err := storage.Open(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.DocIndexName, shard)))
		if err != nil {
			panic(err)
		}
		e.DocStorages = append(e.DocStorages, s)

		//初始化Keys存储
		ks, kerr := storage.Open(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.KeyIdName, shard)))
		if kerr != nil {
			panic(err)
		}
		e.KeyMapperStorages = append(e.KeyMapperStorages, ks)

		//id和keys映射
		iks, ikerr := storage.Open(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.IdKeyName, shard)))
		if ikerr != nil {
			panic(ikerr)
		}
		e.IdKeyMapperStorages = append(e.IdKeyMapperStorages, iks)
	}
	log.Println("初始化完成")

	//初始化完成，自动检测索引并持久化到磁盘
	go e.automaticFlush()
}
func (e *Engine) IndexDocument(doc model.IndexDoc) {
	e.AddDocumentWorkerChan <- doc
}

// DocumentWorkerExec 添加文档队列
func (e *Engine) DocumentWorkerExec() {
	for {
		doc := <-e.AddDocumentWorkerChan
		e.AddDocument(&doc)
	}
}

// 自动保存索引，10秒钟检测一次
func (e *Engine) automaticFlush() {
	ticker := time.NewTicker(time.Second * 10)
	size := e.GetIndexSize()

	for {
		<-ticker.C
		//检查数据是否有变动
		if size != e.GetIndexSize() {
			size = e.GetIndexSize()
			e.FlushIndex()
		}
		//定时GC
		runtime.GC()
	}

}

func (e *Engine) getShard(id uint32) int {
	return int(id % 10)
}

func (e *Engine) InitOption(option *Option) {

	if option == nil {
		//默认值
		option = e.GetOptions()
	}
	e.Option = option

	//初始化其他的
	e.Init()

}

func (e *Engine) getFilePath(fileName string) string {
	return e.IndexPath + string(os.PathSeparator) + fileName
}

func (e *Engine) GetOptions() *Option {
	return &Option{
		KeyIndexName: "key",
		DocIndexName: "doc",
		KeyIdName:    "key_id",
		IdKeyName:    "id_key",
	}
}

// WordCut 分词，只取长度大于2的词
func (e *Engine) WordCut(text string) []string {
	//不区分大小写
	text = strings.ToLower(text)

	var wordMap = make(map[string]int)

	resultChan := seg.CutForSearch(text, true)
	for {
		w, ok := <-resultChan
		if !ok {
			break
		}
		_, found := wordMap[w]
		if !found {
			//去除重复的词
			wordMap[w] = 1
		}
	}

	var wordsSlice []string
	for k, _ := range wordMap {
		wordsSlice = append(wordsSlice, k)
	}

	return wordsSlice
}

// AddDocument 分词索引
func (e *Engine) AddDocument(index *model.IndexDoc) {
	//等待初始化完成
	e.Wait()
	text := index.Text

	words := e.WordCut(text)

	//id对应的词

	keys := make([]uint32, 0)

	for _, word := range words {
		keyValue := utils.StringToInt(word)
		shard := e.getShard(keyValue)
		//添加到内存的tree
		if !e.Indexes[shard].Exists(keyValue) {
			e.Indexes[shard].Insert(keyValue)
		}

		keys = append(keys, keyValue)
		e.addKeyIndex(keyValue, index.Id)
	}

	//添加id索引
	e.addIdIndex(index, keys)
}

func (e *Engine) addKeyIndex(keyValue uint32, id uint32) {
	e.Lock()
	defer e.Unlock()
	ids := make([]uint32, 0)

	k := utils.Uint32ToBytes(keyValue)
	shard := e.getShard(keyValue)

	s := e.KeyMapperStorages[shard]

	//查找是否存在
	found := e.Indexes[shard].Exists(keyValue)

	if found {
		//存在
		//添加到列表
		buf, find := s.Get(k)
		if find {
			//解码
			utils.Decoder(buf, &ids)

			//直接添加，不排序，无序有利于快排
			//判断是否存在
			if !arrays.Exists(ids, id) {
				ids = append(ids, id)
			}
		} else {
			ids = append(ids, id)
		}
	} else {
		ids = append(ids, id)
	}

	err := s.Set(k, utils.Encoder(ids))
	if err != nil {
		return
	}
}

func (e *Engine) addIdIndex(index *model.IndexDoc, keys []uint32) {
	e.Lock()
	defer e.Unlock()
	//gob序列化
	k := utils.Uint32ToBytes(index.Id)
	shard := e.getShard(index.Id)
	s := e.DocStorages[shard]

	//id和key的映射
	iks := e.IdKeyMapperStorages[shard]

	doc := &model.StorageIndexDoc{
		IndexDoc: index,
		Keys:     keys,
	}

	//存储id和key以及文档的映射
	s.Set(k, utils.Encoder(doc))

	//设置到id和key的映射中
	iks.Set(k, utils.Encoder(keys))
}

// MultiSearch 多线程搜索
func (e *Engine) MultiSearch(request *model.SearchRequest) *model.SearchResult {
	//等待搜索初始化完成
	e.Wait()
	//分词搜索
	words := e.WordCut(request.Query)

	totalTime := float64(0)

	fastSort := new(sorts.FastSort)
	var wg sync.WaitGroup
	wg.Add(len(words))

	keys := make([]uint32, len(words))

	for i, word := range words {
		keys[i] = utils.StringToInt(word)
	}

	_time := utils.ExecTime(func() {

		for _, key := range keys {
			go e.SimpleSearch(key, keys, func(values []*model.SliceItem) {
				wg.Done()
				fastSort.Add(values)
			})
		}

		wg.Wait()
	})
	if e.isDebug {
		log.Println("数组查找耗时：", totalTime, "ms")
		log.Println("搜索时间:", _time, "ms")
	}
	// 处理分页
	request = request.GetAndSetDefault()

	//读取文档
	var result = &model.SearchResult{
		Total: fastSort.Count(),
		Time:  float32(_time),
		Page:  request.Page,
		Limit: request.Limit,
		Words: words,
	}

	_time = utils.ExecTime(func() {

		pager := new(pagination.Pagination)
		var resultIds []model.SliceItem
		_tt := utils.ExecTime(func() {
			resultIds = fastSort.GetAll(request.Order)
		})

		if e.isDebug {
			log.Println("处理排序耗时", _tt, "ms")
		}

		pager.Init(request.Limit, len(resultIds))
		//设置总页数
		result.PageCount = pager.PageCount

		//读取单页的id
		if pager.PageCount != 0 {

			start, end := pager.GetPage(request.Page)
			items := resultIds[start:end]

			//只读取前面100个
			for _, item := range items {

				buf := e.GetDocById(item.Id)
				doc := new(model.ResponseDoc)

				doc.Score = item.Score

				if buf != nil {
					//gob解析
					storageDoc := new(model.StorageIndexDoc)
					utils.Decoder(buf, &storageDoc)
					doc.Document = storageDoc.Document
					text := storageDoc.Text
					//处理关键词高亮
					highlight := request.Highlight
					if highlight != nil {
						//全部小写
						text = strings.ToLower(text)
						for _, word := range words {
							text = strings.ReplaceAll(text, word, fmt.Sprintf("%s%s%s", highlight.PreTag, word, highlight.PostTag))
						}
					}
					doc.Text = text
					doc.Id = item.Id
					result.Documents = append(result.Documents, *doc)
				}
			}
		}
	})
	if e.isDebug {
		log.Println("处理数据耗时：", _time, "ms")
	}

	return result
}
func (e *Engine) getRankAsync(keys []uint32, slice *model.SliceItem, call func()) {
	score := e.getRank(keys, slice.Id)
	slice.Score = score
	call()
}
func (e *Engine) getRank(keys []uint32, id uint32) float32 {
	shard := e.getShard(id)
	iks := e.IdKeyMapperStorages[shard]
	score := float32(1)
	if buf, exists := iks.Get(utils.Uint32ToBytes(id)); exists {
		memKeys := make([]uint32, 0)
		utils.Decoder(buf, &memKeys)

		//判断两个数的交集部分，就是得分

		size := len(keys)
		for i, k := range keys {
			//二分法查找
			count := float32(0)
			if arrays.BinarySearch(memKeys, k) {
				//计算基础分，至少1分
				base := float32(size - i)
				//关键词在越前面，分数越高
				score += base
				count++
			}

			//匹配关键词越多，数分越高
			if count != 0 {
				score *= count
			}
		}

	}
	return score
}
func (e *Engine) SimpleSearch(key uint32, keys []uint32, call func(ranks []*model.SliceItem)) {

	shard := e.getShard(key)
	found := e.Indexes[shard].Exists(key)

	if found {
		//读取id
		s := e.KeyMapperStorages[shard]

		kv := utils.Uint32ToBytes(key)

		data, find := s.Get(kv)

		if find {
			array := make([]uint32, 0)
			//解码
			utils.Decoder(data, &array)
			results := make([]*model.SliceItem, len(array))

			var wg sync.WaitGroup
			wg.Add(len(array))

			for index, id := range array {
				rank := &model.SliceItem{}
				rank.Id = id
				go e.getRankAsync(keys, rank, func() {
					wg.Done()
				})
				results[index] = rank
			}
			wg.Wait()
			//放结果集
			call(results)
		} else {
			call(nil)
		}

		//通过关键词匹配度，来计算得分
	} else {
		call(nil)
	}

}

func (e *Engine) GetIndexSize() int {
	size := 0
	for _, index := range e.Indexes {
		size += index.Size()
	}
	return size
}

// FlushIndex 刷新缓存到磁盘
func (e *Engine) FlushIndex() {
	e.Lock()
	defer e.Unlock()

	for i, index := range e.Indexes {
		dump.Write(index.Root, e.getFilePath(fmt.Sprintf("%s.%d", e.Option.KeyIndexName, i)))
	}
}

// GetDocById 通过id获取文档
func (e *Engine) GetDocById(id uint32) []byte {
	shard := e.getShard(id)
	key := utils.Uint32ToBytes(id)
	buf, found := e.DocStorages[shard].Get(key)
	if found {
		return buf
	}

	return nil
}

// RemoveIndex 根据ID移除索引
func (e *Engine) RemoveIndex(id uint32) error {
	//移除
	e.Lock()
	defer e.Unlock()

	shard := e.getShard(id)
	key := utils.Uint32ToBytes(id)

	//关键字和Id映射
	//KeyMapperStorages []*storage.LeveldbStorage
	//ID和key映射，用于计算相关度，一个id 对应多个key
	ik := e.IdKeyMapperStorages[shard]
	keysValue, found := ik.Get(key)
	if !found {
		return errors.New(fmt.Sprintf("没有找到id=%d", id))
	}

	keys := make([]uint32, 0)
	utils.Decoder(keysValue, &keys)

	//符合条件的key，要移除id
	for _, k := range keys {
		kv := utils.Uint32ToBytes(k)
		ks := e.KeyMapperStorages[e.getShard(k)]
		buf, exists := ks.Get(kv)
		if exists {
			ids := make([]uint32, 0)
			utils.Decoder(buf, &ids)
			//如果存在，才移除
			index := arrays.Find(ids, id)
			if index != -1 {
				ids = utils.DeleteArray(ids, index)
				//ids = append(ids[:index], ids[index+1:]...)
				ks.Set(kv, utils.Encoder(ids))
			}
			//如果key映射没有了，删除key-tree节省内存
			if len(ids) == 0 {
				e.Indexes[e.getShard(k)].Remove(k)
			}
		}
	}

	//删除id映射
	err := ik.Delete(key)
	if err != nil {
		return errors.New(err.Error())
	}

	//文档仓
	err = e.DocStorages[shard].Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) Close() {
	e.Lock()
	defer e.Unlock()

	//保存文件
	e.FlushIndex()

}
