package main

import (
	"os"

	"github.com/andygrunwald/go-jira"
)

func NewIssue(issue Issue) jira.Issue {
	return jira.Issue{
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
			Summary:     issue.Title,
			Description: issue.Description,
		},
	}
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
