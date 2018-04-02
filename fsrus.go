// Package fsrus is a file system hook for logrus.
package fsrus

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
)

var defaultTextFormatter = &logrus.TextFormatter{DisableColors: true}
var defaultLevels = logrus.AllLevels

// FsHook represents a filesystem hook.
type FsHook struct {
	lock             *sync.Mutex
	formatter        logrus.Formatter
	levelPathMap     LevelPathMap
	defaultLevelPath string
	levels           []logrus.Level
}

// LevelPathMap maps a file path to a logrus log level.
type LevelPathMap map[logrus.Level]string

// NewFilesystemHook returns a filesystem hook.
// A valid level path map or default path is required.
func NewFilesystemHook(lMap LevelPathMap, defaultPath string, lvls []logrus.Level, fmt logrus.Formatter) (*FsHook, error) {
	hook := &FsHook{
		lock: &sync.Mutex{},
	}

	if lMap == nil && len(defaultPath) <= 0 {
		return hook, errors.New("A default path or a level map must be provided")
	}

	hook.SetLevelPathMap(lMap)
	hook.SetDefaultLevelPath(defaultPath)
	hook.SetLevels(lvls)
	hook.SetFormatter(fmt)

	return hook, nil
}

// SetFormatter sets the filesystem hook's logrus formatter.
// If the provided formatter is nil, it will default to a text formatter
// with colors disabled.
func (hook *FsHook) SetFormatter(formatter logrus.Formatter) {
	if formatter == nil {
		formatter = defaultTextFormatter
	}

	hook.formatter = formatter
}

// SetLevels sets the hook's levels to the specified array of logrus levels.
// These levels determine when the hook should fire.
func (hook *FsHook) SetLevels(levels []logrus.Level) {
	if levels == nil {
		hook.levels = defaultLevels
	} else {
		hook.levels = levels
	}
}

// SetLevelPathMap sets the hook's level path map.
func (hook *FsHook) SetLevelPathMap(levelPathMap LevelPathMap) {
	hook.levelPathMap = levelPathMap
}

// SetDefaultLevelPath sets the default path for all levels.
func (hook *FsHook) SetDefaultLevelPath(defaultLevelPath string) {
	if len(defaultLevelPath) > 0 {
		hook.defaultLevelPath = defaultLevelPath
	}
}

// Fire writes the specified log entry to a file.
func (hook *FsHook) Fire(entry *logrus.Entry) error {
	return hook.writeToFile(entry)
}

// Levels defines which levels should fire the hook.
func (hook *FsHook) Levels() []logrus.Level {
	return hook.levels
}

// writeToFile writes the specified entry to a file
// using the hook's underlying formatter.
func (hook *FsHook) writeToFile(entry *logrus.Entry) error {
	hook.lock.Lock()
	defer hook.lock.Unlock()

	levelPath, ok := hook.levelPathMap[entry.Level]
	if !ok {
		if hook.defaultLevelPath != "" {
			levelPath = hook.defaultLevelPath
		} else {
			return nil
		}
	}

	// Try to make the directory for the corresponding level path if it
	// doesn't exist already and chmod 777.
	dir := filepath.Dir(levelPath)
	_ = os.MkdirAll(dir, os.ModePerm)

	// TODO: DESIGN: Some use-cases may want O_SYNC
	file, err := os.OpenFile(levelPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Printf("Failed to open file. Path: %s, Err: %v\n", levelPath, err)
		return err
	}
	defer file.Close()

	msg, err := hook.formatter.Format(entry)

	if err != nil {
		log.Printf("Failed to format entry. Err: %v\n", err)
		return err
	}

	file.Write(msg)
	return nil
}
