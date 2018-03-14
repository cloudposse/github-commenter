#!/bin/bash

export GITHUB_TOKEN=XXXXXXXXXXXXXXXX
export GITHUB_OWNER=cloudposse
export GITHUB_REPO=github-commenter
export GITHUB_COMMENT_TYPE=pr
export GITHUB_PR_ISSUE_NUMBER=1
export GITHUB_COMMENT_FORMAT="My comment:<br/>{{.}}"
export GITHUB_COMMENT="+1 LGTM"

../dist/bin/github-commenter
