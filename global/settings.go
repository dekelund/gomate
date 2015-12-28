package global

import "log/syslog"

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
