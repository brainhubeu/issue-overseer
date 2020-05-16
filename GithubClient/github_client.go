package GithubClient

import (
	"../Interfaces"
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
	Organization   string
	Token          string
	RequestsNumber int
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

type AddLabelRequestBody struct {
	Labels []string `json:"labels"`
}

func InitGithubClient(organization string, token string) *GithubClient {
	githubClient := &GithubClient{organization, token, 0}
	return githubClient
}

func (githubClient *GithubClient) incrementRequestNumber() {
	githubClient.RequestsNumber++
	fmt.Println("(v 1.0.3) request to GitHub #", githubClient.RequestsNumber)
}

func (githubClient *GithubClient) FindRepos() []string {
	repoNames := []string{}
	client := &http.Client{}
	for page := 1; ; page += 1 {
		githubClient.incrementRequestNumber()
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
		repositories := []Interfaces.Repository{}
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

func (githubClient *GithubClient) FindLabels(repoName string) []Interfaces.Label {
	client := &http.Client{}
	githubClient.incrementRequestNumber()
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
	labels := []Interfaces.Label{}
	err = json.Unmarshal([]byte(string(body)), &labels)
	if err != nil {
		log.Fatalln(err)
	}
	return labels
}

func (githubClient *GithubClient) DeleteLabel(repoName string, labelName string) {
	client := &http.Client{}
	githubClient.incrementRequestNumber()
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

func (githubClient *GithubClient) CreateLabel(repoName string, label Interfaces.Label) {
	client := &http.Client{}
	jsonValue, err := json.Marshal(label)
	if err != nil {
		log.Fatalln(err)
	}
	githubClient.incrementRequestNumber()
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

func (githubClient *GithubClient) RemoveLabel(issueUrl string, labelName string) {
	client := &http.Client{}
	url := strings.Replace(issueUrl, "https://github.com", "https://api.github.com/repos", 1) + "/labels/" + labelName
	fmt.Println("to remove", issueUrl, url, labelName)
	githubClient.incrementRequestNumber()
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

func (githubClient *GithubClient) AddLabel(issueUrl string, labelName string) {
	client := &http.Client{}
	fmt.Println("labelName", labelName)
	requestBody := AddLabelRequestBody{Labels: []string{labelName}}
	jsonValue, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalln(err)
	}
	url := strings.Replace(issueUrl, "https://github.com", "https://api.github.com/repos", 1) + "/labels"
	fmt.Println("to add", issueUrl, url, labelName)
	githubClient.incrementRequestNumber()
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

func (githubClient *GithubClient) FindIssues(repoName string) []Interfaces.Issue {
	client := &http.Client{}
	cursor := (*string)(nil)
	result := []Interfaces.Issue{}
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
		graphqlVariables := GraphqlVariables{Organization: githubClient.Organization, RepoName: repoName, Cursor: cursor}
		graphqlRequestBody := GraphqlRequestBody{Variables: graphqlVariables, Query: query}
		jsonValue, err := json.Marshal(graphqlRequestBody)
		if err != nil {
			log.Fatalln(err)
		}
		githubClient.incrementRequestNumber()
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
		issuesData := Interfaces.Issues{}
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
