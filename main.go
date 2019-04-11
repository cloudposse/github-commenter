package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"text/template"
)

type roundTripper struct {
	accessToken string
}

func (rt roundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("token %s", rt.accessToken))
	return http.DefaultTransport.RoundTrip(r)
}

var (
	token              = flag.String("token", os.Getenv("GITHUB_TOKEN"), "Github access token")
	owner              = flag.String("owner", os.Getenv("GITHUB_OWNER"), "Github repository owner")
	repo               = flag.String("repo", os.Getenv("GITHUB_REPO"), "Github repository name")
	commentType        = flag.String("type", os.Getenv("GITHUB_COMMENT_TYPE"), "Comment type: 'commit', 'pr', 'issue', 'pr-review' or 'pr-file'")
	sha                = flag.String("sha", os.Getenv("GITHUB_COMMIT_SHA"), "Commit SHA")
	number             = flag.String("number", os.Getenv("GITHUB_PR_ISSUE_NUMBER"), "Pull Request or Issue number")
	file               = flag.String("file", os.Getenv("GITHUB_PR_FILE"), "Pull Request File Name")
	position           = flag.String("position", os.Getenv("GITHUB_PR_FILE_POSITION"), "Position in Pull Request File")
	templ              = flag.String("template", os.Getenv("GITHUB_COMMENT_TEMPLATE"), "Template to format comment. Supports `Go` templates: My comment:<br/>{{.}}. Use either `template` or `template_file`")
	templateFile       = flag.String("template_file", os.Getenv("GITHUB_COMMENT_TEMPLATE_FILE"), "The path to a template file to format comment. Supports `Go` templates. Use either `template` or `template_file`")
	format             = flag.String("format", os.Getenv("GITHUB_COMMENT_FORMAT"), "Alias of `template`")
	formatFile         = flag.String("format_file", os.Getenv("GITHUB_COMMENT_FORMAT_FILE"), "Alias of `template_file`")
	comment            = flag.String("comment", os.Getenv("GITHUB_COMMENT"), "Comment text")
	deleteCommentRegex = flag.String("delete-comment-regex", os.Getenv("GITHUB_DELETE_COMMENT_REGEX"), "Regex to find previous comments to delete before creating the new comment. Supported for comment types `commit`, `pr-file`, `issue` and `pr`")
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

func getPullRequestFilePosition(str string) (int, error) {
	if str == "" {
		return 0, errors.New("-position or GITHUB_PR_FILE_POSITION required")
	}

	position, err := strconv.Atoi(str)
	if err != nil {
		return 0, errors.WithMessage(err, "-position or GITHUB_PR_FILE_POSITION must be an integer")
	}

	return position, nil
}

func getComment() (string, error) {
	// Read the comment from the command-line argument or ENV var first
	if *comment != "" {
		return *comment, nil
	}

	// Read from stdin
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return "", errors.WithMessage(err, "Comment must be provided either as command-line argument, ENV variable, or from 'stdin'")
	}

	return string(data), nil
}

func formatComment(comment string) (string, error) {
	if *format == "" && *formatFile == "" && *templ == "" && *templateFile == "" {
		return comment, nil
	}

	var t *template.Template
	var err error
	var templateFinal string
	var templateFileFinal string

	if *format != "" || *templ != "" {
		if *templ != "" {
			templateFinal = *templ
		} else {
			templateFinal = *format
		}
		t = template.New("formatComment").Funcs(sprig.TxtFuncMap())
		t, err = t.Parse(templateFinal)
		if err != nil {
			return "", err
		}
	} else {
		if *templateFile != "" {
			templateFileFinal = *templateFile
		} else {
			templateFileFinal = *formatFile
		}
		name := path.Base(templateFileFinal)
		t = template.New(name).Funcs(sprig.TxtFuncMap())
		t, err = t.ParseFiles(templateFileFinal)
		if err != nil {
			return "", err
		}
	}

	var doc bytes.Buffer

	err = t.Execute(&doc, comment)
	if err != nil {
		return "", err
	}

	// Remove ANSI escape codes
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"
	var re = regexp.MustCompile(ansi)
	var s = doc.String()
	return re.ReplaceAllString(s, ""), nil
}

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
	if *commentType == "" {
		flag.PrintDefaults()
		log.Fatal("-type or GITHUB_COMMENT_TYPE required")
	}
	if *commentType != "commit" && *commentType != "pr" && *commentType != "issue" && *commentType != "pr-review" && *commentType != "pr-file" {
		flag.PrintDefaults()
		log.Fatal("-type or GITHUB_COMMENT_TYPE must be one of 'commit', 'pr', 'issue', 'pr-review' or 'pr-file'")
	}

	http.DefaultClient.Transport = roundTripper{*token}
	githubClient := github.NewClient(http.DefaultClient)

	// https://developer.github.com/v3/guides/working-with-comments
	// https://developer.github.com/v3/repos/comments
	if *commentType == "commit" {
		if *sha == "" {
			flag.PrintDefaults()
			log.Fatal("-sha or GITHUB_COMMIT_SHA required")
		}

		comment, err := getComment()
		if err != nil {
			log.Fatal(err)
		}

		formattedComment, err := formatComment(comment)
		if err != nil {
			log.Fatal(err)
		}

		// Find and delete existing comment(s) before creating the new one
		if *deleteCommentRegex != "" {
			r, err := regexp.Compile(*deleteCommentRegex)
			if err != nil {
				log.Fatal(err)
			}

			listOptions := &github.ListOptions{}
			comments, _, err := githubClient.Repositories.ListCommitComments(context.Background(), *owner, *repo, *sha, listOptions)
			if err != nil {
				log.Println("github-commenter: Error listing commit comments: ", err)
			} else {
				for _, comment := range comments {
					if r.MatchString(*comment.Body) {
						_, err = githubClient.Repositories.DeleteComment(context.Background(), *owner, *repo, *comment.ID)
						if err != nil {
							log.Println("github-commenter: Error deleting commit comment: ", err)
						} else {
							log.Println("github-commenter: Deleted commit comment: ", *comment.ID)
						}
					}
				}
			}
		}

		commitComment := &github.RepositoryComment{Body: &formattedComment}
		commitComment, _, err = githubClient.Repositories.CreateComment(context.Background(), *owner, *repo, *sha, commitComment)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("github-commenter: Created GitHub Commit comment", commitComment.ID)
	} else if *commentType == "pr-review" {
		// https://developer.github.com/v3/pulls/reviews/#create-a-pull-request-review
		num, err := getPullRequestOrIssueNumber(*number)
		if err != nil {
			log.Fatal(err)
		}

		comment, err := getComment()
		if err != nil {
			log.Fatal(err)
		}

		formattedComment, err := formatComment(comment)
		if err != nil {
			log.Fatal(err)
		}

		pullRequestReviewRequest := &github.PullRequestReviewRequest{Body: &formattedComment, Event: github.String("COMMENT")}
		pullRequestReview, _, err := githubClient.PullRequests.CreateReview(context.Background(), *owner, *repo, num, pullRequestReviewRequest)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("github-commenter: Created GitHub PR Review comment", pullRequestReview.ID)
	} else if *commentType == "issue" || *commentType == "pr" {
		// https://developer.github.com/v3/issues/comments
		num, err := getPullRequestOrIssueNumber(*number)
		if err != nil {
			log.Fatal(err)
		}

		comment, err := getComment()
		if err != nil {
			log.Fatal(err)
		}

		formattedComment, err := formatComment(comment)
		if err != nil {
			log.Fatal(err)
		}

		// Find and delete existing comment(s) before creating the new one
		if *deleteCommentRegex != "" {
			r, err := regexp.Compile(*deleteCommentRegex)
			if err != nil {
				log.Fatal(err)
			}

			listOptions := &github.IssueListCommentsOptions{}
			comments, _, err := githubClient.Issues.ListComments(context.Background(), *owner, *repo, num, listOptions)
			if err != nil {
				log.Println("github-commenter: Error listing Issue/PR comments: ", err)
			} else {
				for _, comment := range comments {
					if r.MatchString(*comment.Body) {
						_, err = githubClient.Issues.DeleteComment(context.Background(), *owner, *repo, *comment.ID)
						if err != nil {
							log.Println("github-commenter: Error deleting Issue/PR comment: ", err)
						} else {
							log.Println("github-commenter: Deleted Issue/PR comment: ", *comment.ID)
						}
					}
				}
			}
		}

		issueComment := &github.IssueComment{Body: &formattedComment}
		issueComment, _, err = githubClient.Issues.CreateComment(context.Background(), *owner, *repo, num, issueComment)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("github-commenter: Created GitHub Issue comment", issueComment.ID)
	} else if *commentType == "pr-file" {
		// https://developer.github.com/v3/pulls/comments
		num, err := getPullRequestOrIssueNumber(*number)
		if err != nil {
			log.Fatal(err)
		}

		if *sha == "" {
			flag.PrintDefaults()
			log.Fatal("-sha or GITHUB_COMMIT_SHA required")
		}

		if *file == "" {
			flag.PrintDefaults()
			log.Fatal("-file or GITHUB_PR_FILE required")
		}

		position, err := getPullRequestFilePosition(*position)
		if err != nil {
			log.Fatal(err)
		}

		comment, err := getComment()
		if err != nil {
			log.Fatal(err)
		}

		formattedComment, err := formatComment(comment)
		if err != nil {
			log.Fatal(err)
		}

		// Find and delete existing comment(s) before creating the new one
		if *deleteCommentRegex != "" {
			r, err := regexp.Compile(*deleteCommentRegex)
			if err != nil {
				log.Fatal(err)
			}

			listOptions := &github.PullRequestListCommentsOptions{}
			comments, _, err := githubClient.PullRequests.ListComments(context.Background(), *owner, *repo, num, listOptions)
			if err != nil {
				log.Println("github-commenter: Error listing PR file comments: ", err)
			} else {
				for _, comment := range comments {
					if r.MatchString(*comment.Body) {
						_, err = githubClient.PullRequests.DeleteComment(context.Background(), *owner, *repo, *comment.ID)
						if err != nil {
							log.Println("github-commenter: Error deleting PR file comment: ", err)
						} else {
							log.Println("github-commenter: Deleted PR file comment: ", *comment.ID)
						}
					}
				}
			}
		}

		pullRequestComment := &github.PullRequestComment{Body: &formattedComment, Path: file, Position: &position, CommitID: sha}
		pullRequestComment, _, err = githubClient.PullRequests.CreateComment(context.Background(), *owner, *repo, num, pullRequestComment)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("github-commenter: Created GitHub PR comment on file: ", pullRequestComment.ID)
	}
}
