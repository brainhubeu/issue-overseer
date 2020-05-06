package IssuesTriage

import (
	"../Types"
)

type IssuesTriage struct {
}

func InitIssuesTriage() *IssuesTriage {
	issuesTriage := &IssuesTriage{}
	return issuesTriage
}

func (issuesTriage IssuesTriage) TriageOneIssue(issue Types.Issue) int {
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
			return Types.IssueTypeEnum.OURS
		} else {
			if lastAuthorAssociation == "MEMBER" {
				return Types.IssueTypeEnum.ANSWERED
			} else {
				return Types.IssueTypeEnum.NOT_ANSWERED
			}
		}
	} else {
		if len(comments) == 0 {
			return Types.IssueTypeEnum.NOT_ANSWERED
		} else {
			j := len(comments) - 1
			for ; j >= 0; j-- {
				if comments[j].Node.Author.Login != "issuehunt-app" {
					break
				}
			}
			if comments[j].Node.AuthorAssociation == "MEMBER" {
				return Types.IssueTypeEnum.ANSWERED
			} else {
				return Types.IssueTypeEnum.NOT_ANSWERED
			}
		}
	}
}

func (issuesTriage IssuesTriage) TriageManyIssues(issues []Types.Issue) ([]Types.Issue, []Types.Issue, []Types.Issue) {
	ourIssues := []Types.Issue{}
	answeredIssues := []Types.Issue{}
	notAnsweredIssues := []Types.Issue{}
	for i := 0; i < len(issues); i++ {
		issue := issues[i]
		issueType := issuesTriage.TriageOneIssue(issue)
		if issueType == Types.IssueTypeEnum.OURS {
			ourIssues = append(ourIssues, issue)
		} else if issueType == Types.IssueTypeEnum.ANSWERED {
			answeredIssues = append(answeredIssues, issue)
		} else {
			notAnsweredIssues = append(notAnsweredIssues, issue)
		}
	}
	return ourIssues, answeredIssues, notAnsweredIssues
}
