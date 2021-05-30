package main

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type App struct {
	repository Repository
}

func (a App) run() {
	router := a.initRouter()

	log.Info("Server started")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func (a App) initRouter() *mux.Router {
	router := mux.NewRouter()

	router.Handle("/images", a.handlerWithLogging(a.handlerImages)).Methods("GET", "POST")
	router.Handle("/images/{imageID}", a.handlerWithLogging(a.handlerImageByID)).Methods("GET", "DELETE")
	router.Handle("/images/{imageID}/metadata", a.handlerWithLogging(a.handlerImageMetadata)).Methods("GET")
	router.PathPrefix("/").Handler(a.handlerWithLogging(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))

	return router
}

func (a App) handlerWithLogging(handler func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return handlers.LoggingHandler(
		log.StandardLogger().WriterLevel(log.DebugLevel),
		http.HandlerFunc(handler),
	)
}
