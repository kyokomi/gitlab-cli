package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"unicode/utf8"

	"github.com/codegangsta/cli"
	"github.com/kyokomi/go-gitlab-client/gogitlab"
	"github.com/mitchellh/colorstring"
)

const (
	ProjectIssueUrl = "projects/%d/issues/"
)

var gitlabAppConfig *GitlabCliAppConfig

// Gitlabクライアントを作成する.
func CreateGitlab() (*gogitlab.Gitlab, error) {
	config, err := gitlabAppConfig.ReadGitlabAccessTokenJson()
	if err != nil {
		return nil, err
	}
	flag.Parse()
	return gogitlab.NewGitlab(config.Host, config.ApiPath, config.Token), nil
}

// 対象Projectのissueを作成する.
func PostIssue(gitlab *gogitlab.Gitlab, projectId int, data url.Values) error {
	issue := fmt.Sprintf(ProjectIssueUrl, projectId)
	url := gitlab.ResourceUrl(issue, nil)

	res, err := gitlab.Client.PostForm(url, data)
	if err != nil {
		return err
	}
	fmt.Println(res)

	return nil
}

func ShowIssue(gitlab *gogitlab.Gitlab, projectId int) {
	c := make(chan []*gogitlab.Issue)
	go func(s chan<- []*gogitlab.Issue) {
		page := 1
		for {
			issues, err := gitlab.ProjectIssues(projectId, page)
			if err != nil || len(issues) == 0 {
				break
			}
			page++

			s <- issues
		}
		close(s)
	}(c)

	for {
		issues, ok := <-c
		if !ok {
			break
		}

		for _, issue := range issues {
			titleCount := 90 + ((utf8.RuneCountInString(issue.Title) - len(issue.Title)) / 2)
			nameCount := 16 + ((utf8.RuneCountInString(issue.Assignee.Name) - len(issue.Assignee.Name)) / 2)
			t := fmt.Sprintf("[blue]#%%-4d %%-7s [white]%%-%ds [green]%%-%ds [white]%%-33s / %%-33s", titleCount, nameCount)
			fmt.Println(colorstring.Color(fmt.Sprintf(t,
				issue.LocalId,
				issue.State,
				issue.Title,
				issue.Assignee.Name,
				issue.CreatedAt,
				issue.UpdatedAt)))
		}
	}
}

// issue create task.
func doCreateIssue(c *cli.Context) {

	gitlab, err := CreateGitlab()
	if err != nil {
		log.Fatal("error create gitlab ")
	}

	projectName, err := GetCurrentDirProjectName()
	if err != nil {
		log.Fatal("not gitlab projectName ", err)
	}

	projectId, err := GetProjectId(gitlab, projectName)
	if err != nil {
		log.Fatal("not gitlab projectId ", err)
	}

	PostIssue(gitlab, projectId, url.Values{
		//		"id":           {"1"},
		"title":       {c.String("t")},
		"description": {c.String("d")},
		//		"assignee_id":  {"1"},
		//		"milestone_id": {"1"},
		"labels": {c.String("l")},
	})
}

// project check task.
func doCheckProject(_ *cli.Context) {
	projectName, err := GetCurrentDirProjectName()
	if err != nil {
		log.Fatal("not gitlab projectName ", err)
	}
	fmt.Println("projectName = ", projectName)
}

func doListIssue(_ *cli.Context) {
	gitlab, err := CreateGitlab()
	if err != nil {
		log.Fatal("error create gitlab ")
	}

	projectName, err := GetCurrentDirProjectName()
	if err != nil {
		log.Fatal("not gitlab projectName ", err)
	}

	projectId, err := GetProjectId(gitlab, projectName)
	if err != nil {
		log.Fatal("not gitlab projectId ", err)
	}

	ShowIssue(gitlab, projectId)
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
		cli.BoolFlag{"gitlab.skip-cert-check",
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
