package xp

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/creamlab/revcor/helpers"
)

type Participant struct {
	Id     string `json:"id"`
	Sounds string `json:"sounds"`
}

func LoadParticipant(es ExperimentSettings, participantId string) (p Participant, err error) {
	sounds, err := getParticipantSounds(es, participantId)
	p = Participant{
		Id:     participantId,
		Sounds: strings.Join(sounds, ","),
	}
	return
}

// if participant state is empty, generate the complete list of sounds that compose a run
// the length os this state is 2 (two sounds to be compared for each trial) * TrialsPerBlock * BlockCount
func selectSounds(es ExperimentSettings, participantId string) []string {
	length := 2 * es.TrialsPerBlock * es.BlockCount

	soundsPath := "data/" + es.Id + "/sounds"
	sounds := helpers.FindFilesUnder(soundsPath, ".wav")
	return truncatedInPlaceShuffle(sounds, length)
}

func getParticipantSounds(es ExperimentSettings, participantId string) (sounds []string, err error) {
	folder := "state/" + es.Id
	helpers.EnsureFolder(folder)

	path := folder + "/" + participantId + ".txt"

	if helpers.PathExists(path) {
		// load from state
		sounds, err = helpers.ReadFileLines(path)
	} else {
		// create and save state
		sounds = selectSounds(es, participantId)
		state := strings.Join(sounds[:], "\n")
		err = ioutil.WriteFile(path, []byte(state), 0644)
	}
	return
}

func UpdateParticipantState(es ExperimentSettings, p Participant, stimuli1, stimuli2 string) (err error) {
	path := "state/" + es.Id + "/" + p.Id + ".txt"
	if helpers.PathExists(path) {
		helpers.RemoveLinesFromFile(path, stimuli1, stimuli2)
	} else {
		return errors.New("missing-state")
	}
	return
}
