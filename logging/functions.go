package logging

import (
	"fmt"
	"os"
)

var log interface { //*syslog.Writer
	Err(m string) (err error)
	Crit(m string) (err error)
	Emerg(m string) (err error)
	Debug(m string) (err error)
	Info(m string) (err error)
	Notice(m string) (err error)
}

// Err logs error level messages to gomates default logger
func Err(msg string) {
	log.Err(msg)
}

// Errf logs formatted error level messages to gomates default logger
func Errf(msg string, args ...interface{}) {
	Err(fmt.Sprintf(msg, args...))
}

// Debug logs debug level messages to gomates default logger
func Debug(msg string) {
	log.Debug(msg)
}

// Debugf logs formatted debug level messages to gomates default logger
func Debugf(msg string, args ...interface{}) {
	Debug(fmt.Sprintf(msg, args...))
}

// Info logs info level messages to gomates default logger
func Info(msg string) {
	log.Info(msg)
}

// Infof logs formatted info level messages to gomates default logger
func Infof(msg string, args ...interface{}) {
	Info(fmt.Sprintf(msg, args...))
}

// Notice logs notice level messages to gomates default logger
func Notice(msg string) {
	log.Notice(msg)
}

// Noticef logs formatted notice level messages to gomates default logger
func Noticef(msg string, args ...interface{}) {
	Notice(fmt.Sprintf(msg, args...))
}

// Panicf logs panic level messages to gomates default logger
func Panicf(reason string, args ...interface{}) {
	Panic(fmt.Sprintf(reason, args...))
}

// Panic logs formatted panic level messages to gomates default logger
func Panic(reason string) {
	log.Crit(reason)
	panic(reason)
}

// Fatal logs fatal level messages to gomates default logger
func Fatal(reason string) {
	log.Crit(reason)
	os.Exit(1)
}

// Fatalf logs formatted fatal level messages to gomates default logger
func Fatalf(reason string, args ...interface{}) {
	Fatal(fmt.Sprintf(reason, args...))
}
