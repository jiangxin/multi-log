package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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

	stderr   io.Writer
	exitFunc func(int)
}

// MultiLogger implements Logger interface. It wraps two loggers, one for console, one for file
type MultiLogger struct {
	StdLogger  *logrus.Logger
	FileLogger *logrus.Logger
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
		Formatter: &logrus.TextFormatter{
			DisableTimestamp:       true,
			FullTimestamp:          false,
			DisableLevelTruncation: true,
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
				Formatter: &logrus.TextFormatter{
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
