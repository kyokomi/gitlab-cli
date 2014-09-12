package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var testDir string

func init() {
	currentDir, _ := os.Getwd()
	testDir = strings.Join([]string{currentDir, "_test"}, "/")

	os.MkdirAll(testDir, 0755)

	fmt.Println("init complete")
}

func TestReadFileGitlabAccessTokenJson(t *testing.T) {

	// test用にConfigディレクトリを変更
	ac := NewGitlabCliAppConfig("test")
	ac.ConfigDirPath = testDir

	// テスト前に削除
	if err := ac.RemoveAppConfig(); err != nil {
		t.Error(err)
	}

	// デフォルト書き込み
	if err := ac.WriteDefaultGitlabAccessConfig(); err != nil {
		t.Error(err)
	}

	// 読み込み
	c, err := ac.ReadGitlabAccessTokenJson()
	if err != nil {
		t.Error(err)
	} else {
		if c.Token != defaultConfig.Token {
			t.Error("config token missmatch")
		}
	}

	// テスト後に削除
	if err := ac.RemoveAppConfig(); err != nil {
		t.Error(err)
	}
}
