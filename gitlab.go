package main

import (
	"strconv"

	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/kyokomi/go-gitlab-client/gogitlab"
	color "github.com/mitchellh/colorstring"
)

const (
	titleCount          = 60
	nameCount           = 16
	outTitleReplaceText = " ..."
)

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

func (gitLab *gitLabCli) PostIssue(projectID int, values url.Values) ([]byte, error) {
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

func (gitLab *gitLabCli) PrintIssue(projectID int) {
	c := make(chan []*gogitlab.Issue)
	go func(s chan<- []*gogitlab.Issue) {
		page := 1
		for {
			issues, err := gitLab.ProjectIssues(projectID, page)
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

			title := issue.Title
			if len(title) > titleCount {
				title = trimPrefixIndex(title)
			}

			name := issue.Assignee.Name
			if len(name) > nameCount {
				name = trimPrefixIndex(name)
			}

			titleCount := titleCount + ((utf8.RuneCountInString(title) - len(title)) / 2)
			nameCount := nameCount + ((utf8.RuneCountInString(name) - len(name)) / 2)
			t := fmt.Sprintf("[blue]#%%-4d %%-7s [white]%%-%ds [green]%%-%ds [white]%%-33s / %%-33s", titleCount, nameCount)

			fmt.Println(color.Color(fmt.Sprintf(t,
				issue.LocalID,
				issue.State,
				title,
				name,
				issue.CreatedAt,
				issue.UpdatedAt)))
		}
	}
}

func trimPrefixIndex(s string) string {
	t := ""
	for _, r := range s {
		t += string(r)
		if len(t) >= (titleCount - len(outTitleReplaceText)) {
			break
		}
	}
	return t + outTitleReplaceText
}
