package main

import (
	"./GithubClient"
	"./GithubOperator"
	"./Interfaces"
	"./IssuesTriage"
	"fmt"
	"os"
)

func main() {
	organization := os.Args[1]
	token := os.Getenv("GITHUB_TOKEN")
	OUR_LABEL_TEXT := "answering: reported by " + organization
	const ANSWERED_LABEL_TEXT = "answering: answered"
	const NOT_ANSWERED_LABEL_TEXT = "answering: not answered"
	answeringLabels := []Interfaces.Label{
		Interfaces.Label{Name: OUR_LABEL_TEXT, Color: "a0a000"},
		Interfaces.Label{Name: ANSWERED_LABEL_TEXT, Color: "00a000"},
		Interfaces.Label{Name: NOT_ANSWERED_LABEL_TEXT, Color: "a00000"},
	}

	fmt.Println(token, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT)

	githubClient := GithubClient.InitGithubClient(organization, token)
	issuesTriage := IssuesTriage.InitIssuesTriage()
	githubOperator := GithubOperator.InitGithubOperator(githubClient, issuesTriage, answeringLabels, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT)
	repoNames := githubClient.FindRepos()
	fmt.Println("repoNames", repoNames)
	githubOperator.UpdateRepos(repoNames)
}
