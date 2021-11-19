package helpers

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func EnsureFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0775)
	}
	return nil
}

func FindFilesUnder(path string, match string) (matches []string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
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

func RemoveLinesFromFile(path string, line1, line2 string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	var bs []byte
	buf := bytes.NewBuffer(bs)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != line1 && scanner.Text() != line2 {
			_, err = buf.Write(scanner.Bytes())
			if err != nil {
				return
			}
			_, err = buf.WriteString("\n")
			if err != nil {
				return
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}

	err = os.WriteFile(path, buf.Bytes(), 0644)
	return
}
