package router

import (
	"gofound/web/controller"

	"github.com/gin-gonic/gin"
)

// InitBaseRouter 基础管理路由
func InitBaseRouter(Router *gin.RouterGroup) {

	BaseRouter := Router.Group("")
	{
		BaseRouter.GET("/", controller.Welcome)
		BaseRouter.POST("query", controller.Query)
		BaseRouter.GET("status", controller.Status)
		BaseRouter.GET("gc", controller.GC)

	}
}
