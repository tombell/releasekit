package releasekit

import (
	"context"

	"github.com/google/go-github/github"
)

// GetCommitForTag gets the commit a tag is a reference to.
func GetCommitForTag(c *github.Client, owner, repo, tag string) (*github.RepositoryCommit, error) {
	ref, _, err := c.Git.GetRef(context.Background(), owner, repo, "refs/tags/"+tag)
	if err != nil {
		return nil, err
	}

	sha := *ref.Object.SHA

	if *ref.Object.Type == "tag" {
		tag, _, err := c.Git.GetTag(context.Background(), owner, repo, *ref.Object.SHA)
		if err != nil {
			return nil, err
		}

		sha = *tag.Object.SHA
	}

	commit, _, err := c.Repositories.GetCommit(context.Background(), owner, repo, sha)
	if err != nil {
		return nil, err
	}

	return commit, nil
}

// GetFirstCommit gets the first commit to the repository.
func GetFirstCommit(c *github.Client, owner, repo string) (*github.RepositoryCommit, error) {
	opts := &github.CommitsListOptions{}

	commits, resp, err := c.Repositories.ListCommits(context.Background(), owner, repo, opts)
	if err != nil {
		return nil, err
	}

	if resp.NextPage == 0 {
		return commits[len(commits)-1], nil
	}

	opts.Page = resp.LastPage

	commits, _, err = c.Repositories.ListCommits(context.Background(), owner, repo, opts)
	if err != nil {
		return nil, err
	}

	return commits[len(commits)-1], nil
}

// GetComparison gets the commit comparison for the given base and head range.
func GetComparison(c *github.Client, owner, repo, base, head string) (*github.CommitsComparison, error) {
	comparison, _, err := c.Repositories.CompareCommits(context.Background(), owner, repo, base, head)
	return comparison, err
}
