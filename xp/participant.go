package xp

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"slices"

	"github.com/neuro-team-femto/revcor/helpers"
)

type StrMap map[string]string

type Participant struct {
	Id            string   `json:"id"`
	ExperimentId  string   `json:"experimentId"`
	Todo          []string `json:"todo"`
	InfoCollected bool     `json:"infoCollected"`
	Info          StrMap   `json:"info,omitempty"`
	// not serialized
	infoKeys []string
}

func getStateFolder(experimentId string) string {
	return "data/" + experimentId + "/state/"
}

func (p *Participant) getStateFile() string {
	return getStateFolder(p.ExperimentId) + p.Id + ".json"
}

func (p *Participant) getInfoKeys() []string {
	if len(p.infoKeys) == 0 {
		var keys []string
		for k := range p.Info {
			keys = append(keys, k)
		}
		// ensure order
		slices.Sort(keys)
		p.infoKeys = keys
	}
	return p.infoKeys
}

func (p *Participant) getInfoValues() (values []string) {
	for _, k := range p.getInfoKeys() {
		values = append(values, p.Info[k])
	}
	return
}

func truncatedInPlaceShuffle(input []string, max int) []string {
	if len(input) == 0 {
		return nil
	}
	rand.Shuffle(len(input), func(i, j int) {
		input[i], input[j] = input[j], input[i]
	})
	return input[:max]
}

// if participant state is empty, generate the complete list of assets that compose a run
// the length os this state is:
// - for NInterval==2 -> (two assets to be compared for each trial) * TrialsPerBlock * BlocksPerXp
// - for NInterval==1 -> TrialsPerBlock * BlocksPerXp
func generateTodo(es ExperimentSettings, participantId string) (todos []string) {
	intervalFactor := 2
	if es.NInterval == 1 {
		intervalFactor = 1
	}
	length := intervalFactor * es.TrialsPerBlock * es.BlocksPerXp

	allAssetsPath := "data/" + es.Id + "/assets"
	assets := helpers.FindFilesUnder(allAssetsPath, "."+es.FileExtension)
	todos = truncatedInPlaceShuffle(assets, length)

	if es.AddRepeatBlock {
		// duplicate trials from last block
		repeat := todos[length-(intervalFactor*es.TrialsPerBlock):]
		todos = append(todos, repeat...)
	}

	return
}

func (p *Participant) saveState() (err error) {
	contents, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return
	}

	if err = os.WriteFile(p.getStateFile(), contents, 0644); err != nil {
		log.Printf("[error][saveState] unable to write path '%v' error: %+v\n", p.getStateFile(), err)
		return
	}
	return
}

// API

func DoesParticipantExist(experimentId, participantId string) bool {
	if !helpers.IsIdValid(participantId) {
		return false
	}

	configPath := "data/" + experimentId + "/config/"
	// check participant exists
	participantPaths := helpers.FindFilesUnder(configPath, "participants")
	for _, p := range participantPaths {
		if helpers.IsLineInFile(configPath+p, participantId) {
			return true
		}
	}

	return false
}

// Check if ids sent by client are valid (match a regex + configuration file exists)
func IsParticipantValid(experimentId, participantId string) bool {
	if !IsExperimentValid(experimentId) {
		return false
	}
	return DoesParticipantExist(experimentId, participantId)
}

func InitParticipant(es ExperimentSettings, participantId string) (p Participant, err error) {
	p = Participant{Id: participantId, ExperimentId: es.Id}

	stateFolder := getStateFolder(es.Id)
	err = helpers.EnsureFolder(stateFolder)
	if err != nil {
		log.Printf("[error] could not ensure state folder %v: %+v\n", stateFolder, err)
		return p, err
	}

	stateFile := p.getStateFile()
	// check if participant is new (not considered an error!)
	if _, silentErr := os.Stat(stateFile); errors.Is(silentErr, os.ErrNotExist) {
		p.InfoCollected = false
		p.Todo = generateTodo(es, participantId)
		p.saveState()
		return p, nil
	} else {
		file, err := os.ReadFile(stateFile)
		if err != nil {
			log.Printf("[error][InitParticipant] unable to read path '%v' error: %+v\n", stateFile, err)
			return p, err
		}
		if err = json.Unmarshal([]byte(file), &p); err != nil {
			return p, err
		}
		return p, nil
	}

}

func (p *Participant) UpdateInfo(info StrMap) (err error) {
	p.Info = info
	p.InfoCollected = true
	return p.saveState()
}

func (p *Participant) UpdateTodo(stimuli string) (err error) {
	var newTodo []string
	// remove only once
	removed := false
	for _, t := range p.Todo {
		if !removed && t == stimuli {
			removed = true
		} else {
			newTodo = append(newTodo, t)
		}
	}
	p.Todo = newTodo
	return p.saveState()
}
