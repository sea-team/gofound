package storage

import "github.com/syndtr/goleveldb/leveldb"

type LeveldbStorage struct {
	db   *leveldb.DB
	path string
}

func Open(path string) (*LeveldbStorage, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LeveldbStorage{
		db:   db,
		path: path,
	}, nil
}

func (s *LeveldbStorage) Get(key []byte) ([]byte, bool) {

	buffer, err := s.db.Get(key, nil)
	if err != nil {
		return nil, false
	}
	return buffer, true
}

func (s *LeveldbStorage) Set(key []byte, value []byte) error {
	return s.db.Put(key, value, nil)
}

// Delete 删除
func (s *LeveldbStorage) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

// Close 关闭
func (s *LeveldbStorage) Close() error {
	return s.db.Close()
}
