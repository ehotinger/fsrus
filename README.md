# fsrus

[![Build Status](https://travis-ci.org/ehotinger/fsrus.svg?branch=master)](https://travis-ci.org/ehotinger/fsrus)

A simple filesystem hook for [sirupsen/logrus](https://github.com/sirupsen/logrus).
Log to files with logrus.

## Installing

```golang
go get github.com/ehotinger/fsrus
```

## Building

Build with a simple `make` or `make build`. Invoke `make help` for more information.

### Testing:

`make test`

## Examples

```golang
logger := logrus.New()
logger.Level = logrus.DebugLevel
hook, err := fsrus.NewFilesystemHook(nil, "output.txt", nil, nil)
logger.AddHook(hook)
logger.Info(expectedMsg)
```