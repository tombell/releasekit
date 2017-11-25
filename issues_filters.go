package releasekit

import (
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/google/go-github/github"
)

const (
	closedByPullRequestRegex = `(?i)(close|closes|closed|resolve|resolves|resolved|fix|fixes|fixed) #([0-9]+)`
	mergedPullRequestRegex   = `(?i)merge pull request #([0-9]+)|\(#([0-9]+)\)`
)

// FilterClosedBefore filters out all issues that were closed after the
// specified time.
func FilterClosedBefore(issues []*github.Issue, d time.Time) []*github.Issue {
	var filtered []*github.Issue

	for _, issue := range issues {
		if issue.ClosedAt.After(d) {
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

// FilterClosedAfter filters out all issues that were closed before the
// specified time.
func FilterClosedAfter(issues []*github.Issue, d time.Time) []*github.Issue {
	var filtered []*github.Issue

	// add some leniency for difference between tagging/merging pull request
	// creating the commit
	time := d.Add(2 * time.Second)

	for _, issue := range issues {
		if issue.ClosedAt.Before(time) {
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

// FilterClosedByPull filters out all issues that were closed automatically
// by a pull request.
func FilterClosedByPull(issues []*github.Issue) []*github.Issue {
	r, _ := regexp.Compile(closedByPullRequestRegex)

	var ignore []int

	for _, issue := range issues {
		if !issue.IsPullRequest() {
			continue
		}

		matches := r.FindStringSubmatch(*issue.Body)
		if matches != nil {
			num, _ := strconv.Atoi(matches[len(matches)-1])
			ignore = append(ignore, num)
		}
	}

	var filtered []*github.Issue

	for _, issue := range issues {
		if !contains(ignore, *issue.Number) {
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

// FilterNonMergedPulls filters out all pull requests that were closed and not merged.
func FilterNonMergedPulls(issues []*github.Issue, c *github.Client, owner, repo string) []*github.Issue {
	var ignore []int

	for _, issue := range issues {
		if !issue.IsPullRequest() {
			continue
		}

		pr, err := GetPullRequest(c, owner, repo, *issue.Number)
		if err != nil {
			log.Fatal(err)
		}

		if !*pr.Merged {
			ignore = append(ignore, *issue.Number)
		}
	}

	var filtered []*github.Issue

	for _, issue := range issues {
		if !contains(ignore, *issue.Number) {
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

// FilterClosedByCommits filters out any issues that have not been closed
// by commit messages.
func FilterClosedByCommits(issues []*github.Issue, commits []github.RepositoryCommit) []*github.Issue {
	r, _ := regexp.Compile(closedByPullRequestRegex)

	var ignore []int

	for _, c := range commits {
		matches := r.FindStringSubmatch(*c.Commit.Message)
		if matches != nil {
			num, _ := strconv.Atoi(matches[len(matches)-1])
			ignore = append(ignore, num)
		}
	}

	var filtered []*github.Issue

	for _, issue := range issues {
		if contains(ignore, *issue.Number) || issue.IsPullRequest() {
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

// FilterMergedPullsAfter filters out any issues or pull requests closed
// outside of the commit comparison range.
func FilterMergedPullsAfter(issues []*github.Issue, commits []github.RepositoryCommit) []*github.Issue {
	r, _ := regexp.Compile(mergedPullRequestRegex)

	var prs []int

	for _, issue := range issues {
		if issue.IsPullRequest() {
			prs = append(prs, *issue.Number)
		}
	}

	var merged []int

	for _, c := range commits {
		matches := r.FindStringSubmatch(*c.Commit.Message)
		if matches != nil {
			var pr string

			if matches[1] != "" {
				pr = matches[1]
			} else if matches[2] != "" {
				pr = matches[2]
			}

			num, _ := strconv.Atoi(pr)
			merged = append(merged, num)
		}
	}

	var include []int

	for _, id := range prs {
		if contains(merged, id) {
			include = append(include, id)
		}
	}

	var filtered []*github.Issue

	for _, issue := range issues {
		if contains(include, *issue.Number) || !issue.IsPullRequest() {
			filtered = append(filtered, issue)
		}
	}

	return filtered
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
