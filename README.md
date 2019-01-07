<!-- This file was automatically generated by the `build-harness`. Make all changes to `README.yaml` and run `make readme` to rebuild this file. -->
[![README Header][readme_header_img]][readme_header_link]

[![Cloud Posse][logo]](https://cpco.io/homepage)

# github-commenter [![Build Status](https://travis-ci.org/cloudposse/github-commenter.svg?branch=master)](https://travis-ci.org/cloudposse/github-commenter) [![Latest Release](https://img.shields.io/github/release/cloudposse/github-commenter.svg)](https://github.com/cloudposse/github-commenter/releases/latest) [![Slack Community](https://slack.cloudposse.com/badge.svg)](https://slack.cloudposse.com)


Command line utility for creating GitHub comments on Commits, Pull Request Reviews, Pull Request Files, Issues and Pull Requests.

GitHub API supports these types of comments:

* [Comments on Repos/Commits](https://developer.github.com/v3/repos/comments)
* [Comments on Pull Request Reviews](https://developer.github.com/v3/pulls/reviews/#create-a-pull-request-review)
* [Comments on Pull Request Files](https://developer.github.com/v3/pulls/comments)
* [Comments on Issues](https://developer.github.com/v3/issues/comments)
* [Comments on Pull Requests (in the global section)](https://developer.github.com/v3/issues/comments)

Since GitHub considers Pull Requests as Issues, `Comments on Issues` and `Comments on Pull Requests` use the same API.

The utility supports all these types of comments (`commit`, `pr-review`, `pr-file`, `issue`, `pr`).


---

This project is part of our comprehensive ["SweetOps"](https://cpco.io/sweetops) approach towards DevOps. 
[<img align="right" title="Share via Email" src="https://docs.cloudposse.com/images/ionicons/ios-email-outline-2.0.1-16x16-999999.svg"/>][share_email]
[<img align="right" title="Share on Google+" src="https://docs.cloudposse.com/images/ionicons/social-googleplus-outline-2.0.1-16x16-999999.svg" />][share_googleplus]
[<img align="right" title="Share on Facebook" src="https://docs.cloudposse.com/images/ionicons/social-facebook-outline-2.0.1-16x16-999999.svg" />][share_facebook]
[<img align="right" title="Share on Reddit" src="https://docs.cloudposse.com/images/ionicons/social-reddit-outline-2.0.1-16x16-999999.svg" />][share_reddit]
[<img align="right" title="Share on LinkedIn" src="https://docs.cloudposse.com/images/ionicons/social-linkedin-outline-2.0.1-16x16-999999.svg" />][share_linkedin]
[<img align="right" title="Share on Twitter" src="https://docs.cloudposse.com/images/ionicons/social-twitter-outline-2.0.1-16x16-999999.svg" />][share_twitter]




It's 100% Open Source and licensed under the [APACHE2](LICENSE).











## Screenshots


![PR](images/github-pr-review-comment.png)
*GitHub PR Review Comment*



## Usage

__NOTE__: Create a [GitHub token](https://help.github.com/articles/creating-an-access-token-for-command-line-use) with `repo:status` and `public_repo` scopes.

__NOTE__: The utility accepts parameters as command-line arguments or as ENV variables (or any combination of command-line arguments and ENV vars).
Command-line arguments take precedence over ENV vars.


| Command-line argument |  ENV var                     |  Description                                                                                                                                                                                                |
|:----------------------|:-----------------------------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| token                 | GITHUB_TOKEN                 | Github access [token](https://help.github.com/articles/creating-an-access-token-for-command-line-use)                                                                                                       |
| owner                 | GITHUB_OWNER                 | Github repository owner (_e.g._ `cloudposse`)                                                                                                                                                               |
| repo                  | GITHUB_REPO                  | Github repository name (_e.g._ `github-commenter`)                                                                                                                                                          |
| type                  | GITHUB_COMMENT_TYPE          | Comment type: `commit`, `pr`, `issue`, `pr-review` or `pr-file`                                                                                                                                             |
| sha                   | GITHUB_COMMIT_SHA            | Commit SHA. Required when `type=commit` or `type=pr-file`                                                                                                                                                   |
| number                | GITHUB_PR_ISSUE_NUMBER       | Pull Request or Issue number. Required for all comment types except for `commit`                                                                                                                            |
| file                  | GITHUB_PR_FILE               | Pull Request File Name to comment on. For more info see [create comment](https://developer.github.com/v3/pulls/comments/#create-a-comment)                                                                  |
| position              | GITHUB_PR_FILE_POSITION      | Position in Pull Request File. For more info see [create comment](https://developer.github.com/v3/pulls/comments/#create-a-comment)                                                                         |
| format                | GITHUB_COMMENT_FORMAT        | Comment format template (optional). Supports `Go` [templates](https://golang.org/pkg/text/template). _E.g._ `My comment:<br/>{{.}}`. Use either `format` or `format_file`                                   |
| format_file           | GITHUB_COMMENT_FORMAT_FILE   | The path to a template file to format comment (optional). Supports `Go` templates. Use either `format` or `format_file`                                                                                     |
| template              | GITHUB_COMMENT_TEMPLATE      | Comment format template (optional). This is an alias to `format`. Supports `Go` [templates](https://golang.org/pkg/text/template). _E.g._ `My comment:<br/>{{.}}`. Use either `template` or `template_file` |
| template_file         | GITHUB_COMMENT_TEMPLATE_FILE | The path to a template file to format comment (optional). This is an alias to `format_file`. Supports `Go` templates. Use either `template` or `template_file`                                              |
| comment               | GITHUB_COMMENT               | Comment text. If neither `comment` nor `GITHUB_COMMENT` provided, will read from `stdin`                                                                                                                    |
| delete-comment-regex  | GITHUB_DELETE_COMMENT_REGEX  | Regex to find previous comments to delete before creating the new comment. Supported for comment types `commit`, `pr-file`, `issue` and `pr`                                                                |


__NOTE__: The utility accepts the text of the comment from the command-line argument `comment`, from the ENV variable `GITHUB_COMMENT`, or from the standard input.
Command-line argument takes precedence over ENV var, and ENV var takes precedence over standard input.
Accepting comments from `stdin` allows using Unix pipes to send the output from another program as the input to the tool:

```sh
    cat comment.txt | github-commenter ...
```

```sh
    terraform plan 2>&1 | github-commenter -format "Output from `terraform plan`<br/>```{{.}}```"
```

__NOTE__: The utility supports [sprig functions](http://masterminds.github.io/sprig/) in `Go` templates, allowing to use string replacement and Regular Expressions in the `format` argument.

See [string functions](http://masterminds.github.io/sprig/strings.html) for more details.

For example:

```sh
GITHUB_COMMENT_FORMAT="Helm diff:<br><br><pre>{{regexReplaceAllLiteral `\\n` . `<br>` }}<pre>"
```




## Examples

The utility can be called directly or as a Docker container.

### Build the Go program locally

```sh
go get

CGO_ENABLED=0 go build -v -o "./dist/bin/github-commenter" *.go
```


### Run locally with ENV vars
[run_locally_with_env_vars.sh](examples/run_locally_with_env_vars.sh)

```sh
export GITHUB_TOKEN=XXXXXXXXXXXXXXXX
export GITHUB_OWNER=cloudposse
export GITHUB_REPO=github-commenter
export GITHUB_COMMENT_TYPE=pr
export GITHUB_PR_ISSUE_NUMBER=1
export GITHUB_COMMENT_FORMAT="My comment:<br/>{{.}}"
export GITHUB_COMMENT="+1 LGTM"

./dist/bin/github-commenter
```


### Run locally with command-line arguments
[run_locally_with_command_line_args.sh](examples/run_locally_with_command_line_args.sh)

```sh
./dist/bin/github-commenter \
        -token XXXXXXXXXXXXXXXX \
        -owner cloudposse \
        -repo github-commenter \
        -type pr \
        -number 1 \
        -format "My comment:<br/>{{.}}" \
        -comment "+1 LGTM"
```

### Build the Docker image
__NOTE__: it will download all `Go` dependencies and then build the program inside the container (see [`Dockerfile`](Dockerfile))


```sh
docker build --tag github-commenter  --no-cache=true .
```


### Run in a Docker container with ENV vars
[run_docker_with_env_vars.sh](examples/run_docker_with_env_vars.sh)

```sh
docker run -i --rm \
        -e GITHUB_TOKEN=XXXXXXXXXXXXXXXX \
        -e GITHUB_OWNER=cloudposse \
        -e GITHUB_REPO=github-commenter \
        -e GITHUB_COMMENT_TYPE=pr \
        -e GITHUB_PR_ISSUE_NUMBER=1 \
        -e GITHUB_COMMENT_FORMAT="My comment:<br/>{{.}}" \
        -e GITHUB_COMMENT="+1 LGTM" \
        github-commenter
```


### Run with Docker
Run `github-commenter` in a Docker container with local ENV vars propagated into the container's environment.
[run_docker_with_local_env_vars.sh](examples/run_docker_with_local_env_vars.sh)

```sh
export GITHUB_TOKEN=XXXXXXXXXXXXXXXX
export GITHUB_OWNER=cloudposse
export GITHUB_REPO=github-commenter
export GITHUB_COMMENT_TYPE=pr
export GITHUB_PR_ISSUE_NUMBER=1
export GITHUB_COMMENT_FORMAT="Helm diff:<br><br><pre>{{regexReplaceAllLiteral `\\n` . `<br>` }}<pre>"
export GITHUB_COMMENT="Helm diff comment"

docker run -i --rm \
        -e GITHUB_TOKEN \
        -e GITHUB_OWNER \
        -e GITHUB_REPO \
        -e GITHUB_COMMENT_TYPE \
        -e GITHUB_PR_ISSUE_NUMBER \
        -e GITHUB_COMMENT_FORMAT \
        -e GITHUB_COMMENT \
        github-commenter
```


### Run with Docker using Env File
Run the `github-commenter` in a Docker container with ENV vars declared in a file.
  [run_docker_with_env_vars_file.sh](examples/run_docker_with_env_vars_file.sh)

```sh
docker run -i --rm --env-file ./example.env github-commenter
```


### `delete-comment-regex` example 1
Delete all previous comments on Pull Request #2 that contain the string `test1` in the body of the comments and create a new PR comment

```sh
./dist/bin/github-commenter \
        -token XXXXXXXXXXXXXXXX \
        -owner cloudposse \
        -repo github-commenter \
        -type pr \
        -number 2 \
        -format "{{.}}" \
        -delete-comment-regex "test1" \
        -comment "New Pull Request comment"
```

### `delete-comment-regex` example 2
Delete all previous comments on Issue #1 that contain the string `test2` at the end of the comment's body and create a new Issue comment

```sh
./dist/bin/github-commenter \
        -token XXXXXXXXXXXXXXXX \
        -owner cloudposse \
        -repo github-commenter \
        -type issue \
        -number 1 \
        -format "{{.}}" \
        -delete-comment-regex "test2$" \
        -comment "New Issue comment"
```

### `delete-comment-regex` example 3
Delete all previous commit comments that contain the string `test3` in the body and create a new commit comment

```sh
./dist/bin/github-commenter \
        -token XXXXXXXXXXXXXXXX \
        -owner cloudposse \
        -repo github-commenter \
        -type commit \
        -sha xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
        -format "{{.}}" \
        -delete-comment-regex "test3" \
        -comment "New commit comment"
```


### `delete-comment-regex` example 4
Delete all previous comments on a Pull Request file `doc.txt` that contain the string `test4` in the body of the comments and create a new comment on the file

```sh
./dist/bin/github-commenter \
        -token XXXXXXXXXXXXXXXX \
        -owner cloudposse \
        -repo github-commenter \
        -type pr-file \
        -sha xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
        -number 2 \
        -file doc.txt \
        -position 1 \
        -format "{{.}}" \
        -delete-comment-regex "test4" \
        -comment "New comment on the PR file"
```





## Share the Love 

Like this project? Please give it a ★ on [our GitHub](https://github.com/cloudposse/github-commenter)! (it helps us **a lot**) 

Are you using this project or any of our other projects? Consider [leaving a testimonial][testimonial]. =)


## Related Projects

Check out these related projects.

- [github-status-updater](https://github.com/cloudposse/github-status-updater) - Command line utility for updating GitHub commit statuses and enabling required status checks for pull requests
- [slack-notifier](https://github.com/cloudposse/slack-notifier) - Command line utility to send messages with attachments to Slack channels via Incoming Webhooks



## Help

**Got a question?**

File a GitHub [issue](https://github.com/cloudposse/github-commenter/issues), send us an [email][email] or join our [Slack Community][slack].

[![README Commercial Support][readme_commercial_support_img]][readme_commercial_support_link]

## Commercial Support

Work directly with our team of DevOps experts via email, slack, and video conferencing. 

We provide [*commercial support*][commercial_support] for all of our [Open Source][github] projects. As a *Dedicated Support* customer, you have access to our team of subject matter experts at a fraction of the cost of a full-time engineer. 

[![E-Mail](https://img.shields.io/badge/email-hello@cloudposse.com-blue.svg)][email]

- **Questions.** We'll use a Shared Slack channel between your team and ours.
- **Troubleshooting.** We'll help you triage why things aren't working.
- **Code Reviews.** We'll review your Pull Requests and provide constructive feedback.
- **Bug Fixes.** We'll rapidly work to fix any bugs in our projects.
- **Build New Terraform Modules.** We'll [develop original modules][module_development] to provision infrastructure.
- **Cloud Architecture.** We'll assist with your cloud strategy and design.
- **Implementation.** We'll provide hands-on support to implement our reference architectures. 




## Slack Community

Join our [Open Source Community][slack] on Slack. It's **FREE** for everyone! Our "SweetOps" community is where you get to talk with others who share a similar vision for how to rollout and manage infrastructure. This is the best place to talk shop, ask questions, solicit feedback, and work together as a community to build totally *sweet* infrastructure.

## Newsletter

Signup for [our newsletter][newsletter] that covers everything on our technology radar.  Receive updates on what we're up to on GitHub as well as awesome new projects we discover. 

## Contributing

### Bug Reports & Feature Requests

Please use the [issue tracker](https://github.com/cloudposse/github-commenter/issues) to report any bugs or file feature requests.

### Developing

If you are interested in being a contributor and want to get involved in developing this project or [help out](https://cpco.io/help-out) with our other projects, we would love to hear from you! Shoot us an [email][email].

In general, PRs are welcome. We follow the typical "fork-and-pull" Git workflow.

 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull Request** so that we can review your changes

**NOTE:** Be sure to merge the latest changes from "upstream" before making a pull request!


## Copyright

Copyright © 2017-2019 [Cloud Posse, LLC](https://cpco.io/copyright)



## License 

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) 

See [LICENSE](LICENSE) for full details.

    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.









## Trademarks

All other trademarks referenced herein are the property of their respective owners.

## About

This project is maintained and funded by [Cloud Posse, LLC][website]. Like it? Please let us know by [leaving a testimonial][testimonial]!

[![Cloud Posse][logo]][website]

We're a [DevOps Professional Services][hire] company based in Los Angeles, CA. We ❤️  [Open Source Software][we_love_open_source].

We offer [paid support][commercial_support] on all of our projects.  

Check out [our other projects][github], [follow us on twitter][twitter], [apply for a job][jobs], or [hire us][hire] to help with your cloud strategy and implementation.



### Contributors

|  [![Erik Osterman][osterman_avatar]][osterman_homepage]<br/>[Erik Osterman][osterman_homepage] | [![Andriy Knysh][aknysh_avatar]][aknysh_homepage]<br/>[Andriy Knysh][aknysh_homepage] | [![Igor Rodionov][goruha_avatar]][goruha_homepage]<br/>[Igor Rodionov][goruha_homepage] |
|---|---|---|


  [osterman_homepage]: https://github.com/osterman
  [osterman_avatar]: http://s.gravatar.com/avatar/88c480d4f73b813904e00a5695a454cb?s=144


  [aknysh_homepage]: https://github.com/aknysh/
  [aknysh_avatar]: https://avatars0.githubusercontent.com/u/7356997?v=4&u=ed9ce1c9151d552d985bdf5546772e14ef7ab617&s=144


  [goruha_homepage]: https://github.com/goruha/
  [goruha_avatar]: http://s.gravatar.com/avatar/bc70834d32ed4517568a1feb0b9be7e2?s=144




[![README Footer][readme_footer_img]][readme_footer_link]
[![Beacon][beacon]][website]

  [logo]: https://cloudposse.com/logo-300x69.svg
  [docs]: https://cpco.io/docs
  [website]: https://cpco.io/homepage
  [github]: https://cpco.io/github
  [jobs]: https://cpco.io/jobs
  [hire]: https://cpco.io/hire
  [slack]: https://cpco.io/slack
  [linkedin]: https://cpco.io/linkedin
  [twitter]: https://cpco.io/twitter
  [testimonial]: https://cpco.io/leave-testimonial
  [newsletter]: https://cpco.io/newsletter
  [email]: https://cpco.io/email
  [commercial_support]: https://cpco.io/commercial-support
  [we_love_open_source]: https://cpco.io/we-love-open-source
  [module_development]: https://cpco.io/module-development
  [terraform_modules]: https://cpco.io/terraform-modules
  [readme_header_img]: https://cloudposse.com/readme/header/img?repo=cloudposse/github-commenter
  [readme_header_link]: https://cloudposse.com/readme/header/link?repo=cloudposse/github-commenter
  [readme_footer_img]: https://cloudposse.com/readme/footer/img?repo=cloudposse/github-commenter
  [readme_footer_link]: https://cloudposse.com/readme/footer/link?repo=cloudposse/github-commenter
  [readme_commercial_support_img]: https://cloudposse.com/readme/commercial-support/img?repo=cloudposse/github-commenter
  [readme_commercial_support_link]: https://cloudposse.com/readme/commercial-support/link?repo=cloudposse/github-commenter
  [share_twitter]: https://twitter.com/intent/tweet/?text=github-commenter&url=https://github.com/cloudposse/github-commenter
  [share_linkedin]: https://www.linkedin.com/shareArticle?mini=true&title=github-commenter&url=https://github.com/cloudposse/github-commenter
  [share_reddit]: https://reddit.com/submit/?url=https://github.com/cloudposse/github-commenter
  [share_facebook]: https://facebook.com/sharer/sharer.php?u=https://github.com/cloudposse/github-commenter
  [share_googleplus]: https://plus.google.com/share?url=https://github.com/cloudposse/github-commenter
  [share_email]: mailto:?subject=github-commenter&body=https://github.com/cloudposse/github-commenter
  [beacon]: https://ga-beacon.cloudposse.com/UA-76589703-4/cloudposse/github-commenter?pixel&cs=github&cm=readme&an=github-commenter
