# Couture

![Couture](https://static.thenounproject.com/png/566246-200.png)

## Overview

_TODO_

### Execution

    bin/couture

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

---

* [Source Definitions](internal/pkg/source/source.go) - Where sources are defined.
* [Source CLI Definition](cmd/cli/source.go) - Where CLI arguments and parsing is defined for a source.
* [Source Config Definition](cmd/config/source.go) - Where configuration setup and parsing is defined for a source.

---

* [Sink Definitions](internal/pkg/sink/sink.go) - Where sinks are defined.
* [Sink CLI Definition](cmd/cli/sink.go) - Where CLI arguments and parsing is defined for a sink.
* [Sink Config Definition](cmd/config/sink.go) - Where configuration setup and parsing is defined for a sink.

## Libraries Used

* CLI Parsing: [Kong](https://github.com/alecthomas/kong)
* Configuration: [Configuro](https://github.com/sherifabdlnaby/configuro)
* Event Bus: [EventBus](https://github.com/asaskevich/EventBus)
* AWS: [AWS SDK](https://github.com/aws/aws-sdk-go)
* Colorization: [Aurora](https://github.com/logrusorgru/aurora)
* Terminal Management: [Termdash](https://github.com/mum4k/termdash)

## Limitations

* Can't tail resources from multiple AWS regions or profiles simultaneously.

## Future Ideas

* AWS CodePipeline events?
* ElasticSearch
* Read/Write flat files via [Abstract File System](https://github.com/viant/afs)?
