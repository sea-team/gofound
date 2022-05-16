package controller

import (
	"github.com/gin-gonic/gin"
	"gofound/searcher/model"
)

func Welcome(c *gin.Context) {
	ResponseSuccessWithData(c, "Welcome to GoFound")
}

// Query 查询
func Query(c *gin.Context) {
	var request = &model.SearchRequest{
		Database: c.Query("database"),
	}
	if err := c.ShouldBind(&request); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}
	//调用搜索
	r := srv.Base.Query(request)
	ResponseSuccessWithData(c, r)
}

// GC 释放GC
func GC(c *gin.Context) {
	srv.Base.GC()
	ResponseSuccess(c)
}

// Status 获取服务器状态
func Status(c *gin.Context) {
	r := srv.Base.Status()
	ResponseSuccessWithData(c, r)
}
