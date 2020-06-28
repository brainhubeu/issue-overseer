package githubclient

import (
	"bytes"
	"encoding/json"
	"github.com/brainhubeu/issue-overseer/interfaces"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type githubclient struct {
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

func Initgithubclient(organization string, token string) *githubclient {
	githubClient := &githubclient{organization, token, 0}
	return githubClient
}

func (githubClient *githubclient) incrementRequestNumber() {
	githubClient.RequestsNumber++
	log.Println("(v 1.0.7) request to GitHub #", githubClient.RequestsNumber)
}

func createJson(data interface{}) io.Reader {
	if data == nil {
		return nil
	}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	return bytes.NewBuffer(jsonValue)
}

func isStatusOK(actual int, expectedValues []int) bool {
	for i := 0; i < len(expectedValues); i++ {
		if actual == expectedValues[i] {
			return true
		}
	}
	return false
}

func (githubClient *githubclient) request(method string, url string, source interface{}, expectedStatuses []int, requestBody interface{}) {
	log.Println("request", method, url)
	client := &http.Client{}
	githubClient.incrementRequestNumber()
	jsonReader := createJson(requestBody)
	req, err := http.NewRequest(method, url, jsonReader)
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
	if !isStatusOK(resp.StatusCode, expectedStatuses) {
		log.Fatalln(resp.Status, string(body))
	}
	err = json.Unmarshal([]byte(string(body)), &source)
	if err != nil {
		log.Fatalln(err)
	}
}

func (githubClient *githubclient) FindRepos() []string {
	repoNames := []string{}
	for page := 1; ; page += 1 {
		repositories := []interfaces.Repository{}
		githubClient.request(http.MethodGet, "https://api.github.com/orgs/"+githubClient.Organization+"/repos?page="+strconv.Itoa(page), &repositories, []int{200}, nil)
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

func (githubClient *githubclient) FindLabels(repoName string) []interfaces.Label {
	labels := []interfaces.Label{}
	githubClient.request(http.MethodGet, "https://api.github.com/repos/"+githubClient.Organization+"/"+repoName+"/labels", &labels, []int{200}, nil)
	return labels
}

func (githubClient *githubclient) DeleteLabel(repoName string, labelName string) {
	githubClient.request(http.MethodDelete, "https://api.github.com/repos/"+githubClient.Organization+"/"+repoName+"/labels/"+labelName, nil, []int{204}, nil)
}

func (githubClient *githubclient) CreateLabel(repoName string, label interfaces.Label) {
	labelToCreate := Label{Name: label.Name, Color: label.Color}
	githubClient.request(http.MethodPost, "https://api.github.com/repos/"+githubClient.Organization+"/"+repoName+"/labels", nil, []int{201}, labelToCreate)
}

func (githubClient *githubclient) RemoveLabel(issueUrl string, labelName string) {
	url := strings.Replace(issueUrl, "https://github.com", "https://api.github.com/repos", 1) + "/labels/" + labelName
	githubClient.request(http.MethodDelete, url, nil, []int{200, 404}, nil)
}

func (githubClient *githubclient) AddLabel(issueUrl string, labelName string) {
	requestBody := AddLabelRequestBody{Labels: []string{labelName}}
	url := strings.Replace(issueUrl, "https://github.com", "https://api.github.com/repos", 1) + "/labels"
	githubClient.request(http.MethodPost, url, nil, []int{200}, requestBody)
}

func transformDataIntoIssue(issueData Issue) interfaces.Issue {
	labelsCount := len(issueData.Labels.Edges)
	commentsCount := len(issueData.Comments.Edges)
	labels := make([]interfaces.Label, labelsCount)
	comments := make([]interfaces.Comment, commentsCount)
	for i := 0; i < labelsCount; i++ {
		labelData := issueData.Labels.Edges[i].Node
		labels[i] = interfaces.Label{
			Name:  labelData.Name,
			Color: labelData.Color,
		}
	}
	for i := 0; i < commentsCount; i++ {
		commentData := issueData.Comments.Edges[i].Node
		comments[i] = interfaces.Comment{
			AuthorAssociation: commentData.AuthorAssociation,
			AuthorLogin:       commentData.Author.Login,
		}
	}

	return interfaces.Issue{
		Title:             issueData.Title,
		Url:               issueData.Url,
		Number:            issueData.Number,
		AuthorAssociation: issueData.AuthorAssociation,
		Labels:            labels,
		Comments:          comments,
	}
}

func (githubClient *githubclient) FindIssues(repoName string) []interfaces.Issue {
	cursor := (*string)(nil)
	result := []interfaces.Issue{}
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
		issuesData := Issues{}
		githubClient.request(http.MethodPost, "https://api.github.com/graphql", &issuesData, []int{200}, graphqlRequestBody)
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