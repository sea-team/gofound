package gofound

import (
	"errors"
	"github.com/sea-team/gofound/searcher/model"
)

// AddIndex 添加索引
func (c *Client) AddIndex(dbName string, request *model.IndexDoc) error {
	if request.Text == "" {
		return errors.New("text is empty")
	}
	c.container.GetDataBase(dbName).IndexDocument(request)

	return nil
}

// BatchAddIndex 批次添加索引
func (c *Client) BatchAddIndex(dbName string, documents []*model.IndexDoc) error {
	db := c.container.GetDataBase(dbName)
	// 数据预处理
	for _, doc := range documents {
		if doc.Text == "" {
			return errors.New("text is empty")
		}
		if doc.Document == nil {
			return errors.New("document is empty")
		}
	}
	for _, doc := range documents {
		go db.IndexDocument(doc)
	}
	return nil
}

// RemoveIndex 删除索引
func (c *Client) RemoveIndex(dbName string, data *model.RemoveIndexModel) error {
	db := c.container.GetDataBase(dbName)
	if err := db.RemoveIndex(data.Id); err != nil {
		return err
	}
	return nil
}
