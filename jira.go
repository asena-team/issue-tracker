package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"os"
)

func NewIssue(issue *Issue) *jira.Issue {
	return &jira.Issue{
		Fields: &jira.IssueFields{
			Project: jira.Project{
				Key: Project,
			},
			Type: jira.IssueType{
				Name: issue.Type,
			},
			Priority: &jira.Priority{
				Name: issue.Priority,
			},
			Labels:      []string{"automation"},
			Summary:     issue.Title,
			Description: buildDescription(issue),
		},
	}
}

func buildDescription(issue *Issue) string {
	return fmt.Sprintf("Reporter: %s <%s>\n\n%s", issue.Reporter, issue.Mail, issue.Description)
}

func JiraClient() *jira.Client {
	username := os.Getenv("JIRA_USERNAME")
	token := os.Getenv("JIRA_TOKEN")
	if username == "" || token == "" {
		panic("Please put JIRA_USERNAME and JIRA_TOKEN env vars.")
	}

	tp := jira.BasicAuthTransport{
		Username: username,
		Password: token,
	}

	jiraClient, err := jira.NewClient(tp.Client(), BaseUrl)
	if jiraClient == nil {
		panic(err)
	}

	return jiraClient
}
