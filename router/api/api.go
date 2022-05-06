package api

import (
	"github.com/gin-gonic/gin"
	"gofound/router/result"
	"gofound/searcher"
	"gofound/searcher/model"
	"runtime"
	"runtime/debug"
)

var container *searcher.Container

func SetContainer(c *searcher.Container) {
	container = c
}

func query(c *gin.Context) {

	var request = &model.SearchRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(200, result.Error(err.Error()))
		return
	}

	//调用搜索
	r := container.GetDataBase(c.Query("database")).MultiSearch(request)
	c.JSON(200, result.Success(r))
}

func gc(c *gin.Context) {
	runtime.GC()

	c.JSON(200, result.Success(nil))
}

// status 获取服务器状态
func status(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memory := map[string]any{
		"alloc":         m.Alloc,
		"total":         m.TotalAlloc,
		"sys":           m.Sys,
		"heap":          m.HeapAlloc,
		"heap_sys":      m.HeapSys,
		"heap_idle":     m.HeapIdle,
		"heap_inuse":    m.HeapInuse,
		"heap_released": m.HeapReleased,
		"heap_objects":  m.HeapObjects,
	}
	system := &map[string]any{
		"os":      runtime.GOOS,
		"arch":    runtime.GOARCH,
		"cores":   runtime.NumCPU(),
		"version": runtime.Version(),
	}

	r := gin.H{
		"memory": memory,
		"system": system,
		"status": "ok",
	}
	// 获取服务器状态
	c.JSON(200, result.Success(r))
}

func addIndex(c *gin.Context) {
	document := model.IndexDoc{}
	err := c.BindJSON(&document)
	if err != nil {
		c.JSON(200, result.Error(err.Error()))
		return
	}

	go container.GetDataBase(c.Query("database")).IndexDocument(document)

	c.JSON(200, result.Success(nil))
}

func batchAddIndex(c *gin.Context) {
	documents := make([]model.IndexDoc, 0)
	err := c.BindJSON(&documents)
	if err != nil {
		c.JSON(200, result.Error(err.Error()))
		return
	}

	db := container.GetDataBase(c.Query("database"))
	for _, doc := range documents {
		go db.IndexDocument(doc)
	}

	c.JSON(200, result.Success(nil))
}

// dump 持久化到磁盘
func dump(c *gin.Context) {

	c.JSON(200, result.Error("The interface has been cancelled!"))
}

func wordCut(c *gin.Context) {
	q := c.Query("q")
	r := container.Tokenizer.Cut(q)
	c.JSON(200, result.Success(r))

}

func welcome(c *gin.Context) {
	c.JSON(200, result.Success("Welcome to GoFound"))
}

func removeIndex(c *gin.Context) {
	removeIndexModel := &model.RemoveIndexModel{}
	err := c.BindJSON(&removeIndexModel)
	if err != nil {
		c.JSON(200, result.Error(err.Error()))
		return
	}
	db := container.GetDataBase(c.Query("database"))

	err = db.RemoveIndex(removeIndexModel.Id)
	if err != nil {
		c.JSON(200, result.Error(err.Error()))
		return
	}
	c.JSON(200, result.Success(nil))
}

func dbs(ctx *gin.Context) {
	ctx.JSON(200, result.Success(container.GetDataBases()))
}

//Recover 处理异常
func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			c.JSON(200, result.Error(err.(error).Error()))
		}
		c.Abort()
	}()
	c.Next()
}

func Register(router *gin.Engine) {

	router.GET("/api/", welcome)

	router.GET("/api/dbs", dbs)

	router.POST("/api/query", query)

	router.GET("/api/status", status).POST("/api/status", status)

	router.GET("/api/gc", gc).POST("/api/gc", gc)

	router.GET("/api/word/cut", wordCut)

	router.GET("/api/dump", dump).POST("/api/dump", dump)

	router.POST("/api/index", addIndex)

	router.POST("/api/index/batch", batchAddIndex)

	router.POST("/api/remove", removeIndex)

}
