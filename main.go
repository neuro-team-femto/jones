package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/neuro-team-femto/revcor/front"
	"github.com/neuro-team-femto/revcor/helpers"
	"github.com/neuro-team-femto/revcor/server"
)

var (
	port string
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

	// port
	port = os.Getenv("APP_PORT")
	if len(port) == 0 {
		port = "8100"
	}

	server := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("[main] http listening on port %v\n", port)
	log.Fatal(server.ListenAndServe()) // blocking
}

func main() {
	front.Build()

	if os.Getenv("APP_ENV") != "BUILD_FRONT" {
		// launch http (with websockets) server
		runServer()
	}
}
