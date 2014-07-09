package main

import (
	"github.com/kyokomi/go-gitlab-client/gogitlab"
	"flag"
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
	IconPath string `json:"icon_path"`
}

// Gitlabクライアントを作成する
func CreateGitlab() *gogitlab.Gitlab {
	config := readGitlabAccessTokenJson()

	// TODO: あとで別だし
	// --gitlab.skip-cert-checkを読み込む
	flag.Parse()

	return gogitlab.NewGitlab(config.Host, config.ApiPath, config.Token)
}

// アクセストークンを保存してるローカルファイルを読み込んで返却
func readGitlabAccessTokenJson() GitlabAccessConfig {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
		os.Exit(1)
	}

	// TODO: あとで別だし
	file, err := ioutil.ReadFile(usr.HomeDir + "/.ggn/config.json")
	if err != nil {
		log.Fatalf("Config file error: %v\n", err)
		os.Exit(1)
	}

	var config GitlabAccessConfig
	json.Unmarshal(file, &config)
	// TODO: あとで別だし
	config.IconPath = usr.HomeDir + "/.ggn/logo.png"

	return config
}
