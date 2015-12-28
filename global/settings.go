package global

import "log/syslog"

var Settings struct {
	SysLog      bool
	LogPriority syslog.Priority
	Debug       bool
	PPrint      bool
	CWD         string
	DefPattern  string
}
