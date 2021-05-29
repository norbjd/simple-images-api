package main

import (
	"fmt"

	"github.com/minio/minio-go/v6"
)

type MinioRepository struct {
	bucketName string
	client     *minio.Client
}

func NewMinioClient(endpoint, accessKey, secretKey string) (*minio.Client, error) {
	return minio.New(endpoint, accessKey, secretKey, false)
}

func (r MinioRepository) AddImageAndReturnID(image Image) (string, error) {
	return "", fmt.Errorf("not implemented yet")
}

func (r MinioRepository) GetImages() ([]ImageIDWithMetadata, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r MinioRepository) GetImageByID(imageID string) (*ImageContent, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r MinioRepository) GetImageMetadataByID(imageID string) (*ImageIDWithMetadata, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r MinioRepository) DeleteImageByID(imageID string) error {
	return fmt.Errorf("not implemented yet")
}
