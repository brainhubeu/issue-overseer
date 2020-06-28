package issuestriage

import (
	"github.com/brainhubeu/issue-overseer/interfaces"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	status := m.Run()
	if status == 0 && testing.CoverMode() != "" {
		coverage := testing.Coverage()
		requiredCoverage := 1.0
		if coverage < requiredCoverage {
			log.Println("too low tests coverage:", coverage, ", should be at least", requiredCoverage)
			os.Exit(1)
		}
	}
}

func TestTriageManyIssues(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "issuestriage")
}

var _ = Describe("issuestriage", func() {
	It("triages an empty list", func() {
		issues := []interfaces.Issue{}

		issuesTriage := Initissuestriage()
		ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.TriageManyIssues(issues)

		Expect(ourIssues).To(Equal([]interfaces.Issue{}))
		Expect(answeredIssues).To(Equal([]interfaces.Issue{}))
		Expect(notAnsweredIssues).To(Equal([]interfaces.Issue{}))
	})

	It("triages a non-empty list", func() {
		issues := []interfaces.Issue{
			interfaces.Issue{Title: "title", Url: "url", Number: 122, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{}},
			interfaces.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
				interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
			interfaces.Issue{Title: "title", Url: "url", Number: 124, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
				interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			}},
			interfaces.Issue{Title: "title", Url: "url", Number: 125, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
				interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
			interfaces.Issue{Title: "title", Url: "url", Number: 126, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
				interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
			interfaces.Issue{Title: "title", Url: "url", Number: 127, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{}},
		}

		issuesTriage := Initissuestriage()
		ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.TriageManyIssues(issues)

		Expect(ourIssues).To(Equal([]interfaces.Issue{
			interfaces.Issue{Title: "title", Url: "url", Number: 122, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{}},
			interfaces.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
				interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
		}))
		Expect(answeredIssues).To(Equal([]interfaces.Issue{
			interfaces.Issue{Title: "title", Url: "url", Number: 125, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
				interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
			interfaces.Issue{Title: "title", Url: "url", Number: 126, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
				interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}},
		}))
		Expect(notAnsweredIssues).To(Equal([]interfaces.Issue{
			interfaces.Issue{Title: "title", Url: "url", Number: 124, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
				interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			}},
			interfaces.Issue{Title: "title", Url: "url", Number: 127, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{}},
		}))
	})

	It("returns OURS for an issue created by a member with no comments", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.OURS))
	})

	It("returns OURS for an issue created by a member with a member comment", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.OURS))
	})

	It("returns NOT_ANSWERED for an issue created by a member with an external comment", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("returns ANSWERED for an issue created by a member with an external comment followed by a member comment", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.ANSWERED))
	})

	It("returns NOT_ANSWERED for an issue created by a member with a member comment followed by an external comment", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("returns NOT_ANSWERED for an issue created by a non-member with no comments", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("returns ANSWERED for an issue created by a non-member with a member comment", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.ANSWERED))
	})

	It("returns NOT_ANSWERED for an issue created by a non-member with an external comment", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("returns ANSWERED for an issue created by a non-member with an external comment followed by a non-member comment", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.ANSWERED))
	})

	It("returns NOT_ANSWERED for an issue created by a non-member with a member comment followed by an external comment", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("excludes an issuehunt-app comment for our issue when there are no other comments", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.OURS))
	})

	It("excludes an issuehunt-app comment for an external issue when there are no other comments", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("excludes an issuehunt-app comment when the last comment is ours", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.ANSWERED))
	})

	It("excludes an issuehunt-app comment when the last comment isn't ours", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("works correctly with many comments when the last comment is ours", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.ANSWERED))
	})

	It("works correctly with many comments when the last comment isn't ours", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.NOT_ANSWERED))
	})

	It("works correctly with many comments when the last comment is by issuehunt-app and before it there's our comment", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.ANSWERED))
	})

	It("works correctly with many comments when the last comment is by issuehunt-app and before it there's an external comment", func() {
		issue := interfaces.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []interfaces.Label{}, Comments: []interfaces.Comment{
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			interfaces.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
		}}

		issuesTriage := Initissuestriage()
		issueType := issuesTriage.TriageOneIssue(issue)

		Expect(issueType).To(Equal(interfaces.IssueTypeEnum.NOT_ANSWERED))
	})
})
