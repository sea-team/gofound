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
	Keys []string `json:"keys,omitempty"`
}

type ResponseDoc struct {
	IndexDoc
	OriginalText string   `json:"originalText,omitempty"`
	Score        int      `json:"score,omitempty"` //得分
	Keys         []string `json:"keys,omitempty"`
}

type RemoveIndexModel struct {
	Id uint32 `json:"id,omitempty"`
}

type ResponseDocSort []ResponseDoc

func (r ResponseDocSort) Len() int {
	return len(r)
}

func (r ResponseDocSort) Less(i, j int) bool {
	return r[i].Score < r[j].Score
}

func (r ResponseDocSort) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
