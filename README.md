# The Lambda Layer Inspector (L2I)

This CLI tool allows to inspect one or more AWS Lambda [layer(s)](https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html), using the layer ARN(s) as input.

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

Once installed, all you need to inspect one or more Lambda layer are their ARNs.
For examples of layer ARNs, check out [mthenw/awesome-layers](https://github.com/mthenw/awesome-layers).

### Plain metadata inspection

If you want to inspect a single AWS Lambda layer, provide the ARN to `l2i` using
the `--layers` parameter like so:

```sh
$ l2i --layers arn:aws:lambda:eu-west-1:553035198032:layer:git:10
Name: git
Version: 10
Description: Git 2.25.0 and openssh binaries
Created on: 2020-01-13T20:41:57.917+0000
Size: 17,456 kB
Location: https://awslambda-eu-west-1-layers.s3.eu-west-1.amazonaws.com/snapshots/553035198032/git-c86b3b6b-1ff4-48e2-bdc3-3721ae076147?versionId=YhboGnC0BP6h5jlTaS2jUxyeZxXFBQU3

$ l2i --layers arn:aws:lambda:eu-west-1:464622532012:layer:Datadog-Python37:1
Name: Datadog-Python37
Version: 1
Description: Datadog Lambda Layer for Python
Created on: 2019-05-06T18:48:17.694+0000
Size: 7,657 kB
Location: https://awslambda-eu-west-1-layers.s3.eu-west-1.amazonaws.com/snapshots/464622532012/Datadog-Python37-44fdbec2-a76c-4fb4-a3e2-be9781f035be?versionId=4vdBqcXS41OnGhs5Ml16_G31mK2mBaqM
CompatibleRuntimes: python3.7
```

You can also inspect multiple layers at once, using a comma-separated list and
`l2i` will provide a tabular overview:

```sh
$ l2i --layers "arn:aws:lambda:eu-west-1:464622532012:layer:Datadog-Python37:1,arn:aws:lambda:eu-west-1:553035198032:layer:git:10"
NAME              VERSION  DESCRIPTION                      CREATED ON                    SIZE (kB)
Datadog-Python37  1        Datadog Lambda Layer for Python  2019-05-06T18:48:17.694+0000  7,657
git               10       Git 2.25.0 and openssh binaries  2020-01-13T20:41:57.917+0000  17,456
```

### Content inspection

If you provide the `--export` parameter, `l2i` will not only display metadata of
a layer but also download its contents into the (relative or 
absolute) provided path, for example:

```sh
$ l2i --layers arn:aws:lambda:eu-west-1:553035198032:layer:git:10 --export .
Name: git
Version: 10
Description: Git 2.25.0 and openssh binaries
Created on: 2020-01-13T20:41:57.917+0000
Size: 17,456 kB
Location: https://awslambda-eu-west-1-layers.s3.eu-west-1.amazonaws.com/snapshots/553035198032/git-c86b3b6b-1ff4-48e2-bdc3-3721ae076147?versionId=YhboGnC0BP6h5jlTaS2jUxyeZxXFBQU3

I exported the layer's contents to: /Users/janedoe/serverless/layer-git-content
```

If you now look into the `layer-git-content` directory you should see something
like this:

```
$ tree -d layer-git-content/
layer-git-content/
├── bin
├── etc
│   └── ssh
├── lib
│   └── fipscheck
├── libexec
│   ├── git-core
│   │   └── mergetools
│   └── openssh
└── share
    ├── git-core
    │   └── templates
    └── licenses
        ├── fipscheck-1.3.1
        ├── git-2.25.0
        ├── openssh-7.4p1
        └── pcre2-10.21

17 directories
```

You can also export multiple layers at once, like so:

```sh
$ l2i --layers "arn:aws:lambda:eu-west-1:464622532012:layer:Datadog-Python37:1,arn:aws:lambda:eu-west-1:553035198032:layer:nodejs10:20" --export layers
NAME              VERSION  DESCRIPTION                      CREATED ON                    SIZE (kB)
Datadog-Python37  1        Datadog Lambda Layer for Python  2019-05-06T18:48:17.694+0000  7,657
nodejs10          20       Node.js v10.18.1 custom runtime  2020-01-10T21:32:36.590+0000  12,404

I exported the layers' contents to: /Users/janedoe/serverless/layers/
```

Have fun!