package githubstructures

type issueAnsweringTypeEnum struct {
	OURS         int
	ANSWERED     int
	NOT_ANSWERED int
}

var IssueAnsweringTypeEnum = &issueAnsweringTypeEnum{
	OURS:         1,
	ANSWERED:     2,
	NOT_ANSWERED: 3,
}

type issueManualLabelTypeEnum struct {
	EXISTENT     int
	NON_EXISTENT int
}

var IssueManualLabelTypeEnum = &issueManualLabelTypeEnum{
	EXISTENT:     1,
	NON_EXISTENT: 2,
}

type Label struct {
	Name  string
	Color string
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

type ManualLabelConfig struct {
	Prefix          string
	ParentLabelName string
}
