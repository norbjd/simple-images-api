package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAcceptance(t *testing.T) {
	app := App{repository: MinioRepository{bucketName: "images", client: minioClient}}
	router := app.initRouter()

	// insert an image
	image1Bytes, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAIAAAD/gAIDAAACu0lEQVR42u2cXXaEIAyFgeMiZYW6S/owL+2oGDAJF3vzOJ0ewkf+xYmllECRSSICwiIswiIswiIsIiAswiIswiIswiICwjKRZRZF9xiPH66+86UIMs86ZSEUN2QQsJ6Q8iS4vA/TOy3LmdRaynHFJgMcY1kxxm3Cs1mGaLmFKWWZyNHyBevsdQALOKYs/nDWCh4hx+UpLKuDVL52okqyjn8X6qMTY0M9kBBI9Um5E8laWwh7jFGmdhpLar3emHoVXjmVTab84kMqjygavlivB889bqF+PIuzTeVqGLJu7spZEY9bwft43JPWonKWCYGUc4dYOZh68FJopCWGnTFINbWoR6b/eqzcen7msDKqWXV4TLImFeYk5Z0NfUh9lU6ty30Ck3C+lkxJSbb6W8A9MVlgyki1gmKk50NWpGxIWDXZfH3wdiyDng037WBBN6Qbzim726R0nbxAl0+4dGJWnvO5aesgUC3A52nDNmMWYb146jC7HBMXYZ3H39OwSzeUkgpaN/9aJ3DH0ZV6ApWrJFcGxbJGTf6a1gVyQy1eD4euKA9ZrffZSrx1rQRoI93/bu3LYyyrIwOom+Tp94vDLZpuXpUdfv6kglUr1ZqUDk0q9u1WbkenavQprPaGxcPSySLcXNaWF2vdKmwYs5r2r16UqpPCKh20eNk95tF80elJ5FJxye5sIFQScepwmyh1M4n8OJVfodMyLrcA2qRb8lF07A0ZrdWThQfNMleAuyY51sQUSZnAUu/7QEhZWRYCL3VSwfqF8lsoFqvXF32yYnS+IGs+GDAjFXx+qsChSfQZ1Pj9rkNrnJIopl7lo8Dqi+vyaZSHj/v3Ip6llnIzN7Bxs6am/+AWucsFYYQFS4WaR1rHfGcSDRMuLFjhlSPCIizCIizCohAWYREWYREWYVEIi7AIi7AI66XyA2Yom/1j6CTHAAAAAElFTkSuQmCC")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, insertImageRequest(image1Bytes, "Toto", ""))

	assert.Equal(t, 200, recorder.Code)
	var image1IDWithMetadata ImageIDWithMetadata
	json.Unmarshal(recorder.Body.Bytes(), &image1IDWithMetadata)

	assert.True(t, isValidUUID(image1IDWithMetadata.ID))
	assert.Equal(t, "Toto", image1IDWithMetadata.Metadata.Name)
	assert.Equal(t, "", image1IDWithMetadata.Metadata.Description)

	// insert another image
	image2Bytes, _ := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAGQAAABkCAIAAAD/gAIDAAABqklEQVR42u3dTRLCIAyGYUGO5Wk8mKfxWh3cdqO2JYH8vNk5burjByY41dJ7v1HHqkIAFlhggQUWWBCABdbqatYuqDzv+4f9tZEslqF0rMAiWWBRYIEFFljn+gZTzQTJOvN2GvnC4m+CLMw9bpJVnvflS9JEsi4oLAnaeiyRvMyxqwGkpn1orkyWxitUjVi01kH1c6BGipV2NZjsJsuv1NRkuWaahxWAiUEaLM0sh02WRsMVfBnKesXfswS9UmzwYuP6tEHaQgMxOGbnah0Gd/2aJ1bjFzNjGZrt4M+uSjp4M7NhmKlQPVnBpBSx4klpYYWUEtuzourIJ8up1IVuvt0y1eC402CiKfXWZ1kbocHylixT9yKphksmWf21eSEb4tY7orHZf428qen2LIsnpQzS2cOlsmd5idX6Y2VHC/DspdKUrsOKfbBVkTpenJROxErCdB0rFZD8MvQ+9KXAmnzU0WCKhmXksMzKPdKPUr499Tbzs7OMO2CBBRZYYFFggQUWWGAFLSuD9H4A/DEnkiw3VfgrGZIFFlhggQUWBGCBBRZYYIEFwfH6ABAxbX1ln8i9AAAAAElFTkSuQmCC")

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, insertImageRequest(image2Bytes, "Tree", "A beautiful tree"))

	assert.Equal(t, 200, recorder.Code)
	var image2IDWithMetadata ImageIDWithMetadata
	json.Unmarshal(recorder.Body.Bytes(), &image2IDWithMetadata)

	assert.True(t, isValidUUID(image2IDWithMetadata.ID))
	assert.Equal(t, "Tree", image2IDWithMetadata.Metadata.Name)
	assert.Equal(t, "A beautiful tree", image2IDWithMetadata.Metadata.Description)

	// retrieve all images
	recorder = httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/images", nil)
	router.ServeHTTP(recorder, req)

	var imagesIDWithMetadata []ImageIDWithMetadata
	json.Unmarshal(recorder.Body.Bytes(), &imagesIDWithMetadata)

	expected := []ImageIDWithMetadata{image1IDWithMetadata, image2IDWithMetadata}

	assert.Equal(t, 200, recorder.Code)
	assert.ElementsMatch(t, expected, imagesIDWithMetadata)

	// retrieve first image
	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/images/"+image1IDWithMetadata.ID, nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, image1Bytes, recorder.Body.Bytes())

	// retrieve first image metadata
	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/images/"+image1IDWithMetadata.ID+"/metadata", nil)
	router.ServeHTTP(recorder, req)

	var imageIDWithMetadata ImageIDWithMetadata
	json.Unmarshal(recorder.Body.Bytes(), &imageIDWithMetadata)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, image1IDWithMetadata, imageIDWithMetadata)

	// delete first image
	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/images/"+image1IDWithMetadata.ID, nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)

	// delete second image
	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/images/"+image2IDWithMetadata.ID, nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)

	// get all images
	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/images", nil)
	router.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "[]", recorder.Body.String())
}
