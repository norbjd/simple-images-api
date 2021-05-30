package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (a App) handlerImages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.handlerAddImage(w, r)
	case http.MethodGet:
		a.handlerGetImages(w, r)
	}
}

func (a App) extractImageFromRequest(r *http.Request) (*Image, error) {
	file, _, err := r.FormFile("image")
	if err != nil {
		log.Error(err)
		return nil, fmt.Errorf("cannot retrieve image")
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Error(err)
		return nil, fmt.Errorf("internal error")
	}

	return &Image{
		Content: ImageContent{
			Content: buf.Bytes(),
		},
		Metadata: ImageMetadata{
			Name:        r.Form.Get("name"),
			Description: r.Form.Get("description"),
		},
	}, nil
}

func (a App) handlerAddImage(w http.ResponseWriter, r *http.Request) {
	image, err := a.extractImageFromRequest(r)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}

	imageID, err := a.repository.AddImageAndReturnID(*image)
	if err != nil {
		log.Error("Cannot add image")
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}

	result := ImageIDWithMetadata{
		ID:       imageID,
		Metadata: image.Metadata,
	}

	resultJSON, _ := json.Marshal(result)

	w.Write(resultJSON)
}

func (a App) handlerGetImages(w http.ResponseWriter, r *http.Request) {
	images, err := a.repository.GetImages()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
	}

	resultJSON, _ := json.Marshal(images)

	w.Write(resultJSON)
}
