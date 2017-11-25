package releasekit

import (
	"context"
	"time"

	"github.com/google/go-github/github"
)

// FetchClosedIssuesSince fetches all closed issues since the specified time.
func FetchClosedIssuesSince(c *github.Client, owner, repo string, since time.Time) ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{
		State: "closed",
		Since: since,
	}

	var allIssues []*github.Issue

	for {
		issues, resp, err := c.Issues.ListByRepo(context.Background(), owner, repo, opt)
		if err != nil {
			return nil, err
		}

		allIssues = append(allIssues, issues...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allIssues, nil
}
