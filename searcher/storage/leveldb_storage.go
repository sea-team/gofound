package storage

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
	"sync"
	"time"
)

// LeveldbStorage TODO 要支持事务
type LeveldbStorage struct {
	db       *leveldb.DB
	path     string
	mu       sync.RWMutex //加锁
	closed   bool
	timeout  int64
	lastTime int64
	count    int64
}

func (s *LeveldbStorage) autoOpenDB() {
	if s.isClosed() {
		s.ReOpen()
	}
	s.lastTime = time.Now().Unix()
}

// NewStorage 打开数据库
func NewStorage(path string, timeout int64) (*LeveldbStorage, error) {

	db := &LeveldbStorage{
		path:     path,
		closed:   true,
		timeout:  timeout,
		lastTime: time.Now().Unix(),
	}

	go db.task()

	return db, nil
}

func (s *LeveldbStorage) task() {
	if s.timeout == -1 {
		//不检查
		return
	}
	for {

		if !s.isClosed() && time.Now().Unix()-s.lastTime > s.timeout {
			s.Close()
			//log.Println("leveldb storage timeout", s.path)
		}

		time.Sleep(time.Duration(5) * time.Second)

	}
}

func openDB(path string) (*leveldb.DB, error) {

	////使用布隆过滤器
	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}

	db, err := leveldb.OpenFile(path, o)
	return db, err
}
func (s *LeveldbStorage) ReOpen() {
	if !s.isClosed() {
		log.Println("db is not closed")
		return
	}
	s.mu.Lock()
	db, err := openDB(s.path)
	if err != nil {
		panic(err)
	}
	s.db = db
	s.closed = false
	s.mu.Unlock()
	//计算总条数
	go s.compute()
}

func (s *LeveldbStorage) Get(key []byte) ([]byte, bool) {
	s.autoOpenDB()
	buffer, err := s.db.Get(key, nil)
	if err != nil {
		return nil, false
	}
	return buffer, true
}

func (s *LeveldbStorage) Has(key []byte) bool {
	s.autoOpenDB()
	has, err := s.db.Has(key, nil)
	if err != nil {
		panic(err)
	}
	return has
}

func (s *LeveldbStorage) Set(key []byte, value []byte) {
	s.autoOpenDB()
	err := s.db.Put(key, value, nil)
	if err != nil {
		panic(err)
	}
}

// Delete 删除
func (s *LeveldbStorage) Delete(key []byte) error {
	s.autoOpenDB()
	return s.db.Delete(key, nil)
}

// Close 关闭
func (s *LeveldbStorage) Close() error {
	if s.isClosed() {
		return nil
	}
	s.mu.Lock()
	err := s.db.Close()
	if err != nil {
		return err
	}
	s.closed = true
	s.mu.Unlock()
	return nil
}

func (s *LeveldbStorage) isClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.closed
}

func (s *LeveldbStorage) compute() {
	var count int64
	iter := s.db.NewIterator(nil, nil)
	for iter.Next() {
		count++
	}
	iter.Release()
	s.count = count
}

func (s *LeveldbStorage) GetCount() int64 {
	if s.count == 0 && s.isClosed() {
		s.ReOpen()
		s.compute()
	}
	return s.count
}
