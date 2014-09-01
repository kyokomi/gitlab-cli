package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/codegangsta/cli"
	"github.com/kyokomi/go-gitlab-client/gogitlab"
)

const (
	ProjectIssueUrl = "projects/%d/issues/"
)

// Gitlabクライアントを作成する.
func CreateGitlab() (*gogitlab.Gitlab, error) {
	config, err := ReadGitlabAccessTokenJson()
	if err != nil {
		_, err := WriteDefaultConfig()
		fmt.Println("write file default. ", err)
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
		fmt.Println(url)
		return err
	}
	fmt.Println(res)

	return nil
}

func ShowIssue(gitlab *gogitlab.Gitlab, projectId int, showDetail bool) {
	page := 1
	for {
		issues, err := gitlab.ProjectIssues(projectId, page)
		if err != nil {
			fmt.Println(err)
			break
		}
		if len(issues) == 0 {
			break
		}

		for _, issue := range issues {

			if issue.State != "closed" {
				if showDetail {
					fmt.Printf("[%4d(%d)] %s : [%s] (%s)\n%s\n", issue.Id, issue.LocalId, issue.State, issue.Title, issue.Assignee.Name, issue.Description)
				} else {
					fmt.Printf("[%4d(%d)] %s : [%s] (%s)\n", issue.Id, issue.LocalId, issue.State, issue.Title, issue.Assignee.Name)
				}
			}
		}
		page++
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

func doShowIssue(c *cli.Context) {
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

	ShowIssue(gitlab, projectId, c.Bool("detail"))
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
	if err := WriteAppConfig(c.App.Name, &config); err != nil {
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

	app.Commands = []cli.Command{
		{
			Name:      "create_issue",
			ShortName: "i",
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
			ShortName: "c",
			Usage:     "check project name",
			Action:    doCheckProject,
		},
		{
			Name:      "list-issue",
			ShortName: "l",
			Usage:     "list project issue",
			Action:    doShowIssue,
			Flags: []cli.Flag{
				cli.BoolFlag{"detail, d", "show/hide issue detail.", ""},
			},
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
