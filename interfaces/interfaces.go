package interfaces

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
	Archived bool
	Name     string
}

type Label struct {
	Name  string
	Color string
}

type CommentAuthor struct {
	Login string
}

type Comment struct {
	AuthorAssociation string
	AuthorLogin       string
}

type Issue struct {
	Title             string
	Url               string
	Number            int
	AuthorAssociation string
	Labels            []Label
	Comments          []Comment
}
