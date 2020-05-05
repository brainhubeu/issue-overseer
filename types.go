package main

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
