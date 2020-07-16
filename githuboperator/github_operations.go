package githuboperator

import (
	"github.com/brainhubeu/issue-overseer/githubstructures"
	"log"
	"strings"
)

type GithubClient interface {
	FindRepos() []string
	FindLabels(repoName string) []githubstructures.Label
	DeleteLabel(repoName string, labelName string)
	CreateLabel(repoName string, label githubstructures.Label)
	RemoveLabel(issueUrl string, labelName string)
	AddLabel(issueUrl string, labelName string)
	RenameLabel(repoName string, oldLabelName string, newLabelName string)
	FindIssues(repoName string) []githubstructures.Issue
}

type IssuesTriage interface {
	GroupByAnswering(issues []githubstructures.Issue) ([]githubstructures.Issue, []githubstructures.Issue, []githubstructures.Issue)
	GroupByManualLabel(issues []githubstructures.Issue, config githubstructures.ManualLabelConfig) ([]githubstructures.Issue, []githubstructures.Issue)
}

type githuboperator struct {
	githubclient            GithubClient
	issuestriage            IssuesTriage
	AnsweringLabels         []githubstructures.Label
	OUR_LABEL_TEXT          string
	ANSWERED_LABEL_TEXT     string
	NOT_ANSWERED_LABEL_TEXT string
	DefaultLabels           []githubstructures.Label
	manualLabelConfigs      []githubstructures.ManualLabelConfig
}

func New(githubClient GithubClient, issuesTriage IssuesTriage, answeringLabels []githubstructures.Label, OUR_LABEL_TEXT string, ANSWERED_LABEL_TEXT string, NOT_ANSWERED_LABEL_TEXT string, defaultLabels []githubstructures.Label, manualLabelConfigs []githubstructures.ManualLabelConfig) *githuboperator {
	githubOperator := &githuboperator{githubClient, issuesTriage, answeringLabels, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT, defaultLabels, manualLabelConfigs}
	return githubOperator
}

func (githubOperator githuboperator) createOrUpdateRepoLabels(repoName string) {
	allLabels := githubOperator.githubclient.FindLabels(repoName)
	labelsToDelete := []githubstructures.Label{}
	for i := 0; i < len(githubOperator.DefaultLabels); i++ {
		label := githubOperator.DefaultLabels[i]
		for j := 0; j < len(allLabels); j++ {
			anyLabel := allLabels[j]
			if strings.EqualFold(label.Name, anyLabel.Name) && label.Color != anyLabel.Color {
				labelsToDelete = append(labelsToDelete, label)
			}
		}
	}
	labelsToCreate := append([]githubstructures.Label{}, labelsToDelete...)
	for i := 0; i < len(githubOperator.DefaultLabels); i++ {
		label := githubOperator.DefaultLabels[i]
		j := 0
		for ; j < len(allLabels); j++ {
			anyLabel := allLabels[j]
			if strings.EqualFold(label.Name, anyLabel.Name) {
				break
			}
		}
		if j == len(allLabels) {
			labelsToCreate = append(labelsToCreate, label)
		}
	}
	log.Println(repoName, "labelsToDelete", labelsToDelete)
	log.Println(repoName, "labelsToCreate", labelsToCreate)
	for i := 0; i < len(labelsToDelete); i++ {
		githubOperator.githubclient.DeleteLabel(repoName, labelsToDelete[i].Name)
	}
	for i := 0; i < len(labelsToCreate); i++ {
		githubOperator.githubclient.CreateLabel(repoName, labelsToCreate[i])
	}
}

func (githubOperator githuboperator) updateIssueLabels(issueUrl string, allIssueLabels []githubstructures.Label, labelNameToAdd string) {
	labelsToRemove := []githubstructures.Label{}
	for i := 0; i < len(allIssueLabels); i++ {
		j := 0
		for ; j < len(githubOperator.AnsweringLabels); j++ {
			if githubOperator.AnsweringLabels[j].Name == allIssueLabels[i].Name {
				break
			}
		}
		if j < len(githubOperator.AnsweringLabels) && allIssueLabels[i].Name != labelNameToAdd {
			labelsToRemove = append(labelsToRemove, allIssueLabels[i])
		}
	}
	log.Println(issueUrl, "labelsToRemove", labelsToRemove)
	for i := 0; i < len(labelsToRemove); i++ {
		githubOperator.githubclient.RemoveLabel(issueUrl, labelsToRemove[i].Name)
	}
	githubOperator.githubclient.AddLabel(issueUrl, labelNameToAdd)
}

func (githubOperator githuboperator) updateAnsweringLabelsForRepo(repoName string) {
	issues := githubOperator.githubclient.FindIssues(repoName)
	ourIssues, answeredIssues, notAnsweredIssues := githubOperator.issuestriage.GroupByAnswering(issues)
	log.Println(repoName, "ourIssues", ourIssues)
	log.Println(repoName, "answeredIssues", answeredIssues)
	log.Println(repoName, "notAnsweredIssues", notAnsweredIssues)
	for i := 0; i < len(ourIssues); i++ {
		githubOperator.updateIssueLabels(ourIssues[i].Url, ourIssues[i].Labels, githubOperator.OUR_LABEL_TEXT)
	}
	for i := 0; i < len(answeredIssues); i++ {
		githubOperator.updateIssueLabels(answeredIssues[i].Url, answeredIssues[i].Labels, githubOperator.ANSWERED_LABEL_TEXT)
	}
	for i := 0; i < len(notAnsweredIssues); i++ {
		githubOperator.updateIssueLabels(notAnsweredIssues[i].Url, notAnsweredIssues[i].Labels, githubOperator.NOT_ANSWERED_LABEL_TEXT)
	}
}

func (githubOperator githuboperator) updateMissingManualLabelsForRepo(repoName string) {
	configs := githubOperator.manualLabelConfigs
	for i := 0; i < len(configs); i++ {
		config := configs[i]
		issues := githubOperator.githubclient.FindIssues(repoName)
		issuesWithLabel, issuesWithoutLabel := githubOperator.issuestriage.GroupByManualLabel(issues, config)
		log.Println(repoName, "issues with manual label", config.Prefix, issuesWithLabel)
		log.Println(repoName, "issues without manual label", config.Prefix, issuesWithoutLabel)
		for j := 0; j < len(issuesWithLabel); j++ {
			githubOperator.githubclient.RemoveLabel(issuesWithLabel[j].Url, "missing "+config.Prefix)
		}
		for j := 0; j < len(issuesWithoutLabel); j++ {
			githubOperator.githubclient.AddLabel(issuesWithoutLabel[j].Url, "missing "+config.Prefix)
		}
	}
}

func (githubOperator githuboperator) UpdateRepos(repoNames []string) {
	for i := 0; i < len(repoNames); i++ {
		repoName := repoNames[i]
		githubOperator.createOrUpdateRepoLabels(repoName)
		githubOperator.updateAnsweringLabelsForRepo(repoName)
		githubOperator.updateMissingManualLabelsForRepo(repoName)
	}
}

func (githubOperator githuboperator) RenameLabelInEachRepo(repoNames []string, oldLabelName string, newLabelName string) {
	for i := 0; i < len(repoNames); i++ {
		repoName := repoNames[i]
		githubOperator.githubclient.RenameLabel(repoName, oldLabelName, newLabelName)
	}
}
