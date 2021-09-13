package main

import (
	"antelope-stream-chunking/middlewares"
	"antelope-stream-chunking/pubsub"
	"fmt"
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/toorop/gin-logrus"
	"io"
	"os"
	"time"
)

func InitLogger() {
	formatter := runtime.Formatter{
		ChildFormatter: &log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RubyDate,
		},
	}
	formatter.Line = true
	log.SetFormatter(&formatter)
	if os.Getenv("LOG_FILE") != "" {
		logFile, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file %s for output: %s", os.Getenv("LOG_FILE"), err)
		}
		log.SetOutput(io.MultiWriter(os.Stderr, logFile))
		log.RegisterExitHandler(func() {
			if logFile == nil {
				return
			}
			_ = logFile.Close()
		})
	}
}

func InitSourcesManager() {
	log.Info("Set up stream manager")
	sourceManager := pubsub.GetSourceManager()
	sourceManager.LoadSourcesFromFile(fmt.Sprintf("%s/%s", os.Getenv("RTSP_LEDGER"), "stream_sources"))
	sourceManager.ActivateAllSources()
}

func main() {
	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("error: failed to load the env file")
		}
	}

	InitLogger()

	rootRouter := gin.New()
	rootPort := os.Getenv("PORT")

	rootRouter.Use(middlewares.HttpExceptionHandle())
	rootRouter.Use(ginlogrus.Logger(log.StandardLogger()), gin.Recovery())

	log.Infof("PORT: %s, ENV: %s", rootPort, os.Getenv("ENV"))
	InitRouter(rootRouter)
	InitSourcesManager()

	err := rootRouter.Run("0.0.0.0:" + rootPort)
	if err != nil {
		log.Fatalf("Failed to start server: \n%s\n", err.Error())
	}
}
