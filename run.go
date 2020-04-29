package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type repository struct {
	Archived bool   `json:"archived"`
	Name     string `json:"name"`
}

type label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func findRepos(organization string, token string) []string {
	repoNames := []string{}
	client := &http.Client{}
	for page := 1; ; page += 1 {
		req, err := http.NewRequest("GET", "https://api.github.com/orgs/"+organization+"/repos?page="+strconv.Itoa(page), nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Add("Authorization", "token "+token)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		if resp.StatusCode != 200 {
			log.Fatalln(resp.Status, string(body))
		}
		repositories := []repository{}
		json.Unmarshal([]byte(string(body)), &repositories)
		if len(repositories) == 0 {
			break
		}
		for i := 0; i < len(repositories); i++ {
			repository := repositories[i]
			if !repository.Archived {
				repoNames = append(repoNames, repository.Name)
			}
		}
	}

	sort.Strings(repoNames)
	return repoNames
}

func findLabels(organization string, repoName string, token string) []label {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+organization+"/"+repoName+"/labels", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Authorization", "token "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != 200 {
		log.Fatalln(resp.Status, string(body))
	}
	labels := []label{}
	json.Unmarshal([]byte(string(body)), &labels)
	return labels
}

func deleteLabel(organization string, repoName string, labelName string, token string) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "https://api.github.com/repos/"+organization+"/"+repoName+"/labels/"+labelName, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Authorization", "token "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != 204 {
		log.Fatalln(resp.Status, string(body))
	}
}

func main() {
	organization := os.Args[1]
	token := os.Getenv("GITHUB_TOKEN")
	OUR_LABEL_TEXT := "answering: reported by " + organization
	const ANSWERED_LABEL_TEXT = "answering: answered"
	const NOT_ANSWERED_LABEL_TEXT = "answering: not answered"
	answeringLabels := []label{
		label{OUR_LABEL_TEXT, "a0a000"},
		label{ANSWERED_LABEL_TEXT, "00a000"},
		label{NOT_ANSWERED_LABEL_TEXT, "a00000"},
	}

	fmt.Println(token, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT)

	repoNames := findRepos(organization, token)
	for i := 0; i < len(repoNames); i++ {
		repoName := repoNames[i]
		allLabels := findLabels(organization, repoName, token)
		labelsToDelete := []label{}
		for j := 0; j < len(answeringLabels); j++ {
			label := answeringLabels[j]
			for k := 0; k < len(allLabels); k++ {
				anyLabel := allLabels[k]
				if label.Name == anyLabel.Name && label.Color != anyLabel.Color {
					labelsToDelete = append(labelsToDelete, label)
				}
			}
		}
		labelsToCreate := append([]label{}, labelsToDelete...)
		for j := 0; j < len(answeringLabels); j++ {
			label := answeringLabels[j]
			k := 0
			for ; k < len(allLabels); k++ {
				anyLabel := allLabels[k]
				if label.Name == anyLabel.Name {
					break
				}
			}
			if k == len(allLabels) {
				labelsToCreate = append(labelsToCreate, label)
			}
		}
		fmt.Println(repoName, "allLabels", repoNames[i], allLabels)
		fmt.Println(repoName, "answeringLabels", answeringLabels)
		fmt.Println(repoName, "labelsToDelete", labelsToDelete)
		fmt.Println(repoName, "labelsToCreate", labelsToCreate)
		for j := 0; j < len(labelsToDelete); j++ {
			deleteLabel(organization, repoName, labelsToDelete[j].Name, token)
		}
	}
	fmt.Println("repoNames", repoNames)
}
