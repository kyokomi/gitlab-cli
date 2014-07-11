package main

import (
	"github.com/codegangsta/cli"
	"os"
	"net/url"
	"log"
	"fmt"
	"flag"
	"github.com/kyokomi/go-gitlab-client/gogitlab"
)

func main() {

	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "gitlab-cli"
	app.Usage = "todo:"

	app.Flags = []cli.Flag {
		cli.BoolFlag{"gitlab.skip-cert-check", "If set to true, gitlab client will skip certificate checking for https, possibly exposing your system to MITM attack."},
	}

	flag.Parse()

	app.Commands = []cli.Command{
		{
			Name:      "create_issue",
			ShortName: "i",
			Usage:     "project create issue",
			Flags: []cli.Flag{
				cli.IntFlag{"i", 1, "projectId."},
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

				// TODO: projectIdは今いるディレクトリの.gitから検索したい
				projectId := c.Int("i")

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
	}

	app.Run(os.Args)
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
