package controller

import "github.com/gin-gonic/gin"

// WordCut 分词
func WordCut(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		ResponseErrorWithMsg(c, "请输入关键字")
		return
	}
	r := srv.Word.WordCut(q)
	ResponseSuccessWithData(c, r)
}
