package main

import (
	"strconv"

	"github.com/kyokomi/go-gitlab-client/gogitlab"
	"strings"
	"net/url"
	"io/ioutil"
)

// 対象ProjectのProjectNameを取得する.
func GetProjectID(gitlab *gogitlab.Gitlab, projectName string) (int, error) {
	projects, err := gitlab.Projects()
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

func GetProjectName(gitlab *gogitlab.Gitlab, projectID int) (string, error) {
	project, err := gitlab.Project(strconv.Itoa(projectID))
	if err != nil {
		return "", err
	}
	return project.Name, nil
}

func GetUserName(gitlab *gogitlab.Gitlab, userId int) (string, error) {
	user, err := gitlab.User(strconv.Itoa(userId))
	if err != nil {
		return "", err
	}
	return user.Name, nil
}

//	/projects/:id/milestones/:milestone_id
func GetMilestoneTitle(gitlab *gogitlab.Gitlab, projectID, milestoneID int) (string, error) {
	milestone, err := gitlab.ProjectMilestone(strconv.Itoa(projectID), strconv.Itoa(milestoneID))
	if err != nil {
		return "", err
	}
	return milestone.Title, nil
}

func PostIssue(gitlab *gogitlab.Gitlab, projectID int, values url.Values) ([]byte, error) {
	reader := strings.NewReader(values.Encode())
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	res, err := gitlab.ProjectCreateIssues(projectID, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}
