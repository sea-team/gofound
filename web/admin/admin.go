package admin

import (
	"github.com/gin-gonic/gin"
	"gofound/web/admin/assets"
	"net/http"
	"net/url"
	"os"
)

func adminIndex(ctx *gin.Context) {
	file, err := assets.Static.ReadFile("web/dist/index.html")
	if err != nil && os.IsNotExist(err) {
		ctx.String(http.StatusNotFound, "not found")
		return
	}
	ctx.Data(http.StatusOK, "text/html", file)
}

func handlerStatic(c *gin.Context) {
	staticServer := http.FileServer(http.FS(assets.Static))
	c.Request.URL = &url.URL{Path: "web/dist" + c.Request.RequestURI}
	staticServer.ServeHTTP(c.Writer, c.Request)
}

func Register(router *gin.Engine, handlers ...gin.HandlerFunc) {
	//注册路由
	r := router.Group("/admin", handlers...)
	r.GET("/", adminIndex)
	router.GET("/assets/*filepath", handlerStatic)
}
