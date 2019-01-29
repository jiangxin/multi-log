package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jiangxin/multi-log/formatter"
	"github.com/jiangxin/multi-log/path"
	"github.com/sirupsen/logrus"
)

// Options defines options to initial log
type Options struct {
	Quiet         bool
	Verbose       int
	LogRotateSize int64
	LogFile       string
	LogLevel      string
	ForceColors   bool

	stderr   io.Writer
	exitFunc func(int)
}

// Logger defines our custom basic logger interface
type Logger interface {
	Log(level logrus.Level, args ...interface{})
	Logf(level logrus.Level, format string, args ...interface{})
	Logln(level logrus.Level, args ...interface{})
}

// MultiLogger implements Logger interface. It wraps two loggers, one for console, one for file
type MultiLogger struct {
	StdLogger  *logrus.Logger
	FileLogger *logrus.Logger
	self       Logger
}

const (
	defaultLogRotateSize int64 = 20 * 1024 * 1024
	defaultLogLevel            = "warning"
)

var (
	mLogger = MultiLogger{}
	o       Options
)

// Init must run first to initialize logger
func Init(options Options) {
	var (
		logLevel   logrus.Level
		logFile    string
		err        error
		noExitFunc = func(code int) { return }
	)

	o = options
	if o.LogRotateSize == 0 {
		o.LogRotateSize = defaultLogRotateSize
	}
	if o.stderr == nil {
		o.stderr = os.Stderr
	}
	if o.LogLevel == "" {
		o.LogLevel = defaultLogLevel
	}

	switch o.Verbose {
	case 0:
		logLevel = logrus.WarnLevel
	case 1:
		logLevel = logrus.InfoLevel
	case 2:
		logLevel = logrus.DebugLevel
	default:
		logLevel = logrus.TraceLevel
	}

	mLogger.StdLogger = &logrus.Logger{
		Out: o.stderr,
		Formatter: &formatter.TextFormatter{
			DisableTimestamp:       true,
			FullTimestamp:          false,
			DisableLevelTruncation: true,
			ForceColors:            o.ForceColors,
		},

		Hooks:        make(logrus.LevelHooks),
		Level:        logLevel,
		ExitFunc:     noExitFunc,
		ReportCaller: false,
	}

	if o.LogFile != "" {
		logFile, err = path.Abs(o.LogFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: fail to resolve logfile: %s", err)
		}
	}
	if logFile != "" {
		logLevel = logrus.ErrorLevel
		if o.LogLevel != "" {
			logLevel, err = logrus.ParseLevel(o.LogLevel)
			if err != nil {
				logLevel = logrus.ErrorLevel
			}
		}

		dirname := filepath.Dir(logFile)
		if _, err = os.Stat(dirname); err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(dirname, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: cannot create dir %s for logging\n", dirname)
			}
		}

		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: cannot open %s for logging\n", logFile)
		} else {
			mLogger.FileLogger = &logrus.Logger{
				Out: file,
				Formatter: &formatter.TextFormatter{
					DisableTimestamp:       false,
					FullTimestamp:          true,
					DisableLevelTruncation: false,
				},
				Hooks:        make(logrus.LevelHooks),
				Level:        logLevel,
				ExitFunc:     noExitFunc,
				ReportCaller: false,
			}
		}
	}
}

// Self is used for class override.
// E.g. MultiLoggerWithFields overrides MultiLogger using Self()
func (v *MultiLogger) Self() Logger {
	if v.self != nil {
		return v.self
	}
	return v
}

// Logf is the base function to show message with specific log level
func (v *MultiLogger) Logf(level logrus.Level, format string, args ...interface{}) {
	if v.StdLogger != nil {
		v.StdLogger.Logf(level, format, args...)
	}

	if v.FileLogger != nil {
		v.FileLogger.Logf(level, format, args...)
	}
}

// Tracef is Logf with TraceLevel
func (v *MultiLogger) Tracef(format string, args ...interface{}) {
	v.Self().Logf(logrus.TraceLevel, format, args...)
}

// Debugf is Logf with DebugLevel
func (v *MultiLogger) Debugf(format string, args ...interface{}) {
	v.Self().Logf(logrus.DebugLevel, format, args...)
}

// Infof is Logf with InfoLevel
func (v *MultiLogger) Infof(format string, args ...interface{}) {
	v.Self().Logf(logrus.InfoLevel, format, args...)
}

// Warnf is Logf with WarnLevel
func (v *MultiLogger) Warnf(format string, args ...interface{}) {
	v.Self().Logf(logrus.WarnLevel, format, args...)
}

// Warningf is alias of Warnf
func (v *MultiLogger) Warningf(format string, args ...interface{}) {
	v.Warnf(format, args...)
}

// Errorf is Logf with ErrorLevel
func (v *MultiLogger) Errorf(format string, args ...interface{}) {
	v.Self().Logf(logrus.ErrorLevel, format, args...)
}

// Fatalf is Logf with FatalLevel
func (v *MultiLogger) Fatalf(format string, args ...interface{}) {
	v.Self().Logf(logrus.FatalLevel, format, args...)
	callExitFunc(1)
}

// Panicf is Logf with PanicLevel
func (v *MultiLogger) Panicf(format string, args ...interface{}) {
	// Using logrus.PanicLevel, will run panic and quit directly
	v.Self().Logf(logrus.PanicLevel, format, args...)
}

// Log is the base function to show message with specific log level
func (v *MultiLogger) Log(level logrus.Level, args ...interface{}) {
	if v.StdLogger != nil {
		v.StdLogger.Log(level, args...)
	}

	if v.FileLogger != nil {
		v.FileLogger.Log(level, args...)
	}
}

// Trace is Log with TraceLevel
func (v *MultiLogger) Trace(args ...interface{}) {
	v.Self().Log(logrus.TraceLevel, args...)
}

// Debug is Log with DebugLevel
func (v *MultiLogger) Debug(args ...interface{}) {
	v.Self().Log(logrus.DebugLevel, args...)
}

// Info is Log with InfoLevel
func (v *MultiLogger) Info(args ...interface{}) {
	v.Self().Log(logrus.InfoLevel, args...)
}

// Warn is Log with WarnLevel
func (v *MultiLogger) Warn(args ...interface{}) {
	v.Self().Log(logrus.WarnLevel, args...)
}

// Warning is alias of Warn
func (v *MultiLogger) Warning(args ...interface{}) {
	v.Warn(args...)
}

// Error is Log with ErrorLevel
func (v *MultiLogger) Error(args ...interface{}) {
	v.Self().Log(logrus.ErrorLevel, args...)
}

// Fatal is Log with FatalLevel
func (v *MultiLogger) Fatal(args ...interface{}) {
	v.Self().Log(logrus.FatalLevel, args...)
	callExitFunc(1)
}

// Panic is Log with PanicLevel
func (v *MultiLogger) Panic(args ...interface{}) {
	// Using logrus.PanicLevel, will run panic and quit directly
	v.Self().Log(logrus.PanicLevel, args...)
}

// Logln is the base function to show message with specific log level
func (v *MultiLogger) Logln(level logrus.Level, args ...interface{}) {
	if v.StdLogger != nil {
		v.StdLogger.Logln(level, args...)
	}

	if v.FileLogger != nil {
		v.FileLogger.Logln(level, args...)
	}
}

// Traceln is Logln with TraceLevel
func (v *MultiLogger) Traceln(args ...interface{}) {
	v.Self().Logln(logrus.TraceLevel, args...)
}

// Debugln is Logln with DebugLevel
func (v *MultiLogger) Debugln(args ...interface{}) {
	v.Self().Logln(logrus.DebugLevel, args...)
}

// Infoln is Logln with InfoLevel
func (v *MultiLogger) Infoln(args ...interface{}) {
	v.Self().Logln(logrus.InfoLevel, args...)
}

// Warnln is Logln with WarnLevel
func (v *MultiLogger) Warnln(args ...interface{}) {
	v.Self().Logln(logrus.WarnLevel, args...)
}

// Warningln is alias of Warnln
func (v *MultiLogger) Warningln(args ...interface{}) {
	v.Warnln(args...)
}

// Errorln is Logln with ErrorLevel
func (v *MultiLogger) Errorln(args ...interface{}) {
	v.Self().Logln(logrus.ErrorLevel, args...)
}

// Fatalln is Logln with FatalLevel
func (v *MultiLogger) Fatalln(args ...interface{}) {
	v.Self().Logln(logrus.FatalLevel, args...)
	callExitFunc(1)
}

// Panicln is Logln with PanicLevel
func (v *MultiLogger) Panicln(args ...interface{}) {
	// Using logrus.PanicLevel, will run panic and quit directly
	v.Self().Logln(logrus.PanicLevel, args...)
}

// Logf is the base function to show message with specific log level
func Logf(level logrus.Level, format string, args ...interface{}) {
	if mLogger.StdLogger != nil {
		mLogger.StdLogger.Logf(level, format, args...)
	}

	if mLogger.FileLogger != nil {
		mLogger.FileLogger.Logf(level, format, args...)
	}
}

// Tracef is Logf with TraceLevel
func Tracef(format string, args ...interface{}) {
	Logf(logrus.TraceLevel, format, args...)
}

// Debugf is Logf with DebugLevel
func Debugf(format string, args ...interface{}) {
	Logf(logrus.DebugLevel, format, args...)
}

// Infof is Logf with InfoLevel
func Infof(format string, args ...interface{}) {
	Logf(logrus.InfoLevel, format, args...)
}

// Warnf is Logf with WarnLevel
func Warnf(format string, args ...interface{}) {
	Logf(logrus.WarnLevel, format, args...)
}

// Warningf is alias of Warnf
func Warningf(format string, args ...interface{}) {
	Warnf(format, args...)
}

// Errorf is Logf with ErrorLevel
func Errorf(format string, args ...interface{}) {
	Logf(logrus.ErrorLevel, format, args...)
}

// Fatalf is Logf with FatalLevel
func Fatalf(format string, args ...interface{}) {
	Logf(logrus.FatalLevel, format, args...)
	callExitFunc(1)
}

// Panicf is Logf with PanicLevel
func Panicf(format string, args ...interface{}) {
	// Using logrus.PanicLevel, will run panic and quit directly
	Logf(logrus.PanicLevel, format, args...)
}

// Log is the base function to show message with specific log level
func Log(level logrus.Level, args ...interface{}) {
	if mLogger.StdLogger != nil {
		mLogger.StdLogger.Log(level, args...)
	}

	if mLogger.FileLogger != nil {
		mLogger.FileLogger.Log(level, args...)
	}
}

// Trace is Log with TraceLevel
func Trace(args ...interface{}) {
	Log(logrus.TraceLevel, args...)
}

// Debug is Log with DebugLevel
func Debug(args ...interface{}) {
	Log(logrus.DebugLevel, args...)
}

// Info is Log with InfoLevel
func Info(args ...interface{}) {
	Log(logrus.InfoLevel, args...)
}

// Warn is Log with WarnLevel
func Warn(args ...interface{}) {
	Log(logrus.WarnLevel, args...)
}

// Warning is alias of Warn
func Warning(args ...interface{}) {
	Warn(args...)
}

// Error is Log with ErrorLevel
func Error(args ...interface{}) {
	Log(logrus.ErrorLevel, args...)
}

// Fatal is Log with FatalLevel
func Fatal(args ...interface{}) {
	Log(logrus.FatalLevel, args...)
	callExitFunc(1)
}

// Panic is Log with PanicLevel
func Panic(args ...interface{}) {
	// if using logrus.PanicLevel, will run panic directly
	Log(logrus.PanicLevel, args...)
}

// Logln is the base function to show message with specific log level
func Logln(level logrus.Level, args ...interface{}) {
	if mLogger.StdLogger != nil {
		mLogger.StdLogger.Logln(level, args...)
	}

	if mLogger.FileLogger != nil {
		mLogger.FileLogger.Logln(level, args...)
	}
}

// Traceln is Logln with TraceLevel
func Traceln(args ...interface{}) {
	Logln(logrus.TraceLevel, args...)
}

// Debugln is Logln with DebugLevel
func Debugln(args ...interface{}) {
	Logln(logrus.DebugLevel, args...)
}

// Infoln is Logln with InfoLevel
func Infoln(args ...interface{}) {
	Logln(logrus.InfoLevel, args...)
}

// Warnln is Logln with WarnLevel
func Warnln(args ...interface{}) {
	Logln(logrus.WarnLevel, args...)
}

// Warningln is alias of Warnln
func Warningln(args ...interface{}) {
	Warnln(args...)
}

// Errorln is Logln with ErrorLevel
func Errorln(args ...interface{}) {
	Logln(logrus.ErrorLevel, args...)
}

// Fatalln is Logln with FatalLevel
func Fatalln(args ...interface{}) {
	Logln(logrus.FatalLevel, args...)
	callExitFunc(1)
}

// Panicln is Logln with PanicLevel
func Panicln(args ...interface{}) {
	// if using logrus.PanicLevel, will run panic directly
	Logln(logrus.PanicLevel, args...)
}

func callExitFunc(code int) {
	if o.exitFunc == nil {
		os.Exit(code)
	}
	o.exitFunc(code)
}

func init() {
	Init(Options{})
}
