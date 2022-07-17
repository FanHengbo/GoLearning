package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"gopl.io/ch4/github"
)

var timePtr = flag.String("time", "all", "issues in age categories")
var now = time.Now()
var yearDayNow = now.Year()*365 + now.YearDay()

func ClassifyBasedOnAge(t *string, issue *Issue) bool {
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

	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		if ClassifyBasedOnAge(timePtr, item) {
			fmt.Printf("#%-5d %9.9s %.55s %v\n",
				item.Number, item.User.Login, item.Title, item.CreatedAt)
		}
	}
}
