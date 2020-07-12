package issuestriage

import (
	"github.com/brainhubeu/issue-overseer/githubstructures"
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
	if status != 0 {
		os.Exit(status)
	}
}

func TestIssuesTriage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "issuestriage")
}

var _ = Describe("issuestriage", func() {
	_ = Describe("GroupByAnswering", func() {
		It("triages an empty list", func() {
			issues := []githubstructures.Issue{}

			issuesTriage := New()
			ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.GroupByAnswering(issues)

			Expect(ourIssues).To(Equal([]githubstructures.Issue{}))
			Expect(answeredIssues).To(Equal([]githubstructures.Issue{}))
			Expect(notAnsweredIssues).To(Equal([]githubstructures.Issue{}))
		})

		It("triages a non-empty list", func() {
			issues := []githubstructures.Issue{
				githubstructures.Issue{Title: "title", Url: "url", Number: 122, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
					githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 124, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
					githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 125, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
					githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
					githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 126, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
					githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 127, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{}},
			}

			issuesTriage := New()
			ourIssues, answeredIssues, notAnsweredIssues := issuesTriage.GroupByAnswering(issues)

			Expect(ourIssues).To(Equal([]githubstructures.Issue{
				githubstructures.Issue{Title: "title", Url: "url", Number: 122, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
					githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				}},
			}))
			Expect(answeredIssues).To(Equal([]githubstructures.Issue{
				githubstructures.Issue{Title: "title", Url: "url", Number: 125, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
					githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
					githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 126, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
					githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				}},
			}))
			Expect(notAnsweredIssues).To(Equal([]githubstructures.Issue{
				githubstructures.Issue{Title: "title", Url: "url", Number: 124, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
					githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 127, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{}},
			}))
		})
	})

	_ = Describe("TriageOneIssueByAnswering", func() {
		It("returns OURS for an issue created by a member with no comments", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.OURS))
		})

		It("returns OURS for an issue created by a member with a member comment", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.OURS))
		})

		It("returns NOT_ANSWERED for an issue created by a member with an external comment", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.NOT_ANSWERED))
		})

		It("returns ANSWERED for an issue created by a member with an external comment followed by a member comment", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.ANSWERED))
		})

		It("returns NOT_ANSWERED for an issue created by a member with a member comment followed by an external comment", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.NOT_ANSWERED))
		})

		It("returns NOT_ANSWERED for an issue created by a non-member with no comments", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.NOT_ANSWERED))
		})

		It("returns ANSWERED for an issue created by a non-member with a member comment", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.ANSWERED))
		})

		It("returns NOT_ANSWERED for an issue created by a non-member with an external comment", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.NOT_ANSWERED))
		})

		It("returns ANSWERED for an issue created by a non-member with an external comment followed by a non-member comment", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.ANSWERED))
		})

		It("returns NOT_ANSWERED for an issue created by a non-member with a member comment followed by an external comment", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.NOT_ANSWERED))
		})

		It("excludes an issuehunt-app comment for our issue when there are no other comments", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "MEMBER", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.OURS))
		})

		It("excludes an issuehunt-app comment for an external issue when there are no other comments", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.NOT_ANSWERED))
		})

		It("excludes an issuehunt-app comment when the last comment is ours", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.ANSWERED))
		})

		It("excludes an issuehunt-app comment when the last comment isn't ours", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.NOT_ANSWERED))
		})

		It("works correctly with many comments when the last comment is ours", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.ANSWERED))
		})

		It("works correctly with many comments when the last comment isn't ours", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.NOT_ANSWERED))
		})

		It("works correctly with many comments when the last comment is by issuehunt-app and before it there's our comment", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.ANSWERED))
		})

		It("works correctly with many comments when the last comment is by issuehunt-app and before it there's an external comment", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "NONE", AuthorLogin: "user"},
				githubstructures.Comment{AuthorAssociation: "MEMBER", AuthorLogin: "issuehunt-app"},
			}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByAnswering(issue)

			Expect(issueType).To(Equal(githubstructures.IssueAnsweringTypeEnum.NOT_ANSWERED))
		})
	})

	_ = Describe("TriageOneIssueByManualLabel", func() {
		It("returns NON_EXISTENT when there are no labels", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{}, Comments: []githubstructures.Comment{}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByManualLabel(issue, githubstructures.ManualLabelConfig{Prefix: "severity", ParentLabelName: ""})

			Expect(issueType).To(Equal(githubstructures.IssueManualLabelTypeEnum.NON_EXISTENT))
		})

		It("returns NON_EXISTENT when there are is one non-matching label", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
				githubstructures.Label{Name: "foo", Color: "000000"},
			}, Comments: []githubstructures.Comment{}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByManualLabel(issue, githubstructures.ManualLabelConfig{Prefix: "severity", ParentLabelName: ""})

			Expect(issueType).To(Equal(githubstructures.IssueManualLabelTypeEnum.NON_EXISTENT))
		})

		It("returns NON_EXISTENT when there are many non-matching labels", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
				githubstructures.Label{Name: "foo", Color: "000000"},
				githubstructures.Label{Name: "bar", Color: "000000"},
				githubstructures.Label{Name: "baz", Color: "000000"},
			}, Comments: []githubstructures.Comment{}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByManualLabel(issue, githubstructures.ManualLabelConfig{Prefix: "severity", ParentLabelName: ""})

			Expect(issueType).To(Equal(githubstructures.IssueManualLabelTypeEnum.NON_EXISTENT))
		})

		It("returns EXISTENT when there is a matching label among many labels", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
				githubstructures.Label{Name: "foo", Color: "000000"},
				githubstructures.Label{Name: "severity: major", Color: "000000"},
				githubstructures.Label{Name: "baz", Color: "000000"},
			}, Comments: []githubstructures.Comment{}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByManualLabel(issue, githubstructures.ManualLabelConfig{Prefix: "severity", ParentLabelName: ""})

			Expect(issueType).To(Equal(githubstructures.IssueManualLabelTypeEnum.EXISTENT))
		})

		It("returns EXISTENT when there is only one label and it's matching", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
				githubstructures.Label{Name: "severity: major", Color: "000000"},
			}, Comments: []githubstructures.Comment{}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByManualLabel(issue, githubstructures.ManualLabelConfig{Prefix: "severity", ParentLabelName: ""})

			Expect(issueType).To(Equal(githubstructures.IssueManualLabelTypeEnum.EXISTENT))
		})

		It("returns NON_EXISTENT for exact name", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
				githubstructures.Label{Name: "severity", Color: "000000"},
			}, Comments: []githubstructures.Comment{}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByManualLabel(issue, githubstructures.ManualLabelConfig{Prefix: "severity", ParentLabelName: ""})

			Expect(issueType).To(Equal(githubstructures.IssueManualLabelTypeEnum.NON_EXISTENT))
		})

		It("returns NON_EXISTENT for the prefix following something else than the semicolon and space", func() {
			issue := githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
				githubstructures.Label{Name: "severity-", Color: "000000"},
			}, Comments: []githubstructures.Comment{}}

			issuesTriage := New()
			issueType := issuesTriage.TriageOneIssueByManualLabel(issue, githubstructures.ManualLabelConfig{Prefix: "severity", ParentLabelName: ""})

			Expect(issueType).To(Equal(githubstructures.IssueManualLabelTypeEnum.NON_EXISTENT))
		})
	})

	_ = Describe("GroupByManualLabel", func() {
		It("groups", func() {
			issues := []githubstructures.Issue{
				githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
					githubstructures.Label{Name: "severity: major", Color: "000000"},
				}, Comments: []githubstructures.Comment{}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 122, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
					githubstructures.Label{Name: "foo", Color: "000000"},
				}, Comments: []githubstructures.Comment{}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
					githubstructures.Label{Name: "severity: minor", Color: "000000"},
				}, Comments: []githubstructures.Comment{}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 124, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
					githubstructures.Label{Name: "bar", Color: "000000"},
				}, Comments: []githubstructures.Comment{}},
			}

			issuesTriage := New()
			issuesWithLabel, issuesWithoutLabel := issuesTriage.GroupByManualLabel(issues, githubstructures.ManualLabelConfig{Prefix: "severity", ParentLabelName: ""})

			Expect(issuesWithLabel).To(Equal([]githubstructures.Issue{
				githubstructures.Issue{Title: "title", Url: "url", Number: 121, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
					githubstructures.Label{Name: "severity: major", Color: "000000"},
				}, Comments: []githubstructures.Comment{}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 123, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
					githubstructures.Label{Name: "severity: minor", Color: "000000"},
				}, Comments: []githubstructures.Comment{}},
			}))
			Expect(issuesWithoutLabel).To(Equal([]githubstructures.Issue{
				githubstructures.Issue{Title: "title", Url: "url", Number: 122, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
					githubstructures.Label{Name: "foo", Color: "000000"},
				}, Comments: []githubstructures.Comment{}},
				githubstructures.Issue{Title: "title", Url: "url", Number: 124, AuthorAssociation: "NONE", Labels: []githubstructures.Label{
					githubstructures.Label{Name: "bar", Color: "000000"},
				}, Comments: []githubstructures.Comment{}},
			}))
		})
	})
})
