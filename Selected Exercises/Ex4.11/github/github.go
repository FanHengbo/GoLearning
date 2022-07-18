// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 110.
//!+

// Package github provides a Go API for the GitHub issue tracker.
// See https://developer.github.com/v3/search/#search-issues.
package github

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

const IssuesURL = "https://api.github.com/repos"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}
	//!-
	// For long-term stability, instead of http.Get, use the
	// variant below which adds an HTTP request header indicating
	// that only version 3 of the GitHub API is acceptable.
	//
	//   req, err := http.NewRequest("GET", IssuesURL+"?q="+q, nil)
	//   if err != nil {
	//       return nil, err
	//   }
	//   req.Header.Set(
	//       "Accept", "application/vnd.github.v3.text-match+json")
	//   resp, err := http.DefaultClient.Do(req)
	//!+

	// We must close resp.Body on all execution paths.
	// (Chapter 5 presents 'defer', which makes this simpler.)
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}
func GetIssue(owner, repo, number string) (*Issue, error) {
	resp, err := http.Get(IssuesURL + "/" + owner + "/" + repo + "/issues/" + number)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}
	var result Issue
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// issue number is ommited here
func CreateIssue(owner, repo string) (*Issue, error) {

	vi := "vim"
	tmpDir := os.TempDir()
	tmpFile, tmpFileErr := ioutil.TempFile(tmpDir, "tempFile")
	if tmpFileErr != nil {
		return nil, fmt.Errorf("error %s while creating tempFile", tmpFileErr)
	}
	path, execErr := exec.LookPath(vi)
	if execErr != nil {
		return nil, fmt.Errorf("error %s while looking up for %s", path, vi)
	}
	fmt.Printf("%s is available at %s\nCalling it with file %s \n", vi, path, tmpFile.Name())

	cmd := exec.Command(path, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start failed: %s", err)
	}
	fmt.Printf("Waiting for command to finish.\n")
	if err := cmd.Wait(); err != nil {
		//fmt.Printf("command finished with error: %v", err)
		return nil, fmt.Errorf("command finished with error: %v", err)
	}

	scanner := bufio.NewScanner(tmpFile)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	//concatenate multiple lines of issue into a singe string, which is the issue body
	issueBody := strings.Join(lines, "")
	tmpFile.Close()
	os.Remove(tmpFile.Name())

	//Ask user to input issue title
	fmt.Println("Please enter issue title")
	var issueTitle string
	fmt.Scanln(&issueTitle)

	//Display issue title and body
	fmt.Println("title: ", issueTitle)
	fmt.Println("Body: ", issueBody)
	return nil, fmt.Errorf("not end yet")
}

//!-
