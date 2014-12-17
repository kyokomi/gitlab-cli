// Package appConfig ~/.{appName}/configの書き込み,作成,読み込みを行うパッケージ
package appConfig

import (
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

type AppConfig struct {
	ConfigFileName string
	ConfigDirPath  string
	AppName        string
}

// NewAppConfig create AppConfig.
func NewDefaultAppConfig(appName string) *AppConfig {
	return NewAppConfig(appName, "config")
}

func NewAppConfig(appName, configFileName string) *AppConfig {
	dirPath, err := createAppConfigDirPath(appName)
	if err != nil {
		dirPath = "./"
	}
	return &AppConfig{
		ConfigFileName: configFileName,
		ConfigDirPath:  dirPath,
		AppName:        appName,
	}
}

// configファイルを作成する中身は空.
func (a AppConfig) WriteAppConfig(data []byte) error {
	if err := createAppConfigDir(a.ConfigDirPath); err != nil {
		return err
	}
	return ioutil.WriteFile(a.AppConfigFilePath(), data, os.FileMode(0644))
}

func (a AppConfig) AppConfigFilePath() string {
	return strings.Join([]string{a.ConfigDirPath, a.ConfigFileName}, "/")
}

// configファイルを読み込む[]byte.
func (a AppConfig) ReadAppConfig() ([]byte, error) {
	return ioutil.ReadFile(a.AppConfigFilePath())
}

func (a AppConfig) RemoveAppConfig() error {
	return os.RemoveAll(a.ConfigDirPath)
}

// ~/.{appName}ディレクトリを作成
// すでに存在する場合スルー
func createAppConfigDir(dirPath string) error {
	// check
	if _, err := ioutil.ReadDir(dirPath); err == nil {
		return nil
	}

	// create dir
	return os.Mkdir(dirPath, os.FileMode(0755))
}

func createAppConfigDirPath(appName string) (string, error) {
	// home
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	dirName := "." + appName
	dirPath := strings.Join([]string{usr.HomeDir, dirName}, "/")
	return dirPath, nil
}
