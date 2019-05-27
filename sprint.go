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

// Stracef is sprint with TraceLevel
func Stracef(format string, args ...interface{}) string {
	return sprint("trace", fmt.Sprintf(format, args...))
}

// Sdebugf is sprint with DebugLevel
func Sdebugf(format string, args ...interface{}) string {
	return sprint("debug", fmt.Sprintf(format, args...))
}

// Sinfof is sprint with InfoLevel
func Sinfof(format string, args ...interface{}) string {
	return sprint("info", fmt.Sprintf(format, args...))
}

// Swarnf is sprint with WarnLevel
func Swarnf(format string, args ...interface{}) string {
	return sprint("warn", fmt.Sprintf(format, args...))
}

// Swarningf is alias of Warnf
func Swarningf(format string, args ...interface{}) string {
	return sprint("warn", fmt.Sprintf(format, args...))
}

// Serrorf is sprint with ErrorLevel
func Serrorf(format string, args ...interface{}) string {
	return sprint("error", fmt.Sprintf(format, args...))
}

// Strace is sprint with TraceLevel
func Strace(args ...interface{}) string {
	return sprint("trace", args...)
}

// Sdebug is sprint with DebugLevel
func Sdebug(args ...interface{}) string {
	return sprint("debug", args...)
}

// Sinfo is sprint with InfoLevel
func Sinfo(args ...interface{}) string {
	return sprint("info", args...)
}

// Swarn is sprint with WarnLevel
func Swarn(args ...interface{}) string {
	return sprint("warn", args...)
}

// Swarning is alias of Warn
func Swarning(args ...interface{}) string {
	return sprint("warn", args...)
}

// Serror is sprint with ErrorLevel
func Serror(args ...interface{}) string {
	return sprint("error", args...)
}

// Straceln is sprint with TraceLevel
func Straceln(args ...interface{}) string {
	return sprint("trace", fmt.Sprintln(args...))
}

// Sdebugln is sprint with DebugLevel
func Sdebugln(args ...interface{}) string {
	return sprint("debug", fmt.Sprintln(args...))
}

// Sinfoln is sprint with InfoLevel
func Sinfoln(args ...interface{}) string {
	return sprint("info", fmt.Sprintln(args...))
}

// Swarnln is sprint with WarnLevel
func Swarnln(args ...interface{}) string {
	return sprint("warn", fmt.Sprintln(args...))
}

// Swarningln is alias of Swarnln
func Swarningln(args ...interface{}) string {
	return sprint("warn", fmt.Sprintln(args...))
}

// Serrorln is sprint with ErrorLevel
func Serrorln(args ...interface{}) string {
	return sprint("error", fmt.Sprintln(args...))
}
