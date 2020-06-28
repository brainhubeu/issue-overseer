<br/>
<h1 align="center">
  issue-overseer
</h1>

<p align="center">
  Bot adding labels to GitHub issues which allows easy finding of issues with a missing answer.
<p align="center">
  <strong>
    <a href="https://brainhub.eu/contact/">Hire us</a>
  </strong>
</p>

<div align="center">

  [![CircleCI](https://img.shields.io/circleci/project/github/brainhubeu/issue-overseer.svg)](https://circleci.com/gh/brainhubeu/issue-overseer)
  [![Last commit](https://img.shields.io/github/last-commit/brainhubeu/issue-overseer.svg)](https://github.com/brainhubeu/issue-overseer/commits/master)
  [![license](https://img.shields.io/badge/License-MIT-green)](https://github.com/brainhubeu/issue-overseer/blob/master/LICENSE.md)
  [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

  [![codecov](https://codecov.io/gh/brainhubeu/issue-overseer/branch/master/graph/badge.svg)](https://codecov.io/gh/brainhubeu/issue-overseer)
  [![Activity](https://img.shields.io/github/commit-activity/m/brainhubeu/issue-overseer.svg)](https://github.com/brainhubeu/issue-overseer/commits/master)
</div>

Let's assume `my-acme-org` is your GitHub organization name.

For each open issue (among the comments, it excludes the ones made by **issuehunt-bot**), it:
- puts "**answering: reported by my-acme-org**" label if the issue is created by any member of the my-acme-org organization with no comments by external contributors;
- otherwise, puts "**answering: answered**" label if the last comment is by a member of the organization;
- otherwise, puts "**answering: not answered**"
- removes the remaining answering labels because they are exclusive

## run

Regardless of the way, you choose, you need to export the `GITHUB_TOKEN` environmental variable:

in bash:
```
export GITHUB_TOKEN=my-gh-token
```

in fish:
```
export set GITHUB_TOKEN=my-gh-token
```

### dynamically with go
```
go run . my-acme-org
```

### compile and run an executable file
```
go build -o issue-overseer
./issue-overseer my-acme-org
```

### run with Docker
```
docker-compose up
```
