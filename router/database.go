package router

import (
	"gofound/controller"

	"github.com/gin-gonic/gin"
)

// InitDatabasegoRouter 数据库路由
func InitDatabasegoRouter(Router *gin.RouterGroup) {

	databaseRouter := Router.Group("db")
	{
		databaseRouter.GET("list", controller.DBS)              // 查看数据库
		databaseRouter.GET("drop", controller.DatabaseDrop)     // 删除数据库
		databaseRouter.GET("create", controller.DatabaseCreate) // 添加数据库
	}
}
