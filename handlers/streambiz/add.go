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

type StreamSource struct {
	Source string `json:"source"`
}

func (sb StreamBizRouter) AddStream(ctx *gin.Context) {
	// Check request format
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

	// Get source from request body
	// format:
	// {
	// 	  "source": "some rtsp"
	// }
	source := payload.Source
	sources, fileReader, err := common.OpenFileAndReadLines(fmt.Sprintf("%s/%s", os.Getenv("RTSP_LEDGER"), "stream_sources"))
	if err != nil {
		_ = ctx.Error(&middlewares.HttpException{
			Message: fmt.Sprintf("can not open ledger: %s", err.Error()),
			Code: http.StatusInternalServerError,
		})
		return
	}

	// Add source to save file
	if common.StrArrContains(source, sources) {
		_ = ctx.Error(&middlewares.HttpException{
			Message: "source already exists",
			Code: http.StatusFound,
		})
		return
	}
	common.AppendLineAndClose(fileReader, source)

	go pubsub.GetSourceManager().AddSourceToMem(source)

	ctx.JSON(http.StatusOK, gin.H{"new_sources": append(sources, source)})
}
