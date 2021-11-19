package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/creamlab/revcor/front"
	"github.com/creamlab/revcor/helpers"
	"github.com/creamlab/revcor/ws"
	"github.com/creamlab/revcor/xp"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var (
	port           string
	allowedOrigins = []string{}
	webPrefix      string
	indexTemplate  *template.Template
	upgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			log.Info().Msgf("[server] ws upgrade from origin: %v", origin)
			return helpers.Contains(allowedOrigins, origin)
		},
	}
)

type templateData map[string]interface{}

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

	indexTemplate = template.Must(template.ParseFiles("public/templates/index.html.gtpl"))

	// create state folder
	err := helpers.EnsureFolder("state")
	if err != nil {
		log.Fatal().Err(err)
	}

	// log
	log.Info().Msgf("[server] allowed ws origins: %v", allowedOrigins)
}

// handle incoming websockets
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade HTTP request to Websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("[server] can't upgrade ws")
		return
	}

	ws.RunServer(conn)
}

func basicAuthWith(refLogin, refPassword string) mux.MiddlewareFunc {
	// source https://www.alexedwards.net/blog/basic-authentication-in-go
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			login, password, ok := r.BasicAuth()
			if ok {
				// Calculate SHA-256 hashes for the provided and expected usernames and passwords.
				loginHash := sha256.Sum256([]byte(login))
				passwordHash := sha256.Sum256([]byte(password))
				expectedLoginHash := sha256.Sum256([]byte(refLogin))
				expectedPasswordHash := sha256.Sum256([]byte(refPassword))

				loginMatch := (subtle.ConstantTimeCompare(loginHash[:], expectedLoginHash[:]) == 1)
				passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

				if loginMatch && passwordMatch {
					next.ServeHTTP(w, r)
					return
				}
			}

			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})
	}
}

func soundHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := "data/" + vars["experimentId"] + "/sounds/" + vars["file"]
	http.ServeFile(w, r, path)
}

func xpHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	experimentId := vars["experimentId"]
	participantId := vars["participantId"]

	if xp.IsValid(experimentId, participantId) {
		indexTemplate.Execute(w, templateData{
			"webPrefix":     webPrefix,
			"experimentId":  experimentId,
			"participantId": participantId,
		})
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func runServer() {
	router := mux.NewRouter()
	// websocket handler
	router.HandleFunc(webPrefix+"/ws", websocketHandler)

	// serve assets (js & css) under front/public
	router.PathPrefix(webPrefix + "/scripts/").Handler(http.StripPrefix(webPrefix+"/scripts/", http.FileServer(http.Dir("./public/scripts/"))))
	router.PathPrefix(webPrefix + "/styles/").Handler(http.StripPrefix(webPrefix+"/styles/", http.FileServer(http.Dir("./public/styles/"))))
	// serve assets under data/{experimentId]/sounds, with rewrite
	router.
		PathPrefix(webPrefix + "/xp/{experimentId:[a-zA-Z0-9-_]+}/sounds/{file:.*}").
		HandlerFunc(soundHandler)

	// pages with basic auth
	authedRouter := router.PathPrefix(webPrefix + "/").Subrouter()
	//authedRouter.Use(basicAuthWith("test", "test")
	authedRouter.
		PathPrefix("/xp/{experimentId:[a-zA-Z0-9-_]+}/{participantId:[a-zA-Z0-9-_]+}").
		HandlerFunc(xpHandler)

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

	log.Info().Msgf("[server] http listening on port %v", port)
	log.Fatal().Err(server.ListenAndServe()) // blocking
}

func main() {
	front.Build()

	if os.Getenv("APP_ENV") != "BUILD_FRONT" {
		// launch http (with websockets) server
		runServer()
	}
}
