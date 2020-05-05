package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type GithubClient struct {
	Organization string
	Token        string
}

func (githubClient GithubClient) findRepos() []string {
	repoNames := []string{}
	client := &http.Client{}
	for page := 1; ; page += 1 {
		req, err := http.NewRequest("GET", "https://api.github.com/orgs/"+githubClient.Organization+"/repos?page="+strconv.Itoa(page), nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Add("Authorization", "token "+githubClient.Token)
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
		err = json.Unmarshal([]byte(string(body)), &repositories)
		if err != nil {
			log.Fatalln(err)
		}
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

func (githubClient GithubClient) findLabels(repoName string) []Label {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+githubClient.Organization+"/"+repoName+"/labels", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Authorization", "token "+githubClient.Token)
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
	err = json.Unmarshal([]byte(string(body)), &labels)
	if err != nil {
		log.Fatalln(err)
	}
	return labels
}

func (githubClient GithubClient) deleteLabel(repoName string, labelName string) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "https://api.github.com/repos/"+githubClient.Organization+"/"+repoName+"/labels/"+labelName, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Authorization", "token "+githubClient.Token)
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

func (githubClient GithubClient) createLabel(repoName string, label Label) {
	client := &http.Client{}
	jsonValue, err := json.Marshal(label)
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest("POST", "https://api.github.com/repos/"+githubClient.Organization+"/"+repoName+"/labels", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Authorization", "token "+githubClient.Token)
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

func (githubClient GithubClient) removeLabel(issueUrl string, labelName string) {
	client := &http.Client{}
	url := strings.Replace(issueUrl, "https://github.com", "https://api.github.com/repos", 1) + "/labels/" + labelName
	fmt.Println("to remove", issueUrl, url, labelName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Authorization", "token "+githubClient.Token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != 200 && resp.StatusCode != 404 {
		log.Fatalln(resp.Status, string(body))
	}
	fmt.Println("removed", issueUrl, labelName, resp.StatusCode)
}

func (githubClient GithubClient) addLabel(issueUrl string, labelName string) {
	client := &http.Client{}
	fmt.Println("labelName", labelName)
	requestBody := AddLabelRequestBody{[]string{labelName}}
	jsonValue, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalln(err)
	}
	url := strings.Replace(issueUrl, "https://github.com", "https://api.github.com/repos", 1) + "/labels"
	fmt.Println("to add", issueUrl, url, labelName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Authorization", "token "+githubClient.Token)
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
	fmt.Println("added", issueUrl, labelName, resp.StatusCode)
}

func (githubClient GithubClient) findIssues(repoName string) []Issue {
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
				labels (first:100) {
				  edges {
					node {
					  name
					}
				  }
				}
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
		graphqlVariables := GraphqlVariables{githubClient.Organization, repoName, cursor}
		graphqlRequestBody := GraphqlRequestBody{graphqlVariables, query}
		jsonValue, err := json.Marshal(graphqlRequestBody)
		if err != nil {
			log.Fatalln(err)
		}
		req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Add("Authorization", "token "+githubClient.Token)
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
		err = json.Unmarshal([]byte(string(body)), &issuesData)
		if err != nil {
			log.Fatalln(err)
		}
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
