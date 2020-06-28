package main

import (
	"github.com/brainhubeu/issue-overseer/githubclient"
	"github.com/brainhubeu/issue-overseer/githuboperator"
	"github.com/brainhubeu/issue-overseer/interfaces"
	"github.com/brainhubeu/issue-overseer/issuestriage"
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

	githubClient := githubclient.New(organization, token)
	issuesTriage := issuestriage.New()
	githubOperator := githuboperator.New(githubClient, issuesTriage, answeringLabels, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT)
	repoNames := githubClient.FindRepos()
	log.Println("repoNames", repoNames)
	githubOperator.UpdateRepos(repoNames)
}
