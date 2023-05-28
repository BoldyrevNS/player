package broker

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Reader interface {
	FetchMessage(ctx context.Context, messageCh chan kafka.Message)
	CommitMessage(ctx context.Context, commitCh chan kafka.Message)
}

type readerImpl struct {
	Reader *kafka.Reader
}

func NewBrokerReader(topic string, groupId string, brokers []string) Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupId,
	})
	return &readerImpl{
		Reader: reader,
	}
}

func (k *readerImpl) FetchMessage(ctx context.Context, messageCh chan kafka.Message) {
	for {
		message, _ := k.Reader.FetchMessage(ctx)
		messageCh <- message
	}
}

func (k *readerImpl) CommitMessage(ctx context.Context, commitCh chan kafka.Message) {
	for {
		msg := <-commitCh
		_ = k.Reader.CommitMessages(ctx, msg)
	}
}
