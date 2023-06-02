package episode

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"os"
	"shared/DTO/events"
	"shared/broker"
	"strings"
	"upload-ms/DTO"
	"upload-ms/model"
	"upload-ms/provider/watch_thumbnail"
	"upload-ms/service/file"
)

type SaveWatchInfo struct {
}

type Service interface {
	UploadNewEpisode(data DTO.UploadNewEpisodeDTO) error
	OnUserStartWatch()
	UploadWatchThumbnail(data DTO.UploadWatchThumbnailDTO) error
	SaveWatchThumbnail()
}

type serviceImpl struct {
	fileService                  file.Service
	watchThumbnailProvider       watch_thumbnail.Provider
	uploadEpisodeProducer        broker.Writer
	updateWatchThumbnailProducer broker.Writer
	userWatchConsumer            broker.Reader
	ctx                          context.Context
	saveFrameInfoCh              chan file.SaveFrameInfo
	uploadedFrameInfoCh          chan file.UploadedFrameInfo
	deleteFrameCh                chan string
}

func NewEpisodeService(
	fileService file.Service,
	watchThumbnailProvider watch_thumbnail.Provider,
	saveFrameInfoCh chan file.SaveFrameInfo,
	uploadedFrameInfoCh chan file.UploadedFrameInfo,
	deleteFrameCh chan string,
) Service {
	brokers := strings.Split(os.Getenv("BROKERS"), ",")
	return &serviceImpl{
		fileService:                  fileService,
		watchThumbnailProvider:       watchThumbnailProvider,
		uploadEpisodeProducer:        broker.NewBrokerWriter("upload-episode", os.Getenv("WRITE_BROKER")),
		updateWatchThumbnailProducer: broker.NewBrokerWriter("update-watch_thumbnail", os.Getenv("WRITE_BROKER")),
		userWatchConsumer:            broker.NewBrokerReader("user-start-watch", "upload-ms", brokers),
		ctx:                          context.Background(),
		saveFrameInfoCh:              saveFrameInfoCh,
		uploadedFrameInfoCh:          uploadedFrameInfoCh,
		deleteFrameCh:                deleteFrameCh,
	}
}

func (s *serviceImpl) UploadNewEpisode(data DTO.UploadNewEpisodeDTO) error {
	uploadedVideoInfo, err := s.fileService.UploadMultipartFile(data.Video, file.Video)
	if err != nil {
		return err
	}
	uploadedImageInfo, err := s.fileService.UploadMultipartFile(data.Thumbnail, file.Image)
	if err != nil {
		err = s.fileService.DeleteFile(uploadedVideoInfo.BucketName, uploadedVideoInfo.Name)
		return err
	}
	dataToSend, err := json.Marshal(&events.UploadNewEpisodeDTO{
		EpisodeName:   data.EpisodeName,
		EpisodeNumber: data.EpisodeNumber,
		SeasonId:      data.SeasonId,
		ThumbnailPath: uploadedImageInfo.Path,
		VideoPath:     uploadedVideoInfo.Path,
	})
	if err != nil {
		return err
	}
	s.uploadEpisodeProducer.WriteMessage(context.Background(), string(dataToSend))
	return nil
}

func (s *serviceImpl) OnUserStartWatch() {
	messageCh := make(chan kafka.Message)
	commitCh := make(chan kafka.Message)
	go s.userWatchConsumer.FetchMessage(s.ctx, messageCh)
	go s.userWatchConsumer.CommitMessage(s.ctx, commitCh)
	for {
		var marshalWatch events.UserStartWatchDTO
		msg := <-messageCh
		_ = json.Unmarshal(msg.Value, &marshalWatch)
		_ = s.watchThumbnailProvider.Create(model.WatchThumbnail{
			WatchId:    marshalWatch.WatchId,
			Filename:   "",
			BucketName: "",
			VideoUrl:   marshalWatch.VideoUrl,
		})
		commitCh <- msg
	}
}

func (s *serviceImpl) UploadWatchThumbnail(data DTO.UploadWatchThumbnailDTO) error {
	watch, err := s.watchThumbnailProvider.FindOne(data.WatchId)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	s.saveFrameInfoCh <- file.SaveFrameInfo{
		VideoUrl:           watch.VideoUrl,
		DirName:            "/tmp_thumbnails",
		TimeCode:           data.TimeCode,
		UploadedFilename:   watch.Filename,
		UploadedBucketName: watch.BucketName,
		WatchId:            data.WatchId,
	}
	return nil
}

func (s *serviceImpl) SaveWatchThumbnail() {
	contentType := "image/jpeg"
	for {
		uploadedLocalFrameInfo := <-s.uploadedFrameInfoCh
		uploadInfo, err := s.fileService.UploadLocalFile(contentType, uploadedLocalFrameInfo.FramePath, file.Image)
		if err != nil {
			s.deleteFrameCh <- uploadedLocalFrameInfo.FramePath
			return
		}

		s.deleteFrameCh <- uploadedLocalFrameInfo.FramePath

		if uploadedLocalFrameInfo.UploadedBucketName != "" && uploadedLocalFrameInfo.UploadedFilename != "" {
			err := s.fileService.DeleteFile(uploadedLocalFrameInfo.UploadedBucketName, uploadedLocalFrameInfo.UploadedFilename)
			if err != nil {
				return
			}
		}

		err = s.watchThumbnailProvider.UpdateOne(uploadedLocalFrameInfo.WatchId, watch_thumbnail.UpdateWatch{
			Filename:   uploadInfo.Name,
			BucketName: uploadInfo.BucketName,
		})
		if err != nil {
			s.deleteFrameCh <- uploadedLocalFrameInfo.FramePath
			_ = s.fileService.DeleteFile(uploadInfo.BucketName, uploadInfo.Name)
			return
		}

		updateEventData, err := json.Marshal(events.UpdateThumbnailWatchDTO{
			WatchId:       uploadedLocalFrameInfo.WatchId,
			ThumbnailPath: uploadInfo.Path,
		})
		if err != nil {
			return
		}
		s.updateWatchThumbnailProducer.WriteMessage(context.Background(), string(updateEventData))
	}
}
