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

func main() {
	organization := os.Args[1]
	token := os.Getenv("GITHUB_TOKEN")
	OUR_LABEL_TEXT := "answering: reported by " + organization
	const ANSWERED_LABEL_TEXT = "answering: answered"
	const NOT_ANSWERED_LABEL_TEXT = "answering: not answered"

	fmt.Println(token, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT)

	repoNames := findRepos(organization, token)
	for i := 0; i < len(repoNames); i++ {
		labels := findLabels(organization, repoNames[i], token)
		fmt.Println("labels", repoNames[i], labels)
	}
	fmt.Println("repoNames", repoNames)
}
