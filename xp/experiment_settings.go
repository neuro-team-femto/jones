package xp

import (
	"encoding/json"
	"log"
	"os"

	"github.com/neuro-team-femto/jones/helpers"
)

type Field struct {
	Key       string `json:"key"`
	Label     string `json:"label"`
	InputType string `json:"inputType"`
	Pattern   string `json:"pattern"`
	Min       string `json:"min"`
	Max       string `json:"max"`
}

type ExperimentSettings struct {
	Id             string  `json:"id"`
	Kind           string  `json:"kind"`
	NInterval      int     `json:"nInterval"`
	FileExtension  string  `json:"fileExtension"`
	TrialsPerBlock int     `json:"trialsPerBlock"`
	BlocksPerXp    int     `json:"blocksPerXp"`
	AddRepeatBlock bool    `json:"addRepeatBlock"`
	AllowCreate    bool    `json:"allowCreate"`
	CreatePassword string  `json:"createPassword"`
	AdminPassword  string  `json:"adminPassword"`
	ShowProgress   bool    `json:"showProgress"`
	ForceWidth     string  `json:"forceWidth"`
	CollectInfo    []Field `json:"collectInfo"`
}

// API

func (es ExperimentSettings) CollectsInfo() bool {
	return len(es.CollectInfo) > 0
}

// Check if ids sent by client are valid (match a regex + configuration file exists)
func IsExperimentValid(experimentId string) bool {
	// check IDs format
	if !helpers.IsIdValid(experimentId) {
		return false
	}
	// check config exisis
	return helpers.PathExists("data/" + experimentId + "/config/settings.json")
}

// process with default values
func GetExperimentSettings(experimentId string) (e ExperimentSettings, err error) {
	settingsPath := "data/" + experimentId + "/config/settings.json"
	file, err := os.ReadFile(settingsPath)
	if err != nil {
		log.Printf("[error][GetExperimentSettings] unable to read path '%v' error: %+v\n", settingsPath, err)
		return
	}

	e = ExperimentSettings{}
	if err = json.Unmarshal([]byte(file), &e); err != nil {
		return
	}
	e.Id = experimentId
	if e.Kind != "image" {
		e.Kind = "sound"
	}
	if e.NInterval != 1 {
		e.NInterval = 2
	}
	if len(e.FileExtension) == 0 {
		if e.Kind == "sound" {
			e.FileExtension = "wav"
		} else {
			e.FileExtension = "png"
		}
	}
	return
}

func GetSanitizedExperimentSettings(experimentId string) (e ExperimentSettings, err error) {
	e, err = GetExperimentSettings(experimentId)
	// sanitize
	e.AdminPassword = ""
	e.CreatePassword = ""
	return
}

func GetExperimentWordingRunString(experimentId string) (json string, err error) {
	wordingRunPath := "data/" + experimentId + "/config/wording.run.json"
	return helpers.ReadTrimJSON(wordingRunPath)
}

// no error is returned
func GetExperimentWordingNewMap(experimentId string) (m map[string]string) {
	wordingNewPath := "data/" + experimentId + "/config/wording.new.json"
	file, err := os.ReadFile(wordingNewPath)
	if err != nil {
		log.Printf("[error][GetExperimentWordingNewMap] unable to read path '%v' error: %+v\n", wordingNewPath, err)
		return
	}

	json.Unmarshal([]byte(file), &m)
	return
}
