# issue-label-bot
Bot adding labels to GitHub issues which allows easy finding of issues with a missing answer

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
