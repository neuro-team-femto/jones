package xp

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/creamlab/revcor/helpers"
)

type Participant struct {
	Id   string `json:"id"`
	XpId string `json:"xpId"`
	Todo string `json:"todo"`
	Age  string `json:"age"`
	Sex  string `json:"sex"`
}

func initParticipantWithInfo(es ExperimentSettings, participantId string) (Participant, error) {
	infoPath := "state/" + es.Id + "/" + participantId + "/info.json"
	p := Participant{}

	if _, silentErr := os.Stat(infoPath); errors.Is(silentErr, os.ErrNotExist) {
		// not considered an error
		return p, nil
	}

	file, err := ioutil.ReadFile(infoPath)
	if err != nil {
		return p, err
	}

	if err = json.Unmarshal([]byte(file), &p); err != nil {
		return p, err
	}

	return p, nil
}

func truncatedInPlaceShuffle(input []string, max int) []string {
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(input), func(i, j int) {
		input[i], input[j] = input[j], input[i]
	})
	return input[:max]
}

// if participant state is empty, generate the complete list of sounds that compose a run
// the length os this state is 2 (two sounds to be compared for each trial) * TrialsPerBlock * BlockCount
func genTodo(es ExperimentSettings, participantId string) []string {
	length := 2 * es.TrialsPerBlock * es.BlockCount

	allSoundsPath := "data/" + es.Id + "/sounds"
	sounds := helpers.FindFilesUnder(allSoundsPath, ".wav")
	return truncatedInPlaceShuffle(sounds, length)
}

func getParticipantTodo(es ExperimentSettings, participantId string) (todo []string, err error) {
	folder := "state/" + es.Id + "/" + participantId
	helpers.EnsureFolder(folder)

	todoPath := folder + "/todo.txt"

	if helpers.PathExists(todoPath) {
		// load from state
		todo, err = helpers.ReadFileLines(todoPath)
	} else {
		// create and save state
		todo = genTodo(es, participantId)
		state := strings.Join(todo[:], "\n")
		err = ioutil.WriteFile(todoPath, []byte(state), 0644)
	}
	return
}

// API

func LoadParticipant(es ExperimentSettings, participantId string) (p Participant, err error) {
	p, err = initParticipantWithInfo(es, participantId)
	if err != nil {
		return
	}
	// add fields
	p.Id = participantId
	p.XpId = es.Id
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
		"id":   p.Id,
		"xpId": p.XpId,
		"Age":  age,
		"Sex":  sex,
	}
	contents, err := json.MarshalIndent(toSave, "", "  ")
	if err != nil {
		return
	}

	infoPath := "state/" + p.XpId + "/" + p.Id + "/info.json"
	err = ioutil.WriteFile(infoPath, contents, 0644)
	return
}

func (p *Participant) UpdateTodo(stimuli1, stimuli2 string) (err error) {
	todoPath := "state/" + p.XpId + "/" + p.Id + "/todo.txt"
	if helpers.PathExists(todoPath) {
		helpers.RemoveLinesFromFile(todoPath, stimuli1, stimuli2)
	} else {
		return errors.New("missing-todo")
	}
	return
}
