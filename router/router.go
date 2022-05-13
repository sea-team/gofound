package router

import (
	"gofound/global"
	"gofound/middleware"
	"gofound/web/admin"
	"log"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

// SetupRouter 路由管理
func SetupRouter() *gin.Engine {
	if global.CONFIG.System.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	// 启用GZIP压缩
	if global.CONFIG.System.EnableGzip {
		router.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	var handlers []gin.HandlerFunc
	if global.CONFIG.Auth.Enable {
		handlers = append(handlers, gin.BasicAuth(gin.Accounts{global.CONFIG.Auth.Username: global.CONFIG.Auth.Password}))
		log.Println("Enable Auth:", global.CONFIG.Auth.Enable)
	}

	//注册admin
	if global.CONFIG.Auth.Enable {
		admin.Register(router, handlers...)
		log.Printf("Admin Url: \t http://localhost/:%v/admin", global.CONFIG.System.Addr)
	}

	group := router.Group("/api", handlers...)
	// 分组管理 中间件管理
	group.Use(middleware.Cors(), middleware.Exception())
	{
		InitBaseRouter(group)       // 基础管理
		InitIndexRouter(group)      // 索引管理
		InitDatabasegoRouter(group) // 数据库管理
		InitWordRouter(group)       // 分词管理
	}

	return router
}
