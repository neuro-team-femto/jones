package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/neuro-team-femto/revcor/helpers"
)

var (
	allowedOrigins = []string{}
	webPrefix      string
	upgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			log.Printf("[server] ws upgrade from origin: %v\n", origin)
			return helpers.Contains(allowedOrigins, origin)
		},
	}
)

func init() {
	// environment variables use
	envOrigins := os.Getenv("APP_ORIGINS")
	if len(envOrigins) > 0 {
		allowedOrigins = append(allowedOrigins, strings.Split(envOrigins, ",")...)
	}
	if os.Getenv("APP_ENV") == "DEV" {
		allowedOrigins = append(allowedOrigins, "http://localhost:8100", "https://localhost:8100")
	}

	// web prefix, for instance "/path" if DuckSoup is reachable at https://host/path
	webPrefix = helpers.Getenv("APP_WEB_PREFIX", "")

	// log
	log.Printf("[server] allowed ws origins: %v\n", allowedOrigins)
	log.Printf("[server] APP_WEB_PREFIX: %v\n", webPrefix)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("[server] not found: %v\n", r.URL)
}

func GetRouter() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)
	// public router with no auth
	publicRouter := router.PathPrefix(webPrefix).Subrouter()

	// websocket handler
	publicRouter.HandleFunc("/ws", websocketHandler)

	// serve assets (js & css) under front/public
	publicRouter.PathPrefix("/scripts/").Handler(http.StripPrefix(webPrefix+"/scripts/", http.FileServer(http.Dir("./public/scripts/"))))
	publicRouter.PathPrefix("/styles/").Handler(http.StripPrefix(webPrefix+"/styles/", http.FileServer(http.Dir("./public/styles/"))))

	// serve assets under data/{experimentId]/assets, with path rewrite
	publicRouter.HandleFunc("/xp/{experimentId:[a-zA-Z0-9-_]+}/assets/{file:.*}", soundHandler)

	// run xp
	publicRouter.HandleFunc("/xp/{experimentId:[a-zA-Z0-9-_]+}/run/{participantId:[a-zA-Z0-9-_]+}", runHandler)

	// new (create new participant id)
	publicRouter.HandleFunc("/xp/{experimentId:[a-zA-Z0-9-_]+}/new", createHandler).Methods("POST")
	publicRouter.HandleFunc("/xp/{experimentId:[a-zA-Z0-9-_]+}/new", newHandler).Methods("GET")

	// results router with authentication
	resultsRouter := router.PathPrefix(webPrefix).Subrouter()
	resultsRouter.Use(resultsAuthMiddleware)
	resultsRouter.
		PathPrefix("/xp/{experimentId:[a-zA-Z0-9-_]+}/results").
		Handler(http.StripPrefix(webPrefix+"/xp", http.FileServer(http.Dir("./data/"))))

	return router
}
