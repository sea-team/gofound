package searcher

import (
	"fmt"
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
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

type Engine struct {
	IndexPath string  //索引文件存储目录
	Option    *Option //配置

	invertedIndexStorages []*storage.LeveldbStorage //关键字和Id映射，倒排索引,key=id,value=[]words
	positiveIndexStorages []*storage.LeveldbStorage //ID和key映射，用于计算相关度，一个id 对应多个key，正排索引
	docStorages           []*storage.LeveldbStorage //文档仓

	sync.Mutex                                   //锁
	sync.WaitGroup                               //等待
	addDocumentWorkerChan []chan *model.IndexDoc //添加索引的通道
	IsDebug               bool                   //是否调试模式
	Tokenizer             *words.Tokenizer       //分词器
	DatabaseName          string                 //数据库名

	Shard     int   //分片数
	Timeout   int64 //超时时间,单位秒
	BufferNum int   //分片缓冲数

	documentCount int64 //文档总数量
}

type Option struct {
	InvertedIndexName string //倒排索引
	PositiveIndexName string //正排索引
	DocIndexName      string //文档存储
}

// Init 初始化索引引擎
func (e *Engine) Init() {
	e.Add(1)
	defer e.Done()

	if e.Option == nil {
		e.Option = e.GetOptions()
	}
	if e.Timeout == 0 {
		e.Timeout = 10 * 3 // 默认30s
	}
	//-1代表没有初始化
	e.documentCount = -1
	//log.Println("数据存储目录：", e.IndexPath)
	log.Println("chain num:", e.Shard*e.BufferNum)

	e.addDocumentWorkerChan = make([]chan *model.IndexDoc, e.Shard)
	//初始化文件存储
	for shard := 0; shard < e.Shard; shard++ {

		//初始化chan
		worker := make(chan *model.IndexDoc, e.BufferNum)
		e.addDocumentWorkerChan[shard] = worker

		//初始化chan
		go e.DocumentWorkerExec(worker)

		s, err := storage.NewStorage(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.DocIndexName, shard)), e.Timeout)
		if err != nil {
			panic(err)
		}
		e.docStorages = append(e.docStorages, s)

		//初始化Keys存储
		ks, kerr := storage.NewStorage(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.InvertedIndexName, shard)), e.Timeout)
		if kerr != nil {
			panic(err)
		}
		e.invertedIndexStorages = append(e.invertedIndexStorages, ks)

		//id和keys映射
		iks, ikerr := storage.NewStorage(e.getFilePath(fmt.Sprintf("%s_%d", e.Option.PositiveIndexName, shard)), e.Timeout)
		if ikerr != nil {
			panic(ikerr)
		}
		e.positiveIndexStorages = append(e.positiveIndexStorages, iks)
	}
	go e.automaticGC()
	//log.Println("初始化完成")
}

// 自动保存索引，10秒钟检测一次
func (e *Engine) automaticGC() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		<-ticker.C
		//定时GC
		runtime.GC()
		if e.IsDebug {
			log.Println("waiting:", e.GetQueue())
		}
	}
}

func (e *Engine) IndexDocument(doc *model.IndexDoc) error {
	//数量增加
	e.documentCount++
	e.addDocumentWorkerChan[e.getShard(doc.Id)] <- doc
	return nil
	/*
		select {
		case e.addDocumentWorkerChan[e.getShard(doc.Id)] <- doc:
			e.documentCount++
		default:
			return errors.New("处理缓冲已满")
		}
		return nil
	*/
}

// GetQueue 获取队列剩余
func (e *Engine) GetQueue() int {
	total := 0
	for _, v := range e.addDocumentWorkerChan {
		total += len(v)
	}
	return total
}

// DocumentWorkerExec 添加文档队列
func (e *Engine) DocumentWorkerExec(worker chan *model.IndexDoc) {
	for {
		doc := <-worker
		e.AddDocument(doc)
	}
}

// getShard 计算索引分布在哪个文件块
func (e *Engine) getShard(id uint32) int {
	return int(id % uint32(e.Shard))
}

func (e *Engine) getShardByWord(word string) int {

	return int(utils.StringToInt(word) % uint32(e.Shard))
}

func (e *Engine) InitOption(option *Option) {

	if option == nil {
		//默认值
		option = e.GetOptions()
	}
	e.Option = option
	//shard默认值
	if e.Shard <= 0 {
		e.Shard = 10
	}
	if e.BufferNum <= 0 {
		e.BufferNum = 1000
	}
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
	}
}

// AddDocument 分词索引
func (e *Engine) AddDocument(index *model.IndexDoc) {
	//等待初始化完成
	e.Wait()
	text := index.Text

	splitWords := e.Tokenizer.Cut(text)

	id := index.Id
	// 检查是否需要更新倒排索引 words变更/id不存在
	inserts, needUpdateInverted := e.optimizeIndex(id, splitWords)

	// 将新增的word剔出单独处理，减少I/O操作
	if needUpdateInverted {
		for _, word := range inserts {
			e.addInvertedIndex(word, id)
		}
	}

	// TODO: 是否需要更新正排索引 - 检测document变更
	e.addPositiveIndex(index, splitWords)
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

	if !arrays.ArrayUint32Exists(ids, id) {
		ids = append(ids, id)
	}

	s.Set(key, utils.Encoder(ids))
}

// 移除删去的词
func (e *Engine) optimizeIndex(id uint32, newWords []string) ([]string, bool) {
	// 判断id是否存在
	e.Lock()
	defer e.Unlock()

	// 计算差值
	removes, inserts, changed := e.getDifference(id, newWords)
	if changed {
		if removes != nil && len(removes) > 0 {
			// 移除正排索引
			for _, word := range removes {
				e.removeIdInWordIndex(id, word)
			}
		}
	}
	return inserts, changed
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
// @return []string: 需要删除的词
// @return bool    : words出现变更返回true，否则返回false
func (e *Engine) getDifference(id uint32, newWords []string) ([]string, []string, bool) {
	shard := e.getShard(id)
	wordStorage := e.positiveIndexStorages[shard]
	key := utils.Uint32ToBytes(id)
	buf, found := wordStorage.Get(key)
	if found {
		oldWords := make([]string, 0)
		utils.Decoder(buf, &oldWords)

		// 计算需要移除的
		removes := make([]string, 0)
		for _, word := range oldWords {
			// 旧的在新的里面不存在，就是需要移除的
			if !arrays.ArrayStringExists(newWords, word) {
				removes = append(removes, word)
			}
		}
		// 计算需要新增的
		inserts := make([]string, 0)
		for _, word := range newWords {
			if !arrays.ArrayStringExists(oldWords, word) {
				inserts = append(inserts, word)
			}
		}
		if len(removes) != 0 || len(inserts) != 0 {
			return removes, inserts, true
		}
		// 没有改变
		return removes, inserts, false
	}
	// id不存在，相当于insert
	return nil, newWords, true
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
func (e *Engine) MultiSearch(request *model.SearchRequest) (*model.SearchResult, error) {
	//等待搜索初始化完成
	e.Wait()

	//分词搜索
	words := e.Tokenizer.Cut(request.Query)

	fastSort := &sorts.FastSort{
		IsDebug: e.IsDebug,
		Order:   request.Order,
	}

	_time := utils.ExecTime(func() {

		base := len(words)
		wg := &sync.WaitGroup{}
		wg.Add(base)

		for _, word := range words {
			go e.processKeySearch(word, fastSort, wg)
		}
		wg.Wait()
	})
	if e.IsDebug {
		log.Println("搜索时间:", _time, "ms")
	}
	// 处理分页
	request = request.GetAndSetDefault()

	//计算交集得分和去重
	fastSort.Process()

	wordMap := make(map[string]bool)
	for _, word := range words {
		wordMap[word] = true
	}

	//读取文档
	var result = &model.SearchResult{
		Total: fastSort.Count(),
		Page:  request.Page,
		Limit: request.Limit,
		Words: words,
	}

	t, err := utils.ExecTimeWithError(func() error {

		pager := new(pagination.Pagination)

		pager.Init(request.Limit, fastSort.Count())
		//设置总页数
		result.PageCount = pager.PageCount

		//读取单页的id
		if pager.PageCount != 0 {

			start, end := pager.GetPage(request.Page)
			if request.ScoreExp != "" {
				// 分数表达式不为空,获取所有的数据
				start, end = 0, pager.Total
			}

			var resultItems = make([]model.SliceItem, 0)
			fastSort.GetAll(&resultItems, start, end)

			count := len(resultItems)

			result.Documents = make([]model.ResponseDoc, count)
			//只读取前面100个
			wg := new(sync.WaitGroup)
			wg.Add(count)
			for index, item := range resultItems {
				go e.getDocument(item, &result.Documents[index], request, &wordMap, wg)
			}
			wg.Wait()
			if request.ScoreExp != "" {
				// 生成计算表达式
				exp, err := govaluate.NewEvaluableExpression(request.ScoreExp)
				if err != nil {
					return err
				}
				parameters := make(map[string]interface{})
				// 根据表达式计算分数
				for i, doc := range result.Documents {
					parameters["score"] = doc.Score
					for k, v := range doc.Document {
						parameters["document."+k] = v
					}
					val, err := exp.Evaluate(parameters)
					if err != nil {
						log.Printf("表达式执行'%v'错误: %v 值内容: %v", request.ScoreExp, err, parameters)
					} else {
						result.Documents[i].Score = int(val.(float64))
					}
				}
				if request.Order == "desc" {
					sort.Sort(sort.Reverse(model.ResponseDocSort(result.Documents)))
				} else {
					sort.Sort(model.ResponseDocSort(result.Documents))
				}
				// 取出page
				start, end := pager.GetPage(request.Page)
				result.Documents = result.Documents[start:end]
			}
		}
		return nil
	})
	if e.IsDebug {
		log.Println("处理数据耗时：", _time, "ms")
	}
	if err != nil {
		return nil, err
	}
	result.Time = _time + t

	return result, nil
}

func (e *Engine) getDocument(item model.SliceItem, doc *model.ResponseDoc, request *model.SearchRequest, wordMap *map[string]bool, wg *sync.WaitGroup) {
	buf := e.GetDocById(item.Id)
	defer wg.Done()
	doc.Score = item.Score

	if buf != nil {
		//gob解析
		storageDoc := new(model.StorageIndexDoc)
		utils.Decoder(buf, &storageDoc)
		doc.Document = storageDoc.Document
		doc.Keys = storageDoc.Keys
		text := storageDoc.Text
		//处理关键词高亮
		highlight := request.Highlight
		if highlight != nil {
			//全部小写
			text = strings.ToLower(text)
			//还可以优化，只替换击中的词
			for _, key := range storageDoc.Keys {
				if ok := (*wordMap)[key]; ok {
					text = strings.ReplaceAll(text, key, fmt.Sprintf("%s%s%s", highlight.PreTag, key, highlight.PostTag))
				}
			}
			//放置原始文本
			doc.OriginalText = storageDoc.Text
		}
		doc.Text = text
		doc.Id = item.Id

	}

}

func (e *Engine) processKeySearch(word string, fastSort *sorts.FastSort, wg *sync.WaitGroup) {
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
		fastSort.Add(&ids)
	}

}

// GetIndexCount 获取索引数量
func (e *Engine) GetIndexCount() int64 {
	var size int64
	for i := 0; i < e.Shard; i++ {
		size += e.invertedIndexStorages[i].GetCount()
	}
	return size
}

// GetDocumentCount 获取文档数量
func (e *Engine) GetDocumentCount() int64 {
	if e.documentCount == -1 {
		var count int64
		//使用多线程加速统计
		wg := sync.WaitGroup{}
		wg.Add(e.Shard)
		//这里的统计可能会出现数据错误，因为没加锁
		for i := 0; i < e.Shard; i++ {
			go func(i int) {
				count += e.docStorages[i].GetCount()
				wg.Done()
			}(i)
		}
		wg.Wait()
		e.documentCount = count
	}

	return e.documentCount
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
	//减少数量
	e.documentCount--

	return nil
}

func (e *Engine) Close() {
	e.Lock()
	defer e.Unlock()

	for i := 0; i < e.Shard; i++ {
		e.invertedIndexStorages[i].Close()
		e.positiveIndexStorages[i].Close()
	}
}

// Drop 删除
func (e *Engine) Drop() error {
	e.Lock()
	defer e.Unlock()
	//删除文件
	if err := os.RemoveAll(e.IndexPath); err != nil {
		return err
	}

	//清空内存
	for i := 0; i < e.Shard; i++ {
		e.docStorages = make([]*storage.LeveldbStorage, 0)
		e.invertedIndexStorages = make([]*storage.LeveldbStorage, 0)
		e.positiveIndexStorages = make([]*storage.LeveldbStorage, 0)
	}

	return nil
}
