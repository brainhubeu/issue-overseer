package githuboperator

import (
	"github.com/brainhubeu/issue-overseer/githubstructures"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

type Mockissuestriage struct{}

var mockGroupByAnswering func(issues []githubstructures.Issue) ([]githubstructures.Issue, []githubstructures.Issue, []githubstructures.Issue)

func (issuesTriage Mockissuestriage) GroupByAnswering(issues []githubstructures.Issue) ([]githubstructures.Issue, []githubstructures.Issue, []githubstructures.Issue) {
	return mockGroupByAnswering(issues)
}

type Mockgithubclient struct{}

var mockFindRepos func() []string
var mockFindLabels func(repoName string) []githubstructures.Label
var mockDeleteLabel func(repoName string, labelName string)
var mockCreateLabel func(repoName string, label githubstructures.Label)
var mockRemoveLabel func(issueUrl string, labelName string)
var mockAddLabel func(issueUrl string, labelName string)
var mockFindIssues func(repoName string) []githubstructures.Issue

func (githubClient Mockgithubclient) FindRepos() []string {
	return mockFindRepos()
}
func (githubClient Mockgithubclient) FindLabels(repoName string) []githubstructures.Label {
	return mockFindLabels(repoName)
}
func (githubClient Mockgithubclient) DeleteLabel(repoName string, labelName string) {
	mockDeleteLabel(repoName, labelName)
}
func (githubClient Mockgithubclient) CreateLabel(repoName string, label githubstructures.Label) {
	mockCreateLabel(repoName, label)
}
func (githubClient Mockgithubclient) RemoveLabel(issueUrl string, labelName string) {
	mockRemoveLabel(issueUrl, labelName)
}
func (githubClient Mockgithubclient) AddLabel(issueUrl string, labelName string) {
	mockAddLabel(issueUrl, labelName)
}
func (githubClient Mockgithubclient) FindIssues(repoName string) []githubstructures.Issue {
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
		mockFindLabels = func(repoName string) []githubstructures.Label {
			Fail("mockFindLabels not implemented")
			return nil
		}
		mockDeleteLabel = func(repoName string, labelName string) {
			Fail("mockDeleteLabel not implemented")
		}
		mockCreateLabel = func(repoName string, label githubstructures.Label) {
			Fail("mockCreateLabel not implemented")
		}
		mockRemoveLabel = func(issueUrl string, labelName string) {
			Fail("mockRemoveLabel not implemented")
		}
		mockAddLabel = func(issueUrl string, labelName string) {
			Fail("mockAddLabel not implemented")
		}
		mockFindIssues = func(repoName string) []githubstructures.Issue {
			Fail("mockFindIssues not implemented")
			return nil
		}
	})

	It("triages an empty list", func() {
		repoNames := []string{}
		answeringLabels := []githubstructures.Label{
			githubstructures.Label{Name: "label-1", Color: "color-1"},
			githubstructures.Label{Name: "label-2", Color: "color-2"},
			githubstructures.Label{Name: "label-3", Color: "color-3"},
		}

		githubClient := Mockgithubclient{}
		issuesTriage := Mockissuestriage{}
		githubOperator := New(githubClient, issuesTriage, answeringLabels, "label-1", "label-2", "label-3", answeringLabels)
		githubOperator.UpdateRepos(repoNames)
	})

	It("creates labels", func() {
		mockCreateLabelsParams := []interface{}{}
		repoNames := []string{
			"repo-1",
			"repo-2",
			"repo-3",
		}
		answeringLabels := []githubstructures.Label{
			githubstructures.Label{Name: "label-1", Color: "color-1"},
			githubstructures.Label{Name: "label-2", Color: "color-2"},
			githubstructures.Label{Name: "label-3", Color: "color-3"},
		}

		mockGroupByAnswering = func(issues []githubstructures.Issue) ([]githubstructures.Issue, []githubstructures.Issue, []githubstructures.Issue) {
			return []githubstructures.Issue{}, []githubstructures.Issue{}, []githubstructures.Issue{}
		}
		mockFindLabels = func(repoName string) []githubstructures.Label {
			return []githubstructures.Label{}
		}
		mockCreateLabel = func(repoName string, label githubstructures.Label) {
			mockCreateLabelsParams = append(mockCreateLabelsParams, []interface{}{repoName, label})
		}
		mockFindIssues = func(repoName string) []githubstructures.Issue {
			return []githubstructures.Issue{}
		}
		githubClient := Mockgithubclient{}
		issuesTriage := Mockissuestriage{}
		githubOperator := New(githubClient, issuesTriage, answeringLabels, "label-1", "label-2", "label-3", answeringLabels)

		githubOperator.UpdateRepos(repoNames)

		Expect(mockCreateLabelsParams).To(Equal([]interface{}{
			[]interface{}{"repo-1", githubstructures.Label{Name: "label-1", Color: "color-1"}},
			[]interface{}{"repo-1", githubstructures.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-1", githubstructures.Label{Name: "label-3", Color: "color-3"}},
			[]interface{}{"repo-2", githubstructures.Label{Name: "label-1", Color: "color-1"}},
			[]interface{}{"repo-2", githubstructures.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-2", githubstructures.Label{Name: "label-3", Color: "color-3"}},
			[]interface{}{"repo-3", githubstructures.Label{Name: "label-1", Color: "color-1"}},
			[]interface{}{"repo-3", githubstructures.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-3", githubstructures.Label{Name: "label-3", Color: "color-3"}},
		}))
	})

	It("does not create already created labels", func() {
		mockCreateLabelsParams := []interface{}{}
		repoNames := []string{
			"repo-1",
			"repo-2",
			"repo-3",
		}
		answeringLabels := []githubstructures.Label{
			githubstructures.Label{Name: "label-1", Color: "color-1"},
			githubstructures.Label{Name: "label-2", Color: "color-2"},
			githubstructures.Label{Name: "label-3", Color: "color-3"},
		}

		mockGroupByAnswering = func(issues []githubstructures.Issue) ([]githubstructures.Issue, []githubstructures.Issue, []githubstructures.Issue) {
			return []githubstructures.Issue{}, []githubstructures.Issue{}, []githubstructures.Issue{}
		}
		mockFindLabels = func(repoName string) []githubstructures.Label {
			if repoName == "repo-1" {
				return []githubstructures.Label{
					githubstructures.Label{Name: "label-1", Color: "color-1"},
					githubstructures.Label{Name: "label-2", Color: "color-2"},
				}
			}
			if repoName == "repo-2" {
				return []githubstructures.Label{
					githubstructures.Label{Name: "label-1", Color: "color-1"},
				}
			}
			if repoName == "repo-3" {
				return []githubstructures.Label{
					githubstructures.Label{Name: "label-1", Color: "color-1"},
					githubstructures.Label{Name: "label-2", Color: "color-2"},
					githubstructures.Label{Name: "label-3", Color: "color-3"},
				}
			}
			return []githubstructures.Label{}
		}
		mockCreateLabel = func(repoName string, label githubstructures.Label) {
			mockCreateLabelsParams = append(mockCreateLabelsParams, []interface{}{repoName, label})
		}
		mockFindIssues = func(repoName string) []githubstructures.Issue {
			return []githubstructures.Issue{}
		}
		githubClient := Mockgithubclient{}
		issuesTriage := Mockissuestriage{}
		githubOperator := New(githubClient, issuesTriage, answeringLabels, "label-1", "label-2", "label-3", answeringLabels)

		githubOperator.UpdateRepos(repoNames)

		Expect(mockCreateLabelsParams).To(Equal([]interface{}{
			[]interface{}{"repo-1", githubstructures.Label{Name: "label-3", Color: "color-3"}},
			[]interface{}{"repo-2", githubstructures.Label{Name: "label-2", Color: "color-2"}},
			[]interface{}{"repo-2", githubstructures.Label{Name: "label-3", Color: "color-3"}},
		}))
	})

	It("deletes invalid labels", func() {
		mockDeleteLabelsParams := []interface{}{}
		repoNames := []string{
			"repo-1",
			"repo-2",
			"repo-3",
		}
		answeringLabels := []githubstructures.Label{
			githubstructures.Label{Name: "label-1", Color: "color-1"},
			githubstructures.Label{Name: "label-2", Color: "color-2"},
			githubstructures.Label{Name: "label-3", Color: "color-3"},
		}

		mockGroupByAnswering = func(issues []githubstructures.Issue) ([]githubstructures.Issue, []githubstructures.Issue, []githubstructures.Issue) {
			return []githubstructures.Issue{}, []githubstructures.Issue{}, []githubstructures.Issue{}
		}
		mockFindLabels = func(repoName string) []githubstructures.Label {
			if repoName == "repo-1" {
				return []githubstructures.Label{
					githubstructures.Label{Name: "label-1", Color: "color-1-invalid"},
					githubstructures.Label{Name: "label-2", Color: "color-2"},
				}
			}
			if repoName == "repo-2" {
				return []githubstructures.Label{
					githubstructures.Label{Name: "label-1", Color: "color-1"},
				}
			}
			if repoName == "repo-3" {
				return []githubstructures.Label{
					githubstructures.Label{Name: "label-1", Color: "color-1"},
					githubstructures.Label{Name: "label-2", Color: "color-2-invalid"},
					githubstructures.Label{Name: "label-3", Color: "color-3-invalid"},
					githubstructures.Label{Name: "label-4", Color: "color-4"},
				}
			}
			return []githubstructures.Label{}
		}
		mockCreateLabel = func(repoName string, label githubstructures.Label) {
		}
		mockDeleteLabel = func(repoName string, labelName string) {
			mockDeleteLabelsParams = append(mockDeleteLabelsParams, []interface{}{repoName, labelName})
		}
		mockFindIssues = func(repoName string) []githubstructures.Issue {
			return []githubstructures.Issue{}
		}
		githubClient := Mockgithubclient{}
		issuesTriage := Mockissuestriage{}
		githubOperator := New(githubClient, issuesTriage, answeringLabels, "label-1", "label-2", "label-3", answeringLabels)

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
		answeringLabels := []githubstructures.Label{
			githubstructures.Label{Name: "label-1", Color: "color-1"},
			githubstructures.Label{Name: "label-2", Color: "color-2"},
			githubstructures.Label{Name: "label-3", Color: "color-3"},
		}

		mockGroupByAnswering = func(issues []githubstructures.Issue) ([]githubstructures.Issue, []githubstructures.Issue, []githubstructures.Issue) {
			return []githubstructures.Issue{
					githubstructures.Issue{Url: "url-1"},
					githubstructures.Issue{Url: "url-2"},
				},
				[]githubstructures.Issue{
					githubstructures.Issue{Url: "url-3"},
					githubstructures.Issue{Url: "url-4"},
				},
				[]githubstructures.Issue{
					githubstructures.Issue{Url: "url-5"},
					githubstructures.Issue{Url: "url-6"},
				}
		}
		mockFindLabels = func(repoName string) []githubstructures.Label {
			return []githubstructures.Label{}
		}
		mockCreateLabel = func(repoName string, label githubstructures.Label) {
		}
		mockFindIssues = func(repoName string) []githubstructures.Issue {
			return []githubstructures.Issue{}
		}
		mockAddLabel = func(issueUrl string, labelName string) {
			mockAddLabelParams = append(mockAddLabelParams, []interface{}{issueUrl, labelName})
		}
		githubClient := Mockgithubclient{}
		issuesTriage := Mockissuestriage{}
		githubOperator := New(githubClient, issuesTriage, answeringLabels, "by-ours", "answered", "not-answered", answeringLabels)

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
		answeringLabels := []githubstructures.Label{
			githubstructures.Label{Name: "by-ours", Color: "color-1"},
			githubstructures.Label{Name: "answered", Color: "color-2"},
			githubstructures.Label{Name: "not-answered", Color: "color-3"},
		}

		mockGroupByAnswering = func(issues []githubstructures.Issue) ([]githubstructures.Issue, []githubstructures.Issue, []githubstructures.Issue) {
			return []githubstructures.Issue{
					githubstructures.Issue{Url: "url-1", Labels: []githubstructures.Label{githubstructures.Label{Name: "by-ours"}}},
					githubstructures.Issue{Url: "url-2", Labels: []githubstructures.Label{githubstructures.Label{Name: "not-answered"}}},
				},
				[]githubstructures.Issue{
					githubstructures.Issue{Url: "url-3", Labels: []githubstructures.Label{githubstructures.Label{Name: "answered"}}},
					githubstructures.Issue{Url: "url-4", Labels: []githubstructures.Label{githubstructures.Label{Name: "not-answered"}}},
				},
				[]githubstructures.Issue{
					githubstructures.Issue{Url: "url-5", Labels: []githubstructures.Label{githubstructures.Label{Name: "answered"}}},
					githubstructures.Issue{Url: "url-6", Labels: []githubstructures.Label{githubstructures.Label{Name: "not-answered"}}},
				}
		}
		mockFindLabels = func(repoName string) []githubstructures.Label {
			return []githubstructures.Label{}
		}
		mockCreateLabel = func(repoName string, label githubstructures.Label) {
		}
		mockFindIssues = func(repoName string) []githubstructures.Issue {
			return []githubstructures.Issue{}
		}
		mockAddLabel = func(issueUrl string, labelName string) {
		}
		mockRemoveLabel = func(issueUrl string, labelName string) {
			mockRemoveLabelParams = append(mockRemoveLabelParams, []interface{}{issueUrl, labelName})
		}
		githubClient := Mockgithubclient{}
		issuesTriage := Mockissuestriage{}
		githubOperator := New(githubClient, issuesTriage, answeringLabels, "by-ours", "answered", "not-answered", answeringLabels)

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
