package main

import (
	"fmt"
	"os"
)

func main() {
	organization := os.Args[1]
	token := os.Getenv("GITHUB_TOKEN")
	OUR_LABEL_TEXT := "answering: reported by " + organization
	const ANSWERED_LABEL_TEXT = "answering: answered"
	const NOT_ANSWERED_LABEL_TEXT = "answering: not answered"
	answeringLabels := []Label{
		Label{OUR_LABEL_TEXT, "a0a000"},
		Label{ANSWERED_LABEL_TEXT, "00a000"},
		Label{NOT_ANSWERED_LABEL_TEXT, "a00000"},
	}

	fmt.Println(token, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT)

	githubClient := InitGithubClient(organization, token)
	githubOperator := GithubOperator{githubClient, answeringLabels, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT}
	repoNames := githubClient.findRepos()
	fmt.Println("repoNames", repoNames)
	githubOperator.updateRepos(repoNames)
}
