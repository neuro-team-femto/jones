package server

import (
	"html/template"
	"net/http"
)

var (
	runTemplate *template.Template
	newTemplate *template.Template
)

type templateData map[string]interface{}

func init() {
	runTemplate = template.Must(template.ParseFiles("public/templates/run.html.gtpl"))
	newTemplate = template.Must(template.ParseFiles("public/templates/new.html.gtpl"))
}

func renderNew(w http.ResponseWriter, experimentId, id string, wording map[string]string, errorLabl string) {
	newTemplate.Execute(w, templateData{
		"webPrefix":    webPrefix,
		"experimentId": experimentId,
		"id":           id,
		"wording":      wording,
		"error":        wording[errorLabl],
	})
}

func renderRun(w http.ResponseWriter, experimentId, participantId string) {
	runTemplate.Execute(w, templateData{
		"webPrefix":     webPrefix,
		"experimentId":  experimentId,
		"participantId": participantId,
	})
}
