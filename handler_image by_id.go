package main

import (
	"net/http"
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
	a.handlerNotImplemented(w, r)
}

func (a App) handlerDeleteImage(w http.ResponseWriter, r *http.Request) {
	a.handlerNotImplemented(w, r)
}
