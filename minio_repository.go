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
	images := make([]ImageIDWithMetadata, 0)

	doneCh := make(chan struct{})
	defer close(doneCh)

	for object := range r.client.ListObjects(r.bucketName, "metadata/", false, doneCh) {
		if object.Err != nil {
			log.Error("Error while iterating on objects ", object.Err.Error())
			return nil, fmt.Errorf("cannot iterate on objects")
		}

		objectName := object.Key

		metadataReader, err := r.client.GetObject(r.bucketName, objectName, minio.GetObjectOptions{})
		if err != nil {
			log.Error("Error while reading object ", objectName)
			return nil, fmt.Errorf("error while reading object " + objectName)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(metadataReader)

		var imageIDWithMetadata ImageIDWithMetadata
		json.Unmarshal(buf.Bytes(), &imageIDWithMetadata)

		images = append(images, imageIDWithMetadata)
	}

	return images, nil
}

func (r MinioRepository) GetImageByID(imageID string) (*ImageContent, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r MinioRepository) GetImageMetadataByID(imageID string) (*ImageIDWithMetadata, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func (r MinioRepository) DeleteImageByID(imageID string) error {
	if r.client.RemoveObject(r.bucketName, imageID) != nil {
		return fmt.Errorf("cannot remove image %s", imageID)
	}

	metadataObjectName := "metadata/"+imageID+".json"
	if r.client.RemoveObject(r.bucketName, metadataObjectName) != nil {
		return fmt.Errorf("cannot remove image metadata %s", metadataObjectName)
	}
	
	return nil
}
