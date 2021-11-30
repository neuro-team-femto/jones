package helpers

import "regexp"

var (
	validID = regexp.MustCompile("^[a-zA-Z0-9-_]+$")
)

func IsIdValid(id string) bool {
	// check IDs format
	return validID.MatchString(id)
}
