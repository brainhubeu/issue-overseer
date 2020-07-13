package migrations

import (
	"log"
)

func up_2020_07_13_issue_type(githubOperator GitHubOperator, repoNames []string) {
	githubOperator.RenameLabelInEachRepo(repoNames, "bug", "type: bug")
	githubOperator.RenameLabelInEachRepo(repoNames, "enhancement", "type: enhancement")
	githubOperator.RenameLabelInEachRepo(repoNames, "question", "type: question")
	log.Println("2020-07-13-issue-type migration finished")
}
