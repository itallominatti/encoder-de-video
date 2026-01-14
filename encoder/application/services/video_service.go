package services

import (
	"context"
	"encoder/application/repositories"
	"encoder/domain"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func (v *VideoService) Download() error {
	ctx := context.Background()

	bucketName := os.Getenv("AWS_STORAGE_BUCKET_NAME")
	if bucketName == "" {
		return fmt.Errorf("AWS_STORAGE_BUCKET_NAME not set")
	}

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

func (v *VideoService) Fragment() error {

	err := os.Mkdir(os.Getenv("localStoragePath"+"/"+v.Video.ID), os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("localStoragePath") + "/" + v.Video.ID + "/frag"

	cmd := exec.Command("mp4fragment", source, target)
	err = cmd.Run()
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error fragmenting video: %v, output: %s", err, string(output))
	}
	log.Printf("video %v has been fragmented", v.Video.ID)
	return nil
}
