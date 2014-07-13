package login

import (
	"os/user"
	"log"
	"os"
	"io/ioutil"
	"encoding/json"
)

type GitlabAccessConfig struct {
	Host     string `json:"host"`
	ApiPath  string `json:"api_path"`
	Token    string `json:"token"`
}

// アクセストークンを保存してるローカルファイルを読み込んで返却
func ReadGitlabAccessTokenJson(appName string) GitlabAccessConfig {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
		os.Exit(1)
	}

	filePath := usr.HomeDir + "/." + appName + "/config.json"
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Config file error: %v\n", err)
		os.Exit(1)
	}

	var config GitlabAccessConfig
	json.Unmarshal(file, &config)

	return config
}
