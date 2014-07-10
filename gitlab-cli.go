package main

import (
	"github.com/codegangsta/cli"
	"os"
	"net/http"
	"net/url"
	"log"
	"fmt"
)

func main() {

	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "gitlab-cli"
	app.Usage = "todo:"

	app.Flags = []cli.Flag {
//		cli.BoolFlag{"gitlab.skip-cert-check", "If set to true, gitlab client will skip certificate checking for https, possibly exposing your system to MITM attack."},
	}

	app.Commands = []cli.Command{
		{
			Name:      "create_issue",
			ShortName: "i",
			Usage:     "project create issue",
			Flags: []cli.Flag{
				cli.StringFlag{"t", "", "issue title."},
				cli.StringFlag{"d", "", "issue description."},
			},
			Action: func(c *cli.Context) {

				// TODO: 今いるディレクトリの.gitから検索したい
				projectId := 1
				// TODO: いい感じに保持したい
				accessToken := "14K1ZR6QaH1yznNFWRtw"

				PostIssue(projectId, accessToken, url.Values {
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

func PostIssue(projectId int, accessToken string, data url.Values) {
	// TODO: local test server url
	url := createUrl(projectId, accessToken)
	res, err := http.PostForm(url, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
func createUrl(projectId int, accessToken string) string {
	domain := "http://172.17.8.101:10080/" + "api/v3/"
	issue := fmt.Sprintf("projects/%d/issues/", projectId)
	token := "?private_token=" + accessToken
	return domain + issue + token
}
