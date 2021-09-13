package main

import (
	"antelope-stream-chunking/handlers/streambiz"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"reflect"
)

func InitRouter(router *gin.Engine) {
	streamController := new(streambiz.StreamBizRouter)

	apiV1Group := router.Group("/api/v1")
	{
		apiV1Group.POST("/stream", streamController.AddStream)
		defer log.Infof("register POST /api/v1/stream to %s", reflect.TypeOf(streamController))

		apiV1Group.GET("/stream", streamController.GetAll)
		defer log.Infof("register GET /api/v1/stream to %s", reflect.TypeOf(streamController))

		apiV1Group.DELETE("/stream", streamController.RemoveStream)
		defer log.Infof("register DELETE /api/v1/stream/ to %s", reflect.TypeOf(streamController))

		apiV1Group.GET("/stream/activate", streamController.ActivateStream)
		defer log.Infof("register GET /api/v1/activate/ to %s", reflect.TypeOf(streamController))
	}

}
