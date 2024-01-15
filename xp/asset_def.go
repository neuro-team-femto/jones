package xp

import (
	"bufio"
	"os"
	"strings"
)

func getAssetDefFile(es ExperimentSettings, asset string) (file *os.File, err error) {
	filterFile := strings.TrimSuffix(asset, "."+es.FileExtension) + ".txt"
	filterPath := "data/" + es.Id + "/assets/" + filterFile
	file, err = os.Open(filterPath)
	return
}

func getAssetDefHeaders(es ExperimentSettings, asset string) (headers []string, err error) {
	file, err := getAssetDefFile(es, asset)
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

func getAssetDefAllValues(es ExperimentSettings, asset string) (allValues [][]string, err error) {
	file, err := getAssetDefFile(es, asset)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// skip first line
	scanner.Scan()

	for scanner.Scan() {
		values := strings.Split(scanner.Text(), ",")
		allValues = append(allValues, values)
	}
	return
}
