package pkg

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/pkg/errors"
	"log"
	"time"
)

type Watcher interface {
	Watch() (notifications []*Notification, err error)
}

type JiraWatcher struct {
	JiraClient    *jira.Client
	JiraQuery     string
	LastQueryTime time.Time
}

func NewJiraWatchers(config *Config) ([]*JiraWatcher, error) {
	tp := jira.BasicAuthTransport{
		Username: config.JiraUsername,
		Password: config.JiraPassword,
	}

	jiraClient, err := jira.NewClient(tp.Client(), config.JiraBaseUrl)
	if err != nil {
		return nil, err
	}

	var jiraWatchers []*JiraWatcher

	for _, query := range config.JiraQuery {
		jiraWatchers = append(jiraWatchers, &JiraWatcher{
			JiraClient: jiraClient,
			JiraQuery:  query,
		})
	}
	return jiraWatchers, nil
}

func (jw *JiraWatcher) getQueryWithTime() string {
	formatTime := jw.LastQueryTime.Format("2006-01-02 15:04")
	query := fmt.Sprintf("%s and (created >= \"%s\" or updated >= \"%s\")", jw.JiraQuery, formatTime, formatTime)
	jw.LastQueryTime = time.Now()
	return query
}

func (jw *JiraWatcher) Watch() ([]*Notification, error) {
	maxResults := 100
	if jw.LastQueryTime.Before(time.Now().Add(-time.Hour * 72)) {
		maxResults = 1
	}

	query := jw.getQueryWithTime()
	log.Printf("running: %s", query)
	issues, response, err := jw.JiraClient.Issue.Search(query, &jira.SearchOptions{
		StartAt:    0,
		MaxResults: maxResults,
	})

	if err != nil {
		return nil, errors.Wrapf(err, "got a %d status when making a request to jira with query: %s", response.StatusCode, query)
	}

	var notifications []*Notification
	for _, issue := range issues {
		notifications = append(notifications, jw.convertJiraIssueToNotification(issue))
	}

	return notifications, nil
}

func (jw *JiraWatcher) convertJiraIssueToNotification(issue jira.Issue) *Notification {
	notification := newNotification()
	notification.addArg("subtitle", fmt.Sprintf("update to %s", issue.Key))
	notification.addArg("message", issue.Fields.Description)
	url := jw.JiraClient.GetBaseURL()
	notification.addArg("open", fmt.Sprintf("%sbrowse/%s", url.String(), issue.Key))

	return notification
}
