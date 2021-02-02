package main

import (
	"github.com/flamefatex/config"
	v1 "github.com/flamefatex/example-api/handler/v1"
	"github.com/flamefatex/log"
	"github.com/gin-gonic/gin"
)

// gin
func runGinServer() {
	envGinAddr := config.Config().GetString("gin.addr")
	if envGinAddr != "" {
		ginAddr = envGinAddr
	}

	// handler 初始化
	engine := gin.Default()
	v1Group := engine.Group("/v1")
	v1.NewDownloadHandler(v1Group)

	// gin Run the server
	go func() {
		log.Debugf("start gin server")

		log.Fatal(engine.Run(ginAddr))
	}()
}
