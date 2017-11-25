package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/google/go-github/github"

	"github.com/tombell/releasekit"
)

// printVersion will print the version, and commit SHA for the build.
func printVersion() {
	fmt.Printf("releasekit %s (%s)\n", Version, Commit)
}

// printIfVerbose will print out the formatted string if the verbose flag is
// enabled.
func printIfVerbose(format string, a ...interface{}) {
	if !verbose {
		return
	}

	fmt.Printf(format, a...)
}

// exitIfError will fatal log if the error is not nil.
func exitIfError(err error) {
	if err == nil {
		return
	}

	log.Fatal(err)
}

// generateReleaseBody generates the release body from the slice of issues.
func generateReleaseBody(issues []*github.Issue, changed []string, compare string) string {
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

func main() {
	printVersion()
	parseFlags()

	client := releasekit.CreateGitHubClient(options.Token)

	var since time.Time

	if previous == "" || previous == next {
		printIfVerbose("Fetching first commit...\n")
		commit, err := releasekit.GetFirstCommit(client, owner, repo)
		exitIfError(err)
		sha := *commit.SHA
		previous = sha[:8]
	} else {
		printIfVerbose("Fetching commit for tag (%s)...\n", previous)
		base, err := releasekit.GetCommitForTag(client, owner, repo, previous)
		exitIfError(err)
		since = base.Commit.Author.Date.Add(-24 * time.Hour)
	}

	printIfVerbose("Fetching closed issues...\n")
	issues, err := releasekit.FetchClosedIssuesSince(client, owner, repo, since)
	exitIfError(err)

	printIfVerbose("Fetching commit for tag (%s)...\n", next)
	head, err := releasekit.GetCommitForTag(client, owner, repo, next)
	exitIfError(err)

	printIfVerbose("Fetching commit comparison (%s...%s)...\n", previous, next)
	comparison, err := releasekit.GetComparison(client, owner, repo, previous, next)
	exitIfError(err)

	if !since.IsZero() {
		printIfVerbose("Filtering out issues closed before %s...\n", since)
		issues = releasekit.FilterClosedBefore(issues, since)
	}

	printIfVerbose("Filtering out issues closed after %s...\n", *head.Commit.Author.Date)
	issues = releasekit.FilterClosedAfter(issues, *head.Commit.Author.Date)

	printIfVerbose("Filtering out issues closed by a pull request...\n")
	issues = releasekit.FilterClosedByPull(issues)

	printIfVerbose("Filtering out non-merged pull requests...\n")
	issues = releasekit.FilterNonMergedPulls(issues, client, owner, repo)

	if len(comparison.Commits) > 0 {
		printIfVerbose("Filtering out issues not closed by a commit...\n")
		issues = releasekit.FilterClosedByCommits(issues, comparison.Commits)

		printIfVerbose("Filtering out pull requests merged after tag (%s)...\n", next)
		issues = releasekit.FilterMergedPullsAfter(issues, comparison.Commits)
	}

	var changed []string

	if len(watched) > 0 {
		printIfVerbose("Checking for changes in watched files...\n")

		for _, file := range watched {
			name := filepath.Clean(file)

			for _, commitFile := range comparison.Files {
				if name == *commitFile.Filename {
					changed = append(changed, name)
				}
			}
		}
	}

	printIfVerbose("Generating release body...\n")
	body := generateReleaseBody(issues, changed, *comparison.HTMLURL)

	if options.Dry {
		fmt.Println()
		fmt.Println(body)
		return
	}

	printIfVerbose("Checking for existing release for tag (%s)...\n", next)
	release, err := releasekit.GetReleaseByTag(client, owner, repo, next)
	exitIfError(err)

	if release == nil {
		release = &github.RepositoryRelease{}
	}

	release.TagName = &next
	release.Name = &next
	release.Body = &body

	release.Draft = &draft
	release.Prerelease = &prerelease

	if release.ID != nil {
		fmt.Printf("Updating release (%s)...\n", *release.TagName)
	} else {
		fmt.Printf("Creating release (%s)...\n", *release.TagName)
	}

	release, err = releasekit.CreateOrEditRelease(client, owner, repo, release)
	exitIfError(err)

	if len(attachments) > 0 {
		printIfVerbose("Uploading release assets...\n")
		err = releasekit.UploadReleaseAssets(client, owner, repo, *release.ID, attachments)
		exitIfError(err)
	}

	fmt.Println(*release.HTMLURL)
}
