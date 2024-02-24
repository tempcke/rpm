package log

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tempcke/rpm/internal/configs"
)

type Fields = map[string]any

var (
	logger          logrus.FieldLogger
	defaultLogLevel = logrus.InfoLevel

	EnvAppEnv   = "APP_ENV"
	EnvAppName  = "APP_NAME"
	EnvLogLevel = "LOG_LEVEL"
)

// UseLogger sets the logger to be used
// be aware that doing this will make any logLevel set in Setup
// to be ignored, you are expected to set it in the logger
// passed in
func UseLogger(fieldLogger logrus.FieldLogger) {
	logger = fieldLogger
}

// Logger returns the logger
func Logger() logrus.FieldLogger {
	if logger == nil {
		logger = defaultLogger()
	}
	return logger
}
func Entry(fields ...Fields) *logrus.Entry {
	entry := Logger().WithFields(Fields{
		EnvAppEnv:  conf().GetString(EnvAppEnv),
		EnvAppName: conf().GetString(EnvAppName),
	})
	for _, f := range fields {
		entry = entry.WithFields(f)
	}
	return entry
}

func Debug(msg string, fields ...Fields) { Entry(fields...).Debug(msg) }
func Error(msg string, fields ...Fields) { Entry(fields...).Error(msg) }
func Fatal(args ...interface{}) {
	Entry().Fatal(args...)
}

func WithError(err error) *logrus.Entry {
	return Entry().WithError(err)
}
func WithField(key string, value interface{}) *logrus.Entry {
	return Entry().WithField(key, value)
}

func TerminalLogger() logrus.FieldLogger {
	l := logrus.StandardLogger()
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	l.SetLevel(logLevel(logrus.ErrorLevel))
	return l
}

func defaultLogger() logrus.FieldLogger {
	l := logrus.StandardLogger()
	l.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     true,
	}
	l.SetLevel(logLevel(defaultLogLevel))
	return l
}
func logLevel(fallback logrus.Level) logrus.Level {
	if ll := conf().GetString(EnvLogLevel); ll != "" {
		if lvl, err := logrus.ParseLevel(ll); err == nil {
			return lvl
		}
	}
	return fallback
}
func conf() configs.Config {
	return configs.New()
}
