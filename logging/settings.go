package logging

import "log/syslog"

// Settings contains information about logging level,
// tags and destination address.
type Settings struct {
	Active   bool
	UDP      bool
	RAddr    string
	Tag      string
	Priority syslog.Priority
}
