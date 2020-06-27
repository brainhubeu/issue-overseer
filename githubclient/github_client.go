package githubclient

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

type Repository struct {
	Archived bool   `json:"archived"`
	Name     string `json:"name"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type CommentAuthor struct {
	Login string `json:"login"`
}

type Comment struct {
	AuthorAssociation string        `json:"authorAssociation"`
	Author            CommentAuthor `json:"author"`
}

type LabelEdge struct {
	Node Label `json:"node"`
}

type Labels struct {
	Edges []LabelEdge `json:"edges"`
}

type CommentEdge struct {
	Node Comment `json:"node"`
}

type Comments struct {
	Edges []CommentEdge `json:"edges"`
}

type Issue struct {
	Title             string   `json:"title"`
	Url               string   `json:"url"`
	Number            int      `json:"number"`
	AuthorAssociation string   `json:"authorAssociation"`
	Labels            Labels   `json:"labels"`
	Comments          Comments `json:"comments"`
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

type AddLabelRequestBody struct {
	Labels []string `json:"labels"`
}

func InitGithubClient(organization string, token string) *GithubClient {
	githubClient := &GithubClient{organization, token, 0}
	return githubClient
}

func (githubClient *GithubClient) incrementRequestNumber() {
	githubClient.RequestsNumber++
	fmt.Println("(v 1.0.7) request to GitHub #", githubClient.RequestsNumber)
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
	labelToCreate := Label{Name: label.Name, Color: label.Color}
	client := &http.Client{}
	jsonValue, err := json.Marshal(labelToCreate)
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

func transformDataIntoIssue(issueData Issue) Interfaces.Issue {
	labelsCount := len(issueData.Labels.Edges)
	commentsCount := len(issueData.Comments.Edges)
	labels := make([]Interfaces.Label, labelsCount)
	comments := make([]Interfaces.Comment, commentsCount)
	for i := 0; i < labelsCount; i++ {
		labelData := issueData.Labels.Edges[i].Node
		labels[i] = Interfaces.Label{
			Name:  labelData.Name,
			Color: labelData.Color,
		}
	}
	for i := 0; i < commentsCount; i++ {
		commentData := issueData.Comments.Edges[i].Node
		comments[i] = Interfaces.Comment{
			AuthorAssociation: commentData.AuthorAssociation,
			AuthorLogin:       commentData.Author.Login,
		}
	}

	return Interfaces.Issue{
		Title:             issueData.Title,
		Url:               issueData.Url,
		Number:            issueData.Number,
		AuthorAssociation: issueData.AuthorAssociation,
		Labels:            labels,
		Comments:          comments,
	}
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
			issueData := edges[i].Node
			result = append(result, transformDataIntoIssue(issueData))
		}
	}
	return result
}
