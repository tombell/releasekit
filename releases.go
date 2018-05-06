package releasekit

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/github"
)

// GenerateReleaseBody generates the body of the release notes.
func GenerateReleaseBody(issues []*github.Issue, changed []string, compare string) string {
	if len(issues) == 0 {
		return "New Release"
	}

	output := "## Changes\n"

	for _, issue := range issues {
		output += fmt.Sprintf("* [#%d](%s) - %v (@%v)\n", *issue.Number, *issue.HTMLURL, *issue.Title, *issue.User.Login)
	}

	if len(changed) > 0 {
		output += "\n### Watched File Changes\n"
		output += fmt.Sprintf("Changes: %s\n", compare)

		for _, file := range changed {
			output += fmt.Sprintf("* %s\n", file)
		}
	}

	return output
}

// GetReleaseByTag returns a repository release for the given tag if it exists.
func GetReleaseByTag(c *github.Client, owner, repo, tag string) (*github.RepositoryRelease, error) {
	release, res, err := c.Repositories.GetReleaseByTag(context.Background(), owner, repo, tag)
	if err != nil && res.StatusCode != http.StatusNotFound {
		return nil, err
	}

	return release, nil
}

// CreateOrEditRelease creates a repository release if it doesn't exist, else it
// will edit an existing repository release.
func CreateOrEditRelease(c *github.Client, owner, repo string, release *github.RepositoryRelease) (*github.RepositoryRelease, error) {
	var output *github.RepositoryRelease
	var err error

	if release.ID == nil {
		output, _, err = c.Repositories.CreateRelease(context.Background(), owner, repo, release)
	} else {
		output, _, err = c.Repositories.EditRelease(context.Background(), owner, repo, *release.ID, release)
	}

	if err != nil {
		return nil, err
	}

	return output, nil
}

// UploadReleaseAssets uploads the files to the release as assets.
func UploadReleaseAssets(c *github.Client, owner, repo string, id int, attachments []string) error {
	for _, attachment := range attachments {
		f, err := os.OpenFile(attachment, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		name := filepath.Clean(filepath.Base(f.Name()))
		opt := &github.UploadOptions{Name: name}

		_, _, err = c.Repositories.UploadReleaseAsset(context.Background(), owner, repo, id, opt, f)
		if err != nil {
			return err
		}
	}

	return nil
}
