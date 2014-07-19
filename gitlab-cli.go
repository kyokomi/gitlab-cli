package main

import (
	"github.com/codegangsta/cli"
	"os"
	"net/url"
	"log"
	"fmt"
	"github.com/kyokomi/go-gitlab-client/gogitlab"
	"flag"
)

const (
	ProjectIssueUrl = "projects/%d/issues/"
)

// Gitlabクライアントを作成する
func CreateGitlab() *gogitlab.Gitlab {
	config, err := ReadFileGitlabAccessTokenJson()
	if err != nil {
		filePath := WriteFileDefaultConfig()
		fmt.Println("write file default.", filePath)
		return nil
	}
	flag.Parse()
	return gogitlab.NewGitlab(config.Host, config.ApiPath, config.Token)
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
	issue := fmt.Sprintf(ProjectIssueUrl, projectId)
	url := gitlab.ResourceUrl(issue, nil)

	res, err := gitlab.Client.PostForm(url, data)
	if err != nil {
		fmt.Println(url)
		log.Fatal(err)
	}
	fmt.Println(res)
}

func doCreateIssue(c *cli.Context) {

	gitlab := CreateGitlab()
	if gitlab == nil {
		return
	}

	projectName := GetCurrentDirProjectName()
	projectId := GetProjectId(gitlab, projectName)

	PostIssue(gitlab, projectId, url.Values {
			//		"id":           {"1"},
		"title":        {c.String("t")},
		"description":  {c.String("d")},
			//		"assignee_id":  {"1"},
			//		"milestone_id": {"1"},
		"labels":       {c.String("l")},
	})
}

func doCheckProject(_ *cli.Context) {
	projectName := GetCurrentDirProjectName()
	fmt.Println("projectName = ", projectName)
}

func main() {

	app := cli.NewApp()
	app.Version = Version
	app.Name = AppName
	app.Usage = "todo:"

	app.Flags = []cli.Flag {
		cli.BoolFlag{"gitlab.skip-cert-check",
			"If set to true, gitlab client will skip certificate checking for https, possibly exposing your system to MITM attack."},
	}

	app.Commands = []cli.Command{
		{
			Name:      "create_issue",
			ShortName: "i",
			Usage:     "project create issue",
			Flags: []cli.Flag{
				cli.StringFlag{"title, t", "", "issue title."},
				cli.StringFlag{"description, d", "", "issue description."},
				cli.StringFlag{"label, l", "", "label example hoge,fuga,piyo."},
			},
			Action: doCreateIssue,
		},
		{
			Name:      "check-project",
			ShortName: "c",
			Usage:     "check project name",
			Action: doCheckProject,
		},
	}
	app.Run(os.Args)
}
