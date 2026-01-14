package services

import (
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	getObjInput := &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &v.Video.FilePath,
	}

	result, err := client.GetObject(ctx, getObjInput)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	localPath := filepath.Join(os.Getenv("localstoragePath"), v.Video.ID+".mp4")
	f, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.ReadFrom(result.Body)
	if err != nil {
		return err
	}

	log.Printf("video %v has been stored", v.Video.ID)
	return nil
}
