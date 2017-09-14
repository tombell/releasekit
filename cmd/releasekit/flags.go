package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

var options struct {
	Token string `short:"t" long:"token" description:"GitHub API token" required:"true" value-name:"TOKEN"`
	Owner string `short:"o" long:"owner" description:"GitHub repository owner" required:"true" value-name:"USER/ORG"`
	Repo  string `short:"r" long:"repo" description:"GitHub repository name" required:"true" value-name:"REPO"`

	Prev string `short:"p" long:"previous" description:"Previous release tag" value-name:"GIT_TAG"`
	Next string `short:"n" long:"next" description:"Next release tag" required:"true" value-name:"GIT_TAG"`

	Draft      bool `long:"draft" description:"Mark release as draft"`
	Prerelease bool `long:"prerelease" description:"Mark release as prerelease"`

	Attachments []string `long:"attachment" description:"File path to attach release asset" value-name:"FILE_PATH"`
	Watched     []string `long:"watch" description:"File path to watch for changes" value-name:"FILE_PATH"`

	Verbose bool `short:"v" long:"verbose" description:"Verbose debug output"`
	Print   bool `long:"print" description:"Print the release body"`
}

var (
	verbose     bool
	owner       string
	repo        string
	previous    string
	next        string
	draft       bool
	prerelease  bool
	attachments []string
	watched     []string
)

// parseFlags parses the command line flags.
func parseFlags() {
	parser := flags.NewParser(&options, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}

		os.Exit(1)
	}

	verbose = options.Verbose

	owner = options.Owner
	repo = options.Repo
	previous = options.Prev
	next = options.Next
	prerelease = options.Prerelease
	draft = options.Draft
	attachments = options.Attachments
	watched = options.Watched
}
