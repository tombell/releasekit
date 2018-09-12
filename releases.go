package releasekit

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/go-github/v18/github"
)

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
func UploadReleaseAssets(c *github.Client, owner, repo string, id int64, attachments []string) error {
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
