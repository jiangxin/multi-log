package log

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jiangxin/multi-log/formatter"
	"github.com/sirupsen/logrus"
)

var (
	mu sync.Mutex // lock for printf family methods
)

// If not quiet, always show note message on console
func print(prefix string, args ...interface{}) {
	if o.Quiet {
		return
	}

	msg := sprint(prefix, args...)

	l := mLogger.StdLogger
	mu.Lock()
	fmt.Fprint(l.Out, msg)
	defer mu.Unlock()
}

// If quite, return empty string, otherwize the returned string always ends with "\n".
func sprint(prefix string, args ...interface{}) string {
	var (
		msg string
	)

	switch strings.ToLower(prefix) {
	case "note":
		if o.Quiet {
			return ""
		}
	default:
		level, err := logrus.ParseLevel(prefix)
		if err == nil && !mLogger.StdLogger.IsLevelEnabled(level) {
			return ""
		}
	}

	f, ok := mLogger.StdLogger.Formatter.(*formatter.TextFormatter)
	if !ok {
		f = new(formatter.TextFormatter)
	}

	colorSet, colorReset := f.GetColors(prefix)
	prefix = strings.ToUpper(prefix)
	if !f.DisableLevelTruncation {
		prefix = prefix[0:4]
	}

	if f.DisableTimestamp {
		msg += fmt.Sprintf("%s%s:%s ",
			colorSet,
			prefix,
			colorReset,
		)
	} else if !f.FullTimestamp {
		msg += fmt.Sprintf("%s%s[%04d]:%s ",
			colorSet,
			prefix,
			int(time.Now().Sub(formatter.BaseTimestamp)/time.Second),
			colorReset,
		)
	} else {
		msg += fmt.Sprintf("%s%s[%s]:%s ",
			colorSet,
			prefix,
			time.Now().Format(f.TimestampFormat),
			colorReset,
		)
	}

	msg += fmt.Sprint(args...)
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	return msg
}

// Note will show message on console only if not quiet.
func Note(args ...interface{}) {
	print("NOTE", args...)
}

// Notef is printf version of Note
func Notef(format string, args ...interface{}) {
	print("NOTE", fmt.Sprintf(format, args...))
}

// Noteln is println version of Note
func Noteln(args ...interface{}) {
	print("NOTE", fmt.Sprintln(args...))
}

// Print is alias of Note
func Print(args ...interface{}) {
	print("NOTE", args...)
}

// Printf is alias of Notef
func Printf(format string, args ...interface{}) {
	print("NOTE", fmt.Sprintf(format, args...))
}

// Println is alias of Noteln
func Println(args ...interface{}) {
	print("NOTE", fmt.Sprintln(args...))
}

// Snote will return the output message to display as note
func Snote(args ...interface{}) string {
	return sprint("NOTE", args...)
}

// Snotef is printf version of Snote
func Snotef(format string, args ...interface{}) string {
	return sprint("NOTE", fmt.Sprintf(format, args...))
}

// Snoteln is println version of Snote
func Snoteln(args ...interface{}) string {
	return sprint("NOTE", fmt.Sprintln(args...))
}

// Sprint is alias of Snote
func Sprint(args ...interface{}) string {
	return sprint("NOTE", args...)
}

// Sprintf is alias of Snotef
func Sprintf(format string, args ...interface{}) string {
	return sprint("NOTE", fmt.Sprintf(format, args...))
}

// Sprintln is alias of Snoteln
func Sprintln(args ...interface{}) string {
	return sprint("NOTE", fmt.Sprintln(args...))
}
