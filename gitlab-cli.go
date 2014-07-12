package main

import (
	"github.com/codegangsta/cli"
	"os"
	"net/url"
	"log"
	"fmt"
	"flag"
	"github.com/kyokomi/go-gitlab-client/gogitlab"
	"io/ioutil"
	"strings"
)

func main() {

	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "gitlab-cli"
	app.Usage = "todo:"

	app.Flags = []cli.Flag {
		cli.BoolFlag{"gitlab.skip-cert-check",
			"If set to true, gitlab client will skip certificate checking for https, possibly exposing your system to MITM attack."},
	}

	flag.Parse()

	app.Commands = []cli.Command{
		{
			Name:      "create_issue",
			ShortName: "i",
			Usage:     "project create issue",
			Flags: []cli.Flag{
				cli.StringFlag{"t", "", "issue title."},
				cli.StringFlag{"d", "", "issue description."},
				cli.StringFlag{"a", "ZDesNxuMt5jeCjJ9KSpH", "access token."},
				cli.StringFlag{"u", "http://172.17.8.101:10080/", "url."},
			},
			Action: func(c *cli.Context) {

				domainUrl := c.String("u")
				apiUrl := "api/v3/"
				accessToken := c.String("a")
				gitlab := gogitlab.NewGitlab(domainUrl, apiUrl, accessToken)

				projectName := GetCurrentDirProjectName()
				projectId := GetProjectId(gitlab, projectName)

				PostIssue(gitlab, projectId, url.Values {
//					"id":           {"1"},
					"title":        {c.String("t")},
					"description":  {c.String("d")},
//					"assignee_id":  {"1"},
//					"milestone_id": {"1"},
//					"labels":       {"tag"},
				})
			},
		},
		{
			Name:      "check-project",
			ShortName: "c",
			Usage:     "check project name",
			Action: func(_ *cli.Context) {
				projectName := GetCurrentDirProjectName()
				fmt.Println("projectName = ", projectName)
			},
		},
	}

	app.Run(os.Args)
}

func GetCurrentDirProjectName() string {
	fileData, err := ioutil.ReadFile("./.git/config")
	if err != nil {
		log.Fatal(err)
	}

	var projectName string
	fileText := string(fileData)
	for _, line := range strings.Split(fileText, "\n") {
		if !strings.Contains(line, "url") {
			continue
		}

		// TODO: gitlab check
		if !strings.Contains(line, "gitlab.com:") {
			log.Fatal("It does not support the repository.", line)
			break
		}

		// replace projectName
		idx := strings.LastIndex(line, "/")
		projectName = strings.TrimLeft(line, line[0:idx])
		projectName = strings.Replace(projectName, "/", "", 1)
		projectName = strings.Replace(projectName, ".git", "", 1)
		break
	}
	return projectName
}

func GetProjectId(gitlab *gogitlab.Gitlab, projectName string) int {
	projects, err := gitlab.Projects()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, project := range projects {
		if project.Name == projectName {
			return project.Id
		}
	}
	return 0
}

func PostIssue(gitlab *gogitlab.Gitlab, projectId int, data url.Values) {
	issue := fmt.Sprintf("projects/%d/issues/", projectId)
	url := gitlab.ResourceUrl(issue, nil)

	res, err := gitlab.Client.PostForm(url, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
