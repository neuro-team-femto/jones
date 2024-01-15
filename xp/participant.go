package xp

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/neuro-team-femto/revcor/helpers"
)

type Participant struct {
	Id           string `json:"id"`
	ExperimentId string `json:"experimentId"`
	Todo         string `json:"todo"`
	Age          string `json:"age"`
	Sex          string `json:"sex"`
}

func getStateFolder(experimentId, participantId string) string {
	return "data/" + experimentId + "/state/" + participantId + "/"
}

func initParticipantWithInfo(es ExperimentSettings, participantId string) (Participant, error) {
	p := Participant{Id: participantId, ExperimentId: es.Id}
	stateFolder := getStateFolder(es.Id, participantId)
	helpers.EnsureFolder(stateFolder)
	infoPath := stateFolder + "info.json"

	// check if participant is new (not considered an error!)
	if _, silentErr := os.Stat(infoPath); errors.Is(silentErr, os.ErrNotExist) {
		return p, nil
	}

	file, err := os.ReadFile(infoPath)
	if err != nil {
		return p, err
	}

	if err = json.Unmarshal([]byte(file), &p); err != nil {
		return p, err
	}

	return p, nil
}

func truncatedInPlaceShuffle(input []string, max int) []string {
	if len(input) == 0 {
		return nil
	}
	log.Println(input, max)
	rand.Shuffle(len(input), func(i, j int) {
		input[i], input[j] = input[j], input[i]
	})
	return input[:max]
}

// if participant state is empty, generate the complete list of assets that compose a run
// the length os this state is 2 (two assets to be compared for each trial) * TrialsPerBlock * BlocksPerXp
func generateTodo(es ExperimentSettings, participantId string) (todos []string) {
	length := 2 * es.TrialsPerBlock * es.BlocksPerXp

	allAssetsPath := "data/" + es.Id + "/assets"
	assets := helpers.FindFilesUnder(allAssetsPath, "."+es.FileExtension)
	todos = truncatedInPlaceShuffle(assets, length)

	if es.AddRepeatBlock {
		// duplicate trials from last block
		repeat := todos[length-(2*es.TrialsPerBlock):]
		todos = append(todos, repeat...)
	}

	return
}

func getParticipantTodo(es ExperimentSettings, participantId string) (todo []string, err error) {
	todoPath := getStateFolder(es.Id, participantId) + "todo.txt"

	if helpers.PathExists(todoPath) {
		// load from state
		todo, err = helpers.ReadFileLines(todoPath)
	} else {
		// create and save state
		todo = generateTodo(es, participantId)
		state := strings.Join(todo[:], "\n")
		err = os.WriteFile(todoPath, []byte(state), 0644)
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

func LoadParticipant(es ExperimentSettings, participantId string) (p Participant, err error) {
	p, err = initParticipantWithInfo(es, participantId)
	if err != nil {
		return
	}
	// add fields
	todo, err := getParticipantTodo(es, participantId)
	p.Todo = strings.Join(todo, ",")
	return
}

func (p *Participant) UpdateInfo(age, sex string) (err error) {
	// update p
	p.Age = age
	p.Sex = sex

	// save to file (filter todo field)
	toSave := map[string]string{
		"id":           p.Id,
		"experimentId": p.ExperimentId,
		"age":          age,
		"sex":          sex,
	}
	contents, err := json.MarshalIndent(toSave, "", "  ")
	if err != nil {
		return
	}

	infoPath := getStateFolder(p.ExperimentId, p.Id) + "info.json"
	err = os.WriteFile(infoPath, contents, 0644)
	return
}

func (p *Participant) UpdateTodo(stimuli1, stimuli2 string) (err error) {
	todoPath := getStateFolder(p.ExperimentId, p.Id) + "todo.txt"
	if helpers.PathExists(todoPath) {
		helpers.RemoveOnceFromFile(todoPath, stimuli1, stimuli2)
	} else {
		return errors.New("missing-todo")
	}
	return
}
