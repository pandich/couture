![Couture](docs/couture.png)

Couture combines multiple log streams from disparate sources into friendly output.

[comment]: <> (TODO example output - uses asciinema)

[comment]: <> (TODO working badges)
[![Build Status](https://travis-ci.org/gaggle-net/couture.svg?branch=master)]()
[![Go Report Card](https://goreportcard.com/badge/github.com/gaggle-net/couture)]()
[![Coverage Status](https://coveralls.io/repos/github/gaggle-net/couture/badge.svg?branch=master)]()

### Installation

| Tool | Command | 
| ---: | :------ |
| `go`                          | `go get -u github.com/gaggle-net/couture` |
| [Homebrew](https://brew.sh/)  | `brew ...` |
| `make`                        | `make install` |

### Usage:

For usage run `couture --help`. For shell completions run `eval $(couture shell-completion)`.

### Configuration

Configure Couture by creating a `.couture.yaml` file in `$HOME`. Additionally, each directory can have a configuration
file which is consulted prior to consulting the one in `$HOME`.

Available settings:

_describe alias template behavior_

---

_Project Layout attempts to conform to the
[Standard Go Project Layout](https://github.com/golang-standards/project-layout)_
