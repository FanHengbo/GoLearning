package main

import (
	"Ex4.11/github"
)

func ReadIssue(owner, repo, issueNumber string) {

}
func main() {
	/*
		if argc := len(os.Args); argc != 5 && argc != 4 {
			log.Fatalln("Usage: ./IssueManager [read|create|update|delete] OWNER REPO ISSUE_NUMBER")
		}
	*/
	github.CreateIssue("FanHengbo", "GoLearning")
	//issue, _ := github.GetIssue("golang", "go", "35")
	//fmt.Println(issue.Title, issue.Body)\
	//github.EditIssue("FanHengbo", "GoLearning", "1")
}
