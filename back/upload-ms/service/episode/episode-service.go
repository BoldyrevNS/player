package episode

import (
	"context"
	"encoding/json"
	"shared/DTO/events"
	"shared/broker"
	"upload-ms/DTO"
	"upload-ms/service/file"
)

type Service interface {
	UploadNewEpisode(data DTO.UploadNewEpisodeDTO) error
	UploadWatchThumbnail(data DTO.UploadWatchThumbnailDTO) error
}

type serviceImpl struct {
	fileService  file.Service
	brokerWriter broker.Writer
}

func NewEpisodeService(fileService file.Service, brokerWriter broker.Writer) Service {
	return &serviceImpl{
		brokerWriter: brokerWriter,
		fileService:  fileService,
	}
}

func (s *serviceImpl) UploadNewEpisode(data DTO.UploadNewEpisodeDTO) error {
	uploadedVideoInfo, err := s.fileService.UploadFile(data.Video, data.TitleId, data.SeasonId, file.Video)
	if err != nil {
		return err
	}
	uploadedImageInfo, err := s.fileService.UploadFile(data.Thumbnail, data.TitleId, data.SeasonId, file.Image)
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
	s.brokerWriter.WriteMessage(context.Background(), string(dataToSend))
	return nil
}

func (s *serviceImpl) UploadWatchThumbnail(data DTO.UploadWatchThumbnailDTO) error {
	return nil
}
