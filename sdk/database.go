package gofound

import (
	"github.com/sea-team/gofound/searcher"

	"github.com/syndtr/goleveldb/leveldb/errors"
)

// Show 查看数据库
func (c *Client) Show() (map[string]*searcher.Engine, error) {
	// 保持分格一致
	return c.container.GetDataBases(), nil
}

// Drop 删除数据库
func (c *Client) Drop(dbName string) error {
	if dbName == "" {
		return errors.New("database not exist")
	}
	if err := c.container.DropDataBase(dbName); err != nil {
		return err
	}
	return nil
}

// Create 创建数据库
func (c *Client) Create(dbName string) (*searcher.Engine, error) {
	if dbName == "" {
		return nil, errors.New("database name is empty")
	}
	return c.container.GetDataBase(dbName), nil
}
