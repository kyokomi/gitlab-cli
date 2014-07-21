package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

// ${HOME}/.gitlab-cli/config.json
type GitlabAccessConfig struct {
	Host    string `json:"host"`
	ApiPath string `json:"api_path"`
	Token   string `json:"token"`
}

type ReadConfigError struct {
	Err error
}

func (e *ReadConfigError) Error() string { return e.Err.Error() }

// アクセストークンを保存してるローカルファイルを読み込んで返却する.
func ReadFileGitlabAccessTokenJson() (config GitlabAccessConfig, err error) {
	filePath, err := CreateConfigFilePath()
	if err != nil {
		return
	}
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	json.Unmarshal(file, &config)
	return
}

// ConfigFileのパスを作成して返却する.
func CreateConfigFilePath() (filePath string, err error) {
	usr, err := user.Current()
	if err == nil {
		filePath = usr.HomeDir + "/." + AppName + "/config.json"
	}
	return
}

// デフォルトのConfigFileを作成する.
func WriteFileDefaultConfig() string {
	filePath, err := CreateConfigFilePath()
	if err != nil {
		log.Fatalln(err)
	}

	config := GitlabAccessConfig{
		Host:    "https://gitlab.com/",
		ApiPath: "api/v3/",
		Token:   "aaaaaaaaaaaaaaaaaaaaaaa",
	}

	data, err := json.Marshal(config)
	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile(filePath, data, os.FileMode(0644))
	if err != nil {
		log.Fatalln(err)
	}

	return filePath
}
