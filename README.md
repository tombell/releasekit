# releasekit

[![CircleCI](https://circleci.com/gh/tombell/releasekit/tree/master.svg?style=svg)](https://circleci.com/gh/tombell/releasekit/tree/master)

Command-line app to create GitHub releases.

## Installation

To get the most up to date binaries, check [the releases][releases] for the
pre-built binary for your system.

You can also `go get` to install from source.

    go get -u github.com/tombell/releasekit/cmd/releasekit

On macOS you can use [Homebrew](https://brew.sh) to install.

    brew tap tombell/formulae && brew install tombell/formulae/releasekit

[releases]: https://github.com/tombell/releasekit/releases

## Usage

Use the `-h/--help` flag to see all the available flags when running
**releasekit**.

You will need a [GitHub API token][api-token] when running **releasekit**. It's
advised you create a token specifically for **releasekit**.

[api-token]: https://github.com/settings/tokens

### Generating a new Github API token

Visit the [Developer Settings](https://github.com/settings/tokens/new) and click to generate a new **Personal access token**. If your repository is private, select the `repo` scope, otherwise `public_repo` will do just fine.

### Creating a Release

To create a new release, the tag you want to cut a release for should be pushed
to GitHub.

    releasekit -t $GITHUB_TOKEN -o tombell -r releasekit -p v0.1.0 -n v0.2.0 --dry

This will print then release notes for `v0.2.0`, and will generate the notes
from closed issues and merged pull requests between `v0.1.0` and `v0.2.0`.

If you're happy with the release notes, you can rerun the command omitting the
`--dry` flag. This will go ahead and create the release on the GitHub
repository.

### Updating a Release

To update an existing release, you can rerun the command again, including any
additional flags.

    releasekit -t $GITHUB_TOKEN -o tombell -r releasekit -p v0.1.0 -n v0.2.0 --attachment docs/api.md

This will update the existing `v0.2.0` release created above, it will also
attach the `docs/api.md` file as a release asset when updating.

### Draft and Pre-Release Releases

To mark a release as a draft you can use the `--draft` flag. This will create
the release as a draft, and not associate it with a tag name.

    releasekit -t $GITHUB_TOKEN -o tombell -r releasekit -p v0.1.0 -n v0.2.0 --draft

There are some limitations when using the `--draft` flag, if you rerun the
command without the `--draft` flag, it will create a new release but not remove
the old draft.

To mark a release as a pre-release you can use the `--prerelease` flag. This will
create or update a release as a pre-release.

    releasekit -t $GITHUB_TOKEN -o tombell -r releasekit -p v0.1.0 -n v0.2.0 --prerelease

You can rerun the command without the `--prerelease` flag, to remove the
pre-release mark from the release.

### Including Labels

If you would like to highlight specific labels on pull requests/issues in the
release notes, you can use the `--label` flag.

    releasekit -t $GITHUB_TOKEN -o tombell -r releasekit -p v0.1.0 -n v0.2.0 --label bug

The release on GitHub would then have **bug** next to any item that had the
_bug_ label.

### Attaching Release Assets

When you create or update a release, you can attach any files as release assets
using the `--attachment` flag. This flag can be used multiple times to attach
multiple release assets.

    releasekit -t $GITHUB_TOKEN -o tombell -r releasekit -p v0.1.0 -n v0.2.0 --attachment docs/file1 --attachment docs/file2

The release on GitHub would then have **file1** and **file2** as assets
available to download.

### Watching Specific Files

If you would like to include in the release notes if a specific file has changed
in this release, you can use the `--watch` flag. This flag can be used multiple
times to watch multiple files.

    releasekit -t $GITHUB_TOKEN -o tombell -r releasekit -p v0.1.0 -n v0.2.0 --watch cmd/releasekit/main.go --watch releasekit.go

This will include an additional section at the bottom of the release listing
these files if they've changed, and a link to the compare page on GitHub.
