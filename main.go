package main

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := mux.NewRouter()

	log.Info("Server started")
	log.Fatal(http.ListenAndServe(":8080", router))
}
