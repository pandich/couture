# Example Scripts TODO

## Cofiguration

### Source URL Aliases

* [ ] `~/.user/config/couture/mappings.yaml`
	* [ ] Name
	* [ ] Format
	* [ ] Priority
	* [ ] PredicatesByField
	* [ ] FieldByColumn
	* [ ] TemplateByColumn
	* [ ] Fields
	* [ ] TextPattern
	* [ ] One script per format identification type 
* [ ] `~/.user/config/couture/aliases.yaml`
	* [ ] Simple Alias
	* [ ] Alias to a group of sources
	* [ ] Expansion of children

## Features

### Help

* [ ] Confirm that the help text is correct
* [ ] Confirm examples, and asciicinema examples are correct
* [ ] Updated README.md
* [ ] Updated diagrams
* [ ] Updated example commands
* [ ] Example config and alias files

## Miscellanea

* [ ] Rate limit
	* [ ] --rate-limit=rate-limit
	* [ ] $COUTURE_RATE_LIMIT
* [ ] Info dumper
	* [ ] ADD: --show-themes
	* [ ] ADD: --show-aliases
	* [ ] ADD: --show-url-schemes
	* [ ] ADD: --show-sources
	* [ ] ADD: --show-defaults
		* [ ] With config file
	  	* [ ] Without config file 
	* [ ] --show-mappings

### Formatting

* [ ] Auto-resize: --[no-]auto-resize
	* [ ] Enabled
	* [ ] Disabled
	* [ ] $COUTURE_AUTO_RESIZE
* [ ] Color Mode: --color-mode="auto"
	* [ ] Enabled
	* [ ] Disabled
	* [ ] $COUTURE_COLOR_MODE
* [ ] Force TTY: --tty
	* [ ] Enabled
	* [ ] Disabled
* [ ] Wrap: --[no-]wrap
	* [ ] Enabled
	* [ ] Disabled
	* [ ] $COUTURE_WRAP
* [ ] Wrap Width: --width=width
	* [ ] Enabled
	* [ ] Disabled
	* [ ] Too narrow
	* [ ] Too wide
	* [ ] $COUTURE_WIDTH

### Styling

* [ ] Color Mode: --color-mode="auto"
	* [ ] Auto
	* [ ] Dark
	* [ ] Light
* [ ] Consistent Colors: --[no-]consistent-colors
	* [ ] Enabled
	* [ ] Disabled
	* [ ] $COUTURE_CONSISTENT_COLORS
* [ ] Expansion: --[no-]expand
	* [ ] JSON
	* [ ] $COUTURE_EXPAND
* [ ] Level Meter: --[no-]level-meter
	* [ ] Enabled
	* [ ] Disabled
	* [ ] $COUTURE_LEVEL_METER
* [ ] Multi-line: --[no-]multi-line
	* [ ] Enabled
	* [ ] Disabled
	* [ ] With Expand
	* [ ] With Wrap
	* [ ] With Width
* [ ] Source Style: --source-style="XXX"
	* [ ] Named
	* [ ] Custom
	* [ ] Is there a command to list them?
* [ ] Theme: --theme="XXX"
	* [ ] Named
	* [ ] Custom
	* [ ] Is there a command to list them?
	* [ ] $COUTURE_THEME

### Filtering

* [ ] Filter: --filter="XXX"
	* [ ] Positive
	* [ ] Negative
	* [ ] Trigger
* [ ] Highlight: --highlight="XXX"
	* [ ] Enabled
	* [ ] Disabled
	* [ ] Includes
	* [ ] $COUTURE_HIGHLIGHT
* [ ] Log Level
	* [ ] Gaggle standard files
	* [ ] AWS common files
	* [ ] Files with no levels
	* [ ] $COUTURE_LEVEL
* [ ] Since: --since="XXX"
	* [ ] Time
	* [ ] Duration
	* [ ] Human
	* [ ] Sources without time filtering

### Content

* [ ] Column: --column="XXX"
	* [ ] Timestamp
	* [ ] Level
	* [ ] Message
	* [ ] Application
	* [ ] Action
	* [ ] Line
	* [ ] Context
	* [ ] Entity
	* [ ] Error
	* [ ] $COUTURE_COLUMN_NAMES
* [ ] Time Format: --time-format="XXX"
	* [ ] Named
	* [ ] Custom
	* [ ] Is there a command to list them?
	* [ ] $COUTURE_TIME_FORMAT

## Sources

### CloudWatch

* [ ] Short format
* [ ] Friendly format
* [ ] Long format
* [ ] With profile
* [ ] With region
* [ ] With lookbackTime
* [ ] With subsystem

### CloudFormation

* [ ] Short format
* [ ] Friendly format
* [ ] Long format
* [ ] With profile
* [ ] With region
* [ ] With lookbackTime
* [ ] With resource log discovery for:
	* [ ] apigateway
	* [ ] appsync
	* [ ] cloudwatch
	* [ ] codebuild
	* [ ] codecommit
	* [ ] codepipeline
	* [ ] deploy
	* [ ] ecs
	* [ ] lambda
	* [ ] rds
	* [ ] redshift
	* [ ] xray

### ElasticSearch

* [ ] Short format
* [ ] Friendly format
* [ ] Long format
* [ ] HTTP
* [ ] HTTPS

### S3

* [ ] Key
* [ ] Prefix
* [ ] Bucket

### Local File

* [ ] Path
* [ ] Filename rotates
	* [ ] Follow file by descriptor
	* [ ] Follow file by name

### Remote Files

* [ ] SSH Tunnel

## General

* [ ] Test the pre-log-initialzed state error logging
* [ ] goreleaser
