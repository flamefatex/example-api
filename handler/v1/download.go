package v1

import (
	"os"

	"github.com/flamefatex/config"
	"github.com/flamefatex/log"
	"github.com/gin-gonic/gin"
)

type downloadHandler struct {
}

func NewDownloadHandler(routerGroup *gin.RouterGroup) {

	hdr := &downloadHandler{}
	router := routerGroup.Group("/download")
	router.GET("/:file_type/:date/:filename", hdr.Download)
}

func (hdr *downloadHandler) Download(ctx *gin.Context) {
	file_type := ctx.Param("file_type")
	date := ctx.Param("date")
	filename := ctx.Param("filename")

	rootDir := config.Config().GetString("app.data_root_dir")
	filePath := rootDir + "/" + file_type + "/" + date + "/" + filename

	_, err := os.Stat(filePath)
	if err != nil {
		log.Warnf("找不到文件,filePath:%s,err:%s", filePath, err)
		ctx.Status(404)
		return
	}
	ctx.File(filePath)
	return
}
