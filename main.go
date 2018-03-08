package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	//"github.com/google/go-github/github"
	//"golang.org/x/net/context"
)

type roundTripper struct {
	accessToken string
}

func (rt roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("token %s", rt.accessToken))
	return http.DefaultTransport.RoundTrip(r)
}

var (
	token = flag.String("token", os.Getenv("GITHUB_TOKEN"), "Github access token")
	owner = flag.String("owner", os.Getenv("GITHUB_OWNER"), "Github repository owner")
	repo  = flag.String("repo", os.Getenv("GITHUB_REPO"), "Github repository name")
)

func main() {
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

	http.DefaultClient.Transport = roundTripper{*token}
	//githubClient := github.NewClient(http.DefaultClient)

	fmt.Println("github-commenter: Sent comment to GitHub")
}
