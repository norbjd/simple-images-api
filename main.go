package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	router := InitRouter()

	log.Info("Server started")
	log.Fatal(http.ListenAndServe(":8080", router))
}
