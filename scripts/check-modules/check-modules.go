package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var owner = "ConsenSys"
var requestCounter = 0

// ModData data.
type ModData struct {
	Go      string       `json:"Go"`
	Require []ModRequire `json:"Require"`
	Replace []ModReplace `json:"Replace"`
}

// ModRequire data.
type ModRequire struct {
	Path    string `json:"Path"`
	Version string `json:"Version"`
}

// ModReplace data.
type ModReplace struct {
	Old ModPath `json:"Old"`
	New ModPath `json:"New"`
}

// ModPath data.
type ModPath struct {
	Path string `json:"Path"`
}

// BranchData data.
type BranchData struct {
	Name         string       `json:"name"`
	BranchCommit BranchCommit `json:"commit"`
	Protected    bool         `json:"protected"`
}

// BranchCommit data.
type BranchCommit struct {
	Sha string `json:"sha"`
	URL bool   `json:"url"`
}

// Commit data.
type Commit struct {
	Sha    string       `json:"sha"`
	Detail CommitDetail `json:"commit"`
}

// CommitDetail data.
type CommitDetail struct {
	Author CommitAuthor `json:"author"`
}

// CommitAuthor data.
type CommitAuthor struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

// GithubError data.
type GithubError struct {
	Message string `json:"message"`
}

// getPackageList returns package list
func getPackageList() ModData {
	cmd := exec.Command("go", "mod", "edit", "-json")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	r := bytes.NewReader(output)
	var packages ModData
	json.NewDecoder(r).Decode(&packages)

	return packages
}

// getBranches returns branches list
func getBranches(pkg ModRequire) (string, string, []BranchData) {
	repos := strings.Split(pkg.Path, "/")
	repo := repos[len(repos)-1]
	versions := strings.Split(pkg.Version, "-")
	version := versions[len(versions)-1]

	requestCounter++
	url := "https://api.github.com/repos/ConsenSys/" + repo + "/branches"

	httpClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		catchError(res.Body)
	}

	var branches []BranchData
	json.NewDecoder(res.Body).Decode(&branches)
	return repo, version, branches
}

// getCommits returns commits list
func getCommits(repo string, sha string) []Commit {
	requestCounter++
	url := "https://api.github.com/repos/ConsenSys/" + repo + "/commits?sha=" + sha
	httpClient := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	var commits []Commit
	json.NewDecoder(res.Body).Decode(&commits)
	return commits
}

// findInBranches finds a commit in all branches
func findInBranches(branches []BranchData, version string, repo string) (BranchData, Commit) {
	for _, branch := range branches {
		commits := getCommits(repo, branch.BranchCommit.Sha)
		branchFound, commit := findInBranch(branch, commits, version)
		if branchFound {
			return branch, commit
		}
	}
	return BranchData{}, Commit{}
}

// findInBranch finds a commit in a branch
func findInBranch( branch BranchData, commits []Commit, version string) (bool, Commit) {
	for _, commit := range commits {
		re := regexp.MustCompile(version)
		match := re.FindStringSubmatch(commit.Sha)
		if len(match) > 0 {
			return true, commit
		}
	}
	return false, Commit{}
}

// catchError returns GitHub api error
func catchError(body io.ReadCloser) {
	var error GithubError
	json.NewDecoder(body).Decode(&error)
	fmt.Printf("request counter: %d\n", requestCounter)
	fmt.Printf("GitHub error: %s\n", error.Message)
	os.Exit(1)
}

func main() {
	fmt.Printf("Start checking packages:\n\n")
	packages := getPackageList()
	for _, pkg := range packages.Require {

		re := regexp.MustCompile(owner)
		match := re.FindStringSubmatch(pkg.Path)
		if len(match) == 0 {
			fmt.Printf("skip repo: %s\n", pkg.Path)
			continue
		}

		repo, version, branches := getBranches(pkg)
		fmt.Printf("%s - %s\n", repo, version)

		branch, commit := findInBranches(branches, version, repo)
		if (branch != BranchData{}) {
			fmt.Printf("branch: %s - commit by: %s - on %s\n", branch.Name, commit.Detail.Author.Name, commit.Detail.Author.Date)
			fmt.Printf("url: %s\n", "https://github.com/ConsenSys/"+repo+"/tree/"+commit.Sha)
		} else {
			fmt.Printf("branch not found\n")
		}
		fmt.Println()
	}

	fmt.Printf("\nTotal request counter: %d\n", requestCounter)
}
