package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (a App) handlerImageMetadata(w http.ResponseWriter, r *http.Request) {
	imageID := mux.Vars(r)["imageID"]

	imageIDWithMetadata, err := a.repository.GetImageMetadataByID(imageID)
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

	resultJSON, _ := json.Marshal(imageIDWithMetadata)

	w.Header().Set("Content-Type", "application/json")
	w.Write(resultJSON)
}
