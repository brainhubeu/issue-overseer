package main

import (
	"./githubclient"
	"./githuboperator"
	"./interfaces"
	"./issuestriage"
	"log"
	"os"
)

func main() {
	organization := os.Args[1]
	token := os.Getenv("GITHUB_TOKEN")
	OUR_LABEL_TEXT := "answering: reported by " + organization
	const ANSWERED_LABEL_TEXT = "answering: answered"
	const NOT_ANSWERED_LABEL_TEXT = "answering: not answered"
	answeringLabels := []interfaces.Label{
		interfaces.Label{Name: OUR_LABEL_TEXT, Color: "a0a000"},
		interfaces.Label{Name: ANSWERED_LABEL_TEXT, Color: "00a000"},
		interfaces.Label{Name: NOT_ANSWERED_LABEL_TEXT, Color: "a00000"},
	}

	log.Println(token, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT)

	githubClient := githubclient.Initgithubclient(organization, token)
	issuesTriage := issuestriage.Initissuestriage()
	githubOperator := githuboperator.Initgithuboperator(githubClient, issuesTriage, answeringLabels, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT)
	repoNames := githubClient.FindRepos()
	log.Println("repoNames", repoNames)
	githubOperator.UpdateRepos(repoNames)
}
