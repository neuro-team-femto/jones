package xp

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func assetDefFilePath(es ExperimentSettings, asset string) string {
	fileName := strings.TrimSuffix(asset, "."+es.FileExtension) + ".txt"
	return "data/" + es.Id + "/assets/" + fileName
}

func getAssetDefHeaders(es ExperimentSettings, asset string) (headers []string, err error) {
	path := assetDefFilePath(es, asset)
	file, err := os.Open(path)
	if err != nil {
		log.Printf("[error] unable to read path '%v' error: %+v\n", path, err)
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
	path := assetDefFilePath(es, asset)
	file, err := os.Open(path)
	if err != nil {
		log.Printf("[error] unable to read path '%v' error: %+v\n", path, err)
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
