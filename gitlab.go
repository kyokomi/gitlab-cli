package main

import (
	"strconv"

	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"unicode/utf8"

	"bytes"
	"text/template"

	"github.com/kyokomi/go-gitlab-client/gogitlab"
	color "github.com/mitchellh/colorstring"
)

const (
	titleCount          = 60
	nameCount           = 16
	labelCount           = 20
	outTitleReplaceText = " ..."
)

type templateExec struct {
	Issue *gogitlab.Issue
	Notes []*gogitlab.Note
}

const issueDetailTemplate = `
{{.Issue.Title}}

Author: [green]@{{.Issue.Author.Name}} [blue]{{.Issue.State}} [white]{{.Issue.CreatedAt}}
-------------------------------------------------------------------------------------------
{{.Issue.Description}}
-------------------------------------------------------------------------------------------
{{range $idx, $note := .Notes}}[green]@{{$note.Author.Name}}: [white]{{$note.Body}}
{{end}}`

// 対象ProjectのProjectNameを取得する.
func (gitLab *gitLabCli) GetProjectID(projectName string) (int, error) {
	projects, err := gitLab.Projects()
	if err != nil {
		return 0, err
	}

	for _, project := range projects {
		if project.Name == projectName {
			return project.ID, nil
		}
	}
	return 0, nil
}

func (gitLab *gitLabCli) GetProjectName(projectID int) (string, error) {
	project, err := gitLab.Project(strconv.Itoa(projectID))
	if err != nil {
		return "", err
	}
	return project.Name, nil
}

func (gitLab *gitLabCli) GetUserName(userID int) (string, error) {
	user, err := gitLab.User(strconv.Itoa(userID))
	if err != nil {
		return "", err
	}
	return user.Name, nil
}

//	/projects/:id/milestones/:milestone_id
func (gitLab *gitLabCli) GetMilestoneTitle(projectID, milestoneID int) (string, error) {
	milestone, err := gitLab.ProjectMilestone(strconv.Itoa(projectID), strconv.Itoa(milestoneID))
	if err != nil {
		return "", err
	}
	return milestone.Title, nil
}

func (gitLab *gitLabCli) CreateIssue(projectID int, values url.Values) ([]byte, error) {
	reader := strings.NewReader(values.Encode())
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	res, err := gitLab.ProjectCreateIssues(projectID, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (gitLab *gitLabCli) EditIssue(projectID, issueID int, values url.Values) ([]byte, error) {
	reader := strings.NewReader(values.Encode())
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// localIssueID => issueID
	issue := gitLab.findIssueByID(projectID, issueID)
	if issue == nil {
		return nil, fmt.Errorf("issue not found")
	}

	res, err := gitLab.ProjectEditIssues(projectID, issue.ID, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (gitLab *gitLabCli) PrintIssue(projectID int, state string) {
	c := make(chan []*gogitlab.Issue)
	go gitLab.findIssueState(c, projectID, func(issues []*gogitlab.Issue) bool {
		if state != "" && issues[len(issues)-1].State != state {
			return true
		}
		return false
	})

	for {
		issues, ok := <-c
		if !ok {
			break
		}

		for _, issue := range issues {
			if state != "" && issue.State != state {
				continue
			}

			title := issue.Title
			if checkTrim(title) {
				title = trimPrefixIndex(title)
			}

			name := issue.Assignee.Name
			if checkTrim(name) {
				name = trimPrefixIndex(name)
			}

			labels := fmt.Sprint(issue.Labels)
			if checkTrim(name) {
				labels = trimPrefixIndex(labels)
			}

			titleCount := titleCount + ((utf8.RuneCountInString(title) - len(title)) / 2)
			nameCount := nameCount + ((utf8.RuneCountInString(name) - len(name)) / 2)
			labelCount := labelCount + ((utf8.RuneCountInString(labels) - len(labels)) / 2)
			t := fmt.Sprintf("[blue]#%%-4d %%-7s [white]%%-%ds [green]%%-%ds [red]%%-%ds [white]%%-33s / %%-33s", titleCount, nameCount, labelCount)

			fmt.Println(color.Color(fmt.Sprintf(t,
				issue.LocalID,
				issue.State,
				title,
				name,
				labels,
				issue.CreatedAt,
				issue.UpdatedAt)))
		}
	}
}

func (gitLab *gitLabCli) findIssueState(s chan<- []*gogitlab.Issue, projectID int, findFunc func([]*gogitlab.Issue) bool) {
	page := 1
	for {
		issues, err := gitLab.ProjectIssues(projectID, page)
		if err != nil || len(issues) == 0 {
			break
		}
		page++

		s <- issues

		if findFunc(issues) {
			break
		}
	}
	close(s)
}

func (gitLab *gitLabCli) findIssueByID(projectID, issueID int) *gogitlab.Issue {
	c := make(chan []*gogitlab.Issue)
	go gitLab.findIssueState(c, projectID, func(issues []*gogitlab.Issue) bool {
		for _, issue := range issues {
			if issue.LocalID == issueID {
				return true
			}
		}
		return false
	})

	for {
		issues, ok := <-c
		if !ok {
			break
		}

		for _, issue := range issues {
			if issue.LocalID != issueID {
				continue
			}

			return issue
		}
	}
	return nil
}

func (gitLab *gitLabCli) PrintIssueDetail(projectID, issueID int) error {
	if issueID == 0 {
		fmt.Println("not issue_id")
		return nil
	}

	issue := gitLab.findIssueByID(projectID, issueID)
	if issue == nil {
		fmt.Println("not found issue")
		return nil
	}

	notes, err := gitLab.IssuesNotes(projectID, issue.ID)
	if err != nil {
		return err
	}

	t := template.Must(template.New("issueDetail").Parse(issueDetailTemplate))
	var buf bytes.Buffer
	if err := t.Execute(&buf, templateExec{Issue: issue, Notes: notes}); err != nil {
		return err
	}
	fmt.Println(color.Color(buf.String()))

	return nil
}

func trimPrefixIndex(s string) string {
	t := ""
	for _, r := range s {
		if checkTrim(t + string(r) + outTitleReplaceText) {
			break
		}
		t += string(r)
	}
	return t + outTitleReplaceText
}

func checkTrim(t string) bool {
	diff := ((utf8.RuneCountInString(t) - len(t)) / 2)
	if utf8.RuneCountInString(t) > (titleCount + diff) {
		return true
	}
	return false
}
