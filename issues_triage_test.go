package main

import (
	"./IssuesTriage"
	"./Types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoIssuesTriage(t *testing.T) {
	t.Run("does triage of an empty list", func(t *testing.T) {
		issues := []Types.Issue{}

		issuesTriage := IssuesTriage.InitIssuesTriage()
		ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.DoIssuesTriage(issues)

		assert.Equal(t, ourIssues, []Types.Issue{})
		assert.Equal(t, answeredIssues, []Types.Issue{})
		assert.Equal(t, notAnsweredIssues, []Types.Issue{})
	})

	t.Run("does triage of issues with no comments", func(t *testing.T) {
		issues := []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{}}},
		}

		issuesTriage := IssuesTriage.InitIssuesTriage()
		ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.DoIssuesTriage(issues)

		assert.Equal(t, ourIssues, []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{}}},
		})
		assert.Equal(t, answeredIssues, []Types.Issue{})
		assert.Equal(t, notAnsweredIssues, []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{}}},
		})
	})

	t.Run("does triage of issues with comments", func(t *testing.T) {
		issues := []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
		}

		issuesTriage := IssuesTriage.InitIssuesTriage()
		ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.DoIssuesTriage(issues)

		assert.Equal(t, ourIssues, []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
		})
		assert.Equal(t, answeredIssues, []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
		})
		assert.Equal(t, notAnsweredIssues, []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
		})
	})

	t.Run("excludes issuehunt-app comments", func(t *testing.T) {
		issues := []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
		}

		issuesTriage := IssuesTriage.InitIssuesTriage()
		ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.DoIssuesTriage(issues)

		assert.Equal(t, ourIssues, []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
		})
		assert.Equal(t, answeredIssues, []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
		})
		assert.Equal(t, notAnsweredIssues, []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
			}}},
		})
	})
}
