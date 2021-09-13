package streambiz

import (
	"antelope-stream-chunking/common"
	"antelope-stream-chunking/middlewares"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func (sb StreamBizRouter) GetAll(ctx *gin.Context) {
	lines, f, err := common.OpenFileAndReadLines(fmt.Sprintf("%s/%s", os.Getenv("RTSP_LEDGER"), "stream_sources"))
	if err != nil {
		_ = ctx.Error(&middlewares.HttpException{
			Message: fmt.Sprintf("can not open ledger: %s", err.Error()),
			Code: http.StatusInternalServerError,
		})
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	ctx.JSON(http.StatusOK, gin.H{
		"sources": lines,
	})
}