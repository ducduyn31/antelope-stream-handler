package streambiz

import (
	"antelope-stream-chunking/common"
	"antelope-stream-chunking/middlewares"
	"antelope-stream-chunking/pubsub"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (sb StreamBizRouter) ActivateStream(ctx *gin.Context) {
	source := ctx.Request.Header.Get("Source")

	if source == "" {
		_ = ctx.Error(&middlewares.HttpException{
			Message: "source can not be empty",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if !common.StrArrContains(source, pubsub.GetSourceManager().Sources) {
		_ = ctx.Error(&middlewares.HttpException{
			Message: "source can not be found",
			Code:    http.StatusNotFound,
		})
		return
	}

	go pubsub.GetSourceManager().ActivateSource(source)

	ctx.JSON(http.StatusOK, gin.H{})
}
