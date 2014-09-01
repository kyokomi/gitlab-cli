package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/kyokomi/appConfig"
)

const configFileName = "config.json"

// ${HOME}/.gitlab-cli/config.json
type GitlabAccessConfig struct {
	Host    string `json:"host"`
	ApiPath string `json:"api_path"`
	Token   string `json:"token"`
}

var defaultConfig = GitlabAccessConfig{
	Host:    "https://gitlab.com/",
	ApiPath: "api/v3/",
	Token:   "aaaaaaaaaaaaaaaaaaaaaaa",
}

// ConfigFileのパスを作成して返却する.
func CreateConfigFilePath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	filePath := createConfigFilePath(usr.HomeDir)

	return filePath, nil
}

func CreateConfigDirPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	dirPath := createConfigDirPath(usr.HomeDir)

	return dirPath, nil
}

func createConfigFilePath(baseDir string) string {
	return createConfigDirPath(baseDir) + configFileName
}

func createConfigDirPath(baseDir string) string {
	return baseDir + "/" + "." + AppName + "/"
}

// アクセストークンを保存してるローカルファイルを読み込んで返却する.
func ReadGitlabAccessTokenJson() (GitlabAccessConfig, error) {
	filePath, err := CreateConfigFilePath()
	if err != nil {
		return GitlabAccessConfig{}, err
	}

	return ReadFileGitlabAccessTokenJson(filePath)
}

func ReadFileGitlabAccessTokenJson(filePath string) (GitlabAccessConfig, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return GitlabAccessConfig{}, err
	}

	var config GitlabAccessConfig
	json.Unmarshal(file, &config)

	return config, nil
}

func WriteDefaultConfig() (GitlabAccessConfig, error) {
	dirPath, err := CreateConfigDirPath()
	if err != nil {
		return GitlabAccessConfig{}, err
	}

	return WriteFileDefaultConfig(dirPath + configFileName)
}

// デフォルトのConfigFileを作成する.
func WriteFileDefaultConfig(filePath string) (GitlabAccessConfig, error) {
	if config, err := ReadFileGitlabAccessTokenJson(filePath); err == nil {
		return config, nil
	}
	idx := strings.LastIndex(filePath, "/")
	if err := os.Mkdir(filePath[:idx], os.FileMode(0755)); err != nil {
		return GitlabAccessConfig{}, err
	}

	data, err := json.Marshal(defaultConfig)
	if err != nil {
		return GitlabAccessConfig{}, err
	}

	if err := ioutil.WriteFile(filePath, data, os.FileMode(0644)); err != nil {
		return GitlabAccessConfig{}, err
	}

	return defaultConfig, nil
}

//////////////////////

func WriteAppConfig(appName string, config *GitlabAccessConfig) error {
	a := appConfig.NewAppConfig(appName, "config.json")
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return a.WriteAppConfig(data)
}
