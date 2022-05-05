package model

// IndexDoc 索引实体
type IndexDoc struct {
	Database
	Id       uint32                 `json:"id,omitempty"`
	Text     string                 `json:"text,omitempty"`
	Document map[string]interface{} `json:"document,omitempty"`
}

// StorageIndexDoc 文档对象
type StorageIndexDoc struct {
	*IndexDoc
	Keys []string `json:"keys,omitempty"`
}

type ResponseDoc struct {
	IndexDoc
	Score int `json:"score,omitempty"` //得分
}

type RemoveIndexModel struct {
	Database
	Id uint32 `json:"id,omitempty"`
}
