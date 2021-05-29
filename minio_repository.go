package main

import "fmt"

type MinioRepository struct{}

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
