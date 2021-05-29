package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type App struct {
	repository Repository
}

func (a App) run() {
	router := InitRouter()

	log.Info("Server started")
	log.Fatal(http.ListenAndServe(":8080", router))
}
