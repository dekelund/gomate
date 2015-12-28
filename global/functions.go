package global

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

func Err(msg string) {
	log.Err(msg)
}

func Errf(msg string, args ...interface{}) {
	Err(fmt.Sprintf(msg, args...))
}

func Debug(msg string) {
	log.Debug(msg)
}

func Debugf(msg string, args ...interface{}) {
	Debug(fmt.Sprintf(msg, args...))
}

func Info(msg string) {
	log.Info(msg)
}

func Infof(msg string, args ...interface{}) {
	Info(fmt.Sprintf(msg, args...))
}

func Notice(msg string) {
	log.Notice(msg)
}

func Noticef(msg string, args ...interface{}) {
	Notice(fmt.Sprintf(msg, args...))
}

func Panicf(reason string, args ...interface{}) {
	Panic(fmt.Sprintf(reason, args...))
}

func Panic(reason string) {
	log.Crit(reason)
	panic(reason)
}

func Fatal(reason string) {
	log.Crit(reason)
	os.Exit(1)
}

func Fatalf(reason string, args ...interface{}) {
	Fatal(fmt.Sprintf(reason, args...))
}
