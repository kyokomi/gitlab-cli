package main

import (
	"io/ioutil"
	"strings"
)

// read Current dir gir project name.
func GetCurrentDirProjectName() (string, error) {
	data, err := ioutil.ReadFile("./.git/config")
	if err != nil {
		return "", err
	}

	var projectName string
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.Contains(line, "url") {
			continue
		}

		// replace projectName
		idx := strings.LastIndex(line, "/")
		projectName = strings.Replace(line[idx+1:], ".git", "", 1)
		break
	}

	return projectName, nil
}
