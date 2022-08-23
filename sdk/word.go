package gofound

// WordCut 分词
func (c *Client) WordCut(keyword string) []string {
	return c.container.Tokenizer.Cut(keyword)
}

// BatchWordCut 批量分词
func (c *Client) BatchWordCut(keywords []string) *[][]string {
	res := make([][]string, len(keywords))
	for _, w := range keywords {
		res = append(res, c.container.Tokenizer.Cut(w))
	}
	return &res
}
