package releasekit

import (
	"context"

	"github.com/google/go-github/github"
)

// GetPullRequest gets the pull request with the specified number.
func GetPullRequest(c *github.Client, owner, repo string, number int) (*github.PullRequest, error) {
	pr, _, err := c.PullRequests.Get(context.Background(), owner, repo, number)
	if err != nil {
		return nil, err
	}

	return pr, nil
}
