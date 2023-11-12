package unfoldcpp

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// filepath structure
type hppFile struct {
	path         string
	dependencies []string
}

// constructor
func newHppFile(path string) *hppFile {
	return &hppFile{path: path}
}

// get depe
func (hpp *hppFile) TraceDependencies() ([]string, error) {
	// open file
	file, err := os.Open(hpp.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// read file and extract #include "*"
	// setup regex
	r, err := regexp.Compile(`#include[\s]+\"(.*)\"\s*`)
	if err != nil {
		return nil, err
	}
	// read file line by line
	relative_paths := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if r.MatchString(scanner.Text()) {
			line := scanner.Text()
			start_index := strings.Index(line, "\"")
			end_index, err := findNextQuote(line, start_index+1)
			if err != nil {
				return nil, err
			}
			relative_paths = append(relative_paths, line[start_index+1:end_index])
		}
	}
	// get absolute paths
	absolute_paths := make([]string, 0)
	for _, relative_path := range relative_paths {
		absolute_path, err := filepath.Abs(filepath.Dir(hpp.path) + "/" + relative_path)
		if err != nil {
			return nil, err
		}
		absolute_paths = append(absolute_paths, absolute_path)
	}
	// return
	return absolute_paths, nil
}

func findNextQuote(s string, startIndex int) (int, error) {
	if startIndex < 0 || startIndex >= len(s) {
		return -1, errors.New("startIndex out of range")
	}

	for i := startIndex; i < len(s); i++ {
		if s[i] == '"' {
			if i == 0 || s[i-1] != '\\' {
				return i, nil
			}
		}
	}
	return -1, errors.New("no unescaped quote found")
}
