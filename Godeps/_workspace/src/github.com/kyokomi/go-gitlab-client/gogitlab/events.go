package gogitlab

import (
	"encoding/xml"
	"fmt"
	"time"
	"encoding/json"
)

const (
	project_events_path = "/projects/%d/events.json"
)

type Person struct {
	Name  string `xml:"name"json:"name"`
	Email string `xml:"email"json:"email"`
}

type Link struct {
	Rel  string `xml:"rel,attr,omitempty"json:"rel"`
	Href string `xml:"href,attr"json:"href"`
}

type ActivityFeed struct {
	Title   string        `xml:"title"json:"title"`
	Id      string        `xml:"id"json:"id"`
	Link    []Link        `xml:"link"json:"link"`
	Updated time.Time     `xml:"updated,attr"json:"updated"`
	Entries []*FeedCommit `xml:"entry"json:"entries"`
}

type FeedCommit struct {
	Id      string    `xml:"id"json:"id"`
	Title   string    `xml:"title"json:"title"`
	Link    []Link    `xml:"link"json:"link"`
	Updated time.Time `xml:"updated"json:"updated"`
	Author  Person    `xml:"author"json:"author"`
	Summary string    `xml:"summary"json:"summary"`
	//<media:thumbnail width="40" height="40" url="https://secure.gravatar.com/avatar/7070eab7c6206530d3b7820362227fec?s=40&amp;d=mm"/>
}

type Event struct {
	Title string          `json:"title"`
	ProductId int         `json:"project_id"`
	ActionName string     `json:"action_name"`
	TargetId int          `json:"target_id"`
	TargetType string     `json:"target_type"`
	TargetTitle string    `json:"target_title"`
	Data EventData        `json:"data"`
	CreatedAt time.Time   `json:"created_at"`
}

type EventData struct {
	Before string                   `json:"before"`
	After string                    `json:"after"`
	AuthorId int                    `json:"author_id"`
	Ref string                      `json:"ref"`
	UserId int                      `json:"user_id"`
	UserName string                 `json:"user_name"`
	ProductId int                   `json:"project_id"`
	Repository EventRepository      `json:"repository"`
	Commits []EventCommit           `json:"commits"`
	TotalCommitsCount int           `json:"total_commits_count"`
}

type EventRepository struct {
	Name string                     `json:"name"`
	GitUrl string                   `json:"url"`
	Description string              `json:"description"`
	PageUrl string                  `json:"homepage"`
}

type EventCommit struct {
	Id string                       `json:"id"`
	Message string                  `json:"message"`
	Timestamp string                `json:"timestamp"`
	Url string                      `json:"url"`
	Author EventAuthor              `json:"author"`
}

type EventAuthor struct {
	Name string                     `json:"name"`
	Email string                    `json:"email"`
}

func (g *Gitlab) Activity() (ActivityFeed, error) {

	url := g.BaseUrl + dasboard_feed_path + "?private_token=" + g.Token

	contents, err := g.buildAndExecRequest("GET", url, nil)
	if err != nil {
		fmt.Println("%s", err)
	}

	var activity ActivityFeed
	err = xml.Unmarshal(contents, &activity)
	if err != nil {
		fmt.Println("%s", err)
	}

	return activity, err
}

func (g *Gitlab) RepoActivityFeed(feedPath string) ActivityFeed {

	url := g.BaseUrl + g.RepoFeedPath + "?private_token=" + g.Token

	contents, err := g.buildAndExecRequest("GET", url, nil)
	if err != nil {
		fmt.Println("%s", err)
	}

	var activity ActivityFeed
	err = xml.Unmarshal(contents, &activity)
	if err != nil {
		fmt.Println("%s", err)
	}

	return activity
}

func (g *Gitlab) ProjectEvents(projectId int) []Event {
	url := g.ResourceUrl(fmt.Sprintf(project_events_path, projectId), nil)

	contents, err := g.buildAndExecRequest("GET", url, nil)
	if err != nil {
		fmt.Println("%s", err)
	}

	var events []Event
	err = json.Unmarshal(contents, &events)
	if err != nil {
		fmt.Println("%s", err)
	}

	return events
}
