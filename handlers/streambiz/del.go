package streambiz

import (
	"antelope-stream-chunking/common"
	"antelope-stream-chunking/middlewares"
	"antelope-stream-chunking/pubsub"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

func (sb StreamBizRouter) RemoveStream(ctx *gin.Context) {
	dataBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		_ = ctx.Error(&middlewares.HttpException{
			Message: "request can not be read",
			Code:    http.StatusBadRequest,
		})
		return
	}
	var payload *StreamSource
	err = json.Unmarshal(dataBytes, &payload)
	if err != nil {
		_ = ctx.Error(&middlewares.HttpException{
			Message: "request must be json format",
			Code: http.StatusBadRequest,
		})
		return
	}

	source := payload.Source
	sources, f, err := common.OpenFileAndReadLines(fmt.Sprintf("%s/%s", os.Getenv("RTSP_LEDGER"), "stream_sources"))
	if err != nil {
		_ = ctx.Error(&middlewares.HttpException{
			Message: fmt.Sprintf("can not open ledger: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	if !common.StrArrContains(source, sources) {
		_ = ctx.Error(&middlewares.HttpException{
			Message: "source does not exists",
			Code:    http.StatusBadRequest,
		})
		return
	}

	newSources := common.StrArrRemove(source, sources)

	err = common.WriteOverAndClose(f, newSources)
	if err != nil {
		_ = ctx.Error(&middlewares.HttpException{
			Message: fmt.Sprintf("failed to remove %s from sources", source),
			Code:    http.StatusBadRequest,
		})
		return
	}

	go pubsub.GetSourceManager().RemoveSource(source)

	ctx.JSON(http.StatusOK, gin.H{
		"sources": newSources,
	})
}
