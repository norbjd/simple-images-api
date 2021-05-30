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

func TestAddImageAndReturnID(t *testing.T) {
	imageContent, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAIAAAD/gAIDAAACu0lEQVR42u2cXXaEIAyFgeMiZYW6S/owL+2oGDAJF3vzOJ0ewkf+xYmllECRSSICwiIswiIswiIsIiAswiIswiIswiICwjKRZRZF9xiPH66+86UIMs86ZSEUN2QQsJ6Q8iS4vA/TOy3LmdRaynHFJgMcY1kxxm3Cs1mGaLmFKWWZyNHyBevsdQALOKYs/nDWCh4hx+UpLKuDVL52okqyjn8X6qMTY0M9kBBI9Um5E8laWwh7jFGmdhpLar3emHoVXjmVTab84kMqjygavlivB889bqF+PIuzTeVqGLJu7spZEY9bwft43JPWonKWCYGUc4dYOZh68FJopCWGnTFINbWoR6b/eqzcen7msDKqWXV4TLImFeYk5Z0NfUh9lU6ty30Ck3C+lkxJSbb6W8A9MVlgyki1gmKk50NWpGxIWDXZfH3wdiyDng037WBBN6Qbzim726R0nbxAl0+4dGJWnvO5aesgUC3A52nDNmMWYb146jC7HBMXYZ3H39OwSzeUkgpaN/9aJ3DH0ZV6ApWrJFcGxbJGTf6a1gVyQy1eD4euKA9ZrffZSrx1rQRoI93/bu3LYyyrIwOom+Tp94vDLZpuXpUdfv6kglUr1ZqUDk0q9u1WbkenavQprPaGxcPSySLcXNaWF2vdKmwYs5r2r16UqpPCKh20eNk95tF80elJ5FJxye5sIFQScepwmyh1M4n8OJVfodMyLrcA2qRb8lF07A0ZrdWThQfNMleAuyY51sQUSZnAUu/7QEhZWRYCL3VSwfqF8lsoFqvXF32yYnS+IGs+GDAjFXx+qsChSfQZ1Pj9rkNrnJIopl7lo8Dqi+vyaZSHj/v3Ip6llnIzN7Bxs6am/+AWucsFYYQFS4WaR1rHfGcSDRMuLFjhlSPCIizCIizCohAWYREWYREWYVEIi7AIi7AI66XyA2Yom/1j6CTHAAAAAElFTkSuQmCC")
	image := Image{
		Content: ImageContent{
			Content: imageContent,
		},
		Metadata: ImageMetadata{
			Name:        "Toto",
			Description: "Toto's head",
		},
	}

	imageID, err := minioRepository.AddImageAndReturnID(image)

	assert.NoError(t, err)
	assert.True(t, isValidUUID(imageID))

	imageReader, err := minioClient.GetObject(bucketName, imageID, minio.GetObjectOptions{})
	assert.NoError(t, err)

	buf := new(bytes.Buffer)
	buf.ReadFrom(imageReader)
	assert.Equal(t, imageContent, buf.Bytes())

	imageMetadataReader, err := minioClient.GetObject(bucketName, "metadata/"+imageID + ".json", minio.GetObjectOptions{})
	assert.NoError(t, err)

	buf = new(bytes.Buffer)
	buf.ReadFrom(imageMetadataReader)
	imageMetadata, _ := json.Marshal(ImageIDWithMetadata{ID: imageID, Metadata: image.Metadata})
	assert.Equal(t, imageMetadata, buf.Bytes())
}
