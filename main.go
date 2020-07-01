package main

import (
	"github.com/brainhubeu/issue-overseer/githubclient"
	"github.com/brainhubeu/issue-overseer/githuboperator"
	"github.com/brainhubeu/issue-overseer/githubstructures"
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
	answeringLabels := []githubstructures.Label{
		githubstructures.Label{Name: OUR_LABEL_TEXT, Color: "a0a000"},
		githubstructures.Label{Name: ANSWERED_LABEL_TEXT, Color: "00a000"},
		githubstructures.Label{Name: NOT_ANSWERED_LABEL_TEXT, Color: "a00000"},
	}

	defaultLabels := append([]githubstructures.Label{
		githubstructures.Label{Name: "WIP", Color: "a0a000"},
		githubstructures.Label{Name: "blocked", Color: "000000"},
		githubstructures.Label{Name: "hacktoberfest", Color: "202c99"},
		githubstructures.Label{Name: "in code review", Color: "ccfeff"},
		githubstructures.Label{Name: "needs discussion", Color: "dbf259"},
		githubstructures.Label{Name: "needs testing", Color: "dfdf00"},
		githubstructures.Label{Name: "no reproduction details", Color: "c91eb8"},
		githubstructures.Label{Name: "proposed issuehunt", Color: "2803ba"},
		githubstructures.Label{Name: "severity: blocked", Color: "000000"},
		githubstructures.Label{Name: "severity: critical", Color: "800000"},
		githubstructures.Label{Name: "severity: major", Color: "d00000"},
		githubstructures.Label{Name: "severity: medium", Color: "a0a000"},
		githubstructures.Label{Name: "severity: minor", Color: "40a000"},
		githubstructures.Label{Name: "severity: trivial", Color: "40ff40"},
		githubstructures.Label{Name: "tested & fails", Color: "ff4040"},
		githubstructures.Label{Name: "tested & works", Color: "40ff40"},
	}, answeringLabels...)

	log.Println(token, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT)

	githubClient := githubclient.New(organization, token)
	issuesTriage := issuestriage.New()
	githubOperator := githuboperator.New(githubClient, issuesTriage, answeringLabels, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT, defaultLabels)
	repoNames := githubClient.FindRepos()
	log.Println("repoNames", repoNames)
	githubOperator.UpdateRepos(repoNames)
}
