package global

import "fmt"

import "log/syslog"

type logger struct {
}

func init() {
	log = &logger{}
}

func ReconfigureLogger() {
	if Settings.SysLog.Active == false {
		log = &logger{}
	} else if Settings.SysLog.RAddr == "localhost" {
		log, _ = syslog.New(
			Settings.SysLog.Priority,
			Settings.SysLog.Tag,
		)
	} else {
		// Connection to remote address required
		network := "tcp"

		if Settings.SysLog.UDP {
			network = "udp"
		}

		log, _ = syslog.Dial(
			network, Settings.SysLog.RAddr,
			Settings.SysLog.Priority, Settings.SysLog.Tag,
		)
	}
}

func prioPrint(priority syslog.Priority, m string) {
	if priority <= Settings.SysLog.Priority {
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
