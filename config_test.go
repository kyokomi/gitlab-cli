package main

import (
	"os"
	"testing"
)

func TestCreateConfigFilePath(t *testing.T) {
	filePath := createConfigFilePath("_test")
	if filePath != "_test/.gitlab-cli/config.json" {
		t.Error("filePath error")
	}
}

func TestReadFileGitlabAccessTokenJson(t *testing.T) {

	currentDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}

	err = os.RemoveAll(currentDir + "/_test")
	if err != nil {
		t.Error(err)
	}
	err = os.Mkdir(currentDir+"/_test", os.FileMode(0755))
	if err != nil {
		t.Error(err)
	}

	configDirPath := createConfigDirPath(currentDir + "/_test")

	// no file test
	_, err = ReadFileGitlabAccessTokenJson(configDirPath + configFileName)
	if err == nil {
		t.Error(err)
	}

	// default write
	var wConfig GitlabAccessConfig
	wConfig, err = WriteFileDefaultConfig(configDirPath + configFileName)
	if err != nil {
		t.Error(err)
	}
	if wConfig.Token != defaultConfig.Token {
		t.Error("config token missmatch")
	}

	// read
	var rConfig GitlabAccessConfig
	rConfig, err = ReadFileGitlabAccessTokenJson(configDirPath + configFileName)
	if err != nil {
		t.Error(err)
	}
	if rConfig.Token != wConfig.Token {
		t.Error(err)
	}
}
