package formatter

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	purple = 35
	cyan   = 36
	gray   = 37
)

// Default key names for the default fields
const (
	defaultTimestampFormat = time.RFC3339
)

var baseTimestamp time.Time

func init() {
	baseTimestamp = time.Now()
}

// TextFormatter formats logs into text
type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors.
	DisableColors bool

	// Override coloring based on CLICOLOR and CLICOLOR_FORCE. - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// QuoteEmptyFields will wrap empty fields in quotes if true
	QuoteEmptyFields bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	terminalInitOnce sync.Once
}

func (f *TextFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = checkIfTerminal(entry.Logger.Out)

		if f.isTerminal {
			initTerminal(entry.Logger.Out)
		}
	}

	if f.TimestampFormat == "" {
		f.TimestampFormat = defaultTimestampFormat
	}
}

// IsColored checks capability for color output
func (f *TextFormatter) IsColored() bool {
	isColored := f.ForceColors || (f.isTerminal && (runtime.GOOS != "windows"))

	if f.EnvironmentOverrideColors {
		if force, ok := os.LookupEnv("CLICOLOR_FORCE"); ok && force != "0" {
			isColored = true
		} else if ok && force == "0" {
			isColored = false
		} else if os.Getenv("CLICOLOR") == "0" {
			isColored = false
		}
	}

	return isColored && !f.DisableColors
}

// Format renders a single log entry
func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.terminalInitOnce.Do(func() { f.init(entry) })

	f.printEntry(b, entry)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// GetColors gets colors for set and reset
func (f *TextFormatter) GetColors(levelName string) (colorSet, colorReset string) {
	var (
		levelColor int
		level      logrus.Level
		err        error
	)

	if !f.IsColored() {
		return
	}

	level, err = logrus.ParseLevel(levelName)
	if err == nil {
		switch level {
		case logrus.DebugLevel, logrus.TraceLevel:
			levelColor = blue
		case logrus.WarnLevel:
			levelColor = yellow
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			levelColor = red
		case logrus.InfoLevel:
			levelColor = cyan
		default:
			levelColor = gray
		}
	} else {
		switch strings.ToLower(levelName) {
		case "note":
			levelColor = cyan
		default:
			levelColor = gray
		}
	}

	colorSet = fmt.Sprintf("\x1b[1;%dm", levelColor)
	colorReset = "\x1b[0m"
	return
}

func (f *TextFormatter) printEntry(b *bytes.Buffer, entry *logrus.Entry) {
	var (
		levelText            string
		caller               string
		colorSet, colorReset string
	)

	colorSet, colorReset = f.GetColors(entry.Level.String())

	levelText = strings.ToUpper(entry.Level.String())
	if !f.DisableLevelTruncation {
		levelText = levelText[0:4]
	}

	// Remove a single newline if it already exists in the message to keep
	// the behavior of logrus text_formatter the same as the stdlib log package
	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	if entry.HasCaller() {
		caller = fmt.Sprintf("%s:%d %s() ",
			entry.Caller.File, entry.Caller.Line, entry.Caller.Function)
	}

	w := 0
	if len(entry.Data) > 0 {
		w = -44
	}
	if f.DisableTimestamp {
		fmt.Fprintf(b, "%s%s:%s %s%*s",
			colorSet,
			levelText,
			colorReset,
			caller,
			w,
			entry.Message,
		)
	} else if !f.FullTimestamp {
		fmt.Fprintf(b, "%s%s[%04d]:%s %s%*s",
			colorSet,
			levelText,
			int(entry.Time.Sub(baseTimestamp)/time.Second),
			colorReset,
			caller,
			w,
			entry.Message,
		)
	} else {
		fmt.Fprintf(b, "%s%s[%s]:%s %s%*s",
			colorSet,
			levelText,
			entry.Time.Format(f.TimestampFormat),
			colorReset,
			caller,
			w,
			entry.Message,
		)
	}

	if len(entry.Data) > 0 {
		var keys []string
		for k := range entry.Data {
			keys = append(keys, k)
		}
		if !f.DisableSorting {
			sort.Strings(keys)
		}
		if f.IsColored() {
			for _, k := range keys {
				fmt.Fprintf(b, " %s%s%s=", colorSet, k, colorReset)
				f.appendValue(b, entry.Data[k])
			}
		} else {
			b.WriteString(" (")
			i := 0
			for _, k := range keys {
				fmt.Fprintf(b, "%s=%s", k, entry.Data[k])
				if i != len(entry.Data)-1 {
					b.WriteString(" ")
				}
				i++
			}
			b.WriteString(")")
		}
	}
}

// FormatMessagef formats one message
func FormatMessagef(logLevel logrus.Level, prompt, format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	return fmtMessage(logLevel, prompt, msg)
}

// FormatMessageln formats one message
func FormatMessageln(logLevel logrus.Level, prompt string, args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return fmtMessage(logLevel, prompt, msg)
}

// FormatMessage formats one message
func FormatMessage(logLevel logrus.Level, prompt string, args ...interface{}) string {
	msg := fmt.Sprint(args...)
	return fmtMessage(logLevel, prompt, msg)
}

func fmtMessage(logLevel logrus.Level, prompt, msg string) string {
	var levelColor int

	if logLevel < 0 {
		levelColor = 0
	} else {
		switch logLevel {
		case logrus.DebugLevel, logrus.TraceLevel:
			levelColor = gray
		case logrus.WarnLevel:
			levelColor = yellow
		case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
			levelColor = red
		default:
			levelColor = blue
		}
	}

	if levelColor <= 0 {
		return fmt.Sprintf("%s: %s", prompt, msg)
	}

	if prompt == "" {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", levelColor, msg)
	}
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m: %s", levelColor, prompt, msg)
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}
