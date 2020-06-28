package githuboperator

import (
	"github.com/brainhubeu/issue-overseer/interfaces"
	"log"
)

type githuboperator struct {
	githubclient            interfaces.GithubClient
	issuestriage            interfaces.IssuesTriage
	AnsweringLabels         []interfaces.Label
	OUR_LABEL_TEXT          string
	ANSWERED_LABEL_TEXT     string
	NOT_ANSWERED_LABEL_TEXT string
}

func New(githubClient interfaces.GithubClient, issuesTriage interfaces.IssuesTriage, answeringLabels []interfaces.Label, OUR_LABEL_TEXT string, ANSWERED_LABEL_TEXT string, NOT_ANSWERED_LABEL_TEXT string) *githuboperator {
	githubOperator := &githuboperator{githubClient, issuesTriage, answeringLabels, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT}
	return githubOperator
}

func (githubOperator githuboperator) createOrUpdateRepoLabels(repoName string) {
	allLabels := githubOperator.githubclient.FindLabels(repoName)
	labelsToDelete := []interfaces.Label{}
	for i := 0; i < len(githubOperator.AnsweringLabels); i++ {
		label := githubOperator.AnsweringLabels[i]
		for j := 0; j < len(allLabels); j++ {
			anyLabel := allLabels[j]
			if label.Name == anyLabel.Name && label.Color != anyLabel.Color {
				labelsToDelete = append(labelsToDelete, label)
			}
		}
	}
	labelsToCreate := append([]interfaces.Label{}, labelsToDelete...)
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
	log.Println(repoName, "labelsToDelete", labelsToDelete)
	log.Println(repoName, "labelsToCreate", labelsToCreate)
	for i := 0; i < len(labelsToDelete); i++ {
		githubOperator.githubclient.DeleteLabel(repoName, labelsToDelete[i].Name)
	}
	for i := 0; i < len(labelsToCreate); i++ {
		githubOperator.githubclient.CreateLabel(repoName, labelsToCreate[i])
	}
}

func (githubOperator githuboperator) updateIssueLabels(issueUrl string, allIssueLabels []interfaces.Label, labelNameToAdd string) {
	labelsToRemove := []interfaces.Label{}
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

func (githubOperator githuboperator) UpdateRepos(repoNames []string) {
	for i := 0; i < len(repoNames); i++ {
		repoName := repoNames[i]
		githubOperator.createOrUpdateRepoLabels(repoName)
		issues := githubOperator.githubclient.FindIssues(repoName)
		ourIssues, answeredIssues, notAnsweredIssues := githubOperator.issuestriage.TriageManyIssues(issues)
		log.Println(repoName, "ourIssues", ourIssues)
		log.Println(repoName, "answeredIssues", answeredIssues)
		log.Println(repoName, "notAnsweredIssues", notAnsweredIssues)
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
