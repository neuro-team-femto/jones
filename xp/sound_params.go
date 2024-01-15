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

func getParamFile(es ExperimentSettings, asset string) (file *os.File, err error) {
	filterFile := strings.TrimSuffix(asset, "."+es.FileExtension) + ".txt"
	filterPath := "data/" + es.Id + "/assets/" + filterFile
	file, err = os.Open(filterPath)
	return
}

func ReadParamHeaders(es ExperimentSettings, asset string) (headers []string, err error) {
	file, err := getParamFile(es, asset)
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

func ReadParamValues(es ExperimentSettings, asset string) (fs filters, err error) {
	file, err := getParamFile(es, asset)
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
