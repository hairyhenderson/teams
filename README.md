[![Build Status][circleci-image]][circleci-url]
[![Go Report Card][reportcard-image]][reportcard-url]
[![Codebeat Status][codebeat-image]][codebeat-url]
[![Coverage][gocover-image]][gocover-url]
[![Total Downloads][gh-downloads-image]][gh-downloads-url]
[![CII Best Practices][cii-bp-image]][cii-bp-url]

# teams

CLI application to manage GitHub issues and pull requests for a team's repositories.

## Installing

1. Get the latest `teams` for your platform from the [releases](https://github.com/hairyhenderson/teams/releases) page
2. Store the downloaded binary somewhere in your path as `teams` (or `teams.exe`
  on Windows)
3. Make sure it's executable (on Linux/macOS)
3. Test it out with `teams --help`!

In other words:

```console
$ curl -o /usr/local/bin/teams https://github.com/hairyhenderson/teams/releases/download/<version>/teams_<os>-<arch>
$ chmod 755 /usr/local/bin/teams
$ teams --help
...
```

_Please report any bugs found in the [issue tracker](https://github.com/hairyhenderson/teams/issues/)._

## Usage

Set the `GITHUB_API_TOKEN` environment variable to a [Personal Access Token](https://github.com/settings/tokens) value. The scopes necessary are `read:org` and `repo`.


## Releasing

Right now the release process is semi-automatic.

1. Create a release tag: `git tag -a v0.0.9 -m "Releasing v0.9.9" && git push --tags`
2. Build binaries & compress most of them: `make build-release`
3. Create a release in [github](https://github.com/hairyhenderson/teams/releases)!

## License

[The MIT License](http://opensource.org/licenses/MIT)

Copyright (c) 2016 Dave Henderson

[circleci-image]: https://img.shields.io/circleci/project/hairyhenderson/teams.svg?style=flat
[circleci-url]: https://circleci.com/gh/hairyhenderson/teams
[reportcard-image]: https://goreportcard.com/badge/github.com/hairyhenderson/teams
[reportcard-url]: https://goreportcard.com/report/github.com/hairyhenderson/teams
[codebeat-image]: https://codebeat.co/badges/55ea20ac-5ff5-4c81-8fd9-8f8c4f18b821
[codebeat-url]: https://codebeat.co/projects/github-com-hairyhenderson-teams
[gocover-image]: https://gocover.io/_badge/github.com/hairyhenderson/teams
[gocover-url]: https://gocover.io/github.com/hairyhenderson/teams
[gh-downloads-image]: https://img.shields.io/github/downloads/hairyhenderson/teams/total.svg
[gh-downloads-url]: https://github.com/hairyhenderson/teams/releases

[cii-bp-image]: https://bestpractices.coreinfrastructure.org/projects/379/badge
[cii-bp-url]: https://bestpractices.coreinfrastructure.org/projects/379

[![Analytics](https://ga-beacon.appspot.com/UA-82637990-1/teams/README.md?pixel)](https://github.com/igrigorik/ga-beacon)
