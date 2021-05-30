package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type InMemoryRepository struct {
	images map[string]Image
}

func (r InMemoryRepository) AddImageAndReturnID(image Image) (string, error) {
	imageUUID, _ := uuid.NewUUID()
	imageID := imageUUID.String()
	r.images[imageID] = image
	return imageID, nil
}

func (r InMemoryRepository) GetImages() ([]ImageIDWithMetadata, error) {
	result := make([]ImageIDWithMetadata, 0)
	for imageID, image := range r.images {
		result = append(result, ImageIDWithMetadata{
			ID: imageID,
			Metadata: ImageMetadata{
				Name:        image.Metadata.Name,
				Description: image.Metadata.Description,
			},
		})
	}
	return result, nil
}

func (r InMemoryRepository) GetImageByID(imageID string) (*ImageContent, error) {
	return &ImageContent{Content: r.images[imageID].getBinaryContent()}, nil
}

func (r InMemoryRepository) GetImageMetadataByID(imageID string) (*ImageIDWithMetadata, error) {
	return &ImageIDWithMetadata{
		ID:       imageID,
		Metadata: r.images[imageID].Metadata,
	}, nil
}

func (r InMemoryRepository) DeleteImageByID(imageID string) error {
	delete(r.images, imageID)
	return nil
}

func TestAppGetImages(t *testing.T) {
	app := App{repository: InMemoryRepository{images: map[string]Image{
		"id1": {Content: ImageContent{}, Metadata: ImageMetadata{Name: "image 1", Description: "An image"}},
		"id2": {Content: ImageContent{}, Metadata: ImageMetadata{Name: "image 2", Description: "Another image"}},
	}}}
	router := app.initRouter()

	req, _ := http.NewRequest("GET", "/images", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)

	expectedJSON := `[
		{
			"id": "id1",
			"metadata": {
				"name": "image 1",
				"description": "An image"
			}
		},
		{
			"id": "id2",
			"metadata": {
				"name": "image 2",
				"description": "Another image"
			}
		}
	]`

	assert.JSONEq(t, expectedJSON, recorder.Body.String())
}

func TestAppGetImages_no_images(t *testing.T) {
	app := App{repository: InMemoryRepository{images: map[string]Image{}}}
	router := app.initRouter()

	req, _ := http.NewRequest("GET", "/images", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)

	expectedJSON := `[]`

	assert.JSONEq(t, expectedJSON, recorder.Body.String())
}
