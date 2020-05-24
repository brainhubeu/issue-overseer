package GithubOperator

import (
	"../Interfaces"
	"github.com/stretchr/testify/assert"
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

func TestGithubOperator(t *testing.T) {
	t.Run("triages an empty list", func(t *testing.T) {
		repoNames := []string{}
		answeringLabels := []Interfaces.Label{
			Interfaces.Label{Name: "label-1", Color: "color-1"},
			Interfaces.Label{Name: "label-2", Color: "color-2"},
			Interfaces.Label{Name: "label-3", Color: "color-3"},
		}

		githubClient := MockGithubClient{}
		issuesTriage := MockIssuesTriage{}
		githubOperator := InitGithubOperator(githubClient, issuesTriage, answeringLabels, "label-1", "label-2", "label-3")
		githubOperator.UpdateRepos(repoNames)
	})

	t.Run("creates labels", func(t *testing.T) {
		mockCreateLabelsParams := []interface{}{}
		repoNames := []string{
			"repo-1",
			"repo-2",
			"repo-3",
		}
		answeringLabels := []Interfaces.Label{
			Interfaces.Label{Name: "label-1", Color: "color-1"},
			Interfaces.Label{Name: "label-2", Color: "color-2"},
			Interfaces.Label{Name: "label-3", Color: "color-3"},
		}

		mockTriageManyIssues = func(issues []Interfaces.Issue) ([]Interfaces.Issue, []Interfaces.Issue, []Interfaces.Issue) {
			return []Interfaces.Issue{}, []Interfaces.Issue{}, []Interfaces.Issue{}
		}
		mockFindLabels = func(repoName string) []Interfaces.Label {
			return []Interfaces.Label{}
		}
		mockCreateLabel = func(repoName string, label Interfaces.Label) {
			mockCreateLabelsParams = append(mockCreateLabelsParams, []interface{}{repoName, label})
		}
		mockFindIssues = func(repoName string) []Interfaces.Issue {
			return []Interfaces.Issue{}
		}

		githubClient := MockGithubClient{}
		issuesTriage := MockIssuesTriage{}
		githubOperator := InitGithubOperator(githubClient, issuesTriage, answeringLabels, "label-1", "label-2", "label-3")
		githubOperator.UpdateRepos(repoNames)

		assert.Equal(t, mockCreateLabelsParams, []interface{}{
			[]interface{}{"repo-1", Interfaces.Label{Name: "label-1", Color: "color-1"}},
			[]interface{}{"repo-1", Interfaces.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-1", Interfaces.Label{Name: "label-3", Color: "color-3"}},
			[]interface{}{"repo-2", Interfaces.Label{Name: "label-1", Color: "color-1"}},
			[]interface{}{"repo-2", Interfaces.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-2", Interfaces.Label{Name: "label-3", Color: "color-3"}},
			[]interface{}{"repo-3", Interfaces.Label{Name: "label-1", Color: "color-1"}},
			[]interface{}{"repo-3", Interfaces.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-3", Interfaces.Label{Name: "label-3", Color: "color-3"}},
		})
	})

	t.Run("does not create already created labels", func(t *testing.T) {
		mockCreateLabelsParams := []interface{}{}
		repoNames := []string{
			"repo-1",
			"repo-2",
			"repo-3",
		}
		answeringLabels := []Interfaces.Label{
			Interfaces.Label{Name: "label-1", Color: "color-1"},
			Interfaces.Label{Name: "label-2", Color: "color-2"},
			Interfaces.Label{Name: "label-3", Color: "color-3"},
		}

		mockTriageManyIssues = func(issues []Interfaces.Issue) ([]Interfaces.Issue, []Interfaces.Issue, []Interfaces.Issue) {
			return []Interfaces.Issue{}, []Interfaces.Issue{}, []Interfaces.Issue{}
		}
		mockFindLabels = func(repoName string) []Interfaces.Label {
			if repoName == "repo-1" {
				return []Interfaces.Label{
					Interfaces.Label{Name: "label-1", Color: "color-1"},
					Interfaces.Label{Name: "label-2", Color: "color-2"},
				}
			}
			if repoName == "repo-2" {
				return []Interfaces.Label{
					Interfaces.Label{Name: "label-1", Color: "color-1"},
				}
			}
			if repoName == "repo-3" {
				return []Interfaces.Label{
					Interfaces.Label{Name: "label-1", Color: "color-1"},
					Interfaces.Label{Name: "label-2", Color: "color-2"},
					Interfaces.Label{Name: "label-3", Color: "color-3"},
				}
			}
			return []Interfaces.Label{}
		}
		mockCreateLabel = func(repoName string, label Interfaces.Label) {
			mockCreateLabelsParams = append(mockCreateLabelsParams, []interface{}{repoName, label})
		}
		mockFindIssues = func(repoName string) []Interfaces.Issue {
			return []Interfaces.Issue{}
		}

		githubClient := MockGithubClient{}
		issuesTriage := MockIssuesTriage{}
		githubOperator := InitGithubOperator(githubClient, issuesTriage, answeringLabels, "label-1", "label-2", "label-3")
		githubOperator.UpdateRepos(repoNames)

		assert.Equal(t, mockCreateLabelsParams, []interface{}{
			[]interface{}{"repo-1", Interfaces.Label{Name: "label-3", Color: "color-3"}},
			[]interface{}{"repo-2", Interfaces.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-2", Interfaces.Label{Name: "label-3", Color: "color-3"}},
		})
	})
}
