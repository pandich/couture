![Couture](docs/couture.png)

[![Go Report Card](http://goreportcard.com/badge/github.com/pandich/couture)](https://goreportcard.com/badge/github.com/pandich/couture)
[![goreleaser](http://github.com/pandich/couture/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/pandich/couture/actions/workflows/goreleaser.yml)
[![CodeQL](https://github.com/pandich/couture/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/pandich/couture/actions/workflows/codeql-analysis.yml)

---

Couture combines multiple log streams from disparate sources into human-friendly output.

### Installation
	
	go get -u github.com/pandich/couture	

### Usage:

	couture --help

### Configuration

[comment]: <> (TODO config doc)

---

### Examples:

#### Single-line

	couture \
		--highlight --filter=+distincto --filter=+'\"first_name\"\s*:\s*\"B' --filter=+quinoa \ 
		fake://?style=hacker \
		fake://?style=lorem \
		fake://?style=hipster

![Couture](docs/example-fake-single-line.gif)

#### Multi-line

	couture --expand --multiline \
		--highlight --filter=+distincto --filter=+'\"first_name\"\s*:\s*\"B' --filter=+quinoa \ 
		fake://?style=hacker \
		fake://?style=lorem \
		fake://?style=hipster

![Couture](docs/example-fake-multi-line.gif)
