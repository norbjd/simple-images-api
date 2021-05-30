package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a App) handlerImageByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.handlerGetImage(w, r)
	case http.MethodDelete:
		a.handlerDeleteImage(w, r)
	}
}

func (a App) handlerGetImage(w http.ResponseWriter, r *http.Request) {
	imageID := mux.Vars(r)["imageID"]

	imageContent, err := a.repository.GetImageByID(imageID)
	if err != nil {
		if err == ErrImageDoesNotExist {
			w.WriteHeader(404)
			w.Write([]byte("Image does not exist"))
			return
		}
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Add("Content-Type", "image")
	w.Write(imageContent.Content)
}

func (a App) handlerDeleteImage(w http.ResponseWriter, r *http.Request) {
	a.handlerNotImplemented(w, r)
}
