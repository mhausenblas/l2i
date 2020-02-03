# The Lambda Layer Inspector (L2I)

This CLI tool allows to inspect a [AWS Lambda Layer](https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html) given its ARN.
For example, see the ones listed on [mthenw/awesome-layers](https://github.com/mthenw/awesome-layers).

```sh
$ l2i arn:aws:lambda:eu-west-1:553035198032:layer:git:10
```

```sh
$ l2i arn:aws:lambda:eu-west-1:464622532012:layer:Datadog-Python37:1
```