package main

import (
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
	a.handlerNotImplemented(w, r)
}
