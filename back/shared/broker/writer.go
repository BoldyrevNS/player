package broker

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type Writer interface {
	WriteMessage(ctx context.Context, message string)
}

type writerImpl struct {
	Writer *kafka.Writer
}

func NewBrokerWriter(topic string, addr string) Writer {
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	err = conn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     10,
		ReplicationFactor: 1,
	})
	if err != nil {
		panic(err.Error())
	}
	writer := &kafka.Writer{
		Addr:  kafka.TCP(addr),
		Topic: topic,
	}
	return &writerImpl{
		Writer: writer,
	}
}

func (k *writerImpl) WriteMessage(ctx context.Context, message string) {
	err := k.Writer.WriteMessages(ctx, kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		fmt.Println("Broker write error: ", err)
	}
}
