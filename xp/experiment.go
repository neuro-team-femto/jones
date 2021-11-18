package xp

import (
	"encoding/json"
	"io/ioutil"
)

type Experiment struct {
	Id              string `json:"id"`
	DisplayName     string `json:"displayName"`
	Introduction    string `json:"introduction"`
	TrialQuestion   string `json:"trialQuestion"`
	TrialSoundLabel string `json:"trialSoundLabel"`
	BlockCount      int    `json:"blockCount"`
	TrialsPerBlock  int    `json:"trialsPerBlock"`
}

func LoadExperiment(experimentId string) (e Experiment, err error) {
	configPath := "data/" + experimentId + "/config/xp.json"
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}

	e = Experiment{}
	if err = json.Unmarshal([]byte(file), &e); err != nil {
		return
	}
	e.Id = experimentId
	return
}
