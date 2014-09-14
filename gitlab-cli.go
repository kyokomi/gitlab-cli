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

	return &gitLab, nil
}

// issue create task.
func doCreateIssue(c *cli.Context) {
	gitLab, err := newGitLabCli(c.GlobalBool("skip-cert-check"))
	if err != nil {
		log.Fatal("error create gitlab ")
	}

	user, err := gitLab.CurrentUser()
	if err != nil {
		log.Fatal("error get current user gitlab ")
	}

	values := url.Values{
		"title":       {c.String("t")},
		"description": {c.String("d")},
		"assignee_id": {strconv.Itoa(user.ID)},
		//		"milestone_id": {"1"},
		"labels": {c.String("l")},
	}
	res, err := gitLab.CreateIssue(gitLab.currentProjectID, values)
	if err != nil {
		log.Fatal("project issue create error ", err)
	}
	fmt.Println("done. ", string(res))
}

// issue edit task.
func doEditIssue(c *cli.Context) {
	gitLab, err := newGitLabCli(c.GlobalBool("skip-cert-check"))
	if err != nil {
		log.Fatal("error create gitlab ")
	}

	values := url.Values{}
	if t := c.String("t"); t != "" {
		values.Set("title", t)
	}
	if d := c.String("d"); d != "" {
		values.Set("description", d)
	}
	if l := c.String("l"); l != "" {
		values.Set("labels", l)
	}

	issueID := c.Int("issue-id")
	res, err := gitLab.EditIssue(gitLab.currentProjectID, issueID, values)
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

	gitLab.PrintIssue(projectID, c.String("state"))
}

func doShowIssue(c *cli.Context) {

	issueID := c.Int("issue-id")

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

	if err := gitLab.PrintIssueDetail(projectID, issueID); err != nil {
		log.Fatal("not gitlab projectID ", err)
	}
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
			Usage:     "Creates a new project issue.",
			Flags: []cli.Flag{
				cli.StringFlag{"title, t",       "", "(required) - The title of an issue", ""},
				cli.StringFlag{"description, d", "", "(optional) - The description of an issue", ""},
				cli.StringFlag{"label, l",       "", "(optional) - Comma-separated label names for an issue", ""},
			},
			Action: doCreateIssue,
		},
		{
			Name:      "edit_issue",
			ShortName: "edit",
			Usage:     "Updates an existing project issue.",
			Flags: []cli.Flag{
				cli.IntFlag{"issue-id, id",       0, "(required) - The ID of a project issue", ""},
				cli.StringFlag{"title, t",       "", "(optional) - The title of an issue", ""},
				cli.StringFlag{"description, d", "", "(optional) - The description of an issue", ""},
				cli.StringFlag{"label, l",       "", "(optional) - Comma-separated label names for an issue", ""},
			},
			Action: doEditIssue,
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
			Usage:     "Get a list of project issues.",
			Action:    doListIssue,
			Flags: []cli.Flag{
				cli.StringFlag{"state, s", "", "(optional) - The state event of an issue ('close' to close issue and 'reopen' to reopen it)", ""},
			},
		},
		{
			Name:      "issue",
			ShortName: "",
			Usage:     "Gets a single project issue.",
			Action:    doShowIssue,
			Flags: []cli.Flag{
				cli.IntFlag{"issue-id, id", 0, "(required) - The ID of a project issue", ""},
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
