package main

import (
	"strconv"

	"strings"
	"net/url"
	"io/ioutil"
	"github.com/kyokomi/go-gitlab-client/gogitlab"
	"unicode/utf8"
	"github.com/mitchellh/colorstring"
	"fmt"
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

func (gitLab *gitLabCli) GetUserName(userId int) (string, error) {
	user, err := gitLab.User(strconv.Itoa(userId))
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
			titleCount := 90 + ((utf8.RuneCountInString(issue.Title) - len(issue.Title)) / 2)
			nameCount := 16 + ((utf8.RuneCountInString(issue.Assignee.Name) - len(issue.Assignee.Name)) / 2)
			t := fmt.Sprintf("[blue]#%%-4d %%-7s [white]%%-%ds [green]%%-%ds [white]%%-33s / %%-33s", titleCount, nameCount)
			fmt.Println(colorstring.Color(fmt.Sprintf(t,
				issue.LocalID,
				issue.State,
				issue.Title,
				issue.Assignee.Name,
				issue.CreatedAt,
				issue.UpdatedAt)))
		}
	}
}
