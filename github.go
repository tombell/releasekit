package releasekit

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// CreateGitHubClient creates a new GitHub API client with the specified access
// token for authentication.
func CreateGitHubClient(token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}
