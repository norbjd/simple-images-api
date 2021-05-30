package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v6"

	log "github.com/sirupsen/logrus"
)

type MinioRepository struct {
	bucketName string
	client     *minio.Client
}

func NewMinioClient(endpoint, accessKey, secretKey string) (*minio.Client, error) {
	return minio.New(endpoint, accessKey, secretKey, false)
}

func (r MinioRepository) AddImageAndReturnID(image Image) (string, error) {
	imageUUID, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("cannot generate UUID for image")
	}
	imageID := imageUUID.String()

	imageBytes := image.getBinaryContent()
	imageReader := bytes.NewReader(imageBytes)
	imageSizeInBytes := len(imageBytes)

	uploadInfo, err := r.client.PutObject(
		r.bucketName, imageID,
		imageReader, int64(imageSizeInBytes),
		minio.PutObjectOptions{ContentType: "image"},
	)
	if err != nil {
		return "", fmt.Errorf("cannot upload image : %v", err)
	}

	log.Debug("Image " + imageID + " successfully uploaded : " + fmt.Sprintf("%d", uploadInfo) + " bytes")

	imageIDWithMetadataJSON, _ := json.Marshal(ImageIDWithMetadata{ID: imageID, Metadata: image.Metadata})
	metadataReader := bytes.NewReader(imageIDWithMetadataJSON)
	_, err = r.client.PutObject(
		r.bucketName, "metadata/"+imageID+".json",
		metadataReader, int64(len(imageIDWithMetadataJSON)),
		minio.PutObjectOptions{ContentType: "application/json"},
	)
	if err != nil {
		return "", fmt.Errorf("cannot upload metadata : %v", err)
	}

	return imageID, nil
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
