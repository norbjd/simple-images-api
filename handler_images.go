package main

import (
	"encoding/json"
	"net/http"
)

func (a App) handlerImages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.handlerAddImage(w, r)
	case http.MethodGet:
		a.handlerGetImages(w, r)
	}
}

func (a App) handlerAddImage(w http.ResponseWriter, r *http.Request) {
	a.handlerNotImplemented(w, r)
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
