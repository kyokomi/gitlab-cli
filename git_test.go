package main

import "testing"

func TestGetCurrentDirProjectName(t *testing.T) {
	name := GetCurrentDirProjectName()
	if name != "gitlab-cli" {
		t.Error("current projectName error.")
	}
}
