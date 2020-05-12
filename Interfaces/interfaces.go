package Interfaces

type GithubClient interface {
	FindRepos() []string
	FindLabels(repoName string) []Label
	DeleteLabel(repoName string, labelName string)
	CreateLabel(repoName string, label Label)
	RemoveLabel(issueUrl string, labelName string)
	AddLabel(issueUrl string, labelName string)
	FindIssues(repoName string) []Issue
}

type IssuesTriage interface {
	TriageManyIssues(issues []Issue) ([]Issue, []Issue, []Issue)
}

type issueTypeEnum struct {
	OURS         int
	ANSWERED     int
	NOT_ANSWERED int
}

var IssueTypeEnum = &issueTypeEnum{
	OURS:         1,
	ANSWERED:     2,
	NOT_ANSWERED: 3,
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
