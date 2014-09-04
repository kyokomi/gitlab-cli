package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"strconv"

	"github.com/codegangsta/cli"
	"github.com/kyokomi/go-gitlab-client/gogitlab"
)

var gitlabAppConfig *GitlabCliAppConfig

type gitLabCli struct {
	*gogitlab.Gitlab
	currentUser        gogitlab.User
	currentProjectID   int
	currentProjectName string
}

func newGitLabCli(skipCert bool) (*gitLabCli, error) {
	config, err := gitlabAppConfig.ReadGitlabAccessTokenJson()
	if err != nil {
		return nil, err
	}
	gitLab := gitLabCli{Gitlab: gogitlab.NewGitlabCert(config.Host, config.ApiPath, config.Token, skipCert)}

	projectName, err := GetCurrentDirProjectName()
	if err != nil {
		return nil, err
	}
	gitLab.currentProjectName = projectName

	projectID, err := gitLab.GetProjectID(projectName)
	if err != nil {
		return nil, err
	}
	gitLab.currentProjectID = projectID

	user, err := gitLab.CurrentUser()
	if err != nil {
		return nil, err
	}
	gitLab.currentUser = user

	return &gitLab, nil
}

// issue create task.
func doCreateIssue(c *cli.Context) {
	gitLab, err := newGitLabCli(c.GlobalBool("skip-cert-check"))
	if err != nil {
		log.Fatal("error create gitlab ")
	}

	values := url.Values{
		//		"id":           {"1"},
		"title":       {c.String("t")},
		"description": {c.String("d")},
		"assignee_id": {strconv.Itoa(gitLab.currentUser.ID)},
		//		"milestone_id": {"1"},
		"labels": {c.String("l")},
	}
	res, err := gitLab.PostIssue(gitLab.currentProjectID, values)
	if err != nil {
		log.Fatal("project issue create error ", err)
	}
	fmt.Println("done. ", string(res))
}

// project check task.
func doCheckProject(_ *cli.Context) {
	projectName, err := GetCurrentDirProjectName()
	if err != nil {
		log.Fatal("not gitlab projectName ", err)
	}
	fmt.Println("projectName = ", projectName)
}

func doListIssue(c *cli.Context) {
	gitLab, err := newGitLabCli(c.GlobalBool("skip-cert-check"))
	if err != nil {
		log.Fatal("error create gitlab ")
	}

	projectName, err := GetCurrentDirProjectName()
	if err != nil {
		log.Fatal("not gitlab projectName ", err)
	}

	projectID, err := gitLab.GetProjectID(projectName)
	if err != nil {
		log.Fatal("not gitlab projectID ", err)
	}

	gitLab.PrintIssue(projectID)
}

func doShowIssue(_ *cli.Context) {
	// TODO:
}

func doInitConfig(c *cli.Context) {

	hostName := c.String("host")
	apiPath := c.String("api-path")
	token := c.String("token")

	config := GitlabAccessConfig{
		Host:    hostName,
		ApiPath: apiPath,
		Token:   token,
	}
	if err := gitlabAppConfig.WriteGitlabAccessConfig(&config); err != nil {
		log.Fatal("appConfig write error ", err)
	}
}

// main.
func main() {

	app := cli.NewApp()
	app.Version = Version
	app.Name = AppName
	app.Usage = "todo:"

	app.Flags = []cli.Flag{
		cli.BoolFlag{"skip-cert-check,s",
			"If set to true, gitlab client will skip certificate checking for https, possibly exposing your system to MITM attack.",
			"GITLAB.SKIP_CERT_CHECK"},
	}

	gitlabAppConfig = NewGitlabCliAppConfig(AppName)

	app.Commands = []cli.Command{
		{
			Name:      "add_issue",
			ShortName: "add",
			Usage:     "project create issue",
			Flags: []cli.Flag{
				cli.StringFlag{"title, t", "", "issue title.", ""},
				cli.StringFlag{"description, d", "", "issue description.", ""},
				cli.StringFlag{"label, l", "", "label example hoge,fuga,piyo.", ""},
			},
			Action: doCreateIssue,
		},
		{
			Name:      "check-project",
			ShortName: "check",
			Usage:     "check project name",
			Action:    doCheckProject,
		},
		{
			Name:      "list-issue",
			ShortName: "list",
			Usage:     "list project issue",
			Action:    doListIssue,
		},
		{
			Name:      "issue",
			ShortName: "",
			Usage:     "show project issue",
			Action:    doShowIssue,
		},
		{
			Name:      "init-config",
			ShortName: "init",
			Usage:     "initialize to config",
			Action:    doInitConfig,
			Flags: []cli.Flag{
				cli.StringFlag{"host", "https://gitlab.com/", "host name example [https://gitlab.com/]", ""},
				cli.StringFlag{"api-path", "api/v3/", "api path example [api/v3/]", ""},
				cli.StringFlag{"token", "", "your access token", ""},
			},
		},
	}
	app.Run(os.Args)
}
