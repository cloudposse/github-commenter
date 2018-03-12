package main

import (
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"strconv"
)

type roundTripper struct {
	accessToken string
}

func (rt roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("token %s", rt.accessToken))
	return http.DefaultTransport.RoundTrip(r)
}

var (
	token       = flag.String("token", os.Getenv("GITHUB_TOKEN"), "Github access token")
	owner       = flag.String("owner", os.Getenv("GITHUB_OWNER"), "Github repository owner")
	repo        = flag.String("repo", os.Getenv("GITHUB_REPO"), "Github repository name")
	commentType = flag.String("comment_type", os.Getenv("GITHUB_COMMENT_TYPE"), "Comment type: 'commit', 'pr' or 'issue'")
	sha         = flag.String("sha", os.Getenv("GITHUB_COMMIT_SHA"), "Commit SHA")
	number      = flag.String("number", os.Getenv("GITHUB_PR_ISSUE_NUMBER"), "Pull Request or Issue number")
	comment     = flag.String("comment", os.Getenv("GITHUB_COMMENT"), "Comment text")
)

func getPullRequestOrIssueNumber(str string) (int, error) {
	if str == "" {
		return 0, errors.New("-number or GITHUB_PR_ISSUE_NUMBER required")
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.WithMessage(err, "-number or GITHUB_PR_ISSUE_NUMBER must be an integer")
	}

	return num, nil
}

func getComment() (string, error) {
	if *comment != "" {
		return *comment, nil
	}
	return "", nil
}

func main() {
	//command := newGomplateCmd()
	//initFlags(command)
	//if err := command.Execute(); err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}

	flag.Parse()

	if *token == "" {
		flag.PrintDefaults()
		log.Fatal("-token or GITHUB_TOKEN required")
	}
	if *owner == "" {
		flag.PrintDefaults()
		log.Fatal("-owner or GITHUB_OWNER required")
	}
	if *repo == "" {
		flag.PrintDefaults()
		log.Fatal("-repo or GITHUB_REPO required")
	}
	if *commentType == "" {
		flag.PrintDefaults()
		log.Fatal("-comment_type or GITHUB_COMMENT_TYPE required")
	}
	if *commentType != "commit" && *commentType != "pr" && *commentType != "issue" {
		flag.PrintDefaults()
		log.Fatal("-comment_type or GITHUB_COMMENT_TYPE must be one of 'commit', 'pr' or 'issue'")
	}

	http.DefaultClient.Transport = roundTripper{*token}
	githubClient := github.NewClient(http.DefaultClient)

	if *commentType == "commit" {
		if *sha == "" {
			flag.PrintDefaults()
			log.Fatal("-sha or GITHUB_COMMIT_SHA required")
		}

		comment, err := getComment()
		if err != nil {
			log.Fatal(err)
		}

		// https://developer.github.com/v3/repos/comments
		commitComment := &github.RepositoryComment{Body: &comment}
		commitComment, _, err = githubClient.Repositories.CreateComment(context.Background(), *owner, *repo, *sha, commitComment)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("github-commenter: Created GitHub Commit comment", commitComment.ID)
	} else if *commentType == "pr" {
		num, err := getPullRequestOrIssueNumber(*number)
		if err != nil {
			log.Fatal(err)
		}

		comment, err := getComment()
		if err != nil {
			log.Fatal(err)
		}

		// https://developer.github.com/v3/pulls/comments
		pullRequestComment := &github.PullRequestComment{Body: &comment}
		pullRequestComment, _, err = githubClient.PullRequests.CreateComment(context.Background(), *owner, *repo, num, pullRequestComment)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("github-commenter: Created GitHub PR comment", pullRequestComment.ID)
	} else if *commentType == "issue" {
		num, err := getPullRequestOrIssueNumber(*number)
		if err != nil {
			log.Fatal(err)
		}

		comment, err := getComment()
		if err != nil {
			log.Fatal(err)
		}

		// https://developer.github.com/v3/issues/comments
		issueComment := &github.IssueComment{Body: &comment}
		issueComment, _, err = githubClient.Issues.CreateComment(context.Background(), *owner, *repo, num, issueComment)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("github-commenter: Created GitHub Issue comment", issueComment.ID)
	}
}
