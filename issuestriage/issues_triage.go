package issuestriage

import (
	"github.com/brainhubeu/issue-overseer/interfaces"
)

type issuestriage struct {
}

func Initissuestriage() *issuestriage {
	issuesTriage := &issuestriage{}
	return issuesTriage
}

func (issuesTriage issuestriage) TriageOneIssue(issue interfaces.Issue) int {
	comments := issue.Comments
	if issue.AuthorAssociation == "MEMBER" {
		j := len(comments) - 1
		lastAuthorAssociation := ""
		for ; j >= 0; j-- {
			comment := comments[j]
			if comment.AuthorLogin != "issuehunt-app" && lastAuthorAssociation == "" {
				lastAuthorAssociation = comment.AuthorAssociation
			}
			if comment.AuthorLogin != "issuehunt-app" && comment.AuthorAssociation != "MEMBER" {
				break
			}
		}
		if j == -1 {
			return interfaces.IssueTypeEnum.OURS
		} else {
			if lastAuthorAssociation == "MEMBER" {
				return interfaces.IssueTypeEnum.ANSWERED
			} else {
				return interfaces.IssueTypeEnum.NOT_ANSWERED
			}
		}
	} else {
		j := len(comments) - 1
		for ; j >= 0; j-- {
			if comments[j].AuthorLogin != "issuehunt-app" {
				break
			}
		}
		if j == -1 {
			return interfaces.IssueTypeEnum.NOT_ANSWERED
		} else if comments[j].AuthorAssociation == "MEMBER" {
			return interfaces.IssueTypeEnum.ANSWERED
		} else {
			return interfaces.IssueTypeEnum.NOT_ANSWERED
		}
	}
}

func (issuesTriage issuestriage) TriageManyIssues(issues []interfaces.Issue) ([]interfaces.Issue, []interfaces.Issue, []interfaces.Issue) {
	ourIssues := []interfaces.Issue{}
	answeredIssues := []interfaces.Issue{}
	notAnsweredIssues := []interfaces.Issue{}
	for i := 0; i < len(issues); i++ {
		issue := issues[i]
		issueType := issuesTriage.TriageOneIssue(issue)
		if issueType == interfaces.IssueTypeEnum.OURS {
			ourIssues = append(ourIssues, issue)
		} else if issueType == interfaces.IssueTypeEnum.ANSWERED {
			answeredIssues = append(answeredIssues, issue)
		} else {
			notAnsweredIssues = append(notAnsweredIssues, issue)
		}
	}
	return ourIssues, answeredIssues, notAnsweredIssues
}
