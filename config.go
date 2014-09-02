package main

import (
	"encoding/json"

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

type GitlabCliAppConfig struct {
	appConfig.AppConfig
}

func NewGitlabCliAppConfig(appName string) *GitlabCliAppConfig {
	return &GitlabCliAppConfig{
		AppConfig: *appConfig.NewAppConfig(appName, configFileName),
	}
}

// アクセストークンを保存してるローカルファイルを読み込んで返却する.
func (a GitlabCliAppConfig) ReadGitlabAccessTokenJson() (*GitlabAccessConfig, error) {
	data ,err := a.ReadAppConfig()
	if err != nil {
		return nil, err
	}
	var c GitlabAccessConfig
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

//////////////////////

func (a GitlabCliAppConfig) WriteDefaultGitlabAccessConfig() error {
	return a.WriteGitlabAccessConfig(&defaultConfig)
}

func (a GitlabCliAppConfig) WriteGitlabAccessConfig(config *GitlabAccessConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return a.WriteAppConfig(data)
}
