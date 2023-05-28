package provider

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"time"
)

type FileObj struct {
	Buffer      []byte
	Size        int64
	ContentType string
}

type Provider interface {
	UploadFile(file *multipart.FileHeader, filename string, bucketName string) error
	GetFile(filename string, bucketName string) (string, error)
	DeleteFile(filename string, bucketName string) error
}

type providerImpl struct {
	client *minio.Client
	ctx    context.Context
}

func NewFileService() Provider {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyId := os.Getenv("MINIO_ACCESS_KEY")
	secret := os.Getenv("MINIO_SECRET_KEY")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secret, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Minio init error: %v", err)
	}
	return &providerImpl{ctx: context.Background(), client: minioClient}
}

func (p *providerImpl) createBucket(bucketName string) error {
	bucketExists, err := p.client.BucketExists(p.ctx, bucketName)
	if err != nil {
		return err
	}
	if !bucketExists {
		err = p.client.MakeBucket(p.ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		policy := "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:GetBucketLocation\",\"s3:ListBucket\"],\"Resource\":[\"arn:aws:s3:::" + bucketName + "\"]},{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:GetObject\"],\"Resource\":[\"arn:aws:s3:::" + bucketName + "/*\"]}]}\n"
		err = p.client.SetBucketPolicy(p.ctx, bucketName, policy)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (p *providerImpl) UploadFile(file *multipart.FileHeader, filename string, bucketName string) error {
	contentType := file.Header["Content-Type"][0]
	size := file.Size
	buffer, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		err := buffer.Close()
		if err != nil {
			log.Fatal("buffer close error on file uploading:", err)
		}
	}()
	err = p.createBucket(bucketName)
	if err != nil {
		return err
	}
	_, err = p.client.PutObject(p.ctx, bucketName, filename, buffer, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return err
	}
	return nil
}

func (p *providerImpl) DeleteFile(fileName string, bucketName string) error {
	err := p.client.RemoveObject(p.ctx, bucketName, fileName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (p *providerImpl) GetFile(filename string, bucketName string) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=%s", filename))
	presignedURL, err := p.client.PresignedGetObject(p.ctx, bucketName, filename, time.Second*60*60, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.Path, nil
}
