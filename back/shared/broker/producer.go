package broker

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"time"
)

func ProduceMessage(topic string, partition int, message []byte) error {
	brokerUrl := os.Getenv("BROKER_URL")
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokerUrl, topic, partition)
	if err != nil {
		return err
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal("broker connection close error:", err)
		}
	}()
	err = conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return err
	}
	_, err = conn.WriteMessages(kafka.Message{Value: message})
	if err != nil {
		return err
	}
	return nil
}
