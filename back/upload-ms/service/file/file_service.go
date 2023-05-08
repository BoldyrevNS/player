package file

import (
	"encoding/json"
	"fmt"
	"shared/broker"
	"upload-ms/DTO"
	fileProvider "upload-ms/provider"
)

type Service interface {
	UploadVideo(data DTO.UploadFileDTO) error
	GetVideo(data DTO.GetFileDTO) (string, error)
}

type serviceImpl struct {
	fileProvider fileProvider.Provider
}

func NewUploadService(fileProvider fileProvider.Provider) Service {
	return &serviceImpl{fileProvider: fileProvider}
}

func (s *serviceImpl) UploadVideo(data DTO.UploadFileDTO) error {
	err := s.fileProvider.UploadFile(data.File, fmt.Sprintf("video-%s", data.TitleName))
	if err != nil {
		return err
	}
	msg, err := json.Marshal(DTO.ProduceUploadMsg{EpisodeName: data.File.Filename, TitleId: data.TitleId})
	if err != nil {
		return err
	}
	err = broker.ProduceMessage("video-uploaded", 0, msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *serviceImpl) GetVideo(data DTO.GetFileDTO) (string, error) {
	fileObj, err := s.fileProvider.GetFile(data.Filename, fmt.Sprintf("video-%s", data.Title))
	return fileObj, err
}
