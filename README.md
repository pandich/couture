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

For usage run `couture --help`. For shell completions run `eval $(couture complete)`.

### Configuration

Configure Couture by creating a `.couture.yaml` file in `$HOME`. Additionally, each directory can have a configuration
file which is consulted prior to consulting the one in `$HOME`.

Available settings:

| Field | Description | Example |
| ----: | :---------- | :------ |
| `aliases` | Maps short alias names to full source URLs | `aliases:{work: 'es+http://your-server:9200/some_index/_search'}` |
| `paginator` | The paginator use. Can also be set via the `COUTURE_PAGER` environment variable. | `paginator: less` |

---

_Project Layout attempts to conform to the
[Standard Go Project Layout](https://github.com/golang-standards/project-layout)_
