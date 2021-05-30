package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"mime/multipart"
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
	image, exists := r.images[imageID]
	if !exists {
		return nil, ErrImageDoesNotExist
	}

	return &ImageContent{Content: image.getBinaryContent()}, nil
}

func (r InMemoryRepository) GetImageMetadataByID(imageID string) (*ImageIDWithMetadata, error) {
	image, exists := r.images[imageID]
	if !exists {
		return nil, ErrImageDoesNotExist
	}

	return &ImageIDWithMetadata{
		ID:       imageID,
		Metadata: image.Metadata,
	}, nil
}

func (r InMemoryRepository) DeleteImageByID(imageID string) error {
	_, exists := r.images[imageID]
	if !exists {
		return ErrImageDoesNotExist
	}

	delete(r.images, imageID)
	return nil
}

func insertImageRequest(content []byte, name string, description string) *http.Request {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)

	if name != "" {
		nameField, _ := w.CreateFormField("name")
		nameField.Write([]byte(name))
	}

	if description != "" {
		descriptionField, _ := w.CreateFormField("description")
		descriptionField.Write([]byte(description))
	}

	fw, _ := w.CreateFormFile("image", "image")
	fd := bytes.NewReader(content)

	io.Copy(fw, fd)
	w.Close()

	req, _ := http.NewRequest("POST", "/images", buf)
	req.Header.Set("Content-Type", w.FormDataContentType())

	return req
}

func TestAppAddImage(t *testing.T) {
	app := App{repository: InMemoryRepository{images: make(map[string]Image)}}
	router := app.initRouter()

	imageBytes, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAIAAAD/gAIDAAACu0lEQVR42u2cXXaEIAyFgeMiZYW6S/owL+2oGDAJF3vzOJ0ewkf+xYmllECRSSICwiIswiIswiIsIiAswiIswiIswiICwjKRZRZF9xiPH66+86UIMs86ZSEUN2QQsJ6Q8iS4vA/TOy3LmdRaynHFJgMcY1kxxm3Cs1mGaLmFKWWZyNHyBevsdQALOKYs/nDWCh4hx+UpLKuDVL52okqyjn8X6qMTY0M9kBBI9Um5E8laWwh7jFGmdhpLar3emHoVXjmVTab84kMqjygavlivB889bqF+PIuzTeVqGLJu7spZEY9bwft43JPWonKWCYGUc4dYOZh68FJopCWGnTFINbWoR6b/eqzcen7msDKqWXV4TLImFeYk5Z0NfUh9lU6ty30Ck3C+lkxJSbb6W8A9MVlgyki1gmKk50NWpGxIWDXZfH3wdiyDng037WBBN6Qbzim726R0nbxAl0+4dGJWnvO5aesgUC3A52nDNmMWYb146jC7HBMXYZ3H39OwSzeUkgpaN/9aJ3DH0ZV6ApWrJFcGxbJGTf6a1gVyQy1eD4euKA9ZrffZSrx1rQRoI93/bu3LYyyrIwOom+Tp94vDLZpuXpUdfv6kglUr1ZqUDk0q9u1WbkenavQprPaGxcPSySLcXNaWF2vdKmwYs5r2r16UqpPCKh20eNk95tF80elJ5FJxye5sIFQScepwmyh1M4n8OJVfodMyLrcA2qRb8lF07A0ZrdWThQfNMleAuyY51sQUSZnAUu/7QEhZWRYCL3VSwfqF8lsoFqvXF32yYnS+IGs+GDAjFXx+qsChSfQZ1Pj9rkNrnJIopl7lo8Dqi+vyaZSHj/v3Ip6llnIzN7Bxs6am/+AWucsFYYQFS4WaR1rHfGcSDRMuLFjhlSPCIizCIizCohAWYREWYREWYVEIi7AIi7AI66XyA2Yom/1j6CTHAAAAAElFTkSuQmCC")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, insertImageRequest(imageBytes, "Toto", ""))

	assert.Equal(t, 200, recorder.Code)

	var imageIDWithMetadata ImageIDWithMetadata
	json.Unmarshal(recorder.Body.Bytes(), &imageIDWithMetadata)

	assert.True(t, isValidUUID(imageIDWithMetadata.ID))
	assert.Equal(t, "Toto", imageIDWithMetadata.Metadata.Name)
	assert.Equal(t, "", imageIDWithMetadata.Metadata.Description)
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

	var imagesIDWithMetadata []ImageIDWithMetadata
	json.Unmarshal(recorder.Body.Bytes(), &imagesIDWithMetadata)

	expected := []ImageIDWithMetadata{
		{
			ID: "id1",
			Metadata: ImageMetadata{
				Name:        "image 1",
				Description: "An image",
			},
		},
		{
			ID: "id2",
			Metadata: ImageMetadata{
				Name:        "image 2",
				Description: "Another image",
			},
		},
	}

	assert.ElementsMatch(t, expected, imagesIDWithMetadata)
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

func TestAppGetImage(t *testing.T) {
	imageContent, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAIAAAD/gAIDAAACu0lEQVR42u2cXXaEIAyFgeMiZYW6S/owL+2oGDAJF3vzOJ0ewkf+xYmllECRSSICwiIswiIswiIsIiAswiIswiIswiICwjKRZRZF9xiPH66+86UIMs86ZSEUN2QQsJ6Q8iS4vA/TOy3LmdRaynHFJgMcY1kxxm3Cs1mGaLmFKWWZyNHyBevsdQALOKYs/nDWCh4hx+UpLKuDVL52okqyjn8X6qMTY0M9kBBI9Um5E8laWwh7jFGmdhpLar3emHoVXjmVTab84kMqjygavlivB889bqF+PIuzTeVqGLJu7spZEY9bwft43JPWonKWCYGUc4dYOZh68FJopCWGnTFINbWoR6b/eqzcen7msDKqWXV4TLImFeYk5Z0NfUh9lU6ty30Ck3C+lkxJSbb6W8A9MVlgyki1gmKk50NWpGxIWDXZfH3wdiyDng037WBBN6Qbzim726R0nbxAl0+4dGJWnvO5aesgUC3A52nDNmMWYb146jC7HBMXYZ3H39OwSzeUkgpaN/9aJ3DH0ZV6ApWrJFcGxbJGTf6a1gVyQy1eD4euKA9ZrffZSrx1rQRoI93/bu3LYyyrIwOom+Tp94vDLZpuXpUdfv6kglUr1ZqUDk0q9u1WbkenavQprPaGxcPSySLcXNaWF2vdKmwYs5r2r16UqpPCKh20eNk95tF80elJ5FJxye5sIFQScepwmyh1M4n8OJVfodMyLrcA2qRb8lF07A0ZrdWThQfNMleAuyY51sQUSZnAUu/7QEhZWRYCL3VSwfqF8lsoFqvXF32yYnS+IGs+GDAjFXx+qsChSfQZ1Pj9rkNrnJIopl7lo8Dqi+vyaZSHj/v3Ip6llnIzN7Bxs6am/+AWucsFYYQFS4WaR1rHfGcSDRMuLFjhlSPCIizCIizCohAWYREWYREWYVEIi7AIi7AI66XyA2Yom/1j6CTHAAAAAElFTkSuQmCC")

	app := App{repository: InMemoryRepository{images: map[string]Image{
		"id1": {
			Content:  ImageContent{Content: imageContent},
			Metadata: ImageMetadata{Name: "Toto"},
		},
	}}}
	router := app.initRouter()

	req, _ := http.NewRequest("GET", "/images/id1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, []string{"image"}, recorder.Header()["Content-Type"])
	assert.Equal(t, imageContent, recorder.Body.Bytes())
}

func TestAppGetImage_non_existent_id(t *testing.T) {
	app := App{repository: InMemoryRepository{images: map[string]Image{}}}
	router := app.initRouter()

	req, _ := http.NewRequest("GET", "/images/id1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, 404, recorder.Code)
	assert.Equal(t, "Image does not exist", recorder.Body.String())
}

func TestAppDeleteImage(t *testing.T) {
	app := App{repository: InMemoryRepository{images: map[string]Image{
		"id1": {
			Content:  ImageContent{},
			Metadata: ImageMetadata{Name: "Toto"},
		},
	}}}
	router := app.initRouter()

	req, _ := http.NewRequest("DELETE", "/images/id1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "Image deleted", recorder.Body.String())
}

func TestAppDeleteImage_non_existent_id(t *testing.T) {
	app := App{repository: InMemoryRepository{images: map[string]Image{}}}
	router := app.initRouter()

	req, _ := http.NewRequest("DELETE", "/images/id1", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, 404, recorder.Code)
	assert.Equal(t, "Image does not exist", recorder.Body.String())
}

func TestAppGetImageMetadata(t *testing.T) {
	app := App{repository: InMemoryRepository{images: map[string]Image{
		"id1": {Content: ImageContent{}, Metadata: ImageMetadata{Name: "image 1", Description: "An image"}},
	}}}
	router := app.initRouter()

	req, _ := http.NewRequest("GET", "/images/id1/metadata", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)

	expectedJSON := `{
		"id": "id1",
		"metadata": {
			"name": "image 1",
			"description": "An image"
		}
	}`

	assert.JSONEq(t, expectedJSON, recorder.Body.String())
}

func TestAppGetImageMetadata_non_existent_id(t *testing.T) {
	app := App{repository: InMemoryRepository{images: map[string]Image{}}}
	router := app.initRouter()

	req, _ := http.NewRequest("GET", "/images/id1/metadata", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, 404, recorder.Code)
	assert.Equal(t, "Image does not exist", recorder.Body.String())
}
