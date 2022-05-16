package controller

import "github.com/gin-gonic/gin"

// DatabaseDrop 删除数据库
func DatabaseDrop(c *gin.Context) {
	dbName := c.Query("database")
	if dbName == "" {
		ResponseErrorWithMsg(c, "database is empty")
		return
	}

	if err := srv.Database.Drop(dbName); err != nil {
		ResponseErrorWithMsg(c, err.Error())
		return
	}

	ResponseSuccessWithData(c, "删除成功")
}

// DatabaseCreate 创建数据库
func DatabaseCreate(c *gin.Context) {
	dbName := c.Query("database")
	if dbName == "" {
		ResponseErrorWithMsg(c, "database is empty")
		return
	}

	srv.Database.Create(dbName)
	ResponseSuccessWithData(c, "创建成功")
}

// DBS 查询数据库
func DBS(c *gin.Context) {
	ResponseSuccessWithData(c, srv.Database.Show())
}
