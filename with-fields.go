package log

import (
	"github.com/sirupsen/logrus"
)

// MultiLoggerWithFields extend MultiLogger with fields
type MultiLoggerWithFields struct {
	MultiLogger
	Fields logrus.Fields
	self   Logger
}

// WithFields writes log with fields
func WithFields(fields map[string]interface{}) *MultiLoggerWithFields {
	logger := new(MultiLoggerWithFields)
	logger.MultiLogger = mLogger
	logger.Fields = fields
	logger.MultiLogger.self = logger
	return logger
}

// WithField writes log with only one field
func WithField(key string, value interface{}) *MultiLoggerWithFields {
	return WithFields(logrus.Fields{key: value})
}

// Log defines core log methods for MultiLoggerWithFields
func (v *MultiLoggerWithFields) Log(level logrus.Level, args ...interface{}) {
	if v.StdLogger != nil {
		v.StdLogger.WithFields(v.Fields).Log(level, args...)
	}

	if v.FileLogger != nil {
		v.FileLogger.WithFields(v.Fields).Log(level, args...)
	}
}

// Logf defines core log methods for MultiLoggerWithFields
func (v *MultiLoggerWithFields) Logf(level logrus.Level, format string, args ...interface{}) {
	if v.StdLogger != nil {
		v.StdLogger.WithFields(v.Fields).Logf(level, format, args...)
	}

	if v.FileLogger != nil {
		v.FileLogger.WithFields(v.Fields).Logf(level, format, args...)
	}
}

// Logln defines core log methods for MultiLoggerWithFields
func (v *MultiLoggerWithFields) Logln(level logrus.Level, args ...interface{}) {
	if v.StdLogger != nil {
		v.StdLogger.WithFields(v.Fields).Logln(level, args...)
	}

	if v.FileLogger != nil {
		v.FileLogger.WithFields(v.Fields).Logln(level, args...)
	}
}
