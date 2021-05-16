# Couture

![Couture](https://static.thenounproject.com/png/566246-200.png)

## Overview

Allows for tailing multiple of event sources.

### Execution

#### Help

    couture --help

### Shell Completion

    couture __complete

OR

    couture __complete (bash|fish|powershell|zsh)

Add `eval $(couture __complete)` to the end of your shell init script to enable the completions.

## Development

### Build

Builds to `bin/couture`:

    make

### Install

Installs to `$HOME/bin/couture`:

    make install

### Entry Points

* [Command](cmd/couture.go) - CLI command.
* [Manager](internal/pkg/manager/manager.go) - Couture source/sink manager.
* [Source Definitions](internal/pkg/source/source.go) - Where sources are defined.
* [Sink Definitions](internal/pkg/sink/sink.go) - Where sinks are defined.

## TODO

_Migrate these into GitHub issues_

* Auto-complete for sources
* Testing
* Migrate poll to push and remove poll
* Recovery
* Customizable output
    * Column selection and order
    * Header
    * Prefixes
* Make a real README
* Make a real logo
* Determine license
* Customizable JSON schema? Are we conforming to some open logstash standard?
* Ensure no Gaggle-specific code
* Friendly pagination auto integration w/ less or something?
* JSON pretty dump mode (see [Chroma](https://github.com/alecthomas/chroma))
* Working releases / deployments
* CodePipeline
* Can we opensource?

## Future Ideas

* Read/Write flat files via [Abstract File System](https:// github.com/viant/afs)
* Kinesis stream
