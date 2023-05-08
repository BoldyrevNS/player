package broker

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type ConsumerHandlerParams struct {
	Message kafka.Message
	Err     error
	Reader  *kafka.Reader
	Ctx     context.Context
}

func NewConsumerWorker(topic string, groupId string, handler func(handlerParams ConsumerHandlerParams)) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	brokers := strings.Split(os.Getenv("BROKERS"), ",")
	config := kafka.ReaderConfig{
		Brokers:         brokers,
		Topic:           topic,
		GroupID:         groupId,
		MinBytes:        10e3,
		MaxBytes:        10e6,
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
	}
	reader := kafka.NewReader(config)
	defer func() {
		err := reader.Close()
		if err != nil {
			log.Fatal("reader close err:", err)
		}
	}()

	for {
		msg, err := reader.ReadMessage(context.Background())
		params := ConsumerHandlerParams{
			Message: msg,
			Err:     err,
			Reader:  reader,
			Ctx:     context.Background(),
		}
		handler(params)
	}
}
