package main

import (
	"fmt"
)

type GithubOperator struct {
	GithubClient            GithubClient
	AnsweringLabels         []Label
	OUR_LABEL_TEXT          string
	ANSWERED_LABEL_TEXT     string
	NOT_ANSWERED_LABEL_TEXT string
}

func (githubOperator GithubOperator) createOrUpdateRepoLabels(repoName string) {
	allLabels := githubOperator.GithubClient.findLabels(repoName)
	labelsToDelete := []Label{}
	for i := 0; i < len(githubOperator.AnsweringLabels); i++ {
		label := githubOperator.AnsweringLabels[i]
		for j := 0; j < len(allLabels); j++ {
			anyLabel := allLabels[j]
			if label.Name == anyLabel.Name && label.Color != anyLabel.Color {
				labelsToDelete = append(labelsToDelete, label)
			}
		}
	}
	labelsToCreate := append([]Label{}, labelsToDelete...)
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
		githubOperator.GithubClient.deleteLabel(repoName, labelsToDelete[i].Name)
	}
	for i := 0; i < len(labelsToCreate); i++ {
		githubOperator.GithubClient.createLabel(repoName, labelsToCreate[i])
	}
}

func (githubOperator GithubOperator) updateIssueLabels(issueUrl string, allIssueLabels []LabelEdge, labelNameToAdd string) {
	labelsToRemove := []Label{}
	for i := 0; i < len(allIssueLabels)-1; i++ {
		j := 0
		for ; j < len(githubOperator.AnsweringLabels); j++ {
			if githubOperator.AnsweringLabels[j].Name == allIssueLabels[i].Node.Name {
				break
			}
		}
		if j < len(githubOperator.AnsweringLabels) && allIssueLabels[i].Node.Name != labelNameToAdd {
			labelsToRemove = append(labelsToRemove, allIssueLabels[i].Node)
		}
	}
	fmt.Println(issueUrl, "labelsToRemove", labelsToRemove)
	for i := 0; i < len(labelsToRemove); i++ {
		githubOperator.GithubClient.removeLabel(issueUrl, labelsToRemove[i].Name)
	}
	githubOperator.GithubClient.addLabel(issueUrl, labelNameToAdd)
}

func (githubOperator GithubOperator) updateRepos(repoNames []string) {
	for i := 0; i < len(repoNames); i++ {
		repoName := repoNames[i]
		githubOperator.createOrUpdateRepoLabels(repoName)
		issues := githubOperator.GithubClient.findIssues(repoName)
		ourIssues, answeredIssues, notAnsweredIssues := doIssuesTriage(issues)
		fmt.Println(repoName, "ourIssues", ourIssues)
		fmt.Println(repoName, "answeredIssues", answeredIssues)
		fmt.Println(repoName, "notAnsweredIssues", notAnsweredIssues)
		for j := 0; j < len(ourIssues); j++ {
			githubOperator.updateIssueLabels(ourIssues[j].Url, ourIssues[j].Labels.Edges, githubOperator.OUR_LABEL_TEXT)
		}
		for j := 0; j < len(answeredIssues); j++ {
			githubOperator.updateIssueLabels(answeredIssues[j].Url, answeredIssues[j].Labels.Edges, githubOperator.ANSWERED_LABEL_TEXT)
		}
		for j := 0; j < len(notAnsweredIssues); j++ {
			githubOperator.updateIssueLabels(notAnsweredIssues[j].Url, notAnsweredIssues[j].Labels.Edges, githubOperator.NOT_ANSWERED_LABEL_TEXT)
		}
	}
}
