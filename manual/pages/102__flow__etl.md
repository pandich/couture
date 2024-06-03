# FLOW

## § EVENT IDENTIFICATION AND MAPPING

Events are identified by the comparing the first event in the stream
to a series of predicates mapped to log formats. The first predicate
to match has its log format applied to the stream. If no predicate
matches, the stream is considered unstructured.

#### JSON?

If JSON is detected, the stream is considered structured and each
mapper is asked if it recognizes the JSON structure.

##### Example

```yaml
logstash:
  format: json
  priority: 100
  predicates:
    "@version": "^1$"
  mapping:
    timestamp: "@timestamp"
    level: level
    application: application
    context: thread_name
    entity: class
    action: method
    line: line_number
    message: message
    error: exception.stacktrace
```

#### Semi-Structured?

All other events are now candidates for semi-structured data datection
and mapping. Each mapper applies one or more regular expressions to
detect the structure of the event.

##### Example

```yaml
aws-billing-report:
  format: text
  priority: 90
  predicates:
    _: "^(?P<message>REPORT (?P<entity>RequestId):\\s+(?P<action>\\S+)\\s+.+)$"
aws-billing-start:
  format: text
  priority: 90
  predicates:
    _: "^(?P<message>START (?P<entity>RequestId):(?P<action>\\s+\\S+).+)$"
aws-billing-end:
  format: text
  priority: 90
  predicates:
    _: "^(?P<message>END (?P<entity>RequestId):\\s+(?P<action>\\S+)\\s+.+)$"
```

#### Unstructured

If no mapper recognizes the event, it is considered unstructured and sent to the
unknown event-type channel. This channel can be dumped to `stderr` by using the
undocumented `--dump-unknown` command.

### Filtering

- [ ] Filter: --filter="XXX"
  - [ ] Positive
  - [ ] Negative
  - [ ] Trigger
- [ ] Highlight: --highlight="XXX"
  - [ ] Enabled
  - [ ] Disabled
  - [ ] Includes
  - [ ] $COUTURE_HIGHLIGHT
- [ ] Log Level
  - [ ] Gaggle standard files
  - [ ] AWS common files
  - [ ] Files with no levels
  - [ ] $COUTURE_LEVEL
- [ ] Since: --since="XXX"
  - [ ] Time
  - [ ] Duration
  - [ ] Human
  - [ ] Sources without time filtering
- [ ] `~/.user/config/couture/mappings.yaml`
  - [ ] Name
  - [ ] Format
  - [ ] Priority
  - [ ] PredicatesByField
  - [ ] FieldByColumn
  - [ ] TemplateByColumn
  - [ ] Fields
  - [ ] TextPattern
  - [ ] One script per format identification type
