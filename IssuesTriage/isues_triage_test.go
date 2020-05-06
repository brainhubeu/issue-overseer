package IssuesTriage

import (
	"../Types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTriageManyIssues(t *testing.T) {
	t.Run("triages an empty list", func(t *testing.T) {
		issues := []Types.Issue{}

		issuesTriage := InitIssuesTriage()
		ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.TriageManyIssues(issues)

		assert.Equal(t, ourIssues, []Types.Issue{})
		assert.Equal(t, answeredIssues, []Types.Issue{})
		assert.Equal(t, notAnsweredIssues, []Types.Issue{})
	})

	t.Run("triages a non-empty list", func(t *testing.T) {
		issues := []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 122, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 124, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 125, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 126, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 127, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			}}},
		}

		issuesTriage := InitIssuesTriage()
		ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.TriageManyIssues(issues)

		assert.Equal(t, ourIssues, []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 122, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
		})
		assert.Equal(t, answeredIssues, []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 125, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 126, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
		})
		assert.Equal(t, notAnsweredIssues, []Types.Issue{
			Types.Issue{Title: "title", Url: "url", Number: 124, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
				Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			}}},
			Types.Issue{Title: "title", Url: "url", Number: 127, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			}}},
		})
	})

	t.Run("returns OURS for an issue created by a member with no comments", func(t *testing.T) {
		issue := Types.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
		}}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Types.IssueTypeEnum.OURS)
	})

	t.Run("returns OURS for an issue created by a member with a member comment", func(t *testing.T) {
		issue := Types.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
		}}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Types.IssueTypeEnum.OURS)
	})

	t.Run("returns NOT_ANSWERED for an issue created by a member with an external comment", func(t *testing.T) {
		issue := Types.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
		}}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Types.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("returns ANSWERED for an issue created by a member with an external comment followed by a member comment", func(t *testing.T) {
		issue := Types.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
		}}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Types.IssueTypeEnum.ANSWERED)
	})

	t.Run("returns NOT_ANSWERED for an issue created by a member with a member comment followed by an external comment", func(t *testing.T) {
		issue := Types.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
		}}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Types.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("returns NOT_ANSWERED for an issue created by a non-member with no comments", func(t *testing.T) {
		issue := Types.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
		}}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Types.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("returns ANSWERED for an issue created by a non-member with a member comment", func(t *testing.T) {
		issue := Types.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
		}}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Types.IssueTypeEnum.ANSWERED)
	})

	t.Run("returns NOT_ANSWERED for an issue created by a non-member with an external comment", func(t *testing.T) {
		issue := Types.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
		}}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Types.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("returns ANSWERED for an issue created by a non-member with an external comment followed by a non-member comment", func(t *testing.T) {
		issue := Types.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
		}}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Types.IssueTypeEnum.ANSWERED)
	})

	t.Run("returns NOT_ANSWERED for an issue created by a non-member with a member comment followed by an external comment", func(t *testing.T) {
		issue := Types.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
		}}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Types.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("excludes an issuehunt-app comment when the last comment is ours", func(t *testing.T) {
		issue := Types.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "MEMBER", Author: Types.CommentAuthor{Login: "user"}}},
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
		}}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Types.IssueTypeEnum.ANSWERED)
	})

	t.Run("excludes an issuehunt-app comment when the last comment isn't ours", func(t *testing.T) {
		issue := Types.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: Types.Labels{Edges: []Types.LabelEdge{}}, Comments: Types.Comments{Edges: []Types.CommentEdge{
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "user"}}},
			Types.CommentEdge{Node: Types.Comment{AuthorAssociation: "NONE", Author: Types.CommentAuthor{Login: "issuehunt-app"}}},
		}}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Types.IssueTypeEnum.NOT_ANSWERED)
	})
}
