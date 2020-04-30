package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type Repository struct {
	Archived bool   `json:"archived"`
	Name     string `json:"name"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Issue struct {
	Title             string `json:"title"`
	Url               string `json:"url"`
	Number            string `json:"number"`
	AuthorAssociation string `json:"authorAssociation"`
}

type IssueEdge struct {
	Cursor string `json:"cursor"`
	Node   Issue  `json:"node"`
}

type Issues struct {
	Data struct {
		Repository struct {
			Issues struct {
				Edges []IssueEdge `json:"edges"`
			} `json:"issues"`
		} `json:"repository"`
	} `json:"data"`
}

type GraphqlVariables struct {
	Organization string  `json:"organization"`
	RepoName     string  `json:"repoName"`
	Cursor       *string `json:"cursor"`
}

type GraphqlRequestBody struct {
	Variables GraphqlVariables `json:"variables"`
	Query     string           `json:"query"`
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
		repositories := []Repository{}
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

func findLabels(organization string, repoName string, token string) []Label {
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
	labels := []Label{}
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

func createLabel(organization string, repoName string, label Label, token string) {
	client := &http.Client{}
	jsonValue, err := json.Marshal(label)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest("POST", "https://api.github.com/repos/"+organization+"/"+repoName+"/labels", bytes.NewBuffer(jsonValue))
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
	if resp.StatusCode != 201 {
		log.Fatalln(resp.Status, string(body))
	}
}

func findIssues(organization string, repoName string, token string) []Issue {
	client := &http.Client{}
	cursor := (*string)(nil)
	result := []Issue{}
	for {
		query := `query ($organization: String!, $repoName: String!, $cursor: String) {
	  repository(owner: $organization, name: $repoName) {
		issues(first:20, after: $cursor, states:OPEN) {
		  edges {
			cursor
			node {
			  title
			  url
			  number
			  authorAssociation
			  comments(last:100) {
				edges {
				  node {
					bodyText
					authorAssociation
					author {
					  login
					}
				  }
				}
			  }
			}
		  }
		}
	  }
	}`
		graphqlVariables := GraphqlVariables{organization, repoName, cursor}
		graphqlRequestBody := GraphqlRequestBody{graphqlVariables, query}
		jsonValue, err := json.Marshal(graphqlRequestBody)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(bytes.NewBuffer(jsonValue))
		req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonValue))
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
		issuesData := Issues{}
		json.Unmarshal([]byte(string(body)), &issuesData)
		edges := issuesData.Data.Repository.Issues.Edges
		if len(edges) == 0 {
			break
		}
		cursor = &edges[len(edges)-1].Cursor
		for i := 0; i < len(edges); i++ {
			result = append(result, edges[i].Node)
		}
	}
	return result
}

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

	repoNames := findRepos(organization, token)
	for i := 0; i < len(repoNames); i++ {
		repoName := repoNames[i]
		allLabels := findLabels(organization, repoName, token)
		labelsToDelete := []Label{}
		for j := 0; j < len(answeringLabels); j++ {
			label := answeringLabels[j]
			for k := 0; k < len(allLabels); k++ {
				anyLabel := allLabels[k]
				if label.Name == anyLabel.Name && label.Color != anyLabel.Color {
					labelsToDelete = append(labelsToDelete, label)
				}
			}
		}
		labelsToCreate := append([]Label{}, labelsToDelete...)
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
		for j := 0; j < len(labelsToCreate); j++ {
			createLabel(organization, repoName, labelsToCreate[j], token)
		}
		issues := findIssues(organization, repoName, token)
		fmt.Println(repoName, "issues", issues)
	}
	fmt.Println("repoNames", repoNames)
}
