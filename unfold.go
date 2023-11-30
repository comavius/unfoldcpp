package unfoldcpp

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Unfold(path string) (string, error) {
	// trace dependencies recursively
	current_dependencies := make([]string, 0)
	current_dependencies = append(current_dependencies, path)
	current_code := "#define UNFOLDED\n"
	err := unfold_recursively(&current_dependencies, &current_code)
	if err != nil {
		return "", err
	}
	return current_code, nil
}

func unfold_recursively(current_dependencies *[]string, current_code *string) error {
	// open file
	path := (*current_dependencies)[len(*current_dependencies)-1]
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// setup regex
	r, err := regexp.Compile(`^( |\t)*#include[\s]+\"(.*)\"[\s]*$*`)
	if err != nil {
		return err
	}
	// read lines to check for #include
	// if #include found, add to dependencies and call unfold_recursively
	// if #include not found, add to code
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if r.MatchString(scanner.Text()) {
			line := scanner.Text()
			start_index := strings.Index(line, "\"")
			end_index, err := findNextQuote(line, start_index+1)
			if err != nil {
				return err
			}
			relative_path := line[start_index+1 : end_index]
			absolute_path, err := filepath.Abs(filepath.Dir(path) + "/" + relative_path)
			if err != nil {
				return err
			}
			// check if already in dependencies
			already_in_dependencies := false
			for _, dependency := range *current_dependencies {
				if dependency == absolute_path {
					already_in_dependencies = true
					break
				}
			}
			if !already_in_dependencies {
				*current_dependencies = append(*current_dependencies, absolute_path)
				err = unfold_recursively(current_dependencies, current_code)
				if err != nil {
					return err
				}
			} else {
				// call cyclic dependency error
				return errors.New("Cyclic dependency detected")
			}
		} else {
			*current_code += scanner.Text() + "\n"
		}
	}
	// delete here from dependencies
	*current_dependencies = (*current_dependencies)[:len(*current_dependencies)-1]
	return nil
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
