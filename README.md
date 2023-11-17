# changelog

Parse markdown-esque changelogs (like our example in [History.markdown](History.markdown)), parse out versions, sections, changes & references.
Motivation: automate update of changelogs

Bundled with a command, `changelogger`.

Supports several header formats such as SemVer and [KeepAChangelog](https://keepachangelog.com).

[![Build & test](https://github.com/parkr/changelog/actions/workflows/push-build-test-on-push.yml/badge.svg)](https://github.com/parkr/changelog/actions/workflows/push-build-test-on-push.yml)

## `changelogger` command

### Installation

    $ go get github.com/parkr/changelog/changelogger

### Usage

    $ $GOPATH/bin/changelogger
    $ $GOPATH/bin/changelogger -h

## `changelog` package

### Installation

    $ go get github.com/parkr/changelog

### Usage

    // Parse changelog at a given filename
    changes, err := changelog.NewChangelogFromFile("CHANGELOG.md")

    // Discover the filename of your changelog
    filename := changelog.HistoryFilename()

    // Parse changelog from some io.Reader
    changes, err := changelog.NewChangeLogFromReader(req.Body)

## License

MIT License, Copyright 2015 Parker Moore. See [LICENSE](LICENSE) for details.
