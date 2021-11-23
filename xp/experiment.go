package xp

import (
	"encoding/json"
	"io/ioutil"
	"regexp"

	"github.com/creamlab/revcor/helpers"
)

type ExperimentSettings struct {
	Id             string `json:"id"`
	BlockCount     int    `json:"blockCount"`
	TrialsPerBlock int    `json:"trialsPerBlock"`
}

var (
	validID = regexp.MustCompile("[a-zA-Z0-9-_]+")
)

// Check if ids sent by client are valid (match a regex + configuration file exists)
func IsValid(experimentId, participantId string) bool {
	// check IDs format
	if !validID.MatchString(experimentId) {
		return false
	}
	if !validID.MatchString(participantId) {
		return false
	}

	configPath := "data/" + experimentId + "/config/"
	// check config exisis
	if !helpers.PathExists(configPath + "settings.json") {
		return false
	}

	// check participant exists
	participantPaths := helpers.FindFilesUnder(configPath, "participants")
	for _, p := range participantPaths {
		if helpers.IsLineInFile(configPath+p, participantId) {
			return true
		}
	}

	return false
}

func GetExperimentSettings(experimentId string) (e ExperimentSettings, err error) {
	settingsPath := "data/" + experimentId + "/config/settings.json"
	file, err := ioutil.ReadFile(settingsPath)
	if err != nil {
		return
	}

	e = ExperimentSettings{}
	if err = json.Unmarshal([]byte(file), &e); err != nil {
		return
	}
	e.Id = experimentId
	return
}

func GetExperimentWordingString(experimentId string) (json string, err error) {
	wordingPath := "data/" + experimentId + "/config/wording.json"
	return helpers.ReadTrimJSON(wordingPath)
}
