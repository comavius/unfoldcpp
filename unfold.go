package unfoldcpp

import (
	"os"
)

func Unfold(path string) (string, error) {
	// trace dependencies recursively
	founded_queue := make([]string, 0)
	founded_list := make([]string, 0)
	resolved_queue := make([]string, 0)
	founded_queue = append(founded_queue, path)
	founded_list = append(founded_list, path)
	err := unfold_recursively(&founded_queue, &resolved_queue, &founded_list)
	if err != nil {
		return "", err
	}
	// build single-file
	single_file := "#define UNFOLDED\n"
	for _, path := range resolved_queue {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		single_file += string(data) + "\n"
	}
	return single_file, nil
}

func unfold_recursively(founded_queue *[]string, resolved_queue *[]string, founded_list *[]string) error {
	// pop
	path := (*founded_queue)[0]
	*founded_queue = (*founded_queue)[1:]
	// construct hppFile
	hpp := newHppFile(path)
	// trace dependencies
	dependencies, err := hpp.TraceDependencies()
	if err != nil {
		return err
	}
	// add dependencies to founded_queue if not already in founded_list
	for _, dependency := range dependencies {
		founded_formerly := false
		for _, founded_path := range *founded_list {
			if founded_path == dependency {
				founded_formerly = true
				break
			}
		}
		if !founded_formerly {
			*founded_queue = append(*founded_queue, dependency)
			*founded_list = append(*founded_list, dependency)
			unfold_recursively(founded_queue, resolved_queue, founded_list)
		}
	}
	// add path to resolved_queue
	*resolved_queue = append(*resolved_queue, path)
	return nil
}
