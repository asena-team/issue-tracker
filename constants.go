package main

// Current API Version
const APIVersion = 1

// Atlassian Definitions
const (
	BaseUrl = "https://asena-team.atlassian.net"
	Project = "AS"
)

// Issue Types
var IssueTypes = []string{
	"New Feature",
	"Improvement",
	"Bug",
	"Security",
	"Proposal",
	"Other",
}

// Issue Priorities
var Priorities = []string{
	"Highest",
	"High",
	"Medium",
	"Low",
	"Lowest",
}
