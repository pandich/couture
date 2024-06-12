COUTURE
=======

A purely CLI logging experience.


EXAMPLES
--------

## AWS

### AWS Profile and Region

```shell
couture cloudwatch:///aws/lambda/someLambda?profile=developer-alice&region=us-east-1
```

### CloudWatch Logs

```shell
couture cloudwatch:///aws/...your log group ...
```

### S3

```shell
 couture s3:///my-bucket/my-key
```

### CloudFormation

```shell
couture cloudformation:///my-stack
```

#### Short Forms

```shell
# cloudwatch
couture cw:///...
couture logs:///...
# /aws/lambda/...
couture cf:///my-stack
couture stack:///my-stack
couture lambda:///...
# /aws/appsync/apis/...
couture appsync:///...
# codebuild base
couture codebuild:///...
# /aws/rds/...
couture rds:///...
# /aws/rds/instance/...
couture rdsi:///...
# /aws/rds/cluster/...
couture rdsc:///...
```

### Further Examples

```shell
couture rdsi:///my-db-integration-imp/error
couture appsync://xyvwetblljgkahzew5w5pqeije
couture cf:///my-stack
```
