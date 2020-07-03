package issuestriage

import (
	"github.com/brainhubeu/issue-overseer/githubstructures"
)

type issuestriage struct {
}

func New() *issuestriage {
	issuesTriage := &issuestriage{}
	return issuesTriage
}

func (issuesTriage issuestriage) TriageOneIssueByAnswering(issue githubstructures.Issue) int {
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
			return githubstructures.IssueAnsweringTypeEnum.OURS
		} else {
			if lastAuthorAssociation == "MEMBER" {
				return githubstructures.IssueAnsweringTypeEnum.ANSWERED
			} else {
				return githubstructures.IssueAnsweringTypeEnum.NOT_ANSWERED
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
			return githubstructures.IssueAnsweringTypeEnum.NOT_ANSWERED
		} else if comments[j].AuthorAssociation == "MEMBER" {
			return githubstructures.IssueAnsweringTypeEnum.ANSWERED
		} else {
			return githubstructures.IssueAnsweringTypeEnum.NOT_ANSWERED
		}
	}
}

func (issuesTriage issuestriage) GroupByAnswering(issues []githubstructures.Issue) ([]githubstructures.Issue, []githubstructures.Issue, []githubstructures.Issue) {
	ourIssues := []githubstructures.Issue{}
	answeredIssues := []githubstructures.Issue{}
	notAnsweredIssues := []githubstructures.Issue{}
	for i := 0; i < len(issues); i++ {
		issue := issues[i]
		switch issueType := issuesTriage.TriageOneIssueByAnswering(issue); issueType {
		case githubstructures.IssueAnsweringTypeEnum.OURS:
			ourIssues = append(ourIssues, issue)
		case githubstructures.IssueAnsweringTypeEnum.ANSWERED:
			answeredIssues = append(answeredIssues, issue)
		default:
			notAnsweredIssues = append(notAnsweredIssues, issue)
		}
	}
	return ourIssues, answeredIssues, notAnsweredIssues
}
