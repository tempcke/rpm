package log

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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

func Entry() *logrus.Entry {
	return Logger().
		WithField(EnvAppEnv, viper.GetString(EnvAppEnv)).
		WithField(EnvAppName, viper.GetString(EnvAppName))
}

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
	if ll := viper.GetString(EnvLogLevel); ll != "" {
		if lvl, err := logrus.ParseLevel(ll); err == nil {
			return lvl
		}
	}
	return fallback
}
