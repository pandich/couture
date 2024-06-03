# FLOW

## § SOURCE

A source is an abstraction of a pull or push-drive event system.
Its only requirement is to emit textual events of a consistent
shape.

Not all sources support all configuration and command-line options.
For instance, it is not possible to get lookback history on `stdin`.

No detail usage is provided here. See `--help` for more information.

## Flat Files

Flat files are processed line-by-line. Each line must contain
structured or semi-structure data. Unstructured data is appears
as records with only a message field.

#### Local

Local flat files are read from `file://` URIs. Following by file
path or handle is supported if the OS supports it.

#### Pipe

Piped flat files are read from `stdin`.

#### Remote (SSH)

Remote flat files are read from `ssh://` URIs. The URI's target
must be either in DNS, a hostfile, or SSH configuration. In order
for SSH mode to work, a public key must be installed on the remote
`known_hosts` file. A number of ways exist to do this.

This source invokes a `tail -F` command on the remote host.

See: `man ssh-config`, `man ssh-keygen` and `man ssh-copy-id`.

### AWS Resource

AWS Resources are primary read from CloudWatch Logs. Examples exceptions
are S3 files, and CloudFormation status events.

#### S3

S3 files are read from `s3://` URIs.

#### CloudWatch

CloudWatch logs are read from `cloudwatch://`, `logs://` and `cw://` URIs, with additional
short forms for common AWS services:

| Short Form              | Expansion                             | Example                     |
|-------------------------|---------------------------------------|-----------------------------|
| `rdsc://`               | `cloudwatch://aws/rds/cluster`        | `rdsc://my-cluster`         |
| `rdsi://`               | `cloudwatch://aws/rds/instance`       | `rdsi://my-instance`        |
| `api://`                | `cloudwatch://aws/appsync/apis`       | `api://bookstore`           |

#### CloudFormation

CloudFormation resources are scanned. Any resource capable of producing a CloudWatch log is
added to a multi-tail group. Additionally, CloudFormatione events (e.g., `CREATE_COMPLETE`)
are emitted.

- [ ] `~/.user/config/couture/aliases.yaml`
	- [ ] Simple Alias
	- [ ] Alias to a group of sources
	- [ ] Expansion of children
