package episode

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"os"
	"shared/DTO/events"
	"shared/broker"
	"strings"
	"watch-ms/model"
	"watch-ms/provider/episode"
)

type Service interface {
	OnUploadNewEpisode()
}

type serviceImpl struct {
	episodeProvider episode.Provider
	reader          broker.Reader
	messageCh       chan kafka.Message
	commitCh        chan kafka.Message
	ctx             context.Context
}

func NewEpisodeService(episodeProvider episode.Provider) Service {
	brokers := strings.Split(os.Getenv("BROKERS"), ",")
	return &serviceImpl{
		episodeProvider: episodeProvider,
		reader:          broker.NewBrokerReader("upload-episode", "watch-ms", brokers),
		messageCh:       make(chan kafka.Message),
		commitCh:        make(chan kafka.Message),
		ctx:             context.Background(),
	}
}

func (s *serviceImpl) OnUploadNewEpisode() {
	go s.reader.FetchMessage(s.ctx, s.messageCh)
	go s.reader.CommitMessage(s.ctx, s.commitCh)
	for {
		var unmarshalEpisode events.UploadNewEpisodeDTO
		msg := <-s.messageCh
		_ = json.Unmarshal(msg.Value, &unmarshalEpisode)
		_ = s.episodeProvider.Create(model.Episode{
			SeasonId:      unmarshalEpisode.SeasonId,
			EpisodeName:   unmarshalEpisode.EpisodeName,
			EpisodeNumber: unmarshalEpisode.EpisodeNumber,
			VideoUrl:      unmarshalEpisode.VideoPath,
			ThumbnailUrl:  unmarshalEpisode.ThumbnailPath,
		})
		s.commitCh <- msg
	}
}
