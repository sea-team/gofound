package router

import (
	"github.com/sea-team/gofound/web/controller"

	"github.com/gin-gonic/gin"
)

// InitWordRouter 分词路由
func InitWordRouter(Router *gin.RouterGroup) {

	wordRouter := Router.Group("word")
	{
		wordRouter.GET("cut", controller.WordCut)
	}
}
