package main

import (
	"io/ioutil"
	"log"
	"strings"
)

func GetCurrentDirProjectName() string {
	data, err := ioutil.ReadFile("./.git/config")
	if err != nil {
		log.Fatal(err)
	}

	var projectName string
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.Contains(line, "url") {
			continue
		}

		// replace projectName
		idx := strings.LastIndex(line, "/")
		projectName = strings.TrimLeft(line, line[0:idx])
		projectName = strings.Replace(projectName, "/", "", 1)
		projectName = strings.Replace(projectName, ".git", "", 1)
		break
	}
	return projectName
}

