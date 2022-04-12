package model

// IndexDoc 索引实体
type IndexDoc struct {
	Id       uint32                 `json:"id,omitempty"`
	Text     string                 `json:"text,omitempty"`
	Document map[string]interface{} `json:"document,omitempty"`
}

// StorageIndexDoc 文档对象
type StorageIndexDoc struct {
	*IndexDoc

	Keys []uint32 `json:"keys,omitempty"`
}

type ResponseDoc struct {
	IndexDoc
	Score float32 `json:"score,omitempty"` //得分
}
