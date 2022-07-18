package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"gopl.io/ch4/github"
)

var timePtr = flag.String("time", "all", "issues in age categories")
var issueUrl = flag.String("repo", "repo:golang/go", "repo")
var now = time.Now()
var yearDayNow = now.Year()*365 + now.YearDay()

func AgeFilter(t *string, issue *github.Issue) bool {
	if *t == "all" {
		return true
	}
	issueYearDay := issue.CreatedAt.Year()*365 + issue.CreatedAt.YearDay()
	if *t == "month" {
		if yearDayNow-issueYearDay < 30 {
			return true
		}
	}
	if *t == "year" {
		if yearDayNow-issueYearDay > 365 {
			return true
		}
	}
	return false
}
func main() {
	flag.Parse()
	url := []string{*issueUrl}
	result, err := github.SearchIssues(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	fmt.Println(now, " ", yearDayNow)
	for _, item := range result.Items {
		if AgeFilter(timePtr, item) {
			fmt.Printf("#%-5d %9.9s %.55s %10v\n",
				item.Number, item.User.Login, item.Title, item.CreatedAt)
		}
	}
}
