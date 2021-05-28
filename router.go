package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func handlerWithLogging(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return handlers.LoggingHandler(
		log.StandardLogger().Writer(),
		http.HandlerFunc(handler),
	)
}

func handlerNotImplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
	w.Write([]byte("Not implemented"))
}

func InitRouter() *mux.Router {
	router := mux.NewRouter()

	router.Handle("/images", handlerWithLogging(handlerNotImplemented)).Methods("GET", "POST")
	router.Handle("/images/{imageID}", handlerWithLogging(handlerNotImplemented)).Methods("GET", "DELETE")
	router.Handle("/images/{imageID}/metadata", handlerWithLogging(handlerNotImplemented)).Methods("GET")
	router.PathPrefix("/").Handler(handlerWithLogging(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))

	return router
}
