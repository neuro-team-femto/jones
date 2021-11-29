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

func getParamFile(experimentId, wavFile string) (file *os.File, err error) {
	filterFile := strings.TrimSuffix(wavFile, ".wav") + ".txt"
	filterPath := "data/" + experimentId + "/sounds/" + filterFile
	file, err = os.Open(filterPath)
	return
}

func ReadParamHeaders(experimentId, wavFile string) (headers []string, err error) {
	file, err := getParamFile(experimentId, wavFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// headers in first line
	scanner.Scan()
	headers = strings.Split(scanner.Text(), ",")
	return
}

func ReadParamValues(experimentId, wavFile string) (fs filters, err error) {
	file, err := getParamFile(experimentId, wavFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// skip first line
	scanner.Scan()

	for scanner.Scan() {
		raw := strings.Split(scanner.Text(), ",")
		fs = append(fs, filter{
			Freq: raw[0],
			Gain: raw[1],
		})
	}
	return
}
