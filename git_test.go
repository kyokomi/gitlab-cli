package main

import "testing"

func TestGetCurrentDirProjectName(t *testing.T) {

	projectName, err := GetCurrentDirProjectName()
	if err == nil {
		if projectName != "gitlab-cli" {
			t.Errorf("bad projectName %s", projectName)
		}
	}
}
