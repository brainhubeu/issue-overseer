package GithubOperator

import (
	"../Interfaces"
	"testing"
)

type MockIssuesTriage struct{}

var mockTriageManyIssues func(issues []Interfaces.Issue) ([]Interfaces.Issue, []Interfaces.Issue, []Interfaces.Issue)

func (issuesTriage MockIssuesTriage) TriageManyIssues(issues []Interfaces.Issue) ([]Interfaces.Issue, []Interfaces.Issue, []Interfaces.Issue) {
	return mockTriageManyIssues(issues)
}

type MockGithubClient struct{}

var mockFindRepos func() []string
var mockFindLabels func(repoName string) []Interfaces.Label
var mockDeleteLabel func(repoName string, labelName string)
var mockCreateLabel func(repoName string, label Interfaces.Label)
var mockRemoveLabel func(issueUrl string, labelName string)
var mockAddLabel func(issueUrl string, labelName string)
var mockFindIssues func(repoName string) []Interfaces.Issue

func (githubClient MockGithubClient) FindRepos() []string {
	return mockFindRepos()
}
func (githubClient MockGithubClient) FindLabels(repoName string) []Interfaces.Label {
	return mockFindLabels(repoName)
}
func (githubClient MockGithubClient) DeleteLabel(repoName string, labelName string) {
	mockDeleteLabel(repoName, labelName)
}
func (githubClient MockGithubClient) CreateLabel(repoName string, label Interfaces.Label) {
	mockCreateLabel(repoName, label)
}
func (githubClient MockGithubClient) RemoveLabel(issueUrl string, labelName string) {
	mockRemoveLabel(issueUrl, labelName)
}
func (githubClient MockGithubClient) AddLabel(issueUrl string, labelName string) {
	mockAddLabel(issueUrl, labelName)
}
func (githubClient MockGithubClient) FindIssues(repoName string) []Interfaces.Issue {
	return mockFindIssues(repoName)
}

func TestTriageManyIssues(t *testing.T) {
	t.Run("triages an empty list", func(t *testing.T) {
		repoNames := []string{}
		OUR_LABEL_TEXT := "answering: reported by brainhubeu"
		const ANSWERED_LABEL_TEXT = "answering: answered"
		const NOT_ANSWERED_LABEL_TEXT = "answering: not answered"
		answeringLabels := []Interfaces.Label{
			Interfaces.Label{Name: OUR_LABEL_TEXT, Color: "a0a000"},
			Interfaces.Label{Name: ANSWERED_LABEL_TEXT, Color: "00a000"},
			Interfaces.Label{Name: NOT_ANSWERED_LABEL_TEXT, Color: "a00000"},
		}

		mockTriageManyIssues = func(issues []Interfaces.Issue) ([]Interfaces.Issue, []Interfaces.Issue, []Interfaces.Issue) {
			return []Interfaces.Issue{}, []Interfaces.Issue{}, []Interfaces.Issue{}
		}
		mockFindRepos = func() []string {
			return []string{}
		}
		mockFindLabels = func(repoName string) []Interfaces.Label {
			return []Interfaces.Label{}
		}
		mockDeleteLabel = func(repoName string, labelName string) {
		}
		mockCreateLabel = func(repoName string, label Interfaces.Label) {
		}
		mockRemoveLabel = func(issueUrl string, labelName string) {
		}
		mockAddLabel = func(issueUrl string, labelName string) {
		}
		mockFindIssues = func(repoName string) []Interfaces.Issue {
			return []Interfaces.Issue{}
		}

		githubClient := MockGithubClient{}
		issuesTriage := MockIssuesTriage{}
		githubOperator := InitGithubOperator(githubClient, issuesTriage, answeringLabels, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT)
		githubOperator.UpdateRepos(repoNames)
	})
}
