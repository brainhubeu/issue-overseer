package GithubOperator

import (
	"../Interfaces"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
	RegisterFailHandler(Fail)
	RunSpecs(t, "GithubOperator")
}

var _ = Describe("GithubOperator", func() {
	BeforeEach(func() {
		mockFindRepos = func() []string {
			Fail("mockFindRepos not implemented")
			return nil
		}
		mockFindLabels = func(repoName string) []Interfaces.Label {
			Fail("mockFindLabels not implemented")
			return nil
		}
		mockDeleteLabel = func(repoName string, labelName string) {
			Fail("mockDeleteLabel not implemented")
		}
		mockCreateLabel = func(repoName string, label Interfaces.Label) {
			Fail("mockCreateLabel not implemented")
		}
		mockRemoveLabel = func(issueUrl string, labelName string) {
			Fail("mockRemoveLabel not implemented")
		}
		mockAddLabel = func(issueUrl string, labelName string) {
			Fail("mockAddLabel not implemented")
		}
		mockFindIssues = func(repoName string) []Interfaces.Issue {
			Fail("mockFindIssues not implemented")
			return nil
		}
	})

	It("triages an empty list", func() {
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

	It("creates labels", func() {
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

		Expect(mockCreateLabelsParams).To(Equal([]interface{}{
			[]interface{}{"repo-1", Interfaces.Label{Name: "label-1", Color: "color-1"}},
			[]interface{}{"repo-1", Interfaces.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-1", Interfaces.Label{Name: "label-3", Color: "color-3"}},
			[]interface{}{"repo-2", Interfaces.Label{Name: "label-1", Color: "color-1"}},
			[]interface{}{"repo-2", Interfaces.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-2", Interfaces.Label{Name: "label-3", Color: "color-3"}},
			[]interface{}{"repo-3", Interfaces.Label{Name: "label-1", Color: "color-1"}},
			[]interface{}{"repo-3", Interfaces.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-3", Interfaces.Label{Name: "label-3", Color: "color-3"}},
		}))
	})

	It("does not create already created labels", func() {
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

		Expect(mockCreateLabelsParams).To(Equal([]interface{}{
			[]interface{}{"repo-1", Interfaces.Label{Name: "label-3", Color: "color-3"}},
			[]interface{}{"repo-2", Interfaces.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-2", Interfaces.Label{Name: "label-3", Color: "color-3"}},
		}))
	})

	It("deletes invalid labels", func() {
		mockDeleteLabelsParams := []interface{}{}
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
					Interfaces.Label{Name: "label-1", Color: "color-1-invalid"},
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
					Interfaces.Label{Name: "label-2", Color: "color-2-invalid"},
					Interfaces.Label{Name: "label-3", Color: "color-3-invalid"},
					Interfaces.Label{Name: "label-4", Color: "color-4"},
				}
			}
			return []Interfaces.Label{}
		}
		mockCreateLabel = func(repoName string, label Interfaces.Label) {
		}
		mockDeleteLabel = func(repoName string, labelName string) {
			mockDeleteLabelsParams = append(mockDeleteLabelsParams, []interface{}{repoName, labelName})
		}
		mockFindIssues = func(repoName string) []Interfaces.Issue {
			return []Interfaces.Issue{}
		}
		githubClient := MockGithubClient{}
		issuesTriage := MockIssuesTriage{}
		githubOperator := InitGithubOperator(githubClient, issuesTriage, answeringLabels, "label-1", "label-2", "label-3")

		githubOperator.UpdateRepos(repoNames)

		Expect(mockDeleteLabelsParams).To(Equal([]interface{}{
			[]interface{}{"repo-1", "label-1"},
			[]interface{}{"repo-3", "label-2"},
			[]interface{}{"repo-3", "label-3"},
		}))
	})

	It("adds missing labels", func() {
		mockAddLabelParams := []interface{}{}
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
			return []Interfaces.Issue{
					Interfaces.Issue{Url: "url-1"},
					Interfaces.Issue{Url: "url-2"},
				},
				[]Interfaces.Issue{
					Interfaces.Issue{Url: "url-3"},
					Interfaces.Issue{Url: "url-4"},
				},
				[]Interfaces.Issue{
					Interfaces.Issue{Url: "url-5"},
					Interfaces.Issue{Url: "url-6"},
				}
		}
		mockFindLabels = func(repoName string) []Interfaces.Label {
			return []Interfaces.Label{}
		}
		mockCreateLabel = func(repoName string, label Interfaces.Label) {
		}
		mockFindIssues = func(repoName string) []Interfaces.Issue {
			return []Interfaces.Issue{}
		}
		mockAddLabel = func(issueUrl string, labelName string) {
			mockAddLabelParams = append(mockAddLabelParams, []interface{}{issueUrl, labelName})
		}
		githubClient := MockGithubClient{}
		issuesTriage := MockIssuesTriage{}
		githubOperator := InitGithubOperator(githubClient, issuesTriage, answeringLabels, "by-ours", "answered", "not-answered")

		githubOperator.UpdateRepos(repoNames)

		// TODO investigate this behavior (the same params multiple times)
		Expect(mockAddLabelParams).To(Equal([]interface{}{
			[]interface{}{"url-1", "by-ours"},
			[]interface{}{"url-2", "by-ours"},
			[]interface{}{"url-3", "answered"},
			[]interface{}{"url-4", "answered"},
			[]interface{}{"url-5", "not-answered"},
			[]interface{}{"url-6", "not-answered"},
			[]interface{}{"url-1", "by-ours"},
			[]interface{}{"url-2", "by-ours"},
			[]interface{}{"url-3", "answered"},
			[]interface{}{"url-4", "answered"},
			[]interface{}{"url-5", "not-answered"},
			[]interface{}{"url-6", "not-answered"},
			[]interface{}{"url-1", "by-ours"},
			[]interface{}{"url-2", "by-ours"},
			[]interface{}{"url-3", "answered"},
			[]interface{}{"url-4", "answered"},
			[]interface{}{"url-5", "not-answered"},
			[]interface{}{"url-6", "not-answered"},
		}))
	})

	It("removes labels", func() {
		mockRemoveLabelParams := []interface{}{}
		repoNames := []string{
			"repo-1",
			"repo-2",
			"repo-3",
		}
		answeringLabels := []Interfaces.Label{
			Interfaces.Label{Name: "by-ours", Color: "color-1"},
			Interfaces.Label{Name: "answered", Color: "color-2"},
			Interfaces.Label{Name: "not-answered", Color: "color-3"},
		}

		mockTriageManyIssues = func(issues []Interfaces.Issue) ([]Interfaces.Issue, []Interfaces.Issue, []Interfaces.Issue) {
			return []Interfaces.Issue{
					Interfaces.Issue{Url: "url-1", Labels: []Interfaces.Label{Interfaces.Label{Name: "by-ours"}}},
					Interfaces.Issue{Url: "url-2", Labels: []Interfaces.Label{Interfaces.Label{Name: "not-answered"}}},
				},
				[]Interfaces.Issue{
					Interfaces.Issue{Url: "url-3", Labels: []Interfaces.Label{Interfaces.Label{Name: "answered"}}},
					Interfaces.Issue{Url: "url-4", Labels: []Interfaces.Label{Interfaces.Label{Name: "not-answered"}}},
				},
				[]Interfaces.Issue{
					Interfaces.Issue{Url: "url-5", Labels: []Interfaces.Label{Interfaces.Label{Name: "answered"}}},
					Interfaces.Issue{Url: "url-6", Labels: []Interfaces.Label{Interfaces.Label{Name: "not-answered"}}},
				}
		}
		mockFindLabels = func(repoName string) []Interfaces.Label {
			return []Interfaces.Label{}
		}
		mockCreateLabel = func(repoName string, label Interfaces.Label) {
		}
		mockFindIssues = func(repoName string) []Interfaces.Issue {
			return []Interfaces.Issue{}
		}
		mockAddLabel = func(issueUrl string, labelName string) {
		}
		mockRemoveLabel = func(issueUrl string, labelName string) {
			mockRemoveLabelParams = append(mockRemoveLabelParams, []interface{}{issueUrl, labelName})
		}
		githubClient := MockGithubClient{}
		issuesTriage := MockIssuesTriage{}
		githubOperator := InitGithubOperator(githubClient, issuesTriage, answeringLabels, "by-ours", "answered", "not-answered")

		githubOperator.UpdateRepos(repoNames)

		// TODO investigate this behavior (the same params multiple times)
		Expect(mockRemoveLabelParams).To(Equal([]interface{}{
			[]interface{}{"url-2", "not-answered"},
			[]interface{}{"url-4", "not-answered"},
			[]interface{}{"url-5", "answered"},
			[]interface{}{"url-2", "not-answered"},
			[]interface{}{"url-4", "not-answered"},
			[]interface{}{"url-5", "answered"},
			[]interface{}{"url-2", "not-answered"},
			[]interface{}{"url-4", "not-answered"},
			[]interface{}{"url-5", "answered"},
		}))
	})
})
