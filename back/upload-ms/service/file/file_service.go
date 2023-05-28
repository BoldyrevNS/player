package file

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	fileProvider "upload-ms/provider"
)

type UploadKind string

const (
	Video UploadKind = "video"
	Image UploadKind = "image"
)

type UploadInfo struct {
	Path       string
	Name       string
	BucketName string
}

type Service interface {
	UploadFile(video *multipart.FileHeader, titleId uint, seasonId uint, kind UploadKind) (UploadInfo, error)
	DeleteFile(bucketName string, filename string) error
}

type serviceImpl struct {
	fileProvider   fileProvider.Provider
	ctx            context.Context
	videoMimeTypes map[string]string
	imageMimeTypes map[string]string
}

func NewFileService(fileProvider fileProvider.Provider) Service {
	videoMimeTypes := map[string]string{
		"video/mp4":       ".mp4",
		"video/quicktime": ".mov",
		"video/x-msvideo": ".avi",
		"video/x-ms-wmv":  ".wmv",
	}

	imageMimeTypes := map[string]string{
		"image/png":  ".png",
		"image/jpeg": ".jpg",
		"image/webp": ".webp",
	}

	return &serviceImpl{
		fileProvider:   fileProvider,
		ctx:            context.Background(),
		videoMimeTypes: videoMimeTypes,
		imageMimeTypes: imageMimeTypes,
	}
}

func (s *serviceImpl) UploadFile(video *multipart.FileHeader, titleId uint, seasonId uint, kind UploadKind) (UploadInfo, error) {
	var fileExtension string
	var extensionCheck bool
	uploadMimeType := video.Header.Get("Content-type")
	if kind == Video {
		fileExtension, extensionCheck = s.videoMimeTypes[uploadMimeType]
	}
	if kind == Image {
		fileExtension, extensionCheck = s.imageMimeTypes[uploadMimeType]
	}
	if !extensionCheck {
		return UploadInfo{}, errors.New("wrong file extension")
	}
	name, err := s.genRandName(fileExtension)
	if err != nil {
		return UploadInfo{}, err
	}
	bucketName := fmt.Sprintf("%d-%d-%s", titleId, seasonId, kind)
	err = s.fileProvider.UploadFile(video, name, bucketName)
	if err != nil {
		return UploadInfo{}, err
	}
	path, err := s.fileProvider.GetFile(name, bucketName)
	if err != nil {
		return UploadInfo{}, err
	}
	return UploadInfo{
		Path:       path,
		Name:       name,
		BucketName: bucketName,
	}, nil
}

func (s *serviceImpl) DeleteFile(bucketName string, filename string) error {
	return s.fileProvider.DeleteFile(filename, bucketName)
}

func (s *serviceImpl) genRandName(fileExtension string) (string, error) {
	rand, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	newName := rand.String()
	return fmt.Sprintf("%s%s", newName, fileExtension), nil
}
