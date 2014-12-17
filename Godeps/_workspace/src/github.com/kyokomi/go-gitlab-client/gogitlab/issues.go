package gogitlab

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	issues_url                  = "/issues/"                     // Get a specific issues
	project_issues_url          = "/projects/%d/issues"         // Get a specific issues / Post a create issues
)

type Issue struct {
	Id           int        `json:"id"`
	LocalId      int        `json:"iid"`
	ProjectId    int        `json:"project_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Labels       []string   `json:"labels"`
	Milestone    Milestone  `json:"milestone"`
	Author       Person     `json:"author"`
	Assignee     Person     `json:"assignee"`
	State        string     `json:"state"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatedAt    time.Time  `json:"created_at,omitempty"`
	// AccessLevel int
}

/*
Get a list of issues by the authenticated user.
*/
func (g *Gitlab) Issues() ([]*Issue, error) {

	url := g.ResourceUrl(issues_url, nil)

	var issues []*Issue

	contents, err := g.buildAndExecRequest("GET", url, nil)
	if err == nil {
		err = json.Unmarshal(contents, &issues)
	}

	return issues, err
}

func (g *Gitlab) ProjectIssues(projectId int, pageNo int) ([]*Issue, error) {

	url := g.ResourceUrl(fmt.Sprintf(project_issues_url, projectId), nil)
	url += fmt.Sprintf("&page=%d", pageNo)

	var issues []*Issue

	contents, err := g.buildAndExecRequest("GET", url, nil)
	if err == nil {
		err = json.Unmarshal(contents, &issues)
	}

	return issues, err
}

func (g *Gitlab) ProjectCreateIssues(projectId int, data []byte) ([]*Issue, error) {

	url := g.ResourceUrl(fmt.Sprintf(project_issues_url, projectId), nil)

	var issues []*Issue

	contents, err := g.buildAndExecRequest("POST", url, data)
	if err == nil {
		err = json.Unmarshal(contents, &issues)
	}

	return issues, err
}
