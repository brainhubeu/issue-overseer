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

func findRepos(organization string) []string {
	repoNames := []string{}
	for page := 1; ; page += 1 {
		resp, err := http.Get("https://api.github.com/orgs/" + organization + "/repos?page=" + strconv.Itoa(page))
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		repositories := []repository{}
		json.Unmarshal([]byte(string(body)), &repositories)
		fmt.Println("repositories", repositories)
		if len(repositories) == 0 {
			break
		}
		for i := 0; i < len(repositories); i++ {
			repository := repositories[i]
			if !repository.Archived {
				repoNames = append(repoNames, repository.Name)
				fmt.Println("==repoNames", repoNames)
			}
		}
	}

	sort.Strings(repoNames)
	return repoNames
}

func main() {
	organization := os.Args[1]
	token := os.Getenv("GITHUB_TOKEN")
	OUR_LABEL_TEXT := "answering: reported by " + organization
	const ANSWERED_LABEL_TEXT = "answering: answered"
	const NOT_ANSWERED_LABEL_TEXT = "answering: not answered"

	fmt.Println(token, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT)

	repoNames := findRepos(organization)
	fmt.Println("repoNames", repoNames)
}
