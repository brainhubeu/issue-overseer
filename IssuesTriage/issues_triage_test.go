package IssuesTriage

import (
	"../Interfaces"
	"testing"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestTriageManyIssues(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "IssuesTriage")
}

var _ = Describe("IssuesTriage", func() {
	It("triages an empty list", func() {
		issues := []Interfaces.Issue{}

		issuesTriage := InitIssuesTriage()
		ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.TriageManyIssues(issues)

		Expect(ourIssues).To(Equal([]Interfaces.Issue{}))
		Expect(answeredIssues).To(Equal([]Interfaces.Issue{}))
		Expect(notAnsweredIssues).To(Equal([]Interfaces.Issue{}))
	})

	It("triages a non-empty list", func() {
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

		Expect(ourIssues).To(Equal([]Interfaces.Issue{
			Interfaces.Issue{Title: "title", Url: "url", Number: 122, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{}},
			Interfaces.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
				Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
		}))
		Expect(answeredIssues).To(Equal([]Interfaces.Issue{
			Interfaces.Issue{Title: "title", Url: "url", Number: 125, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
				Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
			Interfaces.Issue{Title: "title", Url: "url", Number: 126, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
				Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
		}))
		Expect(notAnsweredIssues).To(Equal([]Interfaces.Issue{
			Interfaces.Issue{Title: "title", Url: "url", Number: 124, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
				Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			}},
			Interfaces.Issue{Title: "title", Url: "url", Number: 127, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{}},
		}))
	})

	It("returns OURS for an issue created by a member with no comments", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.OURS))
	})

	It("returns OURS for an issue created by a member with a member comment", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.OURS))
	})

	It("returns NOT_ANSWERED for an issue created by a member with an external comment", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("returns ANSWERED for an issue created by a member with an external comment followed by a member comment", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.ANSWERED))
	})

	It("returns NOT_ANSWERED for an issue created by a member with a member comment followed by an external comment", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("returns NOT_ANSWERED for an issue created by a non-member with no comments", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("returns ANSWERED for an issue created by a non-member with a member comment", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.ANSWERED))
	})

	It("returns NOT_ANSWERED for an issue created by a non-member with an external comment", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("returns ANSWERED for an issue created by a non-member with an external comment followed by a non-member comment", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.ANSWERED))
	})

	It("returns NOT_ANSWERED for an issue created by a non-member with a member comment followed by an external comment", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("excludes an issuehunt-app comment for our issue when there are no other comments", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.OURS))
	})

	It("excludes an issuehunt-app comment for an external issue when there are no other comments", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("excludes an issuehunt-app comment when the last comment is ours", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.ANSWERED))
	})

	It("excludes an issuehunt-app comment when the last comment isn't ours", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("works correctly with many comments when the last comment is ours", func() {
		issue := Interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []Interfaces.Label{}, Comments: []Interfaces.Comment{
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			Interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := InitIssuesTriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.ANSWERED))
	})

	It("works correctly with many comments when the last comment isn't ours", func() {
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

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("works correctly with many comments when the last comment is by issuehunt-app and before it there's our comment", func() {
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

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.ANSWERED))
	})

	It("works correctly with many comments when the last comment is by issuehunt-app and before it there's an external comment", func() {
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

		Expect(issueType).To(Equal(Interfaces.IssueTypeEnum.NOT_ANSWERED))
	})
})
