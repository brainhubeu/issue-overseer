package githuboperator

import (
	"../interfaces"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

type Mockissuestriage struct{}

var mockTriageManyIssues func(issues []interfaces.Issue) ([]interfaces.Issue, []interfaces.Issue, []interfaces.Issue)

func (issuesTriage Mockissuestriage) TriageManyIssues(issues []interfaces.Issue) ([]interfaces.Issue, []interfaces.Issue, []interfaces.Issue) {
	return mockTriageManyIssues(issues)
}

type Mockgithubclient struct{}

var mockFindRepos func() []string
var mockFindLabels func(repoName string) []interfaces.Label
var mockDeleteLabel func(repoName string, labelName string)
var mockCreateLabel func(repoName string, label interfaces.Label)
var mockRemoveLabel func(issueUrl string, labelName string)
var mockAddLabel func(issueUrl string, labelName string)
var mockFindIssues func(repoName string) []interfaces.Issue

func (githubClient Mockgithubclient) FindRepos() []string {
	return mockFindRepos()
}
func (githubClient Mockgithubclient) FindLabels(repoName string) []interfaces.Label {
	return mockFindLabels(repoName)
}
func (githubClient Mockgithubclient) DeleteLabel(repoName string, labelName string) {
	mockDeleteLabel(repoName, labelName)
}
func (githubClient Mockgithubclient) CreateLabel(repoName string, label interfaces.Label) {
	mockCreateLabel(repoName, label)
}
func (githubClient Mockgithubclient) RemoveLabel(issueUrl string, labelName string) {
	mockRemoveLabel(issueUrl, labelName)
}
func (githubClient Mockgithubclient) AddLabel(issueUrl string, labelName string) {
	mockAddLabel(issueUrl, labelName)
}
func (githubClient Mockgithubclient) FindIssues(repoName string) []interfaces.Issue {
	return mockFindIssues(repoName)
}

func TestMain(m *testing.M) {
	status := m.Run()
	if status == 0 && testing.CoverMode() != "" {
		coverage := testing.Coverage()
		requiredCoverage := 1.0
		if coverage < requiredCoverage {
			log.Println("too low tests coverage:", coverage, ", should be at least", requiredCoverage)
			os.Exit(1)
		}
	}
}

func TestGithubOperator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "githuboperator")
}

var _ = Describe("githuboperator", func() {
	BeforeEach(func() {
		mockFindRepos = func() []string {
			Fail("mockFindRepos not implemented")
			return nil
		}
		mockFindLabels = func(repoName string) []interfaces.Label {
			Fail("mockFindLabels not implemented")
			return nil
		}
		mockDeleteLabel = func(repoName string, labelName string) {
			Fail("mockDeleteLabel not implemented")
		}
		mockCreateLabel = func(repoName string, label interfaces.Label) {
			Fail("mockCreateLabel not implemented")
		}
		mockRemoveLabel = func(issueUrl string, labelName string) {
			Fail("mockRemoveLabel not implemented")
		}
		mockAddLabel = func(issueUrl string, labelName string) {
			Fail("mockAddLabel not implemented")
		}
		mockFindIssues = func(repoName string) []interfaces.Issue {
			Fail("mockFindIssues not implemented")
			return nil
		}
	})

	It("triages an empty list", func() {
		repoNames := []string{}
		answeringLabels := []interfaces.Label{
			interfaces.Label{Name: "label-1", Color: "color-1"},
			interfaces.Label{Name: "label-2", Color: "color-2"},
			interfaces.Label{Name: "label-3", Color: "color-3"},
		}

		githubClient := Mockgithubclient{}
		issuesTriage := Mockissuestriage{}
		githubOperator := Initgithuboperator(githubClient, issuesTriage, answeringLabels, "label-1", "label-2", "label-3")
		githubOperator.UpdateRepos(repoNames)
	})

	It("creates labels", func() {
		mockCreateLabelsParams := []interface{}{}
		repoNames := []string{
			"repo-1",
			"repo-2",
			"repo-3",
		}
		answeringLabels := []interfaces.Label{
			interfaces.Label{Name: "label-1", Color: "color-1"},
			interfaces.Label{Name: "label-2", Color: "color-2"},
			interfaces.Label{Name: "label-3", Color: "color-3"},
		}

		mockTriageManyIssues = func(issues []interfaces.Issue) ([]interfaces.Issue, []interfaces.Issue, []interfaces.Issue) {
			return []interfaces.Issue{}, []interfaces.Issue{}, []interfaces.Issue{}
		}
		mockFindLabels = func(repoName string) []interfaces.Label {
			return []interfaces.Label{}
		}
		mockCreateLabel = func(repoName string, label interfaces.Label) {
			mockCreateLabelsParams = append(mockCreateLabelsParams, []interface{}{repoName, label})
		}
		mockFindIssues = func(repoName string) []interfaces.Issue {
			return []interfaces.Issue{}
		}
		githubClient := Mockgithubclient{}
		issuesTriage := Mockissuestriage{}
		githubOperator := Initgithuboperator(githubClient, issuesTriage, answeringLabels, "label-1", "label-2", "label-3")

		githubOperator.UpdateRepos(repoNames)

		Expect(mockCreateLabelsParams).To(Equal([]interface{}{
			[]interface{}{"repo-1", interfaces.Label{Name: "label-1", Color: "color-1"}},
			[]interface{}{"repo-1", interfaces.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-1", interfaces.Label{Name: "label-3", Color: "color-3"}},
			[]interface{}{"repo-2", interfaces.Label{Name: "label-1", Color: "color-1"}},
			[]interface{}{"repo-2", interfaces.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-2", interfaces.Label{Name: "label-3", Color: "color-3"}},
			[]interface{}{"repo-3", interfaces.Label{Name: "label-1", Color: "color-1"}},
			[]interface{}{"repo-3", interfaces.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-3", interfaces.Label{Name: "label-3", Color: "color-3"}},
		}))
	})

	It("does not create already created labels", func() {
		mockCreateLabelsParams := []interface{}{}
		repoNames := []string{
			"repo-1",
			"repo-2",
			"repo-3",
		}
		answeringLabels := []interfaces.Label{
			interfaces.Label{Name: "label-1", Color: "color-1"},
			interfaces.Label{Name: "label-2", Color: "color-2"},
			interfaces.Label{Name: "label-3", Color: "color-3"},
		}

		mockTriageManyIssues = func(issues []interfaces.Issue) ([]interfaces.Issue, []interfaces.Issue, []interfaces.Issue) {
			return []interfaces.Issue{}, []interfaces.Issue{}, []interfaces.Issue{}
		}
		mockFindLabels = func(repoName string) []interfaces.Label {
			if repoName == "repo-1" {
				return []interfaces.Label{
					interfaces.Label{Name: "label-1", Color: "color-1"},
					interfaces.Label{Name: "label-2", Color: "color-2"},
				}
			}
			if repoName == "repo-2" {
				return []interfaces.Label{
					interfaces.Label{Name: "label-1", Color: "color-1"},
				}
			}
			if repoName == "repo-3" {
				return []interfaces.Label{
					interfaces.Label{Name: "label-1", Color: "color-1"},
					interfaces.Label{Name: "label-2", Color: "color-2"},
					interfaces.Label{Name: "label-3", Color: "color-3"},
				}
			}
			return []interfaces.Label{}
		}
		mockCreateLabel = func(repoName string, label interfaces.Label) {
			mockCreateLabelsParams = append(mockCreateLabelsParams, []interface{}{repoName, label})
		}
		mockFindIssues = func(repoName string) []interfaces.Issue {
			return []interfaces.Issue{}
		}
		githubClient := Mockgithubclient{}
		issuesTriage := Mockissuestriage{}
		githubOperator := Initgithuboperator(githubClient, issuesTriage, answeringLabels, "label-1", "label-2", "label-3")

		githubOperator.UpdateRepos(repoNames)

		Expect(mockCreateLabelsParams).To(Equal([]interface{}{
			[]interface{}{"repo-1", interfaces.Label{Name: "label-3", Color: "color-3"}},
			[]interface{}{"repo-2", interfaces.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-2", interfaces.Label{Name: "label-3", Color: "color-3"}},
		}))
	})

	It("deletes invalid labels", func() {
		mockDeleteLabelsParams := []interface{}{}
		repoNames := []string{
			"repo-1",
			"repo-2",
			"repo-3",
		}
		answeringLabels := []interfaces.Label{
			interfaces.Label{Name: "label-1", Color: "color-1"},
			interfaces.Label{Name: "label-2", Color: "color-2"},
			interfaces.Label{Name: "label-3", Color: "color-3"},
		}

		mockTriageManyIssues = func(issues []interfaces.Issue) ([]interfaces.Issue, []interfaces.Issue, []interfaces.Issue) {
			return []interfaces.Issue{}, []interfaces.Issue{}, []interfaces.Issue{}
		}
		mockFindLabels = func(repoName string) []interfaces.Label {
			if repoName == "repo-1" {
				return []interfaces.Label{
					interfaces.Label{Name: "label-1", Color: "color-1-invalid"},
					interfaces.Label{Name: "label-2", Color: "color-2"},
				}
			}
			if repoName == "repo-2" {
				return []interfaces.Label{
					interfaces.Label{Name: "label-1", Color: "color-1"},
				}
			}
			if repoName == "repo-3" {
				return []interfaces.Label{
					interfaces.Label{Name: "label-1", Color: "color-1"},
					interfaces.Label{Name: "label-2", Color: "color-2-invalid"},
					interfaces.Label{Name: "label-3", Color: "color-3-invalid"},
					interfaces.Label{Name: "label-4", Color: "color-4"},
				}
			}
			return []interfaces.Label{}
		}
		mockCreateLabel = func(repoName string, label interfaces.Label) {
		}
		mockDeleteLabel = func(repoName string, labelName string) {
			mockDeleteLabelsParams = append(mockDeleteLabelsParams, []interface{}{repoName, labelName})
		}
		mockFindIssues = func(repoName string) []interfaces.Issue {
			return []interfaces.Issue{}
		}
		githubClient := Mockgithubclient{}
		issuesTriage := Mockissuestriage{}
		githubOperator := Initgithuboperator(githubClient, issuesTriage, answeringLabels, "label-1", "label-2", "label-3")

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
		answeringLabels := []interfaces.Label{
			interfaces.Label{Name: "label-1", Color: "color-1"},
			interfaces.Label{Name: "label-2", Color: "color-2"},
			interfaces.Label{Name: "label-3", Color: "color-3"},
		}

		mockTriageManyIssues = func(issues []interfaces.Issue) ([]interfaces.Issue, []interfaces.Issue, []interfaces.Issue) {
			return []interfaces.Issue{
					interfaces.Issue{Url: "url-1"},
					interfaces.Issue{Url: "url-2"},
				},
				[]interfaces.Issue{
					interfaces.Issue{Url: "url-3"},
					interfaces.Issue{Url: "url-4"},
				},
				[]interfaces.Issue{
					interfaces.Issue{Url: "url-5"},
					interfaces.Issue{Url: "url-6"},
				}
		}
		mockFindLabels = func(repoName string) []interfaces.Label {
			return []interfaces.Label{}
		}
		mockCreateLabel = func(repoName string, label interfaces.Label) {
		}
		mockFindIssues = func(repoName string) []interfaces.Issue {
			return []interfaces.Issue{}
		}
		mockAddLabel = func(issueUrl string, labelName string) {
			mockAddLabelParams = append(mockAddLabelParams, []interface{}{issueUrl, labelName})
		}
		githubClient := Mockgithubclient{}
		issuesTriage := Mockissuestriage{}
		githubOperator := Initgithuboperator(githubClient, issuesTriage, answeringLabels, "by-ours", "answered", "not-answered")

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
		answeringLabels := []interfaces.Label{
			interfaces.Label{Name: "by-ours", Color: "color-1"},
			interfaces.Label{Name: "answered", Color: "color-2"},
			interfaces.Label{Name: "not-answered", Color: "color-3"},
		}

		mockTriageManyIssues = func(issues []interfaces.Issue) ([]interfaces.Issue, []interfaces.Issue, []interfaces.Issue) {
			return []interfaces.Issue{
					interfaces.Issue{Url: "url-1", Labels: []interfaces.Label{interfaces.Label{Name: "by-ours"}}},
					interfaces.Issue{Url: "url-2", Labels: []interfaces.Label{interfaces.Label{Name: "not-answered"}}},
				},
				[]interfaces.Issue{
					interfaces.Issue{Url: "url-3", Labels: []interfaces.Label{interfaces.Label{Name: "answered"}}},
					interfaces.Issue{Url: "url-4", Labels: []interfaces.Label{interfaces.Label{Name: "not-answered"}}},
				},
				[]interfaces.Issue{
					interfaces.Issue{Url: "url-5", Labels: []interfaces.Label{interfaces.Label{Name: "answered"}}},
					interfaces.Issue{Url: "url-6", Labels: []interfaces.Label{interfaces.Label{Name: "not-answered"}}},
				}
		}
		mockFindLabels = func(repoName string) []interfaces.Label {
			return []interfaces.Label{}
		}
		mockCreateLabel = func(repoName string, label interfaces.Label) {
		}
		mockFindIssues = func(repoName string) []interfaces.Issue {
			return []interfaces.Issue{}
		}
		mockAddLabel = func(issueUrl string, labelName string) {
		}
		mockRemoveLabel = func(issueUrl string, labelName string) {
			mockRemoveLabelParams = append(mockRemoveLabelParams, []interface{}{issueUrl, labelName})
		}
		githubClient := Mockgithubclient{}
		issuesTriage := Mockissuestriage{}
		githubOperator := Initgithuboperator(githubClient, issuesTriage, answeringLabels, "by-ours", "answered", "not-answered")

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
