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
		cli.BoolFlag{"gitlab.skip-cert-check", "If set to true, gitlab client will skip certificate checking for https, possibly exposing your system to MITM attack."},
	}

	app.Commands = []cli.Command{
		{
			Name:      "create_issue",
			ShortName: "i",
			Usage:     "project create issue",
			Flags: []cli.Flag{
//				cli.IntFlag{"project-id", 1, "projectId."},
			},
			Action: func(_ *cli.Context) {

				// TODO: local test server
				res, err := http.PostForm("http://172.17.8.101:10080//api/v3/projects/1/issues?private_token=14K1ZR6QaH1yznNFWRtw", url.Values {
					"id":           {"1"},
					"title":        {"hogehoge"},
					"description":  {"fuga"},
//					"assignee_id":  {"1"},
//					"milestone_id": {"1"},
//					"labels":       {"tag"},
				})

				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(res)
			},
		},
	}

	app.Run(os.Args)
}
