package IssuesTriage

import (
	"../Types"
)

type IssuesTriage struct {
}

func InitIssuesTriage() *IssuesTriage{
	issuesTriage := &IssuesTriage{}
	return issuesTriage
}

func (issuesTriage IssuesTriage) DoIssuesTriage(issues []Types.Issue) ([]Types.Issue, []Types.Issue, []Types.Issue) {
	ourIssues := []Types.Issue{}
	answeredIssues := []Types.Issue{}
	notAnsweredIssues := []Types.Issue{}
	for i := 0; i < len(issues); i++ {
		issue := issues[i]
		comments := issue.Comments.Edges
		if issue.AuthorAssociation == "MEMBER" {
			j := len(comments) - 1
			lastAuthorAssociation := ""
			for ; j >= 0; j-- {
				comment := comments[j].Node
				if comment.Author.Login != "issuehunt-app" && lastAuthorAssociation == "" {
					lastAuthorAssociation = comment.AuthorAssociation
				}
				if comment.Author.Login != "issuehunt-app" && comment.AuthorAssociation != "MEMBER" {
					break
				}
			}
			if j == -1 {
				ourIssues = append(ourIssues, issue)
			} else {
				if lastAuthorAssociation == "MEMBER" {
					answeredIssues = append(answeredIssues, issue)
				} else {
					notAnsweredIssues = append(notAnsweredIssues, issue)
				}
			}
		} else {
			if len(comments) == 0 {
				notAnsweredIssues = append(notAnsweredIssues, issue)
			} else {
				j := len(comments) - 1
				for ; j >= 0; j-- {
					if comments[j].Node.Author.Login != "issuehunt-app" {
						break
					}
				}
				if comments[j].Node.AuthorAssociation == "MEMBER" {
					answeredIssues = append(answeredIssues, issue)
				} else {
					notAnsweredIssues = append(notAnsweredIssues, issue)
				}
			}
		}
	}
	return ourIssues, answeredIssues, notAnsweredIssues
}
