package main

import (
	"fmt"
	"os"

	"github.com/andygrunwald/go-jira"
)

func main() {
	jiraClient := newJiraClient()

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Description: "Test Issue",
			Type: jira.IssueType{
				Name: "Bug",
			},
			Project: jira.Project{
				Key: Project,
			},
			Summary: "Just a demo issue",
		},
	}
	issue, _, err := jiraClient.Issue.Create(&i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %+v\n", issue.Key, i.Fields.Summary)

}

func newJiraClient() *jira.Client {
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
