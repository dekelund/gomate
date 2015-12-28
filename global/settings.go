package global

import "log/syslog"

func init() {
	Settings.SysLog.Priority = syslog.LOG_INFO
}

var Settings struct {
	SysLog struct {
		Active   bool
		UDP      bool
		RAddr    string
		Tag      string
		Priority syslog.Priority
	}
	Forensic   bool
	PPrint     bool
	CWD        string
	DefPattern string
}
