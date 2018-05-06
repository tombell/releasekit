package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/google/go-github/github"

	"github.com/tombell/releasekit"
)

func printVersion() {
	fmt.Printf("releasekit %s (%s)\n", Version, Commit)
}

func printIfVerbose(format string, a ...interface{}) {
	if !verbose {
		return
	}

	fmt.Printf(format, a...)
}

func exitIfError(err error, msg string) {
	if err == nil {
		return
	}

	log.Fatal(fmt.Sprintf("%s:\n%s", msg, err))
}

func main() {
	printVersion()
	parseFlags()

	client := releasekit.CreateGitHubClient(options.Token)

	var since time.Time

	if previous == "" || previous == next {
		printIfVerbose("Fetching first commit...\n")
		commit, err := releasekit.GetFirstCommit(client, owner, repo)
		exitIfError(err, "Could not fetch first commit")

		sha := *commit.SHA
		previous = sha[:8]
	} else {
		printIfVerbose("Fetching commit for tag (%s)...\n", previous)
		base, err := releasekit.GetCommitForTag(client, owner, repo, previous)
		exitIfError(err, "Could not fetch commit for tag")

		since = base.Commit.Author.Date.Add(-24 * time.Hour)
	}

	printIfVerbose("Fetching closed issues...\n")
	issues, err := releasekit.FetchClosedIssuesSince(client, owner, repo, since)
	exitIfError(err, "Could not fetch closed issues")

	printIfVerbose("Fetching commit for tag (%s)...\n", next)
	head, err := releasekit.GetCommitForTag(client, owner, repo, next)
	exitIfError(err, "Could not fetch commit for tag")

	printIfVerbose("Fetching commit comparison (%s...%s)...\n", previous, next)
	comparison, err := releasekit.GetComparison(client, owner, repo, previous, next)
	exitIfError(err, "Could not fetch commit comparison")

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
	body := releasekit.GenerateReleaseBody(issues, changed, *comparison.HTMLURL)

	if options.Dry {
		fmt.Println()
		fmt.Println(body)
		return
	}

	printIfVerbose("Checking for existing release for tag (%s)...\n", next)
	release, err := releasekit.GetReleaseByTag(client, owner, repo, next)
	exitIfError(err, "Could not check for existing release")

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
	exitIfError(err, "Could not create or update release")

	if len(attachments) > 0 {
		printIfVerbose("Uploading release assets...\n")
		err = releasekit.UploadReleaseAssets(client, owner, repo, *release.ID, attachments)
		exitIfError(err, "Could not upload release assets")
	}

	fmt.Println(*release.HTMLURL)
}
