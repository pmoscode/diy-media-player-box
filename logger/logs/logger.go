package logs

import (
	"fmt"
	"github.com/pmoscode/go-common/mqtt"
	"github.com/pmoscode/go-common/shutdown"
	log2 "log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type Log struct {
	lock                        sync.Mutex
	filename                    string
	logRotationPeriodAfterBytes int
	logFile                     *os.File
	byteCounter                 int
}

func (l *Log) Log(message mqtt.Message) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.byteCounter > l.logRotationPeriodAfterBytes {
		l.rotateLogFile()
		l.byteCounter = 0
	}

	statusMessage := &mqtt.StatusPublishMessage{}
	message.ToStruct(statusMessage)

	logLine := fmt.Sprintf("[%s] # %s # %s -> %s\n", statusMessage.Timestamp, message.Topic, statusMessage.Type, statusMessage.Status)
	if l.logFile == nil {
		log2.Println("No log file open. Trying to rotate...")
		l.rotateLogFile()
	}

	_, err := l.logFile.WriteString(logLine)
	if err != nil {
		return
	}

	l.byteCounter += len(logLine)
}

func (l *Log) rotateLogFile() {
	if l.logFile != nil {
		err := l.logFile.Close()
		l.logFile = nil
		if err != nil {
			log2.Printf("Could not close logfile '%s' for reason:", l.filename)
			log2.Fatal(err)
		}
	}

	_, err := os.Stat(l.filename)
	if err == nil {
		now := strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "_")
		archiveFilename := fmt.Sprintf("%s_%s", l.filename, now)
		err = os.Rename(l.filename, archiveFilename)
		if err != nil {
			log2.Printf("Could not rename logfile '%s' to '%s' for reason:", l.filename, archiveFilename)
			log2.Fatal(err)
		}
	}

	if _, err := os.Stat(l.filename); os.IsNotExist(err) {
		pathStr := path.Dir(l.filename)
		err := os.MkdirAll(pathStr, 0755)
		if err != nil {
			log2.Printf("Could not create directories '%s' for reason:", pathStr)
			log2.Fatal(err)
		}
	}

	file, err := os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log2.Printf("Could not open logfile '%s' for reason:", l.filename)
		log2.Fatal(err)
	}

	l.logFile = file
}

func (l *Log) init() {
	l.rotateLogFile()

	shutdown.GetObserver().AddCommand(func() error {
		err := l.logFile.Sync()
		if err != nil {
			return err
		}

		err = l.logFile.Close()
		if err != nil {
			return err
		}

		log2.Println("Writing log and exiting...")

		return nil
	})
}
