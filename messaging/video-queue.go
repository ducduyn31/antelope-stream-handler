package messaging

import (
	"antelope-stream-chunking/common"
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"net"
	"os"
	"strconv"
)

func establishConnection(topic string, partitions int) error {
	conn, err := kafka.Dial("tcp", os.Getenv("KAFKA_URI"))
	defer func(conn *kafka.Conn) {
		_ = conn.Close()
	}(conn)
	if err != nil {
		log.Fatal(err)
		return err
	}

	controller, err := conn.Controller()
	if err != nil {
		log.Fatal(err)
		return err
	}

	var controllerConn *kafka.Conn

	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer func(controllerConn *kafka.Conn) {
		_ = controllerConn.Close()
	}(controllerConn)

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic: topic,
		NumPartitions: partitions,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
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
		_ = establishConnection("motion", 10)
		kafkaWriter = &kafka.Writer{
			Addr:     kafka.TCP(os.Getenv("KAFKA_URI")),
			Topic:    "motion",
			Balancer: &kafka.Hash{},
		}
	}
	return kafkaWriter
}
