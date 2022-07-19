package main

import (
	"fmt"
	"log"
	"os"

	"Ex4.11/github"
)

func main() {
	argc := len(os.Args)
	if argc != 5 && argc != 4 {
		log.Fatalln("Usage: ./IssueManager [read|create|update|close] OWNER REPO ISSUE_NUMBER")
	}
	var operation = os.Args[1]
	var owner = os.Args[2]
	var repo = os.Args[3]
	var issueNumber string
	if argc == 5 {
		issueNumber = os.Args[4]
	}
	switch operation {
	case "read":
		issue, err := github.GetIssue(owner, repo, issueNumber)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("#%-5d %9.9s %.55s %v\n",
			issue.Number, issue.User.Login, issue.Title, issue.State)
	case "create":
		issue, err := github.CreateIssue(owner, repo)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("#%-5d %9.9s %.55s %v\n",
			issue.Number, issue.User.Login, issue.Title, issue.State)
	case "update":
		issue, err := github.EditIssue(owner, repo, issueNumber)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("#%-5d %9.9s %.55s %v\n",
			issue.Number, issue.User.Login, issue.Title, issue.State)
	case "close":
		issue, err := github.CloseIssue(owner, repo, issueNumber)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("#%-5d %9.9s %.55s %v\n",
			issue.Number, issue.User.Login, issue.Title, issue.State)
	default:
		log.Fatal("Unknown operation!")
	}
}
