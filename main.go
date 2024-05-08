package main

import (
	"log"
	"net/http"
	"time"

	"github.com/neuro-team-femto/jones/config"
	"github.com/neuro-team-femto/jones/front"
	"github.com/neuro-team-femto/jones/helpers"
	"github.com/neuro-team-femto/jones/server"
)

func init() {
	// create data folder if needed
	err := helpers.EnsureFolder("data")
	if err != nil {
		log.Fatal(err)
	}
}

func runServer() {
	router := server.GetRouter()

	server := &http.Server{
		Handler:      router,
		Addr:         ":" + config.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("[server] http listening on port %v\n", config.Port)
	log.Fatal(server.ListenAndServe()) // blocking
}

func main() {
	front.Build()

	if config.Mode != "BUILD_FRONT" {
		// launch http (with websockets) server
		runServer()
	}
}
