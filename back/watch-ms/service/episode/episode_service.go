package episode

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"os"
	"shared/DTO/events"
	"shared/broker"
	"strings"
	"watch-ms/model"
	"watch-ms/provider/episode"
)

type Service interface {
	UploadNewEpisode()
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
		reader:          broker.NewBrokerReader("upload-video", "watch-ms", brokers),
		messageCh:       make(chan kafka.Message),
		commitCh:        make(chan kafka.Message),
		ctx:             context.Background(),
	}
}

func (s *serviceImpl) UploadNewEpisode() {
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
		fmt.Println(unmarshalEpisode)
		s.commitCh <- msg
	}
}
