# The Lambda Layer Inspector (L2I)

This CLI tool allows to inspect one or more AWS Lambda [layer(s)](https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html), using the layer ARN(s) as input.
For examples of layer ARNs, see the ones listed on [mthenw/awesome-layers](https://github.com/mthenw/awesome-layers).

## Install

Download the [latest binary](https://github.com/mhausenblas/l2i/releases/latest) for Linux (Intel or Arm), macOS, or Windows.

For example, to install `l2i` from binary on macOS you could do the following:

```sh
curl -L https://github.com/mhausenblas/l2i/releases/latest/download/l2i_darwin_amd64.tar.gz \
    -o l2i.tar.gz && \
    tar xvzf l2i.tar.gz l2i && \
    mv l2i /usr/local/bin && \
    rm l2i*
```

## Use

If you want to inspect a single AWS Lambda layer, simply provide the ARN to `l2i` like so:

```sh
$ l2i arn:aws:lambda:eu-west-1:553035198032:layer:git:10
Name: git
Version: 10
Description: Git 2.25.0 and openssh binaries
Created on: 2020-01-13T20:41:57.917+0000
Size: 17,456 kB
Location: https://awslambda-eu-west-1-layers.s3.eu-west-1.amazonaws.com/snapshots/553035198032/git-c86b3b6b-1ff4-48e2-bdc3-3721ae076147?versionId=YhboGnC0BP6h5jlTaS2jUxyeZxXFBQU3
```

For multiple layers, `l2i` will provide a tabular overview:

```sh
$ l2i arn:aws:lambda:eu-west-1:464622532012:layer:Datadog-Python37:1 \
      arn:aws:lambda:eu-west-1:553035198032:layer:git:10
NAME              VERSION  DESCRIPTION                      CREATED ON                    SIZE (kB)
Datadog-Python37  1        Datadog Lambda Layer for Python  2019-05-06T18:48:17.694+0000  7,657
git               10       Git 2.25.0 and openssh binaries  2020-01-13T20:41:57.917+0000  17,456
```