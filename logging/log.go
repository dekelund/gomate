package logging

import "fmt"

import "log/syslog"

type logger struct {
	settings Settings
}

func init() {
	log = &logger{}
}

// ReconfigureLogger updated default logger and change priority and tags,
// it might also change destination, for instance start logging to syslog-server
// at remote host via TCP or UDP.
func ReconfigureLogger(settings Settings) {
	if !settings.Active {
		log = &logger{settings}
	} else if settings.RAddr == "localhost" {
		log, _ = syslog.New(
			settings.Priority,
			settings.Tag,
		)
	} else {
		// Connection to remote address required
		network := "tcp"

		if settings.UDP {
			network = "udp"
		}

		log, _ = syslog.Dial(
			network, settings.RAddr,
			settings.Priority, settings.Tag,
		)
	}
}

func prioPrint(priority syslog.Priority, m string, settings Settings) {
	if priority <= settings.Priority {
		fmt.Println(m)
	}
}

func (l *logger) Err(m string) (err error) {
	prioPrint(syslog.LOG_ERR, m, l.settings)
	return
}

func (l *logger) Crit(m string) (err error) {
	prioPrint(syslog.LOG_CRIT, m, l.settings)
	return
}

func (l *logger) Emerg(m string) (err error) {
	prioPrint(syslog.LOG_EMERG, m, l.settings)
	return
}

func (l *logger) Debug(m string) (err error) {
	prioPrint(syslog.LOG_DEBUG, m, l.settings)
	return
}

func (l *logger) Info(m string) (err error) {
	prioPrint(syslog.LOG_INFO, m, l.settings)
	return
}

func (l *logger) Notice(m string) (err error) {
	prioPrint(syslog.LOG_NOTICE, m, l.settings)
	return
}
