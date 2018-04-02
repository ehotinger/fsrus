package fsrus_test

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/ehotinger/fsrus"
	"github.com/sirupsen/logrus"
)

// TestWritingToLog_DefaultPath tests writing all levels to a default path.
func TestWritingToLog_DefaultPath(t *testing.T) {
	expectedMsg := "This is expected output"

	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	tmpFile, err := ioutil.TempFile("", "test.txt")
	fileName := tmpFile.Name()

	defer func() {
		tmpFile.Close()
		os.Remove(fileName)
	}()

	if err != nil {
		t.Errorf("Failed to create tmp file. Err: %v", err)
	}

	hook, err := fsrus.NewFilesystemHook(nil, fileName, nil, nil)
	if err != nil {
		t.Errorf("Failed to create file system hook: %v", err)
	}

	logger.AddHook(hook)

	// Log an expected message and a message which should be filtered out.
	logger.Info(expectedMsg)
	logger.Debug(expectedMsg)
	logger.Warning(expectedMsg)
	logger.Error(expectedMsg)

	scanner := bufio.NewScanner(tmpFile)
	scanner.Split(bufio.ScanLines)

	numReceivedMsgs := 0
	for scanner.Scan() {
		curr := scanner.Text()
		if strings.Contains(curr, expectedMsg) {
			numReceivedMsgs++
		} else {
			t.Errorf("Unexpected message: %s", curr)
		}
	}

	if numReceivedMsgs != 4 {
		t.Errorf("Didn't receive all of the expected messages. Received: %d", numReceivedMsgs)
	}
}

// TestWritingToLog_LevelPathMap tests writing a log entry to
// a file using level paths
func TestWritingToLog_LevelPathMap(t *testing.T) {
	expectedMsg := "This is the expected output"
	filteredMsg := "This message should be filtered"

	logger := logrus.New()
	logger.Level = logrus.DebugLevel

	tmpFile, err := ioutil.TempFile("", "test.txt")
	fileName := tmpFile.Name()

	defer func() {
		tmpFile.Close()
		os.Remove(fileName)
	}()

	if err != nil {
		t.Errorf("Failed to create tmp file. Err: %v", err)
	}

	levelPathMap := fsrus.LevelPathMap{
		logrus.InfoLevel: fileName,
	}

	hook, err := fsrus.NewFilesystemHook(levelPathMap, "", nil, nil)
	if err != nil {
		t.Errorf("Failed to create file system hook: %v", err)
	}

	logger.AddHook(hook)

	// Log an expected message and a message which should be filtered out.
	logger.Info(expectedMsg)
	logger.Debug(filteredMsg)
	logger.Warning(filteredMsg)
	logger.Error(filteredMsg)

	scanner := bufio.NewScanner(tmpFile)
	scanner.Split(bufio.ScanLines)

	receivedMsg := false
	for scanner.Scan() {
		curr := scanner.Text()
		if strings.Contains(curr, expectedMsg) {
			receivedMsg = true
			t.Log("Received the expected message")
		} else {
			t.Errorf("Unexpected message: %s", curr)
		}
	}

	if !receivedMsg {
		t.Errorf("Didn't receive the expected message")
	}
}
