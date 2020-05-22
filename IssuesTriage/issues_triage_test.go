package IssuesTriage

import (
	"../Interfaces"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTriageManyIssues(t *testing.T) {
	t.Run("triages an empty list", func(t *testing.T) {
		issues := []Interfaces.Issue{}

		issuesTriage := InitIssuesTriage()
		ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.TriageManyIssues(issues)

		assert.Equal(t, ourIssues, []Interfaces.Issue{})
		assert.Equal(t, answeredIssues, []Interfaces.Issue{})
		assert.Equal(t, notAnsweredIssues, []Interfaces.Issue{})
	})

	t.Run("triages a non-empty list", func(t *testing.T) {
		issues := []Interfaces.Issue{
			Interfaces.Issue{Title: "title", Url: "url", Number: 122, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{}},
			Interfaces.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
				Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
			Interfaces.Issue{Title: "title", Url: "url", Number: 124, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
				Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			}},
			Interfaces.Issue{Title: "title", Url: "url", Number: 125, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
				Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
			Interfaces.Issue{Title: "title", Url: "url", Number: 126, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
				Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
			Interfaces.Issue{Title: "title", Url: "url", Number: 127, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{}},
		}

		issuesTriage := InitIssuesTriage()
		ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.TriageManyIssues(issues)

		assert.Equal(t, ourIssues, []Interfaces.Issue{
			Interfaces.Issue{Title: "title", Url: "url", Number: 122, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{}},
			Interfaces.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
				Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
		})
		assert.Equal(t, answeredIssues, []Interfaces.Issue{
			Interfaces.Issue{Title: "title", Url: "url", Number: 125, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
				Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
			Interfaces.Issue{Title: "title", Url: "url", Number: 126, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
				Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
		})
		assert.Equal(t, notAnsweredIssues, []Interfaces.Issue{
			Interfaces.Issue{Title: "title", Url: "url", Number: 124, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
				Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			}},
			Interfaces.Issue{Title: "title", Url: "url", Number: 127, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{}},
		})
	})

	t.Run("returns OURS for an issue created by a member with no comments", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.OURS)
	})

	t.Run("returns OURS for an issue created by a member with a member comment", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.OURS)
	})

	t.Run("returns NOT_ANSWERED for an issue created by a member with an external comment", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("returns ANSWERED for an issue created by a member with an external comment followed by a member comment", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.ANSWERED)
	})

	t.Run("returns NOT_ANSWERED for an issue created by a member with a member comment followed by an external comment", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("returns NOT_ANSWERED for an issue created by a non-member with no comments", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("returns ANSWERED for an issue created by a non-member with a member comment", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.ANSWERED)
	})

	t.Run("returns NOT_ANSWERED for an issue created by a non-member with an external comment", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("returns ANSWERED for an issue created by a non-member with an external comment followed by a non-member comment", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.ANSWERED)
	})

	t.Run("returns NOT_ANSWERED for an issue created by a non-member with a member comment followed by an external comment", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("excludes an issuehunt-app comment for our issue when there are no other comments", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.OURS)
	})

	t.Run("excludes an issuehunt-app comment for an external issue when there are no other comments", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("excludes an issuehunt-app comment when the last comment is ours", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.ANSWERED)
	})

	t.Run("excludes an issuehunt-app comment when the last comment isn't ours", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("works correctly with many comments when the last comment is ours", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.ANSWERED)
	})

	t.Run("works correctly with many comments when the last comment isn't ours", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.NOT_ANSWERED)
	})

	t.Run("works correctly with many comments when the last comment is by issuehunt-app and before it there's our comment", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.ANSWERED)
	})

	t.Run("works correctly with many comments when the last comment is by issuehunt-app and before it there's an external comment", func(t *testing.T) {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		assert.Equal(t, issueType, Interfaces.IssueTypeEnum.NOT_ANSWERED)
	})
}
