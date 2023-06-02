package watch

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"os"
	"shared/DTO/events"
	"shared/broker"
	"strings"
	watchDto "watch-ms/DTO"
	"watch-ms/model"
	"watch-ms/provider/watch"
)

type Service interface {
	UserStartWatch(data watchDto.UserStartWatchDTO, userId uint) error
	OnUpdateWatchThumbnail()
}

type serviceImpl struct {
	watchProvider                watch.Provider
	userStartWatchProducer       broker.Writer
	updateWatchThumbnailConsumer broker.Reader
	ctx                          context.Context
}

func NewWatchService(watchProvider watch.Provider) Service {
	brokers := strings.Split(os.Getenv("BROKERS"), ",")
	return &serviceImpl{
		watchProvider:                watchProvider,
		userStartWatchProducer:       broker.NewBrokerWriter("user-start-watch", os.Getenv("WRITE_BROKER")),
		updateWatchThumbnailConsumer: broker.NewBrokerReader("update-watch_thumbnail", "watch-ms", brokers),
		ctx:                          context.Background(),
	}
}

func (s *serviceImpl) UserStartWatch(data watchDto.UserStartWatchDTO, userId uint) error {
	createdWatch, err := s.watchProvider.CreateAndReturnJoinVideoUrl(model.Watch{
		UserId:    userId,
		EpisodeId: data.EpisodeId,
	})
	if err != nil {
		return err
	}
	marshalWatch, err := json.Marshal(events.UserStartWatchDTO{
		WatchId:  createdWatch.WatchId,
		VideoUrl: createdWatch.VideoUrl,
	})
	if err != nil {
		return err
	}
	s.userStartWatchProducer.WriteMessage(s.ctx, string(marshalWatch))
	return nil
}

func (s *serviceImpl) OnUpdateWatchThumbnail() {
	messageCh := make(chan kafka.Message)
	commitCh := make(chan kafka.Message)

	go s.updateWatchThumbnailConsumer.FetchMessage(s.ctx, messageCh)
	go s.updateWatchThumbnailConsumer.CommitMessage(s.ctx, commitCh)
	for {
		var updateData events.UpdateThumbnailWatchDTO
		msg := <-messageCh
		err := json.Unmarshal(msg.Value, &updateData)
		if err != nil {
			return
		}
		err = s.watchProvider.UpdateThumbnailUrl(updateData.WatchId, updateData.ThumbnailPath)
		if err != nil {
			return
		}
	}
}
