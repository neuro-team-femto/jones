package xp

import (
	"encoding/json"
	"io/ioutil"

	"github.com/creamlab/revcor/helpers"
)

type ExperimentSettings struct {
	Id             string `json:"id"`
	BlockCount     int    `json:"blockCount"`
	TrialsPerBlock int    `json:"trialsPerBlock"`
}

type experimentWording interface{}

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
