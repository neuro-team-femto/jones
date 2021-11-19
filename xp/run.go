package xp

import (
	"regexp"

	"github.com/creamlab/revcor/helpers"
)

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
