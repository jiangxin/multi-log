package log

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jiangxin/multi-log/formatter"
)

var (
	mu sync.Mutex // lock for printf family methods
)

// If not quiet, always show note message on console
func note(prefix string, showPrefix bool, args ...interface{}) {
	if o.Quiet {
		return
	}

	l := mLogger.StdLogger
	f, ok := l.Formatter.(*formatter.TextFormatter)
	if !ok {
		f = new(formatter.TextFormatter)
	}

	mu.Lock()
	defer mu.Unlock()

	colorSet, colorReset := f.GetColors(prefix)
	prefix = strings.ToUpper(prefix)
	if !f.DisableLevelTruncation {
		prefix = prefix[0:4]
	}

	if showPrefix {
		if f.DisableTimestamp {
			fmt.Fprintf(l.Out, "%s%s:%s ",
				colorSet,
				prefix,
				colorReset,
			)
		} else if !f.FullTimestamp {
			fmt.Fprintf(l.Out, "%s%s[%04d]:%s ",
				colorSet,
				prefix,
				int(time.Now().Sub(formatter.BaseTimestamp)/time.Second),
				colorReset,
			)
		} else {
			fmt.Fprintf(l.Out, "%s%s[%s]:%s ",
				colorSet,
				prefix,
				time.Now().Format(f.TimestampFormat),
				colorReset,
			)
		}
	}

	msg := fmt.Sprint(args...)
	msg = strings.TrimSuffix(msg, "\n")
	fmt.Fprintln(l.Out, msg)
}

// Note will show message on console only if not quiet.
func Note(args ...interface{}) {
	note("NOTE", true, args...)
}

// Notef is printf version of Note
func Notef(format string, args ...interface{}) {
	note("NOTE", true, fmt.Sprintf(format, args...))
}

// Noteln is println version of Note
func Noteln(args ...interface{}) {
	note("NOTE", true, fmt.Sprintln(args...))
}

// Print is alias of Note
func Print(args ...interface{}) {
	note("NOTE", true, args...)
}

// Printf is alias of Notef
func Printf(format string, args ...interface{}) {
	note("NOTE", true, fmt.Sprintf(format, args...))
}

// Println is alias of Noteln
func Println(args ...interface{}) {
	note("NOTE", true, fmt.Sprintln(args...))
}
