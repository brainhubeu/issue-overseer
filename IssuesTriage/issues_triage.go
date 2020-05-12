package IssuesTriage

import (
	"../Interfaces"
)

type IssuesTriage struct {
}

func InitIssuesTriage() *IssuesTriage {
	issuesTriage := &IssuesTriage{}
	return issuesTriage
}

func (issuesTriage IssuesTriage) TriageOneIssue(issue Interfaces.Issue) int {
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
			return Interfaces.IssueTypeEnum.OURS
		} else {
			if lastAuthorAssociation == "MEMBER" {
				return Interfaces.IssueTypeEnum.ANSWERED
			} else {
				return Interfaces.IssueTypeEnum.NOT_ANSWERED
			}
		}
	} else {
		j := len(comments) - 1
		for ; j >= 0; j-- {
			if comments[j].Node.Author.Login != "issuehunt-app" {
				break
			}
		}
		if j == -1 {
			return Interfaces.IssueTypeEnum.NOT_ANSWERED
		} else if comments[j].Node.AuthorAssociation == "MEMBER" {
			return Interfaces.IssueTypeEnum.ANSWERED
		} else {
			return Interfaces.IssueTypeEnum.NOT_ANSWERED
		}
	}
}

func (issuesTriage IssuesTriage) TriageManyIssues(issues []Interfaces.Issue) ([]Interfaces.Issue, []Interfaces.Issue, []Interfaces.Issue) {
	ourIssues := []Interfaces.Issue{}
	answeredIssues := []Interfaces.Issue{}
	notAnsweredIssues := []Interfaces.Issue{}
	for i := 0; i < len(issues); i++ {
		issue := issues[i]
		issueType := issuesTriage.TriageOneIssue(issue)
		if issueType == Interfaces.IssueTypeEnum.OURS {
			ourIssues = append(ourIssues, issue)
		} else if issueType == Interfaces.IssueTypeEnum.ANSWERED {
			answeredIssues = append(answeredIssues, issue)
		} else {
			notAnsweredIssues = append(notAnsweredIssues, issue)
		}
	}
	return ourIssues, answeredIssues, notAnsweredIssues
}
