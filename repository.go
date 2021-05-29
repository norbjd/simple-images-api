package main

type Repository interface {
	AddImageAndReturnID(image Image) (string, error)
	GetImages() ([]ImageIDWithMetadata, error)
	GetImageByID(imageID string) (*ImageContent, error)
	GetImageMetadataByID(imageID string) (*ImageIDWithMetadata, error)
	DeleteImageByID(imageID string) error
}
