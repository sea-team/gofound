package router

import (
	"github.com/gin-gonic/gin"
)

func login(ctx *gin.Context) {
	ctx.HTML(200, "login.html", gin.H{})

}

func Register(router *gin.Engine) {
	//注册路由
	router.GET("/admin/login", login)
}
