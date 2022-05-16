package service

import (
	"gofound/global"
	"gofound/searcher"
)

type Database struct {
	Container *searcher.Container
}

func NewDatabase() *Database {
	return &Database{
		Container: global.Container,
	}
}

// Show 查看数据库
func (d *Database) Show() map[string]*searcher.Engine {
	return d.Container.GetDataBases()
}

// Drop 删除数据库
func (d *Database) Drop(dbName string) error {
	if err := d.Container.DropDataBase(dbName); err != nil {
		return err
	}
	return nil
}

// Create 创建数据库
func (d *Database) Create(dbName string) *searcher.Engine {
	return d.Container.GetDataBase(dbName)
}
