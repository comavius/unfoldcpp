package unfoldcpp

import (
	"bufio"
	"regexp"
	"strings"
)

func deleteIncludes(single_file string) (string, error) {
	guarded_single_file := ""
	r, err := regexp.Compile(`#include[\s]+\"(.*)\"\s*`)
	if err != nil {
		return "", err
	}
	// read string line by line
	scanner := bufio.NewScanner(strings.NewReader(single_file))
	for scanner.Scan() {
		line := scanner.Text()
		if !r.MatchString(line) {
			guarded_single_file += line + "\n"
		}
	}
	return guarded_single_file, nil
}
