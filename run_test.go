package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDoIssuesTriage(t *testing.T) {
	t.Run("does triage of an empty list", func(t *testing.T) {
		issues := []Issue{}

		ourIssues, answeredIssues, notAnsweredIssues := doIssuesTriage(issues)

		assert.Equal(t, ourIssues, []Issue{})
		assert.Equal(t, answeredIssues, []Issue{})
		assert.Equal(t, notAnsweredIssues, []Issue{})
	})

	t.Run("does triage of issues with no comments", func(t *testing.T) {
		issues := []Issue{
			Issue{"title", "url", "number", "MEMBER", Labels{[]LabelEdge{}}, Comments{[]CommentEdge{}}},
			Issue{"title", "url", "number", "NONE", Labels{[]LabelEdge{}}, Comments{[]CommentEdge{}}},
		}

		ourIssues, answeredIssues, notAnsweredIssues := doIssuesTriage(issues)

		assert.Equal(t, ourIssues, []Issue{
			Issue{"title", "url", "number", "MEMBER", Labels{[]LabelEdge{}}, Comments{[]CommentEdge{}}},
			Issue{"title", "url", "number", "NONE", Labels{[]LabelEdge{}}, Comments{[]CommentEdge{}}},
		})
		assert.Equal(t, answeredIssues, []Issue{})
		assert.Equal(t, notAnsweredIssues, []Issue{
		})
	})
}
