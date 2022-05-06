package searcher

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"gofound/searcher/arrays"
	"gofound/searcher/model"
	"gofound/searcher/pagination"
	"gofound/searcher/sorts"
	"gofound/searcher/storage"
	"gofound/searcher/utils"
	"gofound/searcher/words"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type Engine struct {
	IndexPath string  //索引文件存储目录
	Option    *Option //配置

	invertedIndexStorages []*storage.LeveldbStorage //关键字和Id映射，倒排索引,key=id,value=[]words
	positiveIndexStorages []*storage.LeveldbStorage //ID和key映射，用于计算相关度，一个id 对应多个key，正排索引
	docStorages           []*storage.LeveldbStorage //文档仓

	sync.Mutex                                  //锁
	sync.WaitGroup                              //等待
	addDocumentWorkerChan []chan model.IndexDoc //添加索引的通道
	IsDebug               bool                  //是否调试模式
	Tokenizer             *words.Tokenizer      //分词器
	DatabaseName          string                //数据库名
}

type Option struct {
	InvertedIndexName string //倒排索引
	PositiveIndexName string //正排索引
	DocIndexName      string //文档存储
	Shard             int    //分片数，默认为5
}

// Init 初始化索引引擎
func (e *Engine) Init() {
	e.Add(1)
	defer e.Done()

	if e.Option == nil {
		e.Option = e.GetOptions()
	}
	log.Println("数据存储目录：", e.IndexPath)

	e.addDocumentWorkerChan = make([]chan model.IndexDoc, e.Option.Shard)
	//初始化文件存储
	for shard := 0; shard < e.Option.Shard; shard++ {

		//初始化chan
		worker := make(chan model.IndexDoc, 1000)
		e.addDocumentWorkerChan[shard] = worker

		//初始化chan
		go e.DocumentWorkerExec(worker)

		s, err := storage.Open(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.DocIndexName, shard)))
		if err != nil {
			panic(err)
		}
		e.docStorages = append(e.docStorages, s)

		//初始化Keys存储
		ks, kerr := storage.Open(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.InvertedIndexName, shard)))
		if kerr != nil {
			panic(err)
		}
		e.invertedIndexStorages = append(e.invertedIndexStorages, ks)

		//id和keys映射
		iks, ikerr := storage.Open(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.PositiveIndexName, shard)))
		if ikerr != nil {
			panic(ikerr)
		}
		e.positiveIndexStorages = append(e.positiveIndexStorages, iks)
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
	//根据ID来判断，使用多线程，提速
	e.addDocumentWorkerChan[e.getShard(doc.Id)] <- doc
}

// DocumentWorkerExec 添加文档队列
func (e *Engine) DocumentWorkerExec(worker chan model.IndexDoc) {
	for {
		doc := <-worker
		e.AddDocument(&doc)
	}
}

// getShard 计算索引分布在哪个文件块
func (e *Engine) getShard(id uint32) int {
	return int(id % uint32(e.Option.Shard))
}

func (e *Engine) getShardByWord(word string) int {

	return int(utils.StringToInt(word) % uint32(e.Option.Shard))
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
		DocIndexName:      "docs",
		InvertedIndexName: "inverted_index",
		PositiveIndexName: "positive_index",
		Shard:             5,
	}
}

// AddDocument 分词索引
func (e *Engine) AddDocument(index *model.IndexDoc) {
	//等待初始化完成
	e.Wait()
	text := index.Text

	words := e.Tokenizer.Cut(text)

	//id对应的词

	//判断ID是否存在，如果存在，需要计算两次的差值，然后更新
	id := index.Id
	isUpdate := e.optimizeIndex(id, words)

	//没有更新
	if !isUpdate {
		return
	}

	for _, word := range words {
		e.addInvertedIndex(word, id)
	}

	//添加id索引
	e.addPositiveIndex(index, words)
}

// 添加倒排索引
func (e *Engine) addInvertedIndex(word string, id uint32) {
	e.Lock()
	defer e.Unlock()

	shard := e.getShardByWord(word)

	s := e.invertedIndexStorages[shard]

	//string作为key
	key := []byte(word)

	//存在
	//添加到列表
	buf, find := s.Get(key)
	ids := make([]uint32, 0)
	if find {
		utils.Decoder(buf, &ids)
	}

	if !arrays.BinarySearch(ids, id) {
		ids = append(ids, id)
	}

	s.Set(key, utils.Encoder(ids))
}

//	移除没有的词
func (e *Engine) optimizeIndex(id uint32, newWords []string) bool {
	//判断id是否存在
	e.Lock()
	defer e.Unlock()

	//计算差值
	removes, found := e.getDifference(id, newWords)
	if found && len(removes) > 0 {
		//从这些词中移除当前ID
		for _, word := range removes {
			e.removeIdInWordIndex(id, word)
		}
	}

	// 有没有更新
	return !found || len(removes) > 0

}

func (e *Engine) removeIdInWordIndex(id uint32, word string) {

	shard := e.getShardByWord(word)

	wordStorage := e.invertedIndexStorages[shard]

	//string作为key
	key := []byte(word)

	buf, found := wordStorage.Get(key)
	if found {
		ids := make([]uint32, 0)
		utils.Decoder(buf, &ids)

		//移除
		index := arrays.Find(ids, id)
		if index != -1 {
			ids = utils.DeleteArray(ids, index)
			if len(ids) == 0 {
				err := wordStorage.Delete(key)
				if err != nil {
					panic(err)
				}
			} else {
				wordStorage.Set(key, utils.Encoder(ids))
			}
		}
	}

}

// 计算差值
func (e *Engine) getDifference(id uint32, newWords []string) ([]string, bool) {

	shard := e.getShard(id)
	wordStorage := e.positiveIndexStorages[shard]
	key := utils.Uint32ToBytes(id)
	buf, found := wordStorage.Get(key)
	if found {
		oldWords := make([]string, 0)
		utils.Decoder(buf, &oldWords)

		//计算需要移除的
		removes := make([]string, 0)
		for _, word := range oldWords {

			//旧的在新的里面不存在，就是需要移除的
			if !arrays.ArrayStringExists(newWords, word) {
				removes = append(removes, word)
			}
		}
		return removes, true
	}

	return nil, false
}

// 添加正排索引 id=>keys id=>doc
func (e *Engine) addPositiveIndex(index *model.IndexDoc, keys []string) {
	e.Lock()
	defer e.Unlock()

	key := utils.Uint32ToBytes(index.Id)
	shard := e.getShard(index.Id)
	docStorage := e.docStorages[shard]

	//id和key的映射
	positiveIndexStorage := e.positiveIndexStorages[shard]

	doc := &model.StorageIndexDoc{
		IndexDoc: index,
		Keys:     keys,
	}

	//存储id和key以及文档的映射
	docStorage.Set(key, utils.Encoder(doc))

	//设置到id和key的映射中
	positiveIndexStorage.Set(key, utils.Encoder(keys))
}

// MultiSearch 多线程搜索
func (e *Engine) MultiSearch(request *model.SearchRequest) *model.SearchResult {
	//等待搜索初始化完成
	e.Wait()
	//分词搜索
	words := e.Tokenizer.Cut(request.Query)

	totalTime := float64(0)

	fastSort := &sorts.FastSort{
		IsDebug: e.IsDebug,
	}

	_time := utils.ExecTime(func() {

		base := len(words)
		wg := &sync.WaitGroup{}
		wg.Add(base)

		for _, word := range words {
			go e.processKeySearch(word, fastSort, wg, base)
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
		Page:  request.Page,
		Limit: request.Limit,
		Words: words,
	}

	_time += utils.ExecTime(func() {

		pager := new(pagination.Pagination)
		var resultItems []model.SliceItem
		_tt := utils.ExecTime(func() {
			resultItems = fastSort.GetAll(request.Order)
		})

		if e.IsDebug {
			log.Println("处理排序耗时", _tt, "ms")
		}

		pager.Init(request.Limit, fastSort.Count())
		//设置总页数
		result.PageCount = pager.PageCount

		//读取单页的id
		if pager.PageCount != 0 {

			start, end := pager.GetPage(request.Page)
			items := resultItems[start:end]
			count := len(items)
			result.Documents = make([]model.ResponseDoc, count)
			//只读取前面100个
			wg := new(sync.WaitGroup)
			wg.Add(count)
			for index, item := range items {

				go e.getDocument(item, &result.Documents[index], request, &words, wg)
			}
			wg.Wait()
		}
	})
	if e.IsDebug {
		log.Println("处理数据耗时：", _time, "ms")
	}
	result.Time = _time

	return result
}

func (e *Engine) getDocument(item model.SliceItem, doc *model.ResponseDoc, request *model.SearchRequest, words *[]string, wg *sync.WaitGroup) {
	buf := e.GetDocById(item.Id)
	defer wg.Done()
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
			//还可以优化，只替换击中的词
			for _, word := range *words {
				text = strings.ReplaceAll(text, word, fmt.Sprintf("%s%s%s", highlight.PreTag, word, highlight.PostTag))
			}
		}
		doc.Text = text
		doc.Id = item.Id

	}

}

func (e *Engine) processKeySearch(word string, fastSort *sorts.FastSort, wg *sync.WaitGroup, base int) {
	defer wg.Done()

	shard := e.getShardByWord(word)
	//读取id
	invertedIndexStorage := e.invertedIndexStorages[shard]
	key := []byte(word)

	buf, find := invertedIndexStorage.Get(key)
	if find {
		ids := make([]uint32, 0)
		//解码
		utils.Decoder(buf, &ids)
		//ids越多，说明这个词频越高，这个词越重要
		frequency := (len(ids) % base) + 1
		fastSort.Add(ids, frequency)
	}

}

func (e *Engine) GetIndexSize() int64 {
	var size int64
	for i := 0; i < e.Option.Shard; i++ {
		size += e.invertedIndexStorages[i].Size()
	}
	return size
}

// GetDocById 通过id获取文档
func (e *Engine) GetDocById(id uint32) []byte {
	shard := e.getShard(id)
	key := utils.Uint32ToBytes(id)
	buf, found := e.docStorages[shard].Get(key)
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
	//invertedIndexStorages []*storage.LeveldbStorage
	//ID和key映射，用于计算相关度，一个id 对应多个key
	ik := e.positiveIndexStorages[shard]
	keysValue, found := ik.Get(key)
	if !found {
		return errors.New(fmt.Sprintf("没有找到id=%d", id))
	}

	keys := make([]string, 0)
	utils.Decoder(keysValue, &keys)

	//符合条件的key，要移除id
	for _, word := range keys {
		e.removeIdInWordIndex(id, word)
	}

	//删除id映射
	err := ik.Delete(key)
	if err != nil {
		return errors.New(err.Error())
	}

	//文档仓
	err = e.docStorages[shard].Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) Close() {
	e.Lock()
	defer e.Unlock()

	for i := 0; i < e.Option.Shard; i++ {
		e.invertedIndexStorages[i].Close()
		e.positiveIndexStorages[i].Close()
	}
}
