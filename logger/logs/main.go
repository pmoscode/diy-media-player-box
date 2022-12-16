package logs

import (
	"fmt"
	"gitlab.com/pmoscodegrp/common/mqtt"
	log2 "log"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"
)

type log struct {
	lock                        sync.Mutex
	filename                    string
	logRotationPeriodAfterBytes int
	logFile                     *os.File
	byteCounter                 int
}

func (l *log) Log(message mqtt.Message) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.byteCounter > l.logRotationPeriodAfterBytes {
		l.rotateLogFile()
		l.byteCounter = 0
	}

	statusMessage := &mqtt.StatusPublishMessage{}
	message.ToStruct(statusMessage)

	logLine := fmt.Sprintf("[%s] # %s # %s -> %s\n", statusMessage.Timestamp, message.Topic, statusMessage.Type, statusMessage.Status)
	if l.logFile != nil {
		l.logFile.WriteString(logLine)
		l.byteCounter += len(logLine)
	} else {
		log2.Println("No log file open. Trying to rotate...")
		l.rotateLogFile()
	}
}

func (l *log) rotateLogFile() {
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
		os.MkdirAll(path.Dir(l.filename), 0755) // Create your file
	}
	file, err := os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log2.Printf("Could not open logfile '%s' for reason:", l.filename)
		log2.Fatal(err)
	}

	l.logFile = file
}

func (l *log) init() {
	l.rotateLogFile()

	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		l.logFile.Sync()
		l.logFile.Close()
		log2.Println("Writing log and exiting...")
		os.Exit(1)
	}()
}

func New(filename string, filesize int) *log {
	logger := &log{
		logRotationPeriodAfterBytes: filesize,
		filename:                    filename,
		logFile:                     nil,
		byteCounter:                 0,
	}

	logger.init()

	return logger
}
