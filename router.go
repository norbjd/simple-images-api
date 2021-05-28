package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func handlerNotImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	w.Write([]byte("Not implemented"))
}

func InitRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/images", handlerNotImplemented).Methods("GET", "POST")
	router.HandleFunc("/images/{imageID}", handlerNotImplemented).Methods("GET", "DELETE")
	router.HandleFunc("/images/{imageID}/metadata", handlerNotImplemented).Methods("GET")

	return router
}
