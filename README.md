# issue-label-bot
Bot adding labels to GitHub issues which allows easy finding of issues with a missing answer

[![CircleCI](https://img.shields.io/circleci/project/github/brainhubeu/issue-label-bot.svg)](https://circleci.com/gh/brainhubeu/issue-label-bot)
[![Last commit](https://img.shields.io/github/last-commit/brainhubeu/issue-label-bot.svg)](https://github.com/brainhubeu/issue-label-bot/commits/master)
[![license](https://img.shields.io/badge/License-MIT-green)](https://github.com/brainhubeu/issue-label-bot/blob/master/LICENSE.md)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

[![Activity](https://img.shields.io/github/commit-activity/m/brainhubeu/issue-label-bot.svg)](https://github.com/brainhubeu/issue-label-bot/commits/master)

For each open issue (among the comments, it excludes the ones made by **issuehunt-bot**), it:
- puts "**answering: reported by \_\_organization_name\_\_**" label if the issue is created by any member of the \_\_organization_name\_\_ organization with no comments by external contributors;
- otherwise, puts "**answering: answered**" label if the last comment is by a member of the organization;
- otherwise, puts "**answering: not answered**"
- removes the remaining answering labels because they are exclusive

## run
```
export set GITHUB_TOKEN=__your_token__
go run run.go brainhubeu
```
