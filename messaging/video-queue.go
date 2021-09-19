package messaging

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"net"
	"os"
	"strconv"
	"time"
)

type Message struct {
	Order     int64  `json:"order"`
	Source    string `json:"source"`
	Frame     string `json:"frame"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Timestamp int64  `json:"timestamp"`
}

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
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: 1,
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func SendToKafka(payload *gocv.Mat, order int64, source string) {
	bufferImage, err := gocv.IMEncode(gocv.JPEGFileExt, *payload)
	if err != nil {
		log.Error(err)
		return
	}
	b64Image := b64.StdEncoding.EncodeToString(bufferImage)

	m := Message{
		Order:     order,
		Source:    source,
		Frame:     b64Image,
		Width:     payload.Cols(),
		Height:    payload.Rows(),
		Timestamp: time.Now().Unix(),
	}

	jsonMessage, err := json.Marshal(m)

	if err != nil {
		log.Error(err)
		return
	}

	err = GetWriter().WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(source),
		Value: jsonMessage,
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
