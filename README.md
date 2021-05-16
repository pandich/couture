# Couture

![Couture](https://static.thenounproject.com/png/566246-200.png)

Combine multiple log streams from disparate sources.

    couture --level=debug \
            --include='\d+ trousers?' --include='[A-Z]{3}-slacks' \
            --excdlue='(green|brown) culotte' \
            --since=1h \
            'es+http://logging.example.zz:9200/log/?application=pants-service' \
            'file:///var/log/shirts-service/logstash.log' \
            'cloudformation://clothing-service-stack?region=us-west-2&profile=production' \
            'cloudwatch:///aws/lambda/monitor-lambda' \
            'lambda://suits-lambda'

## Overview

Allows for tailing multiple of event sources.

For usage run `couture --help`

For shell completions run `eval $(couture complete)`

To build and install Couture into `$GOPATH` run `make install`. (See [Makefile](Makefile))

##### Important Files

* [Makefile](Makefile)
* [CLI Command](cmd/couture.go) - CLI command.
* [Event Source](internal/pkg/source/source.go) - Where sources are defined.
* [Even Sink](internal/pkg/sink/sink.go) - Where sinks are defined.
* [Source(s) -> Sink Manager](internal/pkg/manager/manager.go) - Couture source/sink manager.
