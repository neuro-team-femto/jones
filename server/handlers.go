package server

import (
	"crypto/sha256"
	"crypto/subtle"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/neuro-team-femto/revcor/helpers"
	"github.com/neuro-team-femto/revcor/xp"
)

// handle incoming websockets
func websocketHandler(w http.ResponseWriter, r *http.Request) {
	// upgrade HTTP request to Websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[error] can't upgrade ws: %v\n", err)
		return
	}

	runWs(conn)
}

func assetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := "data/" + vars["experimentId"] + "/assets/" + vars["file"]
	http.ServeFile(w, r, path)
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	experimentId := vars["experimentId"]
	participantId := vars["participantId"]

	if xp.IsParticipantValid(experimentId, participantId) {
		renderRun(w, experimentId, participantId)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	experimentId := mux.Vars(r)["experimentId"]

	es, err := xp.GetExperimentSettings(experimentId)
	if err == nil && es.AllowCreate {
		wording := xp.GetExperimentWordingNewMap(experimentId)
		renderNew(w, experimentId, "", wording, "")
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	experimentId := mux.Vars(r)["experimentId"]
	id := r.FormValue("id")
	password := r.FormValue("password")

	wording := xp.GetExperimentWordingNewMap(experimentId)
	var errorLabel string

	es, err := xp.GetExperimentSettings(experimentId)
	if err != nil { // check if settings exists
		errorLabel = "serverError"
	} else if password != es.CreatePassword { // check if participant knows create password
		errorLabel = "wrongPassword"
	} else if !helpers.IsIdValid(id) { // check if new id is valid
		errorLabel = "invalidId"
	} else if !xp.DoesParticipantExist(experimentId, id) { // create participant ID only if not already exists
		path := "data/" + experimentId + "/config/new-participants.txt"
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			errorLabel = "serverError"
		} else if _, err := file.WriteString(id + "\n"); err != nil {
			errorLabel = "serverError"
		}
	}

	if len(errorLabel) > 0 {
		renderNew(w, experimentId, id, wording, errorLabel)
	} else {
		http.Redirect(w, r, webPrefix+"/xp/"+experimentId+"/run/"+id, http.StatusFound)
	}
}

func resultsAuthMiddleware(next http.Handler) http.Handler {
	// inspired from https://www.alexedwards.net/blog/basic-authentication-in-go
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		experimentId := mux.Vars(r)["experimentId"]
		if len(experimentId) > 0 {
			settings, err := xp.GetExperimentSettings(experimentId)
			if err == nil {
				login, password, ok := r.BasicAuth()
				if ok {
					// Calculate SHA-256 hashes for the provided and expected usernames and passwords.
					loginHash := sha256.Sum256([]byte(login))
					passwordHash := sha256.Sum256([]byte(password))
					expectedLoginHash := sha256.Sum256([]byte("admin"))
					expectedPasswordHash := sha256.Sum256([]byte(settings.AdminPassword))
					loginMatch := (subtle.ConstantTimeCompare(loginHash[:], expectedLoginHash[:]) == 1)
					passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

					if loginMatch && passwordMatch {
						// authentication succeeded
						next.ServeHTTP(w, r)
						return
					}
				}
			}
		}
		// authentication failed
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
