package githuboperator

import (
	"../Interfaces"
	"fmt"
)

type GithubOperator struct {
	GithubClient            Interfaces.GithubClient
	IssuesTriage            Interfaces.IssuesTriage
	AnsweringLabels         []Interfaces.Label
	OUR_LABEL_TEXT          string
	ANSWERED_LABEL_TEXT     string
	NOT_ANSWERED_LABEL_TEXT string
}

func InitGithubOperator(githubClient Interfaces.GithubClient, issuesTriage Interfaces.IssuesTriage, answeringLabels []Interfaces.Label, OUR_LABEL_TEXT string, ANSWERED_LABEL_TEXT string, NOT_ANSWERED_LABEL_TEXT string) *GithubOperator {
	githubOperator := &GithubOperator{githubClient, issuesTriage, answeringLabels, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT}
	return githubOperator
}

func (githubOperator GithubOperator) createOrUpdateRepoLabels(repoName string) {
	allLabels := githubOperator.GithubClient.FindLabels(repoName)
	labelsToDelete := []Interfaces.Label{}
	for i := 0; i < len(githubOperator.AnsweringLabels); i++ {
		label := githubOperator.AnsweringLabels[i]
		for j := 0; j < len(allLabels); j++ {
			anyLabel := allLabels[j]
			if label.Name == anyLabel.Name && label.Color != anyLabel.Color {
				labelsToDelete = append(labelsToDelete, label)
			}
		}
	}
	labelsToCreate := append([]Interfaces.Label{}, labelsToDelete...)
	for i := 0; i < len(githubOperator.AnsweringLabels); i++ {
		label := githubOperator.AnsweringLabels[i]
		j := 0
		for ; j < len(allLabels); j++ {
			anyLabel := allLabels[j]
			if label.Name == anyLabel.Name {
				break
			}
		}
		if j == len(allLabels) {
			labelsToCreate = append(labelsToCreate, label)
		}
	}
	fmt.Println(repoName, "labelsToDelete", labelsToDelete)
	fmt.Println(repoName, "labelsToCreate", labelsToCreate)
	for i := 0; i < len(labelsToDelete); i++ {
		githubOperator.GithubClient.DeleteLabel(repoName, labelsToDelete[i].Name)
	}
	for i := 0; i < len(labelsToCreate); i++ {
		githubOperator.GithubClient.CreateLabel(repoName, labelsToCreate[i])
	}
}

func (githubOperator GithubOperator) updateIssueLabels(issueUrl string, allIssueLabels []Interfaces.Label, labelNameToAdd string) {
	labelsToRemove := []Interfaces.Label{}
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
	fmt.Println(issueUrl, "labelsToRemove", labelsToRemove)
	for i := 0; i < len(labelsToRemove); i++ {
		githubOperator.GithubClient.RemoveLabel(issueUrl, labelsToRemove[i].Name)
	}
	githubOperator.GithubClient.AddLabel(issueUrl, labelNameToAdd)
}

func (githubOperator GithubOperator) UpdateRepos(repoNames []string) {
	for i := 0; i < len(repoNames); i++ {
		repoName := repoNames[i]
		githubOperator.createOrUpdateRepoLabels(repoName)
		issues := githubOperator.GithubClient.FindIssues(repoName)
		ourIssues, answeredIssues, notAnsweredIssues := githubOperator.IssuesTriage.TriageManyIssues(issues)
		fmt.Println(repoName, "ourIssues", ourIssues)
		fmt.Println(repoName, "answeredIssues", answeredIssues)
		fmt.Println(repoName, "notAnsweredIssues", notAnsweredIssues)
		for j := 0; j < len(ourIssues); j++ {
			githubOperator.updateIssueLabels(ourIssues[j].Url, ourIssues[j].Labels, githubOperator.OUR_LABEL_TEXT)
		}
		for j := 0; j < len(answeredIssues); j++ {
			githubOperator.updateIssueLabels(answeredIssues[j].Url, answeredIssues[j].Labels, githubOperator.ANSWERED_LABEL_TEXT)
		}
		for j := 0; j < len(notAnsweredIssues); j++ {
			githubOperator.updateIssueLabels(notAnsweredIssues[j].Url, notAnsweredIssues[j].Labels, githubOperator.NOT_ANSWERED_LABEL_TEXT)
		}
	}
}
