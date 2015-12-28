package global

import "fmt"
import "log/syslog"

// FIXME: Don't expose struct
type logger struct {
	syslog.Priority
}

func init() {
	log = &logger{syslog.LOG_INFO}
}

func ReconfigureLogger() {
	if Settings.SysLog {
		log, _ = syslog.New(Settings.LogPriority, "unbrokenwing")
	} else {
		log = &logger{Settings.LogPriority}
	}
}

func prioPrint(priority syslog.Priority, m string) {
	if priority <= Settings.LogPriority {
		fmt.Println(m)
	}
}

func (l *logger) Err(m string) (err error) {
	prioPrint(syslog.LOG_ERR, m)
	return
}

func (l *logger) Crit(m string) (err error) {
	prioPrint(syslog.LOG_CRIT, m)
	return
}

func (l *logger) Emerg(m string) (err error) {
	prioPrint(syslog.LOG_EMERG, m)
	return
}

func (l *logger) Debug(m string) (err error) {
	prioPrint(syslog.LOG_DEBUG, m)
	return
}

func (l *logger) Info(m string) (err error) {
	prioPrint(syslog.LOG_INFO, m)
	return
}

func (l *logger) Notice(m string) (err error) {
	prioPrint(syslog.LOG_NOTICE, m)
	return
}
