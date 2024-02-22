package helpers

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		log.Printf("[error][PathExists] path '%v' error: %+v\n", path, err)
	}
	return true
}

func EnsureFolder(path string) error {
	if !PathExists(path) {
		return os.MkdirAll(path, 0775)
	}
	return nil
}

func FindFilesUnder(path string, match string) (matches []string) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("[error][FindFilesUnder] path '%v' error: %+v\n", path, err)
		return
	}
	for _, f := range files {
		if !f.IsDir() && strings.Contains(f.Name(), match) {
			matches = append(matches, f.Name())
		}
	}
	return
}

func ReadFileLines(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("[error][ReadFileLines] path '%v' error: %+v\n", path, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return
}

func ReadTrimJSON(path string) (json string, err error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("[error][ReadTrimJSON] path '%v' error: %+v\n", path, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		json += strings.TrimSpace(scanner.Text())
	}
	return
}

func IsLineInFile(path string, line string) bool {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("[error][IsLineInFile] path '%v' error: %+v\n", path, err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		trimmed := strings.TrimSpace(scanner.Text())
		if trimmed == line {
			return true
		}
	}

	return false
}

// remove line1 and line2 from file once
// (if line1 appears several times, only the first occurrence is removed, same for line2)
// not currently used
func RemoveOnceFromFile(path string, line1, line2 string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("[error][RemoveOnceFromFile][read] path '%v' error: %+v\n", path, err)
		return
	}
	defer file.Close()

	var bs []byte
	buf := bytes.NewBuffer(bs)

	scanner := bufio.NewScanner(file)
	line1RemovedOnce := false
	line2RemovedOnce := false
	for scanner.Scan() {
		currentLine := scanner.Text()
		if currentLine == line1 && !line1RemovedOnce {
			line1RemovedOnce = true
		} else if currentLine == line2 && !line2RemovedOnce {
			line2RemovedOnce = true
		} else {
			_, err = buf.WriteString(scanner.Text() + "\n")
			if err != nil {
				log.Printf("[error][RemoveOnceFromFile][string] path '%v' error: %+v\n", path, err)
				return
			}
		}
	}
	if err = scanner.Err(); err != nil {
		log.Printf("[error][RemoveOnceFromFile][scanner] path '%v' error: %+v\n", path, err)
		return
	}

	if err = os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		log.Printf("[error][RemoveOnceFromFile][wrote] path '%v' error: %+v\n", path, err)
	}
	return
}
