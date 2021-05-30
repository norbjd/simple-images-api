package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v6"

	"testing"

	"github.com/stretchr/testify/assert"
)

const bucketName = "images"

var minioClient, _ = NewMinioClient("localhost:9001", "accessKey", "secretKey")

var minioRepository = MinioRepository{
	bucketName: bucketName,
	client:     minioClient,
}

func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func generateImage1() Image {
	imageContent, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAIAAAD/gAIDAAACu0lEQVR42u2cXXaEIAyFgeMiZYW6S/owL+2oGDAJF3vzOJ0ewkf+xYmllECRSSICwiIswiIswiIsIiAswiIswiIswiICwjKRZRZF9xiPH66+86UIMs86ZSEUN2QQsJ6Q8iS4vA/TOy3LmdRaynHFJgMcY1kxxm3Cs1mGaLmFKWWZyNHyBevsdQALOKYs/nDWCh4hx+UpLKuDVL52okqyjn8X6qMTY0M9kBBI9Um5E8laWwh7jFGmdhpLar3emHoVXjmVTab84kMqjygavlivB889bqF+PIuzTeVqGLJu7spZEY9bwft43JPWonKWCYGUc4dYOZh68FJopCWGnTFINbWoR6b/eqzcen7msDKqWXV4TLImFeYk5Z0NfUh9lU6ty30Ck3C+lkxJSbb6W8A9MVlgyki1gmKk50NWpGxIWDXZfH3wdiyDng037WBBN6Qbzim726R0nbxAl0+4dGJWnvO5aesgUC3A52nDNmMWYb146jC7HBMXYZ3H39OwSzeUkgpaN/9aJ3DH0ZV6ApWrJFcGxbJGTf6a1gVyQy1eD4euKA9ZrffZSrx1rQRoI93/bu3LYyyrIwOom+Tp94vDLZpuXpUdfv6kglUr1ZqUDk0q9u1WbkenavQprPaGxcPSySLcXNaWF2vdKmwYs5r2r16UqpPCKh20eNk95tF80elJ5FJxye5sIFQScepwmyh1M4n8OJVfodMyLrcA2qRb8lF07A0ZrdWThQfNMleAuyY51sQUSZnAUu/7QEhZWRYCL3VSwfqF8lsoFqvXF32yYnS+IGs+GDAjFXx+qsChSfQZ1Pj9rkNrnJIopl7lo8Dqi+vyaZSHj/v3Ip6llnIzN7Bxs6am/+AWucsFYYQFS4WaR1rHfGcSDRMuLFjhlSPCIizCIizCohAWYREWYREWYVEIi7AIi7AI66XyA2Yom/1j6CTHAAAAAElFTkSuQmCC")
	return Image{
		Content: ImageContent{
			Content: imageContent,
		},
		Metadata: ImageMetadata{
			Name:        "Toto",
			Description: "Toto's head",
		},
	}
}

func generateImage2() Image {
	imageContent, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAIAAAD/gAIDAAABqklEQVR42u3dTRLCIAyGYUGO5Wk8mKfxWh3cdqO2JYH8vNk5burjByY41dJ7v1HHqkIAFlhggQUWWBCABdbqatYuqDzv+4f9tZEslqF0rMAiWWBRYIEFFljn+gZTzQTJOvN2GvnC4m+CLMw9bpJVnvflS9JEsi4oLAnaeiyRvMyxqwGkpn1orkyWxitUjVi01kH1c6BGipV2NZjsJsuv1NRkuWaahxWAiUEaLM0sh02WRsMVfBnKesXfswS9UmzwYuP6tEHaQgMxOGbnah0Gd/2aJ1bjFzNjGZrt4M+uSjp4M7NhmKlQPVnBpBSx4klpYYWUEtuzourIJ8up1IVuvt0y1eC402CiKfXWZ1kbocHylixT9yKphksmWf21eSEb4tY7orHZf428qen2LIsnpQzS2cOlsmd5idX6Y2VHC/DspdKUrsOKfbBVkTpenJROxErCdB0rFZD8MvQ+9KXAmnzU0WCKhmXksMzKPdKPUr499Tbzs7OMO2CBBRZYYFFggQUWWGAFLSuD9H4A/DEnkiw3VfgrGZIFFlhggQUWBGCBBRZYYIEFwfH6ABAxbX1ln8i9AAAAAElFTkSuQmCC")
	return Image{
		Content: ImageContent{
			Content: imageContent,
		},
		Metadata: ImageMetadata{
			Name:        "Tree",
			Description: "",
		},
	}
}

func TestAddImageAndReturnID(t *testing.T) {
	image := generateImage1()

	imageID, err := minioRepository.AddImageAndReturnID(image)

	assert.NoError(t, err)
	assert.True(t, isValidUUID(imageID))

	imageReader, err := minioClient.GetObject(bucketName, imageID, minio.GetObjectOptions{})
	assert.NoError(t, err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(imageReader)
	assert.Equal(t, image.getBinaryContent(), buf.Bytes())

	imageMetadataReader, err := minioClient.GetObject(bucketName, "metadata/"+imageID+".json", minio.GetObjectOptions{})
	assert.NoError(t, err)

	buf = new(bytes.Buffer)
	buf.ReadFrom(imageMetadataReader)
	imageMetadata, _ := json.Marshal(ImageIDWithMetadata{ID: imageID, Metadata: image.Metadata})
	assert.Equal(t, imageMetadata, buf.Bytes())

	minioRepository.DeleteImageByID(imageID)
}

func TestGetImages(t *testing.T) {
	imageID1, _ := minioRepository.AddImageAndReturnID(generateImage1())
	imageID2, _ := minioRepository.AddImageAndReturnID(generateImage2())

	images, err := minioRepository.GetImages()

	assert.NoError(t, err)
	assert.Len(t, images, 2)

	expected := []ImageIDWithMetadata{
		{
			ID: imageID1,
			Metadata: ImageMetadata{
				Name:        "Toto",
				Description: "Toto's head",
			},
		},
		{
			ID: imageID2,
			Metadata: ImageMetadata{
				Name:        "Tree",
				Description: "",
			},
		},
	}

	assert.Equal(t, expected, images)

	minioRepository.DeleteImageByID(imageID1)
	minioRepository.DeleteImageByID(imageID2)
}

func TestGetImageByID(t *testing.T) {
	image := generateImage1()

	imageID, _ := minioRepository.AddImageAndReturnID(image)

	imageContent, err := minioRepository.GetImageByID(imageID)
	assert.NoError(t, err)
	assert.Equal(t, image.getBinaryContent(), imageContent.Content)

	minioRepository.DeleteImageByID(imageID)
}

func TestGetImageMetadataByID(t *testing.T) {
	image := generateImage1()

	imageID, _ := minioRepository.AddImageAndReturnID(image)

	imageMetadata, err := minioRepository.GetImageMetadataByID(imageID)
	assert.NoError(t, err)

	expected := &ImageIDWithMetadata{
		ID: imageID,
		Metadata: ImageMetadata{
			Name:        image.Metadata.Name,
			Description: image.Metadata.Description,
		},
	}
	assert.Equal(t, expected, imageMetadata)

	minioRepository.DeleteImageByID(imageID)
}

func TestDeleteImageByID(t *testing.T) {
	imageID, _ := minioRepository.AddImageAndReturnID(generateImage1())

	err := minioRepository.DeleteImageByID(imageID)

	assert.NoError(t, err)

	images, _ := minioRepository.GetImages()
	assert.Len(t, images, 0)
}
