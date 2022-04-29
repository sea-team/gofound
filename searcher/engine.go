package searcher

import (
	"bytes"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/wangbin/jiebago"
	"gofound/searcher/arrays"
	"gofound/searcher/model"
	"gofound/searcher/pagination"
	"gofound/searcher/sorts"
	"gofound/searcher/storage"
	"gofound/searcher/utils"
	"gofound/tree"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Engine struct {
	IndexPath string

	Option *Option

	//关键字和Id映射，正排索引
	InvertedIndexStorages []*storage.LeveldbStorage

	//ID和key映射，用于计算相关度，一个id 对应多个key，倒排索引
	PositiveIndexStorages []*storage.LeveldbStorage

	//文档仓
	DocStorages []*storage.LeveldbStorage

	//锁
	sync.Mutex
	//等待
	sync.WaitGroup

	//添加索引的通道
	AddDocumentWorkerChan chan model.IndexDoc

	//是否调试模式
	IsDebug bool
}

type Option struct {
	KeyIndexName      string
	InvertedIndexName string
	PositiveIndexName string
	DocIndexName      string
	Dictionary        string //词典路径
	Shard             int    //分片数，默认为5
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
	if e.Option == nil {
		e.Option = e.GetOptions()
	}
	log.Println("数据存储目录：", e.IndexPath)

	seg.LoadDictionary(e.Option.Dictionary)

	//初始化chan
	e.AddDocumentWorkerChan = make(chan model.IndexDoc, 1000)

	//初始化文件存储
	for shard := 0; shard < e.Option.Shard; shard++ {
		//初始化chan
		go e.DocumentWorkerExec()

		s, err := storage.Open(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.DocIndexName, shard)))
		if err != nil {
			panic(err)
		}
		e.DocStorages = append(e.DocStorages, s)

		//初始化Keys存储
		ks, kerr := storage.Open(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.InvertedIndexName, shard)))
		if kerr != nil {
			panic(err)
		}
		e.InvertedIndexStorages = append(e.InvertedIndexStorages, ks)

		//id和keys映射
		iks, ikerr := storage.Open(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.PositiveIndexName, shard)))
		if ikerr != nil {
			panic(ikerr)
		}
		e.PositiveIndexStorages = append(e.PositiveIndexStorages, iks)
	}
	go e.automaticGC()
	log.Println("初始化完成")
}

// 自动保存索引，10秒钟检测一次
func (e *Engine) automaticGC() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		<-ticker.C
		//定时GC
		runtime.GC()
	}
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

// getShard 计算索引分布在哪个文件块
func (e *Engine) getShard(id uint32) int {
	return int(id % uint32(e.Option.Shard))
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
		KeyIndexName:      "key",
		DocIndexName:      "docs",
		InvertedIndexName: "inverted_index",
		PositiveIndexName: "positive_index",
		Dictionary:        "./data/dictionary.txt",
		Shard:             5,
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

	keys := make([]uint32, len(words))

	for i, word := range words {
		keyValue := utils.StringToInt(word)
		keys[i] = keyValue
		e.addInvertedIndex(keyValue, index.Id)
	}

	//添加id索引
	e.addPositiveIndex(index, keys)
}

// 添加倒排索引
func (e *Engine) addInvertedIndex(keyValue uint32, id uint32) {
	e.Lock()
	defer e.Unlock()

	key := utils.Uint32ToBytes(keyValue)
	shard := e.getShard(keyValue)

	s := e.InvertedIndexStorages[shard]

	//存在
	//添加到列表
	buf, find := s.Get(key)
	bitmap := roaring.New()

	if find {
		//解码
		_, err := bitmap.FromBuffer(buf)
		if err != nil {
			panic(err)
		}

	}

	bitmap.Add(id)

	value := new(bytes.Buffer)
	_, err := bitmap.WriteTo(value)
	if err != nil {
		panic(err)
	}
	s.Set(key, value.Bytes())
}

// 添加正排索引 id=>keys id=>doc
func (e *Engine) addPositiveIndex(index *model.IndexDoc, keys []uint32) {
	e.Lock()
	defer e.Unlock()

	key := utils.Uint32ToBytes(index.Id)
	shard := e.getShard(index.Id)
	s := e.DocStorages[shard]

	//id和key的映射
	iks := e.PositiveIndexStorages[shard]

	doc := &model.StorageIndexDoc{
		IndexDoc: index,
		Keys:     keys,
	}

	//存储id和key以及文档的映射
	s.Set(key, utils.Encoder(doc))

	//设置到id和key的映射中
	iks.Set(key, utils.Encoder(keys))
}

// MultiSearch 多线程搜索
func (e *Engine) MultiSearch(request *model.SearchRequest) *model.SearchResult {
	//等待搜索初始化完成
	e.Wait()
	//分词搜索
	words := e.WordCut(request.Query)

	totalTime := float64(0)

	keys := make([]uint32, len(words))

	for i, word := range words {
		keys[i] = utils.StringToInt(word)
	}

	fastSort := &sorts.FastSort{
		IsDebug: e.IsDebug,
		Keys:    keys,
		Call:    e.getRank,
	}

	_time := utils.ExecTime(func() {

		wg := &sync.WaitGroup{}
		wg.Add(len(keys))

		for _, key := range keys {
			go e.processKeySearch(key, fastSort, wg)
		}
		wg.Wait()
	})
	if e.IsDebug {
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
		var resultItems []model.SliceItem
		_tt := utils.ExecTime(func() {
			resultItems = fastSort.GetAll(request.Order)
		})

		if e.IsDebug {
			log.Println("处理排序耗时", _tt, "ms")
		}

		pager.Init(request.Limit, len(resultItems))
		//设置总页数
		result.PageCount = pager.PageCount

		//读取单页的id
		if pager.PageCount != 0 {

			start, end := pager.GetPage(request.Page)
			items := resultItems[start:end]

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
	if e.IsDebug {
		log.Println("处理数据耗时：", _time, "ms")
	}

	return result
}

func (e *Engine) processKeySearch(key uint32, fastSort *sorts.FastSort, wg *sync.WaitGroup) {
	defer wg.Done()

	shard := e.getShard(key)
	//读取id
	iis := e.InvertedIndexStorages[shard]
	kv := utils.Uint32ToBytes(key)

	bitmap := roaring.New()
	buf, find := iis.Get(kv)
	if find {
		//解码
		_, err := bitmap.FromBuffer(buf)
		if err != nil {
			panic(err)
		}

	}
	fastSort.Add(bitmap)

}

func (e *Engine) getRank(keys []uint32, id uint32) float32 {
	shard := e.getShard(id)
	iks := e.PositiveIndexStorages[shard]
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

func (e *Engine) GetIndexSize() int64 {
	var size int64
	for i := 0; i < e.Option.Shard; i++ {
		size += e.InvertedIndexStorages[i].Size()
	}
	return size
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
	//InvertedIndexStorages []*storage.LeveldbStorage
	//ID和key映射，用于计算相关度，一个id 对应多个key
	ik := e.PositiveIndexStorages[shard]
	keysValue, found := ik.Get(key)
	if !found {
		return errors.New(fmt.Sprintf("没有找到id=%d", id))
	}

	keys := make([]uint32, 0)
	utils.Decoder(keysValue, &keys)

	//符合条件的key，要移除id
	for _, k := range keys {
		kv := utils.Uint32ToBytes(k)
		ks := e.InvertedIndexStorages[e.getShard(k)]
		buf, exists := ks.Get(kv)
		if exists {
			ids := make([]uint32, 0)
			utils.Decoder(buf, &ids)
			//如果存在，才移除
			index := arrays.Find(ids, id)
			if index != -1 {
				ids = utils.DeleteArray(ids, index)
				ks.Set(kv, utils.Encoder(ids))
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

	for i := 0; i < e.Option.Shard; i++ {
		e.InvertedIndexStorages[i].Close()
		e.PositiveIndexStorages[i].Close()
	}
}
