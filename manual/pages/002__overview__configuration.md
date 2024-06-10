COUTURE
=======

OVERVIEW
--------

## CONFIGURATION

### Defaults

	`$HOME/.config/couture/config.yaml`

```yaml
auto_resize: true
color: true
columns:
	- source
	- timestamp
	- application
	- context
	# unified class name, method name, line number
	- caller
	- level
	- message

	# method name
	#  - action

	# line number
	#  - line

	# class/type/struct name
	#  - entity

	# error message
	#  - error

consistent_colors: true
expand: false
highlight: false
multi_line: false
level_meter: false
theme: prince
time_format: human
# automatic by default
# width: 80 				 
``` 

### Aliases

Aliass and group aliases allow for URL shortening. Each alias is actually a
mustache template with the following implicit and user-definable values.

An alias is references as `alias://name` (or `@name`) where `name` is the alias name.
Groups are referenced `@@name`  where `name` is the group name.

| Variable                 | Description                                        |
|--------------------------|----------------------------------------------------|
| `epoch`                  | Epoch millis.                                      |
| `yyyy`                   | 4-digit year.                                      |
| `yy`                     | 2-digit year.                                      |
| `mm`                     | 2-digit month.                                     |
| `m`                      | 1- or 2-digit month.                               |
| `dd`                     | 2-digit day of month.                              |
| `d`                      | 1- or 2-digit day of month.                        |
| `hh`                     | 2-digit hour from 00 to 23.                        |
| `h`                      | 1-digit hour from 0 to 23.                         |
| `MM`                     | 2-digit minute 00 to 59.                           |
| `M`                      | 1- oe 2-digit minute 0 to 59.                      |
| `ss`                     | 2-digit second 00 to 59.                           |
| `s`                      | 1- oe 2-digit second 0 to 59.                      |
| `_name`                  | Anything passed into the URL's host.               |
| `_path`                  | Any additional URL path.                           |
| `_user`                  | Anything passed into the URL's username parameter. |
| `_password`              | Anything passed into the URL's password parameter. |
| `$COUTURE_CONTEXT_{xyz}` | The value `xyz`.                                   | 

`$HOME/.config/couture/aliases.yaml`

```yaml
aliases:
	simulate-1: "lambda://{{shortenv}}-couture-simulator-Simulate-1"
	simulate-2: "lambda://{{shortenv}}-couture-simulator-Simulate-2"
	simulate-3: "lambda://{{shortenv}}-couture-simulator-Simulate-3"
	simulate-4: "lambda://{{shortenv}}-couture-simulator-Simulate-4"
	simulate-5: "lambda://{{shortenv}}-couture-simulator-Simulate-5"
	simulate-6: "lambda://{{shortenv}}-couture-simulator-Simulate-6"
	simulate-7: "lambda://{{shortenv}}-couture-simulator-Simulate-7"
	simulate-8: "lambda://{{shortenv}}-couture-simulator-Simulate-8"
	simulate-9: "lambda://{{shortenv}}-couture-simulator-Simulate-9"
	bob: "cf:///couture-simulator-Stack"

groups:
	simulate: [ "simulate-1", "simulate-2", "simulate-3", "simulate-4", "simulate-5", "simulate-6", "simulate-7", "simulate-8", "simulate-9" ]

```
