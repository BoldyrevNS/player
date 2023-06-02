package file

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"os/exec"
	fileProvider "upload-ms/provider/file"
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

type SaveFrameInfo struct {
	VideoUrl           string
	DirName            string
	TimeCode           string
	UploadedFilename   string
	UploadedBucketName string
	WatchId            uint
}

type UploadedFrameInfo struct {
	FramePath          string
	FrameName          string
	UploadedFilename   string
	UploadedBucketName string
	WatchId            uint
}

type Service interface {
	UploadMultipartFile(file *multipart.FileHeader, kind UploadKind) (UploadInfo, error)
	UploadLocalFile(contentType string, filePath string, kind UploadKind) (UploadInfo, error)
	DeleteFile(bucketName string, filename string) error
	SaveVideoFrameLocal()
	DeleteVideoFrameLocal()
	GenRandName() (string, error)
	GenRandFilename(fileExtension string) (string, error)
	GetFileContentType(file *os.File) (string, error)
	GetFileHeader(file *os.File) (*multipart.FileHeader, error)
}

type serviceImpl struct {
	fileProvider        fileProvider.Provider
	ctx                 context.Context
	videoMimeTypes      map[string]string
	imageMimeTypes      map[string]string
	saveFrameInfoCh     chan SaveFrameInfo
	uploadedFrameInfoCh chan UploadedFrameInfo
	deleteFrameCh       chan string
}

func NewFileService(
	fileProvider fileProvider.Provider,
	saveFrameInfoCh chan SaveFrameInfo,
	uploadedFrameInfoCh chan UploadedFrameInfo,
	deleteFrameCh chan string,
) Service {
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
		fileProvider:        fileProvider,
		ctx:                 context.Background(),
		videoMimeTypes:      videoMimeTypes,
		imageMimeTypes:      imageMimeTypes,
		saveFrameInfoCh:     saveFrameInfoCh,
		uploadedFrameInfoCh: uploadedFrameInfoCh,
		deleteFrameCh:       deleteFrameCh,
	}
}

func (s *serviceImpl) UploadMultipartFile(video *multipart.FileHeader, kind UploadKind) (UploadInfo, error) {
	contentType := video.Header.Get("Content-type")

	fileExtension, valid := s.validateContentType(contentType, kind)
	if !valid {
		return UploadInfo{}, errors.New("wrong upload extension")
	}
	name, err := s.GenRandFilename(fileExtension)
	if err != nil {
		return UploadInfo{}, err
	}
	bucketName, err := s.GenRandName()
	if err != nil {
		return UploadInfo{}, err
	}
	err = s.fileProvider.UploadMultipartFile(video, name, bucketName)
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

func (s *serviceImpl) UploadLocalFile(contentType string, filePath string, kind UploadKind) (UploadInfo, error) {
	fileExtension, valid := s.validateContentType(contentType, kind)
	if !valid {
		return UploadInfo{}, errors.New("wrong upload extension")
	}
	name, err := s.GenRandFilename(fileExtension)
	if err != nil {
		return UploadInfo{}, err
	}
	bucketName, err := s.GenRandName()
	if err != nil {
		return UploadInfo{}, err
	}
	err = s.fileProvider.UploadLocalFile(filePath, contentType, name, bucketName)
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

func (s *serviceImpl) SaveVideoFrameLocal() {
	width := 640
	height := 360
	size := fmt.Sprintf("%dx%d", width, height)
	var imageBuffer bytes.Buffer

	for {
		data := <-s.saveFrameInfoCh
		if _, err := os.Stat(data.DirName); os.IsNotExist(err) {
			err = os.Mkdir(data.DirName, os.ModePerm)
		}

		frameName, err := s.GenRandFilename(".jpg")
		framePath := fmt.Sprintf("%s/%s", data.DirName, frameName)
		videoUrl := fmt.Sprintf("%s%s", os.Getenv("DOWNLOAD_ENDPOINT"), data.VideoUrl)
		if err != nil {
			return
		}
		cmd := exec.Command("ffmpeg", "-ss", data.TimeCode, "-i", videoUrl, "-vframes", "1", "-s", size, framePath)
		cmd.Stdout = &imageBuffer
		err = cmd.Run()
		if err != nil {
			panic("generate frame error")
		}
		_ = cmd.Wait()
		s.uploadedFrameInfoCh <- UploadedFrameInfo{
			FramePath:          framePath,
			FrameName:          frameName,
			UploadedFilename:   data.UploadedFilename,
			UploadedBucketName: data.UploadedBucketName,
			WatchId:            data.WatchId,
		}
	}
}

func (s *serviceImpl) DeleteVideoFrameLocal() {
	for {
		deletePath := <-s.deleteFrameCh
		_ = os.Remove(deletePath)
	}
}

func (s *serviceImpl) GenRandName() (string, error) {
	rand, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return rand.String(), nil
}

func (s *serviceImpl) GenRandFilename(fileExtension string) (string, error) {
	rand, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	newName := rand.String()
	return fmt.Sprintf("%s%s", newName, fileExtension), nil
}

func (s *serviceImpl) GetFileHeader(file *os.File) (*multipart.FileHeader, error) {
	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	contentType, err := s.GetFileContentType(file)
	if err != nil {
		return nil, err
	}
	return &multipart.FileHeader{
		Filename: fileStat.Name(),
		Header: textproto.MIMEHeader{
			"Content-type": []string{contentType},
		},
		Size: fileStat.Size(),
	}, nil
}

func (s *serviceImpl) GetFileContentType(file *os.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func (s *serviceImpl) validateContentType(contentType string, kind UploadKind) (string, bool) {
	var fileExtension string
	var extensionCheck bool
	if kind == Video {
		fileExtension, extensionCheck = s.videoMimeTypes[contentType]
	}
	if kind == Image {
		fileExtension, extensionCheck = s.imageMimeTypes[contentType]
	}
	return fileExtension, extensionCheck
}
