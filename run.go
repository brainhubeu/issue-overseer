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
	"strings"
)

type Repository struct {
	Archived bool   `json:"archived"`
	Name     string `json:"name"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Comment struct {
	AuthorAssociation string `json:"authorAssociation"`
	Author            struct {
		Login string `json:"login"`
	} `json:"author"`
}

type LabelEdge struct {
	Node Label `json:"node"`
}

type Issue struct {
	Title             string `json:"title"`
	Url               string `json:"url"`
	Number            string `json:"number"`
	AuthorAssociation string `json:"authorAssociation"`
	Labels            struct {
		Edges []LabelEdge `json:"edges"`
	} `json:"labels"`
	Comments struct {
		Edges []struct {
			Node Comment `json:"node"`
		} `json:"edges"`
	} `json:"comments"`
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
	err = json.Unmarshal([]byte(string(body)), &labels)
	if err != nil {
		log.Fatalln(err)
	}
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

func removeLabel(issueUrl string, labelName string, token string) {
	client := &http.Client{}
	url := strings.Replace(issueUrl, "https://github.com", "https://api.github.com/repos", 1) + "/labels/" + labelName
	fmt.Println("to remove", issueUrl, url, labelName)
	req, err := http.NewRequest("DELETE", url, nil)
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
	if resp.StatusCode != 200 && resp.StatusCode != 404 {
		log.Fatalln(resp.Status, string(body))
	}
	fmt.Println("removed", issueUrl, labelName, resp.StatusCode)
}

func addLabel(issueUrl string, labelName string, token string) {
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
	fmt.Println("added", issueUrl, labelName, resp.StatusCode)
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
		graphqlVariables := GraphqlVariables{organization, repoName, cursor}
		graphqlRequestBody := GraphqlRequestBody{graphqlVariables, query}
		jsonValue, err := json.Marshal(graphqlRequestBody)
		if err != nil {
			log.Fatalln(err)
		}
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

func doIssuesTriage(issues []Issue) ([]Issue, []Issue, []Issue) {
	ourIssues := []Issue{}
	answeredIssues := []Issue{}
	notAnsweredIssues := []Issue{}
	for i := 0; i < len(issues); i++ {
		issue := issues[i]
		comments := issue.Comments.Edges
		if issue.AuthorAssociation == "MEMBER" {
			j := len(comments) - 1
			lastAuthorAssociation := ""
			for ; j >= 0; j-- {
				comment := comments[j].Node
				if comment.Author.Login != "issuehunt-app" && lastAuthorAssociation == "" {
					lastAuthorAssociation = comment.AuthorAssociation
				}
				if comment.Author.Login != "issuehunt-app" && comment.AuthorAssociation != "MEMBER" {
					break
				}
			}
			if j == -1 {
				ourIssues = append(ourIssues, issue)
			} else {
				if lastAuthorAssociation == "MEMBER" {
					answeredIssues = append(answeredIssues, issue)
				} else {
					notAnsweredIssues = append(notAnsweredIssues, issue)
				}
			}
		} else {
			if len(comments) == 0 {
				notAnsweredIssues = append(notAnsweredIssues, issue)
			} else {
				j := len(comments) - 1
				for ; j >= 0; j-- {
					if comments[j].Node.Author.Login != "issuehunt-app" {
						break
					}
				}
				if comments[j].Node.AuthorAssociation == "MEMBER" {
					answeredIssues = append(answeredIssues, issue)
				} else {
					notAnsweredIssues = append(notAnsweredIssues, issue)
				}
			}
		}
	}
	return ourIssues, answeredIssues, notAnsweredIssues
}

func createOrUpdateRepoLabels(organization string, repoName string, token string, answeringLabels []Label) {
	allLabels := findLabels(organization, repoName, token)
	labelsToDelete := []Label{}
	for i := 0; i < len(answeringLabels); i++ {
		label := answeringLabels[i]
		for j := 0; j < len(allLabels); j++ {
			anyLabel := allLabels[j]
			if label.Name == anyLabel.Name && label.Color != anyLabel.Color {
				labelsToDelete = append(labelsToDelete, label)
			}
		}
	}
	labelsToCreate := append([]Label{}, labelsToDelete...)
	for i := 0; i < len(answeringLabels); i++ {
		label := answeringLabels[i]
		j := 0
		for ; j < len(allLabels); j++ {
			anyLabel := allLabels[j]
			if label.Name == anyLabel.Name {
				break
			}
		}
		if j == len(allLabels) {
			labelsToCreate = append(labelsToCreate, label)
		}
	}
	fmt.Println(repoName, "labelsToDelete", labelsToDelete)
	fmt.Println(repoName, "labelsToCreate", labelsToCreate)
	for i := 0; i < len(labelsToDelete); i++ {
		deleteLabel(organization, repoName, labelsToDelete[i].Name, token)
	}
	for i := 0; i < len(labelsToCreate); i++ {
		createLabel(organization, repoName, labelsToCreate[i], token)
	}
}

func updateIssueLabels(issueUrl string, allIssueLabels []LabelEdge, answeringLabels []Label, labelNameToAdd string, token string) {
	labelsToRemove := []Label{}
	for i := 0; i < len(allIssueLabels)-1; i++ {
		j := 0
		for ; j < len(answeringLabels); j++ {
			if answeringLabels[j].Name == allIssueLabels[i].Node.Name {
				break
			}
		}
		if j < len(answeringLabels) && allIssueLabels[i].Node.Name != labelNameToAdd {
			labelsToRemove = append(labelsToRemove, allIssueLabels[i].Node)
		}
	}
	fmt.Println(issueUrl, "labelsToRemove", labelsToRemove)
	for i := 0; i < len(labelsToRemove); i++ {
		removeLabel(issueUrl, labelsToRemove[i].Name, token)
	}
	addLabel(issueUrl, labelNameToAdd, token)
}

func updateRepos(organization string, repoNames []string, token string, OUR_LABEL_TEXT string, ANSWERED_LABEL_TEXT string, NOT_ANSWERED_LABEL_TEXT string, answeringLabels []Label) {
	for i := 0; i < len(repoNames); i++ {
		repoName := repoNames[i]
		createOrUpdateRepoLabels(organization, repoName, token, answeringLabels)
		issues := findIssues(organization, repoName, token)
		ourIssues, answeredIssues, notAnsweredIssues := doIssuesTriage(issues)
		fmt.Println(repoName, "ourIssues", ourIssues)
		fmt.Println(repoName, "answeredIssues", answeredIssues)
		fmt.Println(repoName, "notAnsweredIssues", notAnsweredIssues)
		for j := 0; j < len(ourIssues); j++ {
			updateIssueLabels(ourIssues[j].Url, ourIssues[j].Labels.Edges, answeringLabels, OUR_LABEL_TEXT, token)
		}
		for j := 0; j < len(answeredIssues); j++ {
			updateIssueLabels(answeredIssues[j].Url, answeredIssues[j].Labels.Edges, answeringLabels, ANSWERED_LABEL_TEXT, token)
		}
		for j := 0; j < len(notAnsweredIssues); j++ {
			updateIssueLabels(notAnsweredIssues[j].Url, notAnsweredIssues[j].Labels.Edges, answeringLabels, NOT_ANSWERED_LABEL_TEXT, token)
		}
	}
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
	fmt.Println("repoNames", repoNames)
	updateRepos(organization, repoNames, token, OUR_LABEL_TEXT, ANSWERED_LABEL_TEXT, NOT_ANSWERED_LABEL_TEXT, answeringLabels)
}
