package web

import (
	"github.com/gin-gonic/gin"
	"gofound/searcher"
	"gofound/searcher/model"
	"gofound/searcher/system"
	"os"
	"runtime"
)

type Api struct {
	Container *searcher.Container
	Callback  func() map[string]interface{}
}

func (a *Api) query(c *gin.Context) {

	var request = &model.SearchRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(200, Error(err.Error()))
		return
	}

	//调用搜索
	r := a.Container.GetDataBase(c.Query("database")).MultiSearch(request)
	c.JSON(200, Success(r))
}

func (a *Api) gc(c *gin.Context) {
	runtime.GC()

	c.JSON(200, Success(nil))
}

// status 获取服务器状态
func (a *Api) status(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	s := a.Callback()

	r := gin.H{
		"memory": system.GetMemStat(),
		"cpu":    system.GetCPUStatus(),
		"disk":   system.GetDiskStat(),
		"system": s,
	}
	// 获取服务器状态
	c.JSON(200, Success(r))
}

func (a *Api) addIndex(c *gin.Context) {
	document := model.IndexDoc{}
	err := c.BindJSON(&document)
	if err != nil {
		c.JSON(200, Error(err.Error()))
		return
	}

	go a.Container.GetDataBase(c.Query("database")).IndexDocument(document)

	c.JSON(200, Success(nil))
}

func (a *Api) batchAddIndex(c *gin.Context) {
	documents := make([]model.IndexDoc, 0)
	err := c.BindJSON(&documents)
	if err != nil {
		c.JSON(200, Error(err.Error()))
		return
	}

	db := a.Container.GetDataBase(c.Query("database"))
	for _, doc := range documents {
		go db.IndexDocument(doc)
	}

	c.JSON(200, Success(nil))
}

func (a *Api) wordCut(c *gin.Context) {
	q := c.Query("q")
	r := a.Container.Tokenizer.Cut(q)
	c.JSON(200, Success(r))

}

func welcome(c *gin.Context) {
	c.JSON(200, Success("Welcome to GoFound"))
}

func (a *Api) removeIndex(c *gin.Context) {
	removeIndexModel := &model.RemoveIndexModel{}
	err := c.BindJSON(&removeIndexModel)
	if err != nil {
		c.JSON(200, Error(err.Error()))
		return
	}
	db := a.Container.GetDataBase(c.Query("database"))

	err = db.RemoveIndex(removeIndexModel.Id)
	if err != nil {
		c.JSON(200, Error(err.Error()))
		return
	}
	c.JSON(200, Success(nil))
}

func (a *Api) dbs(ctx *gin.Context) {
	ctx.JSON(200, Success(a.Container.GetDataBases()))
}

func (a *Api) restart(c *gin.Context) {

	os.Exit(0)
}
func (a *Api) databaseDrop(c *gin.Context) {
	db := c.Query("database")
	if db == "" {
		c.JSON(200, Error("database is empty"))
	} else {
		err := a.Container.DropDataBase(db)
		if err != nil {
			c.JSON(200, Error(err.Error()))
		} else {
			c.JSON(200, Success("删除成功"))
		}
	}
}
func (a *Api) databaseCreate(c *gin.Context) {
	db := c.Query("database")
	if db == "" {
		c.JSON(200, Error("database is empty"))
	} else {
		a.Container.GetDataBase(db)
		c.JSON(200, Success("创建成功"))
	}

}

func (a *Api) Register(router *gin.Engine, handlers ...gin.HandlerFunc) {

	group := router.Group("/api", handlers...)

	group.GET("/", welcome)

	group.POST("/query", a.query)

	group.GET("/status", a.status)

	group.GET("/gc", a.gc)

	group.GET("/db/list", a.dbs)
	group.GET("/db/drop", a.databaseDrop)
	group.GET("/db/create", a.databaseCreate)

	group.GET("/word/cut", a.wordCut)

	group.POST("/index", a.addIndex)

	group.POST("/index/batch", a.batchAddIndex)

	group.POST("/remove", a.removeIndex)

	group.GET("/restart", a.restart)

}
