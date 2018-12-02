package app

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type Loggable interface {
	Debug(line string)
	Info(line string)
	Warn(line string)
	Error(line string)
	Debugf(fmt string, args ...interface{})
	Infof(fmt string, args ...interface{})
	Warnf(fmt string, args ...interface{})
	Errorf(fmt string, args ...interface{})
}

type Logger struct {
	LogFileHandle *os.File
	MirrorStdout  bool
	IsDebugLevel  bool
	IsInfoLevel   bool
	IsWarnLevel   bool
	IsErrorLevel  bool
	ThreadSafe    *sync.Mutex
}

func NewLogger(cmdVerbose string, logFileHandle *os.File) *Logger {

	verbose, verboseErr := strconv.Atoi(cmdVerbose)
	if verboseErr != nil {
		panic(fmt.Sprintf("Invalid verbose setting: %v", verboseErr))
	} else if verbose < 0 || verbose > 5 {
		panic(fmt.Sprintf("Invalid verbose setting. Must be between 1 and 5."))
	}

	l := new(Logger)
	l.ThreadSafe = new(sync.Mutex)

	// Set to all to false
	l.IsDebugLevel = false
	l.IsInfoLevel = false
	l.IsWarnLevel = false
	l.IsErrorLevel = false

	// Set to true based on verbose level (0:Quiet,1:ERROR,2:WARN,3:INFO,4:DEBUG,5:ALL)
	switch verbose {
	case 0:
		break
	case 1:
		l.IsErrorLevel = true
		break
	case 2:
		l.IsErrorLevel = true
		l.IsWarnLevel = true
		break
	case 3:
		l.IsErrorLevel = true
		l.IsWarnLevel = true
		l.IsInfoLevel = true
		break
	case 4:
		l.IsErrorLevel = true
		l.IsWarnLevel = true
		l.IsInfoLevel = true
		l.IsDebugLevel = true
		break
	case 5:
		l.IsErrorLevel = true
		l.IsWarnLevel = true
		l.IsInfoLevel = true
		l.IsDebugLevel = true
		break
	}

	// TODO: Figure out how we want to implment this better
	// Unless we are 'quietmode' we echo to STDOUT
	// l.MirrorStdout = !config.QuietMode
	l.LogFileHandle = logFileHandle

	return l
}

func (l *Logger) Write(level string, line string) {
	s := time.Now().UTC().Format("2006-01-02T15:04:05.000") + " [" + level + "] " + line

	l.ThreadSafe.Lock()
	defer l.ThreadSafe.Unlock()

	fmt.Fprintln(l.LogFileHandle, s)
	if l.MirrorStdout && (l.LogFileHandle != os.Stderr) {
		fmt.Fprintln(os.Stdout, s)
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.IsDebugLevel {
		line := fmt.Sprintf(format, args...)
		l.Debug(line)
	}
	return
}

func (l *Logger) Debug(line string) {
	if l.IsDebugLevel {
		l.Write("DEBUG", line)
	}
	return
}

func (l *Logger) Info(line string) {
	if l.IsInfoLevel {
		l.Write("INFO", line)
	}
	return
}
func (l *Logger) Infof(format string, args ...interface{}) {
	if l.IsInfoLevel {
		line := fmt.Sprintf(format, args...)
		l.Write("INFO", line)
	}
	return
}

func (l *Logger) Warn(line string) {
	if l.IsWarnLevel {
		l.Write("WARN", line)
	}
	return
}
func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.IsWarnLevel {
		line := fmt.Sprintf(format, args...)
		l.Write("WARN", line)
	}
	return
}

func (l *Logger) Error(line string) {
	if l.IsErrorLevel {
		l.Write("ERROR", line)
	}
	return
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	if l.IsErrorLevel {
		line := fmt.Sprintf(format, args...)
		l.Write("ERROR", line)
	}
	return
}
