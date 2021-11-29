package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/creamlab/revcor/front"
	"github.com/creamlab/revcor/helpers"
	"github.com/creamlab/revcor/server"
)

var (
	port string
)

func init() {
	// create state folder
	err := helpers.EnsureFolder("state")
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
