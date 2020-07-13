package migrations

import (
	"log"
)

type GitHubOperator interface {
	RenameLabelInEachRepo(repoNames []string, oldLabelName string, newLabelName string)
}

func Up(githubOperator GitHubOperator, repoNames []string) {
	up_2020_07_13_issue_type(githubOperator, repoNames)
	log.Println("all migrations finished")
}
