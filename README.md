# fsrus

A simple filesystem hook for [sirupsen/logrus](https://github.com/sirupsen/logrus).
Log to files with logrus.

## Examples:


```golang
logger := logrus.New()
logger.Level = logrus.DebugLevel
hook, err := fsrus.NewFilesystemHook(nil, "output.txt", nil, nil)
logger.AddHook(hook)
logger.Info(expectedMsg)
```