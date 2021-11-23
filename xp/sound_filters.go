package xp

import (
	"bufio"
	"os"
	"strings"
)

type filter struct {
	Freq string
	Gain string
}

type filters []filter

func LoadFilters(experimentId, wavFile string) (fs filters, err error) {
	filterFile := strings.TrimSuffix(wavFile, ".wav") + ".txt"
	filterPath := "data/" + experimentId + "/sounds/" + filterFile
	file, err := os.Open(filterPath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		raw := strings.Split(scanner.Text(), ",")
		fs = append(fs, filter{
			Freq: raw[0],
			Gain: raw[1],
		})
	}
	return
}
