package messaging

import (
	"antelope-stream-chunking/common"
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"os"
)

func establishConnection(topic string, partition int) *kafka.Conn {
	conn, err := kafka.DialLeader(context.Background(), "tcp", os.Getenv("KAFKA_URI"), topic, partition)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return conn
}

func SendToKafka(payload *gocv.Mat, order int, source string) {
	bufferImage, err := gocv.IMEncode(gocv.JPEGFileExt, *payload)
	if err != nil {
		log.Error(err)
		return
	}

	err = GetWriter().WriteMessages(context.Background(), kafka.Message{
		Key: []byte(source),
		Value: common.ConcatByteArr([]byte(fmt.Sprintf("order=%d&image=", order)), bufferImage),
	})
}

var kafkaWriter *kafka.Writer

func GetWriter() *kafka.Writer {
	if kafkaWriter == nil {
		establishConnection("motion", 0)
		kafkaWriter = &kafka.Writer{
			Addr:     kafka.TCP(os.Getenv("KAFKA_URI")),
			Topic:    "motion",
			Balancer: &kafka.LeastBytes{},
		}
	}
	return kafkaWriter
}
